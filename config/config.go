package config

import (
	"RSSHub/pkg/envzilla"
	"RSSHub/pkg/postgres"
	"fmt"
)

type (
	Config struct {
		Postgres postgres.Config
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
