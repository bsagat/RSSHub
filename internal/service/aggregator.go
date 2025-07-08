package service

import (
	"RSSHub/internal/adapters/repo"
	"RSSHub/internal/domain/models"
	"context"
	"log/slog"
	"sync"
	"time"
)

type RssAggregator struct {
	repo *repo.PostgresRepo
	log  *slog.Logger
	mu   sync.Mutex
}

func NewRssAggregator(repo *repo.PostgresRepo, log *slog.Logger) *RssAggregator {
	return &RssAggregator{
		repo: repo,
		log:  log,
	}
}

func (a *RssAggregator) Start(ctx context.Context) error {
	return nil
}

func (a *RssAggregator) Stop() error {
	return nil
}

func (a *RssAggregator) SetInterval(d time.Duration) error {
	return nil
}

func (a *RssAggregator) SetWorkers(count int) error {
	return nil
}

func (a *RssAggregator) Resize(workers int) error {
	return nil
}

func (a *RssAggregator) GetArticles(feedName string, num int) ([]models.RSSItem, error) {
	return nil, nil
}

func (a *RssAggregator) ListFeeds(num int) ([]models.RSSFeed, error) {
	return nil, nil
}

func (a *RssAggregator) AddFeed(name, url string) error {
	return nil
}

func (a *RssAggregator) DeleteFeed(name string) error {
	return nil
}
