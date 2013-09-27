package sf

import (
	"time"
)

type Clock struct {
	startTime time.Time
}

func NewClock() *Clock {
	return &Clock{time.Now()}
}

func (c *Clock) ElapsedTime() time.Duration {
	return time.Since(c.startTime)
}

func (c *Clock) Restart() time.Duration {
	elapsed := time.Since(c.startTime)
	c.startTime = time.Now()

	return elapsed
}
