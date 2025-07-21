package cmd

import (
	"RSSHub/config"
	"RSSHub/internal/app"
	"RSSHub/pkg/logger"
	"RSSHub/pkg/utils"
	"context"
	"flag"
	"os"
)

var helpFlag = flag.Bool("help", false, "Prints help message")

func Run() {
	flag.Parse()

	// Printing help flag
	if *helpFlag {
		utils.PrintHelp()
		return
	}

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
	if err := application.Run(ctx); err != nil {
		logger.Error(ctx, "failed to run application", "error", err)
		os.Exit(1)
	}
}
