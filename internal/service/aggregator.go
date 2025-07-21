package service

import (
	"RSSHub/internal/adapters/httpadapter"
	"RSSHub/internal/adapters/repo"
	"RSSHub/internal/domain/models"
	"RSSHub/pkg/logger"
	"RSSHub/pkg/utils"
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	ErrConfigNotFound        = errors.New("RSS config not found")
	ErrProcessAlreadyRunning = errors.New("background process already running")
	ErrFailedToReadConfig    = errors.New("failed to read config")
	ErrFailedToUpdateStatus  = errors.New("failed to update aggregator status")
)

// RssAggregator is the main service that manages the RSS feed aggregation process.
// It handles the configuration, ticker, and worker controllers.
type RssAggregator struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	log     logger.Logger
	cleanDb func()

	articleRepo *repo.ArticleRepo
	feedRepo    *repo.FeedRepo
	configRepo  *repo.ConfigRepo

	tc *TickerController
	wc *WorkerController
}

func NewRssAggregator(articleRepo *repo.ArticleRepo, feedRepo *repo.FeedRepo, configRepo *repo.ConfigRepo, log logger.Logger, cleanDb func()) *RssAggregator {
	ctx, cancel := context.WithCancel(context.Background())
	return &RssAggregator{
		ctx:         ctx,
		cancel:      cancel,
		cleanDb:     cleanDb,
		log:         log,
		articleRepo: articleRepo,
		feedRepo:    feedRepo,
		configRepo:  configRepo,
	}
}

// Start launches the RSS aggregator, ticker controller, and worker controller with the configuration parameters.
func (a *RssAggregator) Start(ctx context.Context) error {
	const op = "RssAggregator.Start"
	log := a.log.GetSlogLogger().With("op", op)

	// Reading configuration
	cfg, err := a.GetConfig(ctx)
	if err != nil {
		log.Error("failed to load config", "error", err)
		return err
	}

	// Checking if process is already runnning
	if cfg.Run {
		return ErrProcessAlreadyRunning
	}

	// Updating application status
	if err := a.configRepo.UpdateRunStatus(ctx, true); err != nil {
		log.Error("Failed to update aggregator status", "error", err)
		return ErrFailedToUpdateStatus
	}

	a.wc = NewWorkerController(cfg.WorkerCount, a.log)

	// Initialize rss fetcher.
	rssFetcher := httpadapter.NewClient(time.Second * 5)
	a.tc = NewTickerController(cfg.TimerInterval, a.feedRepo, a.articleRepo, rssFetcher, a.log)

	a.wg.Add(2)
	go a.tc.Run(a.ctx, &a.wg, a.wc)
	go a.wc.Run(a.ctx, cfg.WorkerCount, &a.wg)

	a.wg.Add(2)
	go a.intervalUpdater(a.ctx, cfg.TimerInterval)
	go a.countUpdater(a.ctx, cfg.WorkerCount)

	msg := fmt.Sprintf("The background process for fetching feeds has started (interval = %s, workers = %d)", utils.PrettyDuration(cfg.TimerInterval), cfg.WorkerCount)
	a.log.Notify(msg)
	a.ListenShutdown(a.ctx)
	return nil
}

// Stop gracefully shuts down the RSS aggregator and updates the run status in the repository.
func (a *RssAggregator) Stop() error {
	const op = "RssAggregator.Stop"
	log := a.log.GetSlogLogger().With("op", op)

	a.cancel()
	a.wg.Wait()

	if err := a.configRepo.UpdateRunStatus(context.Background(), false); err != nil {
		log.Error("Failed to change app status", "error", err)
		return fmt.Errorf("%s:%w", op, ErrFailedToUpdateStatus)
	}

	msg := "Graceful shutdown: aggregator stopped"
	a.log.Notify(msg)
	return nil
}

// ---------------- TickerController ----------------

type RssFetcher interface {
	FetchRSSFeed(ctx context.Context, url string) (*models.RSSFeed, error)
}

type TickerController struct {
	t           *VarTicker
	intervalCh  chan time.Duration
	feedRepo    *repo.FeedRepo
	articleRepo *repo.ArticleRepo
	rssFethcer  RssFetcher
	log         logger.Logger
}

func NewTickerController(interval time.Duration, feedRepo *repo.FeedRepo, articleRepo *repo.ArticleRepo, rssFethcer RssFetcher, log logger.Logger) *TickerController {
	return &TickerController{
		t:           NewVarTicker(interval),
		intervalCh:  make(chan time.Duration, 1),
		feedRepo:    feedRepo,
		articleRepo: articleRepo,
		rssFethcer:  rssFethcer,
		log:         log,
	}
}

// Starts the ticker loop, periodically fetching stale feeds and dispatching jobs to the worker controller.
func (c *TickerController) Run(ctx context.Context, wg *sync.WaitGroup, wc *WorkerController) {
	defer wg.Done()
	defer c.t.Stop()
	defer close(wc.wp.jobCh)

	for {
		select {
		case <-ctx.Done():
			c.log.Debug(ctx, "ticker controller has been stopped")
			return
		case <-c.t.ticker.C:
			c.processFeeds(ctx, wc)
		case newInterval := <-c.intervalCh:
			oldInterval := c.t.GetDuration()
			c.t.Reset(newInterval)

			msg := fmt.Sprintf("Interval of fetching feeds changed from %s to %s", utils.PrettyDuration(oldInterval), utils.PrettyDuration(newInterval))
			c.log.Notify(msg)
		}
	}
}

