package modules

import (
	"time"

	"github.com/f24-cse535/pbft/internal/config/node/bft"
	"github.com/f24-cse535/pbft/pkg/enum"
	"github.com/f24-cse535/pbft/pkg/models"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"
)

// Validator is function type used to validate input messages.
type Validator func(*pbft.AckMsg) *pbft.AckMsg

// Waiter is a module that handles wait processes on consensus demand.
type Waiter struct {
	cfg *bft.Config
}

// NewWaiter returns an instance of waiter module.
func NewWaiter(cfg *bft.Config) *Waiter {
	return &Waiter{
		cfg: cfg,
	}
}

// NewPrePreparedWaiter takes packets from a channel until it gets 2f+1 matching preprepared messages.
func (w *Waiter) NewPrePreparedWaiter(channel chan *models.Packet, validator Validator) int {
	// create a list of messages
	messages := make(map[string]int)

	// create a new timer, time start flag, and a target holder
	var (
		timer        *time.Timer
		target       string
		timerStarted bool = false
	)

	for {
		if timerStarted {
			select {
			case intr := <-channel:
				// ignore messages that are not preprepared
				if intr.Type != enum.PktPP {
					continue
				}

				// validate the message
				msg := intr.Payload.(*pbft.AckMsg)
				if msg = validator(msg); msg == nil {
					continue
				}

				// extract the digest to count the messages
				digest := msg.GetDigest()
				if _, ok := messages[digest]; !ok {
					messages[digest] = 1
				} else {
					messages[digest]++
				}
			case <-timer.C:
				return messages[target]
			}
		} else {
			// the timer is not started, so keep waiting
			intr := <-channel

			// ignore messages that are not preprepared
			if intr.Type != enum.PktPP {
				continue
			}

			// extract the digest to count the messages
			digest := intr.Payload.(*pbft.AckMsg).GetDigest()
			if _, ok := messages[digest]; !ok {
				messages[digest] = 1
			} else {
				messages[digest]++
			}

			// check for having 2f+1 match messages
			for key, value := range messages {
				if value >= w.cfg.Majority-1 {
					target = key
					timerStarted = true
					timer = time.NewTimer(time.Duration(w.cfg.MajorityTimeout) * time.Millisecond)
				}
			}
		}
	}
}

// NewPreparedWaiter takes packets from a channel until it gets 2f+1 matching prepared messages.
func (w *Waiter) NewPreparedWaiter(channel chan *models.Packet, validator Validator) int {
	// create a list of messages
	messages := make(map[string]int)

	for {
		// get raw interrupts
		intr := <-channel

		// ignore messages that are not prepared
		if intr.Type != enum.PktP {
			continue
		}

		// validate the message
		msg := intr.Payload.(*pbft.AckMsg)
		if msg = validator(msg); msg == nil {
			continue
		}

		// extract the digest to count the messages
		digest := msg.GetDigest()
		if _, ok := messages[digest]; !ok {
			messages[digest] = 1
		} else {
			messages[digest]++
		}

		// check for having 2f+1 match messages
		for _, value := range messages {
			if value >= w.cfg.Majority-1 {
				return value
			}
		}
	}
}
