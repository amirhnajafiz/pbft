package consensus

import (
	"time"

	"github.com/f24-cse535/pbft/pkg/enum"
	"github.com/f24-cse535/pbft/pkg/models"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"
)

// waitForPrePrepareds takes interrupts from the handler channel and waits
// until it gets 2f+1 matching preprepared messages.
func (c *Consensus) waitForPrePrepareds(channel chan *models.InterruptMsg) int {
	// create a list of messages
	messages := make(map[string]int)

	for {
		// get raw interrupts
		intr := <-channel

		// ignore messages that are not preprepared
		if intr.Type != enum.IntrPrePrepared {
			continue
		}

		// extract the digest to count the messages
		digest := intr.Payload.(*pbft.PrePreparedMsg).GetDigest()
		if _, ok := messages[digest]; !ok {
			messages[digest] = 1
		} else {
			messages[digest]++
		}

		// check for having 2f+1 match messages
		for _, value := range messages {
			if value >= c.BFTCfg.Majority {
				return value
			}
		}
	}
}

// waitForPrepareds takes interrupts from the handler channel and waits
// until it gets 2f+1 matching prepared messages.
func (c *Consensus) waitForPrepareds(channel chan *models.InterruptMsg) int {
	// create a list of messages
	messages := make(map[string]int)

	for {
		// get raw interrupts
		intr := <-channel

		// ignore messages that are not prepared
		if intr.Type != enum.IntrPrepared {
			continue
		}

		// extract the digest to count the messages
		digest := intr.Payload.(*pbft.PreparedMsg).GetDigest()
		if _, ok := messages[digest]; !ok {
			messages[digest] = 1
		} else {
			messages[digest]++
		}

		// check for having 2f+1 match messages
		for _, value := range messages {
			if value >= c.BFTCfg.Majority {
				return value
			}
		}
	}
}

// waitForReplys takes interrupts from the handler channel and waits
// until it gets f+1 matching replys messages.
func (c *Consensus) waitReplys(channel chan *models.InterruptMsg) *pbft.ReplyMsg {
	// create a list of messages
	messages := make(map[int64]int)
	replys := make(map[int64]*pbft.ReplyMsg)

	// create a new timer
	timer := time.NewTimer(time.Duration(c.BFTCfg.RequestTimeout) * time.Millisecond)
	defer timer.Stop()

	for {
		select {
		case intr := <-channel:
			// ignore messages that are not preprepared
			if intr.Type != enum.IntrReply {
				continue
			}

			// extract the digest to count the messages
			payload := intr.Payload.(*pbft.ReplyMsg)
			ts := payload.GetTimestamp()

			if _, ok := messages[ts]; !ok {
				messages[ts] = 1
				replys[ts] = payload
			} else {
				messages[ts]++
			}

			// check for having f+1 match messages
			for key, value := range messages {
				if value >= c.BFTCfg.Responses {
					return replys[key]
				}
			}
		case <-timer.C: // timeout
			return nil
		}
	}
}
