package modules

import (
	"sync"
	"time"
)

// Timer is a generic timer struct that is used in consensus module.
type Timer struct {
	duration time.Duration
	clock    *time.Timer
	lock     sync.Mutex
}

// NewTimer returns a timer instance.
func NewTimer(period int, unit time.Duration) *Timer {
	du := time.Duration(period) * unit

	return &Timer{
		lock:     sync.Mutex{},
		clock:    time.NewTimer(du),
		duration: du,
	}
}

// Start the timer.
func (t *Timer) Start() {
	t.lock.Lock()
	t.clock.Reset(t.duration)
	t.lock.Unlock()
}

// Stop the timer.
func (t *Timer) Stop() {
	t.lock.Lock()
	t.clock.Stop()
	t.lock.Unlock()
}

// Notify returns the timer channel.
func (t *Timer) Notify() <-chan time.Time {
	return t.clock.C
}
