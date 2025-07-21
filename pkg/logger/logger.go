package logger

import (
	"context"
	"fmt"
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

	Notify(msg string)
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

func (l *logger) Notify(msg string) {
	l.Debug(context.TODO(), msg)
	fmt.Println(msg)
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

// InitLogger создаёт логгер на основе slog и пишет в файл .log
func InitLogger(logLevel string) Logger {
	// Открываем файл для логов
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("cannot open log file: %v", err))
	}

	// Устанавливаем уровень логирования
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

	// Настраиваем JSON-хендлер с записью в файл
	handler := slog.NewJSONHandler(file, l.opts)
	l.slog = slog.New(handler)

	return &l
}
