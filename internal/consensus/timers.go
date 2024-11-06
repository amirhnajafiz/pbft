package consensus

import (
	"fmt"
	"time"
)

// viewTimer is a timer that expires if the leader does not response during execution waits.
func (c *Consensus) viewTimer() {
	var (
		timer   *time.Timer = time.NewTimer(time.Duration(c.BFTCfg.ViewTimeout) * time.Millisecond)
		hashMap             = make(map[int]bool)
	)

	// in a while loop capture events
	for {
		if c.viewTimerStatus {
			c.processTimerStart(timer, hashMap)
		} else {
			c.processTimerStop(timer, hashMap)
		}
	}
}

// processTimerStart captures the input of a timer when it is started.
func (c *Consensus) processTimerStart(timer *time.Timer, hashMap map[int]bool) {
	select {
	case <-c.viewTimerStart:
		c.viewTimerStatus = true
		timer.Reset(time.Duration(c.BFTCfg.ViewTimeout) * time.Millisecond)
	case <-c.viewTimerHalt:
		c.viewTimerStatus = false
		timer.Stop()
	case seq := <-c.viewTimerIn:
		hashMap[seq] = true
		timer.Reset(time.Duration(c.BFTCfg.ViewTimeout) * time.Millisecond)
	case seq := <-c.viewTimerRm:
		delete(hashMap, seq)
		timer.Reset(time.Duration(c.BFTCfg.ViewTimeout) * time.Millisecond)
	case <-timer.C:
		fmt.Println("timer timeout")
		c.viewTimerStatus = false
	}
}

// processTimerStop captures the input of a timer when it is stopped.
func (c *Consensus) processTimerStop(timer *time.Timer, hashMap map[int]bool) {
	select {
	case <-c.viewTimerStart:
		c.viewTimerStatus = true
		timer.Reset(time.Duration(c.BFTCfg.ViewTimeout) * time.Millisecond)
	case <-c.viewTimerHalt:
		c.viewTimerStatus = false
		timer.Stop()
	case seq := <-c.viewTimerIn:
		hashMap[seq] = true
		c.viewTimerStatus = true
		timer.Reset(time.Duration(c.BFTCfg.ViewTimeout) * time.Millisecond)
	case seq := <-c.viewTimerRm:
		delete(hashMap, seq)
	}
}
