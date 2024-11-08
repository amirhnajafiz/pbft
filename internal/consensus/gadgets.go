package consensus

import (
	"fmt"

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
	digest := hashing.MD5Req(message)

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
	// change the view to stop processing requests
	c.memory.IncView()
	view := c.memory.GetView()

	// send a view change message
	if count := c.communication.SendViewChangeMsg(view, c.logs.GetSequenceNumber()); count < c.cfg.Majority {
		c.logger.Info("not enough available servers to start view change", zap.Int("live servers", count))

		return
	}

	c.logger.Debug("view change started")

	// wait for 2f+1 messages
	for {
		msg := <-c.viewChangeGadgetChannel
		c.logs.AppendViewChange(view, msg)

		if len(c.logs.GetViewChanges(view)) >= c.cfg.Majority-1 {
			break
		}
	}

	c.logger.Debug("view change end")

	// close our channel
	c.inViewChangeMode = false
	c.viewChangeGadgetChannel = nil

	// if the node is the leader, run a new leader gadget
	if c.getCurrentLeader() == c.memory.GetNodeId() {
		c.logger.Debug("new leader", zap.String("id", c.memory.GetNodeId()))
		c.newLeaderGadget()
	}
}

// newLeaderGadget performs the procedure of new leader.
func (c *Consensus) newLeaderGadget() {
	// get all previous messages
	messages := c.logs.GetViewChanges(c.memory.GetView())

	minSequence := c.logs.GetSequenceNumber()
	maxSequence := c.logs.GetSequenceNumber()

	// loop in all messages
	for _, msg := range messages {
		sequence := int(msg.GetSequenceNumber())
		if sequence <= minSequence {
			minSequence = sequence
		}

		if sequence >= maxSequence {
			maxSequence = sequence
		}
	}

	// create an array to store sequences
	requests := make([]int, 0)

	// collect all requets that are prepared
	for i := minSequence; i <= maxSequence; i++ {
		if tmp := c.logs.GetRequest(i); tmp != nil && tmp.GetStatus().Number() > pbft.RequestStatus_REQUEST_STATUS_PP.Number() {
			requests = append(requests, i)
		}
	}

	// broadcase a new-view message
	fmt.Println(requests)
}
