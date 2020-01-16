package scheduler

import (
	"time"
)

type Span struct {
	Name     string
	Start    time.Time
	Duration time.Duration
}
