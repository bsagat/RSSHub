package models

import (
	"time"
)

// RssConfig represent config stored in database
type RssConfig struct {
	Run           bool
	WorkerCount   int
	TimerInterval time.Duration
}
