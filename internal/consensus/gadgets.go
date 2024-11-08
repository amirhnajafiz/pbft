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
		return nil
	}

	// get the digest of input request
	digest := hashing.MD5(message)

	// validate the message
	if !c.validateMsg(digest, msg.GetDigest(), msg.GetView()) {
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
		c.executeRequest(msg)                                               // execute request
		c.logs.SetRequestStatus(index, pbft.RequestStatus_REQUEST_STATUS_E) // update the request and set the status of prepare

		c.communication.SendReplyMsg(msg, c.memory.GetView()) // send the reply message using helper functions
		c.logger.Debug("request executed", zap.Int("sequence number", index))

		index++
		if msg = c.logs.GetRequest(index); msg == nil || msg.GetStatus() != pbft.RequestStatus_REQUEST_STATUS_C {
			return
		}
	}
}

// newViewChangeGadget gets a count number of view-change messages and starts view change procedure.
func (c *Consensus) newViewChangeGadget() {
	// create our channel
	c.viewChangeGadgetChannel = make(chan *pbft.ViewChangeMsg)

	// change the view to stop processing requests
	c.memory.IncView()

	// send a view change message
	c.communication.SendViewChangeMsg(c.memory.GetView(), c.logs.GetSequenceNumber())

	// wait for 2f+1 messages

	// close our channel
	c.inViewChangeMode = false
	c.viewChangeGadgetChannel = nil

	// if the node is the leader, run a new leader gadget
	if c.getCurrentLeader() == c.memory.GetNodeId() {
		c.newLeaderGadget()
	}
}

func (c *Consensus) newLeaderGadget() {
	// send new-view messages
}
