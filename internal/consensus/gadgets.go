package consensus

import (
	"github.com/f24-cse535/pbft/internal/utils/hashing"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"

	"go.uber.org/zap"
)

// newAckGadget validates a preprepared or prepared ack message.
func (c *Consensus) newAckGadget(msg *pbft.AckMsg) *pbft.AckMsg {
	// get the message from our datastore
	message := c.Logs.GetRequest(int(msg.GetSequenceNumber()))
	if message == nil {
		c.Logger.Debug(
			"request not found",
			zap.Int64("sequence number", msg.GetSequenceNumber()),
		)
		return nil
	}

	// get the digest of input request
	digest := hashing.MD5(message)

	// validate the message
	if !c.validateAckMessage(msg, digest) {
		c.Logger.Debug(
			"ack message is not valid",
			zap.Int64("sequence number", msg.GetSequenceNumber()),
		)
		return nil
	}

	return msg
}
