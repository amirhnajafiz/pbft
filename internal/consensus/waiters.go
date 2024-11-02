package consensus

import (
	"github.com/f24-cse535/pbft/pkg/enum"
	"github.com/f24-cse535/pbft/pkg/models"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"
)

// waitForPrePrepareds takes interrupts from the handler channel and waits
// until it gets 2f+1 matching preprepared messages.
func (c *Consensus) waitForPrePrepareds(channel chan *models.InterruptMsg) {
	// create a list of messages
	messages := make(map[string][]*pbft.PrePreparedMsg)

	for {
		// get raw interrupts
		intr := <-channel

		// ignore messages that are not preprepared
		if intr.Type != enum.IntrPrePrepared {
			continue
		}

		// extract the payload
		payload := intr.Payload.(*pbft.PrePreparedMsg)
		if _, ok := messages[payload.GetDigest()]; !ok { // build an array of map
			messages[payload.GetDigest()] = make([]*pbft.PrePreparedMsg, 0)
		}

		// append the payload
		messages[payload.GetDigest()] = append(messages[payload.GetDigest()], payload)

		// check for 2f+1 messages
		for _, value := range messages {
			if len(value) >= c.BFTCfg.Responses {
				return
			}
		}
	}
}
