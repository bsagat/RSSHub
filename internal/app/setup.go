package app

import (
	"RSSHub/config"
	"RSSHub/internal/adapters/cli"
	"RSSHub/internal/adapters/repo"
	"RSSHub/internal/pkg/logger"
	"RSSHub/internal/service"
	"context"
	"log/slog"
	"os"

	_ "github.com/lib/pq"
)

type App struct {
	cliHandler    *cli.CLIHandler
	postrgresRepo *repo.PostgresRepo
	log           *slog.Logger
}

func New(ctx context.Context, cfg *config.Config) *App {
	// Logger setup
	log := logger.New(cfg.Env)
	log.Info("Logger setup has been finished...")

	// Database connection
	log.Info("Connecting to Database...")
	postrgresRepo, err := repo.New(ctx, cfg.DB)
	if err != nil {
		log.Error("Failed to connect repo", "error", err)
		os.Exit(1)
	}
	log.Info("Database connection estabilished")

	aggregator := service.NewRssAggregator(postrgresRepo, log)
	cliHandler := cli.NewCLIHandler(aggregator, log)

	return &App{
		postrgresRepo: postrgresRepo,
		cliHandler:    cliHandler,
		log:           log,
	}
}
