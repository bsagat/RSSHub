package service

import (
	"context"
	"sync"
	"time"
)

type RssAggregator struct {
	mu sync.Mutex
}

func NewRssAggregator() *RssAggregator {
	return &RssAggregator{}
}

func (a *RssAggregator) Start(ctx context.Context) error {
	return nil
}
func (a *RssAggregator) Stop() error {
	return nil
}

func (a *RssAggregator) SetInterval(d time.Duration) {

}
func (a *RssAggregator) Resize(workers int) error {
	return nil
}
