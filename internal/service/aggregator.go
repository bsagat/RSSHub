package service

import (
	"RSSHub/internal/adapters/repo"
	"RSSHub/pkg/logger"
	"context"
	"sync"
)

type RssAggregator struct {
	articleRepo *repo.ArticleRepo
	feedRepo    *repo.FeedRepo

	log logger.Logger
	mu  sync.Mutex
}

func NewRssAggregator(articleRepo *repo.ArticleRepo, feedRepo *repo.FeedRepo, log logger.Logger) *RssAggregator {
	return &RssAggregator{
		articleRepo: articleRepo,
		feedRepo:    feedRepo,
		log:         log,
		mu:          sync.Mutex{},
	}
}

func (a *RssAggregator) Start(ctx context.Context) error {
	return nil
}

func (a *RssAggregator) Stop() error {
	return nil
}
