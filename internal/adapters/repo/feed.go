package repo

import (
	"RSSHub/internal/domain/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrFeedNotFound = errors.New("feed not found")

type FeedRepo struct {
	db *sql.DB
}

func NewFeedRepo(db *sql.DB) *FeedRepo {
	return &FeedRepo{
		db: db,
	}
}

// Creates new feed in database and sets his UUID
func (repo *FeedRepo) Create(ctx context.Context, feed *models.RSSFeed) error {
	const op = "FeedRepo.Create"
	query := `
	INSERT INTO 
		feeds(Name,description, URL)
	VALUES
		($1, $2, $3)
	RETURNING
		ID;
	`
	if err := repo.db.QueryRowContext(ctx, query, feed.Channel.Title, feed.Channel.Description, feed.Channel.Link).
		Scan(&feed.ID); err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	return nil
}

// Deletes feed by his name and his articles by cascade
func (repo *FeedRepo) Delete(ctx context.Context, name string) error {
	const op = "FeedRepo.Delete"
	query := `
		DELETE 
			FROM feeds
		WHERE 
			Name = $1;
	`
	res, err := repo.db.ExecContext(ctx, query, name)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	aff, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	if aff == 0 {
		return ErrFeedNotFound
	}
	return nil
}

// ListAll fetches all recent feeds from the database.
func (repo *FeedRepo) ListAll(ctx context.Context) ([]models.RSSFeed, error) {
	const op = "FeedRepo.ListAll"
	query := `
		SELECT 
			ID, Name, Description, URL, Created_at FROM feeds
		ORDER BY
			 Created_at;`

	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	defer rows.Close()

	var feeds []models.RSSFeed
	for rows.Next() {
		var feed models.RSSFeed
		if err := rows.Scan(&feed.ID, &feed.Channel.Title, &feed.Channel.Description, &feed.Channel.Link, &feed.CreatedAt); err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}
		feeds = append(feeds, feed)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return feeds, nil
}

// List fetches up to limit recent feeds from the database.
func (repo *FeedRepo) List(ctx context.Context, limit int) ([]models.RSSFeed, error) {
	const op = "FeedRepo.List"
	query := `
		SELECT 
			ID, Name, Description, URL, Created_at FROM feeds
		ORDER BY 
			Created_at
		LIMIT $1;`

	rows, err := repo.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	defer rows.Close()

	var feeds []models.RSSFeed
	for rows.Next() {
		var feed models.RSSFeed
		if err := rows.Scan(&feed.ID, &feed.Channel.Title, &feed.Channel.Description, &feed.Channel.Link, &feed.CreatedAt); err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}
		feeds = append(feeds, feed)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return feeds, nil
}
