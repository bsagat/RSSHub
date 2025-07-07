package main

import (
	"RSSHub/config"
	"RSSHub/internal/app"
	"context"
)

func main() {
	ctx := context.Background()

	// Config file parse
	cfg := config.New()

	// Main app configuration setup
	application := app.New(ctx, cfg)

	// Run app
	application.Run()
}
