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
	"os"
	"os/signal"
	"syscall"
)

const serviceName = "rsshub"

type App struct {
	cliHandler *cli.CLIHandler
	postgresDB *postgres.PostgreDB
	aggregator *service.RssAggregator
	log        logger.Logger
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

	// Services
	aggregator := service.NewRssAggregator(articleRepo, logger)

	// CLI Handler
	cliHandler := cli.NewCLIHandler(aggregator, logger)

	return &App{
		cliHandler: cliHandler,
		postgresDB: db,
		aggregator: aggregator,
		log:        logger,
	}, nil
}

func (app *App) close(_ context.Context) {
	// Closing database connection
	app.postgresDB.Close()

	// Closing CLI handler if needed
	app.cliHandler.Close()
}

func (app *App) Run() error {
	const fn = "app.Run"
	log := app.log.GetSlogLogger().With("fn", fn)

	errCh := make(chan error, 1)
	ctx := context.Background()

	// Running CLI
	code := app.cliHandler.ParseFlags()
	if code != 0 {
		return fmt.Errorf("cli exited with code %d", code)
	}

	log.InfoContext(ctx, "application started", "name", serviceName)

	// Waiting for signal or error
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case errRun := <-errCh:
		return errRun
	case s := <-shutdownCh:
		log.InfoContext(ctx, "shutting down application", "signal", s.String())

		app.close(ctx)
		log.InfoContext(ctx, "graceful shutdown completed!")
	}

	return nil
}
