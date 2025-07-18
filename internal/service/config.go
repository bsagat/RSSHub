package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

func (a *RssAggregator) SetInterval(changed time.Duration) error {
	const op = "RssAggregator.SetInterval"
	log := a.log.GetSlogLogger().With(
		slog.String("op", op),
		slog.Duration("new duration", changed),
	)
	ctx := context.TODO()

	last, err := a.configRepo.UpdateTimerInterval(ctx, changed)
	if err != nil {
		log.Error("Failed to update timer interval", "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info(fmt.Sprintf("Interval of fetching feeds changed from %d minutes to %d minutes", int(last.Minutes()), int(changed.Minutes())))
	return nil
}

func (a *RssAggregator) Resize(workers int) error {
	const op = "RssAggregator.Resize"
	log := a.log.GetSlogLogger().With(
		slog.String("op", op),
		slog.Int("worker count", workers),
	)
	ctx := context.TODO()

	if workers > 10000 {
		return errors.New("max goroutine count is 10000")
	}

	if err := a.configRepo.UpdateWorkerCount(ctx, workers); err != nil {
		log.Error("Failed to update worker count", "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
