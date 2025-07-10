package config

import (
	"RSSHub/pkg/envzilla"
	"RSSHub/pkg/postgres"
	"fmt"
	"time"
)

type (
	Config struct {
		Postgres postgres.Config
		App      CLI_APP
	}

	CLI_APP struct {
		TimerInterval time.Duration `env:"CLI_APP_TIMER_INTERVAL" default:"3m"`
		WorkerCount   int           `env:"CLI_APP_WORKERS_COUNT" default:"3"`
	}
)

func New() (*Config, error) {
	cfg := new(Config)

	if err := envzilla.Loader(".env"); err != nil {
		return nil, fmt.Errorf("failed to load config file: %w", err)
	}

	if err := envzilla.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return cfg, nil
}
