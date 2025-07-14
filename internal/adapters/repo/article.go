package repo

import (
	"RSSHub/internal/domain/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ArticleRepo struct {
	pool *pgxpool.Pool
}

func NewArticleRepo(pool *pgxpool.Pool) *ArticleRepo {
	return &ArticleRepo{
		pool: pool,
	}
}

// Create atomically inserts multiple articles for a feed using batch operations
func (r *ArticleRepo) Create(ctx context.Context, feedID string, articles []models.RSSItem) error {
	const op = "ArticleRepo.Create"

	batch := &pgx.Batch{}
	query := `
		INSERT INTO articles(
			title, 
			link, 
			description, 
			published_at, 
			feed_id
		) VALUES (
			$1, $2, $3, $4, $5
		) `

	for _, article := range articles {
		batch.Queue(query,
			article.Title,
			article.Link,
			article.Description,
			article.PubDate,
			feedID,
		)
	}

	br := r.pool.SendBatch(ctx, batch)
	defer br.Close()

	_, err := br.Exec()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// List fetches recent articles by feed name with limit
func (r *ArticleRepo) List(ctx context.Context, feedName string, limit int) ([]*models.RSSItem, error) {
	const op = "ArticleRepo.List"

	query := `
		SELECT 
			a.title, 
			a.link, 
			a.description, 
			a.published_at::TEXT 
		FROM 
			articles a
		JOIN feeds f ON a.feed_id = f.id
		WHERE 
			f.name = $1
		ORDER BY 
			a.published_at DESC
		LIMIT $2`

	rows, err := r.pool.Query(ctx, query, feedName, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	articles, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*models.RSSItem, error) {
		var item models.RSSItem
		err := row.Scan(
			&item.Title,
			&item.Link,
			&item.Description,
			&item.PubDate,
		)
		if err != nil {
			return nil, err
		}
		return &item, nil
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return articles, nil
}

// ListAll fetches all articles for a feed
func (r *ArticleRepo) ListAll(ctx context.Context, feedName string) ([]*models.RSSItem, error) {
	const op = "ArticleRepo.ListAll"

	query := `
		SELECT 
			a.title, 
			a.link, 
			a.description, 
			a.published_at 
		FROM 
			articles a
		JOIN feeds f ON a.feed_id = f.id
		WHERE 
			f.name = $1
		ORDER BY 
			a.published_at DESC`

	rows, err := r.pool.Query(ctx, query, feedName)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	articles, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*models.RSSItem, error) {
		var item models.RSSItem
		err := row.Scan(
			&item.Title,
			&item.Link,
			&item.Description,
			&item.PubDate,
		)
		if err != nil {
			return nil, err
		}
		return &item, nil
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return articles, nil
}
