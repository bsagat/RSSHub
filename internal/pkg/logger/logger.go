package logger

import (
	"log/slog"
	"os"
)

const (
	LocalArea = "local"
	DebugArea = "debug"
	ProdArea  = "prod"
)

func New(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case LocalArea:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case DebugArea:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case ProdArea:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}
	return log
}
