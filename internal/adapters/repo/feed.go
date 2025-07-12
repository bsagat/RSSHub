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

var ErrFeedNotFound = errors.New("feed not found")

type FeedRepo struct {
	db *pgxpool.Pool
}

func NewFeedRepo(pool *pgxpool.Pool) *FeedRepo {
	return &FeedRepo{
		db: pool,
	}
}

// Create new feed in database.
func (r *FeedRepo) Create(ctx context.Context, feed *models.Feed) error {
	const op = "FeedRepo.Create"

	query := `
	INSERT INTO feeds(name, description, url)
	VALUES($1, $2, $3)
	RETURNING id, created_at;
	`

	err := r.db.QueryRow(ctx, query,
		feed.Name,
		feed.Description,
		feed.URL,
	).Scan(&feed.ID, &feed.CreatedAt)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Delete feed by name and its articles by cascade
func (r *FeedRepo) Delete(ctx context.Context, name string) error {
	const op = "FeedRepo.Delete"

	query := `
	DELETE FROM feeds
	WHERE name = $1;
	`

	tag, err := r.db.Exec(ctx, query, name)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, ErrFeedNotFound)
	}

	return nil
}

// ListAll fetches all recent feeds from the database
func (r *FeedRepo) ListAll(ctx context.Context) ([]*models.Feed, error) {
	const op = "FeedRepo.ListAll"

	query := `
		SELECT id, name, description, url, created_at, updated_at
		FROM feeds
		ORDER BY created_at DESC;
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	feeds, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*models.Feed, error) {
		var feed models.Feed
		err := row.Scan(
			&feed.ID,
			&feed.Name,
			&feed.Description,
			&feed.URL,
			&feed.CreatedAt,
			&feed.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		return &feed, nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return feeds, nil
}

// List fetches up to `limit` recent feeds from the database
func (r *FeedRepo) List(ctx context.Context, limit int) ([]*models.Feed, error) {
	const op = "FeedRepo.List"

	query := `
		SELECT id, name, description, url, created_at, updated_at
		FROM feeds
		ORDER BY created_at DESC
		LIMIT $1;
	`

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	feeds, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*models.Feed, error) {
		var feed models.Feed
		err := row.Scan(
			&feed.ID,
			&feed.Name,
			&feed.Description,
			&feed.URL,
			&feed.CreatedAt,
			&feed.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		return &feed, nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return feeds, nil
}

// Exist
func (f *FeedRepo) Exist(ctx context.Context, name string) (bool, error) {
	const op = "FeedRepo.IsUnique"

	query := `
	SELECT COUNT(*) != 0  as exist
	FROM feeds
	WHERE Name = $1
	`
	var exist bool
	if err := f.db.QueryRow(ctx, query, name).
		Scan(&exist); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return exist, nil
}

// GetStaleFeeds returns all feeds that haven't been updated in the given period.
func (f *FeedRepo) GetStaleFeeds(ctx context.Context, period time.Duration) ([]*models.Feed, error) {
	const op = "FeedRepo.GetStaleFeeds"

	query := `
		SELECT id, name, description, url, created_at, updated_at
		FROM feeds
		WHERE updated_at IS NULL OR updated_at < $1
	`

	cutoff := time.Now().Add(-period)
	rows, err := f.db.Query(ctx, query, cutoff)
	if err != nil {
		return nil, fmt.Errorf("%s: query error: %w", op, err)
	}
	defer rows.Close()

	var feeds []*models.Feed
	for rows.Next() {
		feed := new(models.Feed)
		err := rows.Scan(
			feed.ID,
			feed.Name,
			feed.Description,
			feed.URL,
			feed.CreatedAt,
			feed.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: scan error: %w", op, err)
		}
		feeds = append(feeds, feed)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows error: %w", op, err)
	}

	return feeds, nil
}

// UpdateUpdatedAt updates updated_at field
func (f *FeedRepo) UpdateUpdatedAt(ctx context.Context, name string) error {
	const op = "FeedRepo.UpdateUpdatedAt"

	query := `
		UPDATE feeds 
		SET updated_at = NOW()
		WHERE name = $1;
	`

	_, err := f.db.Exec(ctx, query, name)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
