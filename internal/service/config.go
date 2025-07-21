package service

import (
	"RSSHub/internal/domain/models"
	"RSSHub/pkg/logger"
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

	msg := fmt.Sprintf("Interval of fetching feeds changed from %d minutes to %d minutes", int(last.Minutes()), int(changed.Minutes()))
	logger.Notify(log, msg)
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

	oldCount, err := a.configRepo.UpdateWorkerCount(ctx, workers)
	if err != nil {
		log.Error("Failed to update worker count", "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	msg := fmt.Sprintf("Number of workers changed from %d to %d ", oldCount, workers)
	logger.Notify(log, msg)
	return nil
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
		msg := "Background process already running"
		logger.Notify(log, msg)
		return nil, ErrProcessAlreadyRunning
	}

	if err := a.configRepo.UpdateRunStatus(ctx, true); err != nil {
		log.Error("Failed to update aggregator status", "error", err)
		return nil, fmt.Errorf("%s: %w", op, ErrFailedToUpdateStatus)
	}
	return cfg, nil
}
