package service

import (
	"sync"
	"time"
)

// VarTicker wraps time.Ticker, allowing safe concurrent
// reset of the ticker interval during runtime.
type VarTicker struct {
	ticker   *time.Ticker // Underlying ticker that emits ticks at a set interval.
	duration time.Duration
	mu       sync.Mutex // Mutex to protect ticker operations during concurrent access.
}

// NewVarTicker creates a new VarTicker with the specified delay interval.
func NewVarTicker(delay time.Duration) *VarTicker {
	return &VarTicker{
		ticker:   time.NewTicker(delay),
		duration: delay,
	}
}

// Stop safely stops the underlying ticker.
func (t *VarTicker) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.ticker.Stop()
}

// Reset safely resets the ticker interval to a new delay.
func (t *VarTicker) Reset(delay time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.ticker.Reset(delay)
	t.duration = delay
}

func (t *VarTicker) GetDuration() time.Duration {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.duration
}
