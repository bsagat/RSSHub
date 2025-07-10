package logger

import (
	"context"
	"log"
	"log/slog"
	"os"
)

const (
	LevelDebug string = "DEBUG"
	LevelInfo  string = "INFO"
	LevelWarn  string = "WARN"
	LevelError string = "ERROR"
)

// Logger интерфейс логгера
type Logger interface {
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)

	GetSlogLogger() *slog.Logger
	GetLogLogger(level string, prefix string) *log.Logger
}

var l = logger{
	opts: &slog.HandlerOptions{
		Level: slog.LevelDebug,
	},
}

type logger struct {
	opts *slog.HandlerOptions
	slog *slog.Logger
}

// Debug логирование уровня debug
func (l *logger) Debug(ctx context.Context, msg string, args ...any) {
	l.slog.DebugContext(ctx, msg, args...)
}

// Info логирование уровня info
func (l *logger) Info(ctx context.Context, msg string, args ...any) {
	l.slog.InfoContext(ctx, msg, args...)
}

// Warn логирование уровня warn
func (l *logger) Warn(ctx context.Context, msg string, args ...any) {
	l.slog.WarnContext(ctx, msg, args...)
}

// Error логирование уровня error
func (l *logger) Error(ctx context.Context, msg string, args ...any) {
	l.slog.ErrorContext(ctx, msg, args...)
}

// GetSlogLogger метод для возврата объекта slog
func (l *logger) GetSlogLogger() *slog.Logger {
	return l.slog
}

// GetLogLogger метод для возврата объекта log.Logger. Данный лог не разделяется на уровни логирования.
// Поэтому необходимо завать отдельно с каким уровнем будет работать данный лог
// level - уровень с которым будет логироваться log.Logger
// prefix - префикс, который добавляется в сообщения
func (l *logger) GetLogLogger(level string, prefix string) *log.Logger {
	slogLevel := slog.LevelDebug
	switch level {
	case LevelDebug:
		slogLevel = slog.LevelDebug
	case LevelInfo:
		slogLevel = slog.LevelInfo
	case LevelWarn:
		slogLevel = slog.LevelWarn
	case LevelError:
		slogLevel = slog.LevelError
	}

	return log.New(&handlerWriter{l.slog.Handler(), slogLevel, true}, prefix, 0)
}

// InitLogger создаение логгера на основе slog
func InitLogger(logLevel string) Logger {
	switch logLevel {
	case LevelDebug:
		l.opts.Level = slog.LevelDebug
	case LevelInfo:
		l.opts.Level = slog.LevelInfo
	case LevelWarn:
		l.opts.Level = slog.LevelWarn
	case LevelError:
		l.opts.Level = slog.LevelError
	}

	l.slog = slog.New(slog.NewJSONHandler(os.Stdout, l.opts))
	return &l
}
