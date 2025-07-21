package ports

import (
	"RSSHub/internal/domain/models"
	"context"
	"time"
)

type Aggregator interface {
	// Core lifecycle
	Start(ctx context.Context) error // Starts background feed polling
	Stop() error                     // Graceful shutdown

	// Dynamic configuration
	SetInterval(d time.Duration) error // Dynamically changes fetch interval
	Resize(workers int) error          // Dynamically resizes worker pool
	GetConfig(ctx context.Context) (*models.RssConfig, error)

	// Feed management
	AddFeed(name, desc, url string) error      // Adds a new feed
	DeleteFeed(name string) error              // Deletes feed by name
	ListFeeds(num int) ([]*models.Feed, error) // Lists all feeds

	// Article retrieval
	GetArticles(feedName string, num int) ([]*models.RSSItem, error) // Gets latest 'num' articles for the feed
}
