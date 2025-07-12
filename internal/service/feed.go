package service

import (
	"RSSHub/internal/domain/models"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

// Shows <num> recent articles for the given feed.
func (a *RssAggregator) GetArticles(feedName string, num int) ([]models.RSSItem, error) {
	const op = "RssAggregator.GetArticles"
	log := a.log.GetSlogLogger().With(
		slog.String("op:%s", op),
		slog.String("feed name", feedName),
		slog.Int("article count", num),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var (
		articles []models.RSSItem
		err      error
	)
	articles, err = a.articleRepo.List(ctx, feedName, num)
	if err != nil {
		log.Error("Failed to get articles list", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(articles) == 0 {
		return nil, fmt.Errorf("articles by feed %s are not found", feedName)
	}

	return articles, nil
}

// Shows the <num> most recently added feeds.
func (a *RssAggregator) ListFeeds(num int) ([]models.RSSFeed, error) {
	const op = "RssAggregator.ListFeeds"
	log := a.log.GetSlogLogger().With(
		slog.String("op", op),
		slog.Int("feeds count", num),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var (
		feeds []models.RSSFeed
		err   error
	)
	switch num {
	case 0:
		feeds, err = a.feedRepo.ListAll(ctx)
		if err != nil {
			log.Error("Failed to get all feed list", "error", err)
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	default:
		feeds, err = a.feedRepo.List(ctx, num)
		if err != nil {
			log.Error("Failed to get feed list", "error", err)
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	if len(feeds) == 0 {
		return nil, errors.New("feeds are not found")
	}

	return feeds, nil
}

func (a *RssAggregator) AddFeed(name, desc, url string) error {
	const op = "RssAggregator.AddFeed"
	log := a.log.GetSlogLogger().With(
		slog.String("op", op),
		slog.String("feed name", name),
		slog.String("URL", url),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Feed existense check
	exist, err := a.feedRepo.Exist(ctx, name)
	if err != nil {
		log.Error("Failed to check feed existence", "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	if exist {
		return fmt.Errorf("feed name must be unique")
	}

	// Creating a new feed
	feed := models.RSSFeed{
		Channel: models.Channel{
			Title:       name,
			Description: desc,
			Link:        url,
		},
	}
	if err := a.feedRepo.Create(ctx, &feed); err != nil {
		log.Error("Failed to create new feed", "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *RssAggregator) DeleteFeed(name string) error {
	const op = "RssAggregator.DeleteFeed"
	log := a.log.GetSlogLogger().With(
		slog.String("op: %s", op),
		slog.String("feed name", name),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Feed existense check
	exist, err := a.feedRepo.Exist(ctx, name)
	if err != nil {
		log.Error("Failed to check feed existence", "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	if !exist {
		return fmt.Errorf("feed is not exist")
	}

	if err := a.feedRepo.Delete(ctx, name); err != nil {
		log.Error("Failed to delete feed", "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
