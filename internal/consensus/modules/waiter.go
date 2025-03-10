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
func (w *Waiter) StartWaiting(channel chan *models.Packet, targetType enum.PacketType, validator Validator) (int, int) {
	// create a list of messages and nodes
	messages := make(map[string]int)
	nodes := make(map[string]bool)
	signed := 0

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
				if intr.Type != targetType {
					continue
				}

				msg := intr.Payload.(*pbft.AckMsg)
				if msg = validator(msg); msg == nil {
					continue
				}

				if _, ok := nodes[msg.GetNodeId()]; !ok {
					nodes[msg.GetNodeId()] = true
				} else {
					continue
				}

				if msg.GetSign() != nil && len(msg.GetSign()) > 0 {
					signed++
				}

				digest := msg.GetDigest()
				if _, ok := messages[digest]; !ok {
					messages[digest] = 1
				} else {
					messages[digest]++
				}
			case <-timer.C:
				return messages[target], signed
			}
		} else {
			intr := <-channel

			if intr.Type != targetType {
				continue
			}

			msg := intr.Payload.(*pbft.AckMsg)
			if msg = validator(msg); msg == nil {
				continue
			}

			if _, ok := nodes[msg.GetNodeId()]; !ok {
				nodes[msg.GetNodeId()] = true
			} else {
				continue
			}

			if msg.GetSign() != nil && len(msg.GetSign()) > 0 {
				signed++
			}

			digest := msg.GetDigest()
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
