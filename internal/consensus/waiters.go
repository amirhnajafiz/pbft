package consensus

import (
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
			if value >= c.BFTCfg.Responses {
				return value
			}
		}
	}
}
