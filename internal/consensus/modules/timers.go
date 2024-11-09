package modules

import (
	"sync"
	"time"
)

// Timer is a generic timer that is used in consensus module.
type Timer struct {
	duration time.Duration
	clock    *time.Timer
	lock     sync.Mutex
	counter  int
}

// NewTimer returns a timer instance.
func NewTimer(period int, unit time.Duration) *Timer {
	du := time.Duration(period) * unit

	return &Timer{
		lock:     sync.Mutex{},
		counter:  0,
		clock:    time.NewTimer(du),
		duration: du,
	}
}

// Start the timer again.
func (t *Timer) Start(flag bool) {
	t.lock.Lock()

	if flag {
		t.counter++
	}

	t.clock.Reset(t.duration)

	t.lock.Unlock()
}

// Stop the timer.
func (t *Timer) Stop(flag bool) {
	t.lock.Lock()

	if flag {
		t.counter--
	}

	if t.counter == 0 {
		t.clock.Stop()
	}

	t.lock.Unlock()
}

// Notify when timer expires.
func (t *Timer) Notify() <-chan time.Time {
	return t.clock.C
}
