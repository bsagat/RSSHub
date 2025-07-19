package service

import (
	"RSSHub/internal/adapters/httpadapter"
	"RSSHub/internal/adapters/repo"
	"RSSHub/internal/domain/models"
	"RSSHub/pkg/logger"
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
	ErrProcessAlreadyRunning = errors.New("process already running")
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

// loadConfig retrieves the RSS aggregator configuration, checks running state, and updates run status in the repository.
func (a *RssAggregator) loadConfig(ctx context.Context) (*models.RssConfig, error) {
	const op = "RssAggregator.loadConfig"
	log := a.log.GetSlogLogger().With("op", op)

	cfg, err := a.configRepo.Get(ctx)
	if err != nil {
		log.Error("Failed to read config", "error", err)
		return nil, fmt.Errorf("%s: %w", op, ErrFailedToReadConfig)
	}
	if cfg == nil {
		log.Error("Config is not found")
		return nil, ErrConfigNotFound
	}
	if cfg.Run {
		log.Error("Background process already running")
		return nil, ErrProcessAlreadyRunning
	}

	if err := a.configRepo.UpdateRunStatus(ctx, true); err != nil {
		log.Error("Failed to update aggregator status", "error", err)
		return nil, fmt.Errorf("%s: %w", op, ErrFailedToUpdateStatus)
	}
	return cfg, nil
}

// Start launches the RSS aggregator, ticker controller, and worker controller with the configuration parameters.
func (a *RssAggregator) Start(ctx context.Context) error {
	const op = "RssAggregator.Start"
	log := a.log.GetSlogLogger().With("op", op)

	// Читаем конфиг
	cfg, err := a.loadConfig(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.wc = NewWorkerController(cfg.WorkerCount, a.log)
	a.tc = NewTickerController(cfg.TimerInterval, a.feedRepo, a.articleRepo, a.log)

	a.wg.Add(2)
	go a.tc.Run(a.ctx, &a.wg, a.wc)
	go a.wc.Run(a.ctx, cfg.WorkerCount, &a.wg)

	a.wg.Add(2)
	go a.intervalUpdater(a.ctx, cfg.TimerInterval)
	go a.countUpdater(a.ctx, cfg.WorkerCount)

	log.Info(fmt.Sprintf("The background process for fetching feeds has started (interval = %s minutes, workers = %d)", cfg.TimerInterval.String(), cfg.WorkerCount))
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

	log.Info("Graceful shutdown: aggregator stopped")
	return nil
}

// ---------------- TickerController ----------------

type TickerController struct {
	t           *VarTicker
	intervalCh  chan time.Duration
	feedRepo    *repo.FeedRepo
	articleRepo *repo.ArticleRepo
	log         logger.Logger
}

func NewTickerController(interval time.Duration, feedRepo *repo.FeedRepo, articleRepo *repo.ArticleRepo, log logger.Logger) *TickerController {
	return &TickerController{
		t:           NewVarTicker(interval),
		intervalCh:  make(chan time.Duration, 1),
		feedRepo:    feedRepo,
		articleRepo: articleRepo,
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
			fmt.Println("closing ticker controler")
			return
		case <-c.t.ticker.C:
			c.processFeeds(ctx, wc)
		case newInterval := <-c.intervalCh:
			oldInterval := c.t.GetDuration()
			c.t.Reset(newInterval)
			fmt.Printf("Interval of fetching feeds changed from %s minutes to %s minutes", oldInterval, newInterval)
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
			fetched, err := httpadapter.NewClient(time.Second*5).FetchRSSFeed(ctx, feed.URL)
			if err != nil {
				c.log.Error(ctx, "Failed to fetch RSS feed", "feed_URL", feed.URL, "error", err)
				return
			}
			if len(fetched.Channel.Item) == 0 {
				c.log.Error(ctx, "There is no items in the feed", "feed_URL", feed.URL)
				return
			}
			if err := c.articleRepo.Create(ctx, feed.ID, fetched.Channel.Item); err != nil {
				c.log.Error(ctx, "Failed to save feed items", "feed_id", feed.ID, "articles", fetched.Channel.Item, "error", err)
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
			fmt.Println("closing worker controler")
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
			fmt.Println("interval updater stopped")
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
			fmt.Println("worker updater stopped")
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
	a.log.Info(ctx, "shutting down application", "signal", s.String())
	if err := a.Stop(); err != nil {
		a.log.Error(ctx, "Failed to stop aggregator", "error", err)
	}
	a.cleanDb()

	a.log.Info(ctx, "graceful shutdown completed!")
}
