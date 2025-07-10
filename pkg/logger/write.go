package logger

import (
	"context"
	"log/slog"
	"runtime"
	"time"
)

// handlerWriter это io.Writer который вызывает slog.Handler.
type handlerWriter struct {
	h         slog.Handler
	level     slog.Leveler
	capturePC bool
}

// Write записывает информацию в логи (скопирован из slog.NewLogLogger)
func (w *handlerWriter) Write(buf []byte) (int, error) {
	level := w.level.Level()
	if !w.h.Enabled(context.Background(), level) {
		return 0, nil
	}
	var pc uintptr
	var pcs [1]uintptr
	runtime.Callers(4, pcs[:])
	pc = pcs[0]

	// Remove final newline.
	origLen := len(buf) // Report that the entire buf was written.
	if len(buf) > 0 && buf[len(buf)-1] == '\n' {
		buf = buf[:len(buf)-1]
	}
	r := slog.NewRecord(time.Now(), level, string(buf), pc)
	return origLen, w.h.Handle(context.Background(), r)
}
