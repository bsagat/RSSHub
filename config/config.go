package config

import (
	"RSSHub/internal/pkg/envzilla"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"
)

type (
	Config struct {
		Env           string
		TimerInterval time.Duration
		WorkerCount   int
		DB
	}

	DB struct {
		DB_HOST     string
		DB_PORT     string
		DB_USER     string
		DB_PASSWORD string
		DB_NAME     string
	}
)

func New() *Config {
	if err := envzilla.Loader(".env"); err != nil {
		slog.Error("Failed to load config file", "error", err)
		os.Exit(1)
	}

	timerInterval, err := time.ParseDuration(os.Getenv("CLI_APP_TIMER_INTERVAL"))
	if err != nil {
		slog.Error("Failed to parse CLI timer interval", "error", err)
		os.Exit(1)
	}

	workerCount, err := strconv.Atoi(os.Getenv("CLI_APP_WORKERS_COUNT"))
	if err != nil {
		slog.Error("Failed to parse worker count", "error", err)
		os.Exit(1)
	}

	return &Config{
		Env:           os.Getenv("CLI_APP_ENV"),
		TimerInterval: timerInterval,
		WorkerCount:   workerCount,
		DB: DB{
			DB_HOST:     os.Getenv("DB_HOST"),
			DB_PORT:     os.Getenv("DB_PORT"),
			DB_USER:     os.Getenv("DB_USER"),
			DB_PASSWORD: os.Getenv("DB_PASSWORD"),
			DB_NAME:     os.Getenv("DB_NAME"),
		},
	}
}

func (c *DB) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DB_HOST, c.DB_PORT, c.DB_USER, c.DB_PASSWORD, c.DB_NAME,
	)
}
