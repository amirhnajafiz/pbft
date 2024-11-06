package consensus

import (
	"github.com/f24-cse535/pbft/internal/utils/hashing"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"

	"go.uber.org/zap"
)

// newAckGadget validates a preprepared or prepared ack message.
func (c *Consensus) newAckGadget(msg *pbft.AckMsg) *pbft.AckMsg {
	// get the message from our datastore
	message := c.logs.GetRequest(int(msg.GetSequenceNumber()))
	if message == nil {
		c.logger.Debug(
			"request not found",
			zap.Int64("sequence number", msg.GetSequenceNumber()),
		)
		return nil
	}

	// get the digest of input request
	digest := hashing.MD5(message)

	// validate the message
	if !c.validateAckMessage(msg, digest) {
		c.logger.Debug(
			"ack message is not valid",
			zap.Int64("sequence number", msg.GetSequenceNumber()),
		)
		return nil
	}

	return msg
}

// newExecutionGadget gets a sequence number and performs the execution logic.
func (c *Consensus) newExecutionGadget(sequence int) {
	if !c.memory.GetByzantine() && !c.canExecuteRequest(sequence) {
		c.logger.Debug(
			"cannot execute the request yet",
			zap.Int("sequence number", sequence),
		)
		return
	}

	// follow sequence until one is not committed, execute them
	index := sequence
	msg := c.logs.GetRequest(index)

	for {
		c.executeRequest(msg) // execute request

		// update the request and set the status of prepare
		c.logs.SetRequestStatus(index, pbft.RequestStatus_REQUEST_STATUS_E)

		c.communication.SendReplyMsg(msg, c.memory.GetView()) // send the reply message using helper functions

		c.logger.Info("request executed", zap.Int("sequence number", index))

		index++
		if msg = c.logs.GetRequest(index); msg == nil || msg.GetStatus() != pbft.RequestStatus_REQUEST_STATUS_C {
			return
		}
	}
}
