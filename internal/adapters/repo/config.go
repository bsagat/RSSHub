package repo

import (
	"RSSHub/internal/domain/models"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrConfigNotFound = errors.New("config not found in database")

type ConfigRepo struct {
	pool *pgxpool.Pool
}

func NewConfigRepo(pool *pgxpool.Pool) *ConfigRepo {
	return &ConfigRepo{
		pool: pool,
	}
}

// Get retrieves the current configuration from the database
func (r *ConfigRepo) Get(ctx context.Context) (*models.RssConfig, error) {
	const op = "ConfigRepo.Get"
	const query = `
		SELECT 
			run, 
			worker_count, 
			timer_interval 
		FROM 
			config 
		LIMIT 1`

	var config models.RssConfig

	err := r.pool.QueryRow(ctx, query).Scan(
		&config.Run,
		&config.WorkerCount,
		&config.TimerInterval,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrConfigNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &config, nil
}

// Update replaces the existing configuration (single row) with new values
func (r *ConfigRepo) Update(ctx context.Context, config *models.RssConfig) error {
	const op = "ConfigRepo.Update"
	const query = `
		UPDATE config 
		SET 
			run = $1,
			worker_count = $2,
			timer_interval = $3`

	_, err := r.pool.Exec(ctx, query,
		config.Run,
		config.WorkerCount,
		config.TimerInterval,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// UpdateRunStatus updates only the run status in the configuration
func (r *ConfigRepo) UpdateRunStatus(ctx context.Context, run bool) error {
	const op = "ConfigRepo.UpdateRunStatus"
	const query = `
		UPDATE config 
		SET run = $1`

	_, err := r.pool.Exec(ctx, query, run)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// UpdateWorkerCount updates only the worker count in the configuration
func (r *ConfigRepo) UpdateWorkerCount(ctx context.Context, count int) (int, error) {
	const op = "ConfigRepo.UpdateWorkerCount"
	const query = `
		WITH old_count AS (
        	SELECT worker_count FROM config 
    	)
        UPDATE config 
        SET worker_count = $1
		FROM old_count
		RETURNING old_count.worker_count`

	var oldCount int
	if err := r.pool.QueryRow(ctx, query, count).Scan(&oldCount); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return oldCount, nil
}

// UpdateTimerInterval updates only the timer interval in the configuration
func (r *ConfigRepo) UpdateTimerInterval(ctx context.Context, interval time.Duration) (*time.Duration, error) {
	const op = "ConfigRepo.UpdateTimerInterval"
	const query = `
	WITH old_interval AS (
		SELECT timer_interval FROM config
	)
	
	UPDATE config 
	SET timer_interval = $1 
	FROM old_interval
	RETURNING old_interval.timer_interval`

	var lastInterval time.Duration
	if err := r.pool.QueryRow(ctx, query, interval).Scan(&lastInterval); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &lastInterval, nil
}
