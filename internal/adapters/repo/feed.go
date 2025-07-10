package repo

import (
	"RSSHub/internal/domain/models"
	"context"
	"errors"
	"fmt"

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

// Create new feed in database and sets its UUID
func (r *FeedRepo) Create(ctx context.Context, feed *models.RSSFeed) error {
	const op = "FeedRepo.Create"

	query := `
	INSERT INTO feeds(name, description, url)
	VALUES($1, $2, $3)
	RETURNING id, created_at;
	`

	err := r.db.QueryRow(ctx, query,
		feed.Channel.Title,
		feed.Channel.Description,
		feed.Channel.Link,
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
func (r *FeedRepo) ListAll(ctx context.Context) ([]models.RSSFeed, error) {
	const op = "FeedRepo.ListAll"

	query := `
	SELECT id, name, description, url, created_at 
	FROM feeds
	ORDER BY created_at DESC;
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	feeds, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.RSSFeed, error) {
		var feed models.RSSFeed
		err := row.Scan(
			&feed.ID,
			&feed.Channel.Title,
			&feed.Channel.Description,
			&feed.Channel.Link,
			&feed.CreatedAt,
		)
		return feed, err
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return feeds, nil
}

// List fetches up to limit recent feeds from the database
func (r *FeedRepo) List(ctx context.Context, limit int) ([]models.RSSFeed, error) {
	const op = "FeedRepo.List"

	query := `
	SELECT id, name, description, url, created_at 
	FROM feeds
	ORDER BY created_at DESC
	LIMIT $1;
	`

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	feeds, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.RSSFeed, error) {
		var feed models.RSSFeed
		err := row.Scan(
			&feed.ID,
			&feed.Channel.Title,
			&feed.Channel.Description,
			&feed.Channel.Link,
			&feed.CreatedAt,
		)
		return feed, err
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return feeds, nil
}
