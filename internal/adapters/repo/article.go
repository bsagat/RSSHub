package repo

import (
	"RSSHub/internal/domain/models"
	"context"
	"database/sql"
	"fmt"
)

type ArticleRepo struct {
	db *sql.DB
}

// Creates new articles in database by feed ID, atomically
func (repo *ArticleRepo) Create(ctx context.Context, feedID string, articles []models.RSSItem) error {
	const op = "ArticleRepo.Create"
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	query := `
	INSERT INTO 
		articles(title, link, description, Published_at, feed_id)
	VALUES 
		($1, $2, $3, $4, $5);
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s:%w", op, err)
	}
	defer stmt.Close()

	for _, article := range articles {
		_, err := stmt.ExecContext(ctx, article.Title, article.Link, article.Description, article.PubDate, feedID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("%s:%w", op, err)
		}
	}
	return tx.Commit()
}

// List fetches up to limit recent articles by feed name from the database.
func (repo *ArticleRepo) List(ctx context.Context, feedName string, limit int) ([]models.RSSItem, error) {
	const op = "ArticleRepo.List"
	query := `
		SELECT 
			a.title, a.link, a.description, a.Published_at 
		FROM 
			articles a
		INNER JOIN
			feeds f 
		ON 
			a.feed_id=f.ID
		WHERE 
			f.Name=$1
		LIMIT $2;`

	rows, err := repo.db.QueryContext(ctx, query, feedName, limit)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	defer rows.Close()

	var articles []models.RSSItem
	for rows.Next() {
		var article models.RSSItem
		if err := rows.Scan(&article.Title, &article.Link, &article.Description, &article.PubDate); err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}
		articles = append(articles, article)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return articles, nil
}
