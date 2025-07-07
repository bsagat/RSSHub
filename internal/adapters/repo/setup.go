package repo

import (
	"RSSHub/config"
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type PostgresRepo struct {
	Db *sql.DB
}

func New(ctx context.Context, DBconfig config.DB) (*PostgresRepo, error) {
	db, err := sql.Open("postgres", DBconfig.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database conn: %w", err)
	}

	return &PostgresRepo{
		Db: db,
	}, nil
}