// processFeeds retrieves stale feeds and submits jobs for each feed to fetch and save RSS items.
func (c *TickerController) processFeeds(ctx context.Context, wc *WorkerController) {
	feeds, err := c.feedRepo.GetStaleFeeds(ctx, c.t.GetDuration())
	if err != nil {
		c.log.Error(ctx, "Failed to get stale feeds", "error", err)
		return
	}
	if feeds == nil {
		c.log.Error(ctx, "feeds list is empty")
		return
	}
	for _, feed := range feeds {
		wc.SubmitJob(func() {
			fetched, err := c.rssFethcer.FetchRSSFeed(ctx, feed.URL)
			if err != nil {
				c.log.Error(ctx, "Failed to fetch RSS feed", "feed_URL", feed.URL, "error", err)
				return
			}
			if len(fetched.Channel.Item) == 0 {
				c.log.Error(ctx, "There is no items in the feed", "feed_URL", feed.URL)
				return
			}

			if err := c.articleRepo.CreateOrUpdate(ctx, feed.ID, fetched.Channel.Item); err != nil {
				c.log.Error(ctx, "Failed to save feed items", "feed_id", feed.ID, "articles", fetched.Channel.Item, "error", err)
				return
			}

			if err := c.feedRepo.UpdateUpdatedAt(ctx, feed.Name); err != nil {
				c.log.Error(ctx, "Failed to update updated_at", "error", err)
				return
			}
		})
	}

}

// ---------------- WorkerController -----------------

type WorkerController struct {
	wp      *WorkerParty
	countCh chan int
	log     logger.Logger
}

func NewWorkerController(initialCount int, log logger.Logger) *WorkerController {
	wc := &WorkerController{
		wp:      NewWorkerParty(),
		countCh: make(chan int),
		log:     log,
	}
	return wc
}

// Run starts the worker party and listens for changes in worker count or context cancellation for graceful shutdown.
func (wc *WorkerController) Run(ctx context.Context, initCount int, wg *sync.WaitGroup) {
	defer wg.Done()
	go wc.wp.Start(ctx)
	wc.wp.Scale(initCount)

	for {
		select {
		case <-ctx.Done():
			wc.log.Debug(ctx, "worker controller has been stopped")
			return
		case newCount := <-wc.countCh:
			wc.wp.Scale(newCount)
		}
	}
}

func (wc *WorkerController) SubmitJob(job func()) {
	wc.wp.jobCh <- job
}

// ---------------- Updaters ----------------

// intervalUpdater periodically checks for updated timer intervals in the configuration and updates the ticker if needed.
func (a *RssAggregator) intervalUpdater(ctx context.Context, currentInterval time.Duration) {
	defer a.wg.Done()

	t := time.NewTicker(time.Second * 2)
	defer t.Stop()
	defer close(a.tc.intervalCh)

	for {
		select {
		case <-ctx.Done():
			a.log.Debug(ctx, "interval updater has been stopped")
			return
		case <-t.C:
			cfg, err := a.configRepo.Get(ctx)
			if err != nil {
				a.log.Error(ctx, "Failed to read config", "error", err)
				continue
			}
			if currentInterval != cfg.TimerInterval {
				select {
				case a.tc.intervalCh <- cfg.TimerInterval:
					currentInterval = cfg.TimerInterval
				default:
				}
			}
		}
	}
}

// countUpdater periodically checks for updated worker count in the configuration and scales the worker pool accordingly.
func (a *RssAggregator) countUpdater(ctx context.Context, currentCount int) {
	defer a.wg.Done()

	t := time.NewTicker(time.Second * 2)
	defer t.Stop()

	defer close(a.wc.countCh)

	for {
		select {
		case <-ctx.Done():
			a.log.Debug(ctx, "worker count updater has been stopped")
			return
		case <-t.C:
			cfg, err := a.configRepo.Get(ctx)
			if err != nil {
				a.log.Error(ctx, "Failed to read config", "error", err)
				continue
			}
			if currentCount != cfg.WorkerCount {
				select {
				case a.wc.countCh <- cfg.WorkerCount:
					currentCount = cfg.WorkerCount
				default:
				}
			}
		}
	}
}

func (a *RssAggregator) ListenShutdown(ctx context.Context) {
	shutdownCh := make(chan os.Signal, 1)

	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	s := <-shutdownCh
	msg := fmt.Sprintf("catched shutdown signal %s", s.String())
	a.log.Notify(msg)

	if err := a.Stop(); err != nil {
		a.log.Error(ctx, "Failed to stop aggregator", "error", err)
	}
	a.cleanDb()

	msg = "graceful shutdown completed!"
	a.log.Notify(msg)
}
