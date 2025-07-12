package models

import "time"

// RssConfig represent config stored in database
type RssConfig struct {
	Run           bool
	WorkerCount   int
	TimerInterval time.Duration
}

var (
	DefaultRssConfig *RssConfig = &RssConfig{
		Run:           false,
		WorkerCount:   3,
		TimerInterval: 3 * time.Minute,
	}
)
