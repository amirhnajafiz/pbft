package modules

import (
	"time"
)

// Timer is a generic timer that is used in consensus module.
type Timer struct {
	duration time.Duration
	clock    *time.Timer
}

// NewTimer returns a timer instance.
func NewTimer(period int, unit time.Duration) *Timer {
	du := time.Duration(period) * unit

	return &Timer{
		clock:    time.NewTimer(du),
		duration: du,
	}
}

// Start the timer again.
func (t *Timer) Start() {
	t.clock.Reset(t.duration)
}

// Stop the timer.
func (t *Timer) Stop() {
	t.clock.Stop()
}

// Notify when timer expires.
func (t *Timer) Notify() {
	<-t.clock.C
}
