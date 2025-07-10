package cmd

import (
	"RSSHub/config"
	"RSSHub/internal/app"
	"RSSHub/pkg/logger"
	"context"
	"os"
)

func Run() {
	ctx := context.Background()

	logger := logger.InitLogger(logger.LevelDebug)

	// Config file parse
	cfg, err := config.New()
	if err != nil {
		logger.Error(ctx, "failed to init config", "error", err)
		os.Exit(1)
	}

	// Creating application
	application, err := app.New(ctx, cfg, logger)
	if err != nil {
		logger.Error(ctx, "failed to init application", "error", err)
		os.Exit(1)
	}

	// Running the apllication
	if err := application.Run(); err != nil {
		logger.Error(ctx, "failed to run application", "error", err)
		os.Exit(1)
	}
}
