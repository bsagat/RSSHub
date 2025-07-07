package ports

import (
	"context"
	"time"
)

type Aggregator interface {
	Start(ctx context.Context) error // Starts background feed polling
	Stop() error                     // Graceful shutdown

	SetInterval(d time.Duration) // Dynamically changes fetch interval
	Resize(workers int) error    // Dynamically resizes worker pool
}
