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

// newExecutionGadget gets a sequence number and performs the execution logic.
func (c *Consensus) newExecutionGadget(sequence int) {
	if !c.Memory.GetByzantine() && !c.canExecute(sequence) {
		c.Logger.Debug(
			"cannot execute the request yet",
			zap.Int("sequence number", sequence),
		)
		return
	}

	// follow sequence until one is not committed, execute them
	index := sequence
	msg := c.Logs.GetRequest(index)

	for {
		c.executeRequest(msg) // execute request

		// update the request and set the status of prepare
		c.Logs.SetRequestStatus(index, pbft.RequestStatus_REQUEST_STATUS_E)

		c.promiseReply(msg) // send the reply message using helper functions

		c.Logger.Info(
			"request executed",
			zap.Int64("sequence number", msg.GetSequenceNumber()),
		)

		index++

		if msg = c.Logs.GetRequest(index); msg == nil || msg.GetStatus() != pbft.RequestStatus_REQUEST_STATUS_C {
			return
		}
	}
}
