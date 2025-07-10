package service

import (
	"RSSHub/internal/adapters/repo"
	"RSSHub/internal/domain/models"
	"RSSHub/pkg/logger"
	"context"
	"sync"
	"time"
)

type RssAggregator struct {
	repo *repo.ArticleRepo
	log  logger.Logger
	mu   sync.Mutex
}

func NewRssAggregator(repo *repo.ArticleRepo, log logger.Logger) *RssAggregator {
	return &RssAggregator{
		repo: repo,
		log:  log,
		mu:   sync.Mutex{},
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
