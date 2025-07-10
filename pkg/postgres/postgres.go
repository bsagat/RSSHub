package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	DB_HOST     string `env:"DB_HOST" default:"localhost"`
	DB_PORT     string `env:"DB_PORT" default:"5432"`
	DB_USER     string `env:"DB_USER" default:"postgres"`
	DB_PASSWORD string `env:"DB_PASSWORD" default:""`
	DB_NAME     string `env:"DB_NAME" default:"rsshub"`

	MaxOpenConns int32  `env:"POSTGRES_MAX_OPEN_CONN" default:"25"`
	MaxIdleTime  string `env:"POSTGRES_MAX_IDLE_TIME" default:"15m"`
}

func (c Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DB_HOST, c.DB_PORT, c.DB_USER, c.DB_PASSWORD, c.DB_NAME,
	)
}

type PostgreDB struct {
	Pool     *pgxpool.Pool
	DBConfig *pgxpool.Config
}

func New(ctx context.Context, config Config) (*PostgreDB, error) {
	dsn := config.DSN()
	dbConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}

	// Setting maxOpenConns
	dbConfig.MaxConns = config.MaxOpenConns

	// Parse idle timeout duration
	duration, err := time.ParseDuration(config.MaxIdleTime)
	if err != nil {
		return nil, fmt.Errorf("failed to parse MaxIdleTime: %w", err)
	}
	dbConfig.MaxConnIdleTime = duration

	pool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Ping the database
	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgreDB{
		Pool:     pool,
		DBConfig: dbConfig,
	}, nil
}

func (db *PostgreDB) Close() {
	db.Pool.Close()
}

// Health checks the database connection health using Ping
func (db *PostgreDB) Health(ctx context.Context) (bool, error) {
	if db.Pool == nil {
		return false, fmt.Errorf("database pool is not initialized")
	}

	// Create a context with timeout for the health check
	healthCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// Simply ping the database
	if err := db.Pool.Ping(healthCtx); err != nil {
		return false, fmt.Errorf("database ping failed: %w", err)
	}

	return true, nil
}
