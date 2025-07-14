package app

import (
	"RSSHub/config"
	"RSSHub/internal/adapters/cli"
	"RSSHub/internal/adapters/repo"
	"RSSHub/internal/service"
	"RSSHub/pkg/logger"
	"RSSHub/pkg/postgres"
	"context"
	"fmt"
)

type App struct {
	cliHandler *cli.CLIHandler
	postgresDB *postgres.PostgreDB
	aggregator *service.RssAggregator

	cfg *config.Config
	log logger.Logger
}

func New(ctx context.Context, cfg *config.Config, logger logger.Logger) (*App, error) {
	const fn = "app.NewApplication"
	log := logger.GetSlogLogger().With("fn", fn)

	// Database connection
	db, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		log.Error("failed to connect postgres", "error", err)
		return nil, fmt.Errorf("failed to connect postgres")
	}

	// Repository
	articleRepo := repo.NewArticleRepo(db.Pool)
	feedRepo := repo.NewFeedRepo(db.Pool)
	configRepo := repo.NewConfigRepo(db.Pool)

	// Services
	aggregator := service.NewRssAggregator(articleRepo, feedRepo, configRepo, logger, func() {
		db.Close()
	})

	// CLI Handler
	cliHandler := cli.NewCLIHandler(aggregator, cfg.App, logger)

	return &App{
		cliHandler: cliHandler,
		postgresDB: db,
		aggregator: aggregator,

		cfg: cfg,
		log: logger,
	}, nil
}

func (app *App) Run(ctx context.Context) error {
	// Running CLI
	code := app.cliHandler.ParseFlags()
	if code != 0 {
		return fmt.Errorf("cli exited with code %d", code)
	}

	return nil
}
