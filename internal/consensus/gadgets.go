package consensus

import (
	"github.com/f24-cse535/pbft/internal/utils/hashing"
	"github.com/f24-cse535/pbft/pkg/models"
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

func (c *Consensus) newProcessingGadet(sequence int, msg *pbft.RequestMsg) {
	// open a communication channel
	channel := make(chan *models.Packet, c.cfg.Total*2)
	c.requestsHandlersTable[sequence] = channel

	// send preprepare messages
	go c.communication.SendPreprepareMsg(msg, c.memory.GetView())

	// update our own status
	c.logs.SetRequestStatus(sequence, pbft.RequestStatus_REQUEST_STATUS_PP)

	// wait for 2f+1 preprepared messages (count our own)
	c.waiter.NewPrePreparedWaiter(channel, c.newAckGadget)

	// broadcast to all using prepare
	go c.communication.SendPrepareMsg(msg, c.memory.GetView())

	// update our own status
	if !c.memory.GetByzantine() {
		c.logs.SetRequestStatus(sequence, pbft.RequestStatus_REQUEST_STATUS_P)
	}

	// wait for 2f+1 prepared messages (count our own)
	c.waiter.NewPreparedWaiter(channel, c.newAckGadget)

	// broadcast to all using commit, make sure everyone get's it
	go c.communication.SendCommitMsg(msg, c.memory.GetView())

	// delete our input channel as soon as possible
	delete(c.requestsHandlersTable, sequence)

	// update our own status
	c.logs.SetRequestStatus(sequence, pbft.RequestStatus_REQUEST_STATUS_C)

	// send the sequence to the execute handler
	c.executionChannel <- sequence
}

// newExecutionGadget gets a sequence number and performs the execution logic.
func (c *Consensus) newExecutionGadget(sequence int) {
	if !c.memory.GetByzantine() && !c.canExecuteRequest(sequence) {
		return
	}

	// follow sequence until one is not committed, execute them
	index := sequence
	msg := c.logs.GetRequest(index)

	// don't reexecute a request
	if msg.GetStatus() == pbft.RequestStatus_REQUEST_STATUS_E {
		return
	}

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
	sequence := c.logs.GetSequenceNumber()
	preprepares := c.logs.GetPreprepares(sequence)

	c.logs.AppendViewChange(view, &pbft.ViewChangeMsg{
		View:           int64(view),
		SequenceNumber: int64(sequence),
		Preprepares:    preprepares,
	})

	// send a view change message
	if count := c.communication.SendViewChangeMsg(view, sequence, preprepares); count < c.cfg.Majority {
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
	// get current view
	view := c.memory.GetView()

	// get all previous messages
	messages := c.logs.GetViewChanges(view)

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
	requests := make([]*pbft.PrePrepareMsg, 0)

	// collect all requets that are prepared
	for i := minSequence; i <= maxSequence; i++ {
		for _, msg := range messages {
			for _, pp := range msg.GetPreprepares() {
				if int(pp.GetView()) == i {
					requests = append(requests, pp)
				}
			}
		}
	}

	// send new view
	c.communication.SendNewViewMsg(c.memory.GetView(), requests)

	// start the protocol for every request
	for _, req := range requests {
		c.newProcessingGadet(int(req.GetSequenceNumber()), req.Request)
	}
}
