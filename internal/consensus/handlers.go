package consensus

import (
	"github.com/f24-cse535/pbft/internal/utils/hashing"
	"github.com/f24-cse535/pbft/pkg/enum"
	"github.com/f24-cse535/pbft/pkg/models"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"
	"go.uber.org/zap"
)

// preprepareHandler gets gRPC packets of type PP and handles them.
func (c *Consensus) preprepareHandler() {
	for {
		// get raw PP packets and cast it
		raw := <-c.consensusHandlersTable[enum.PktPP]
		msg := raw.Payload.(*pbft.PrePrepareMsg)

		digest := hashing.MD5(msg.GetRequest()) // get the digest of request

		if !c.validateMsg(digest, msg.GetDigest(), msg.GetView()) {
			continue
		}

		// update the request and set the status of preprepared
		c.logs.SetRequest(raw.Sequence, msg.GetRequest())
		c.logs.SetRequestStatus(raw.Sequence, pbft.RequestStatus_REQUEST_STATUS_PP)

		// call preprepared RPC to notify the sender
		c.communication.Client().PrePrepared(msg.GetNodeId(), &pbft.AckMsg{
			View:           int64(c.memory.GetView()),
			SequenceNumber: int64(raw.Sequence),
			Digest:         msg.GetDigest(),
		})
	}
}

// prepareHandler gets gRPC packets of type P and handles them.
func (c *Consensus) prepareHandler() {
	for {
		// get raw P packets and cast them
		raw := <-c.consensusHandlersTable[enum.PktP]
		msg := raw.Payload.(*pbft.AckMsg)

		// get the message from our datastore
		message := c.logs.GetRequest(raw.Sequence)
		if message == nil {
			continue
		}

		digest := hashing.MD5(message) // get the digest of input request

		if !c.memory.GetByzantine() { // byzantine nodes don't prepare messages
			if !c.validateMsg(digest, msg.GetDigest(), msg.GetView()) {
				continue
			}

			// update the request and set the status of prepare
			c.logs.SetRequestStatus(raw.Sequence, pbft.RequestStatus_REQUEST_STATUS_P)
		}

		// call prepared RPC to notify the sender
		c.communication.Client().Prepared(msg.GetNodeId(), &pbft.AckMsg{
			View:           int64(c.memory.GetView()),
			SequenceNumber: int64(raw.Sequence),
			Digest:         digest,
		})
	}
}

// commitHandler gets gRPC packets of type C and handles them.
func (c *Consensus) commitHandler() {
	for {
		// get raw C packets
		raw := <-c.consensusHandlersTable[enum.PktCmt]

		// update the request and set the status of prepare
		c.logs.SetRequestStatus(raw.Sequence, pbft.RequestStatus_REQUEST_STATUS_C)

		// send the sequence to the execute handler
		c.executionChannel <- raw.Sequence
	}
}

// executeHandler gets sequence numbers from handlers and executes a request.
func (c *Consensus) executeHandler() {
	for {
		seq := <-c.executionChannel

		// execute the request
		c.newExecutionGadget(seq)
	}
}

// requestHandler gets a request message and performs the request handling logic.
func (c *Consensus) requestHandler(pkt *models.Packet) {
	// parse the input message
	msg := pkt.Payload.(*pbft.RequestMsg)

	// check if we had a request with the given timestamp
	if req := c.checkRequestExecution(msg.GetTransaction().GetTimestamp()); req != nil {
		c.communication.SendReplyMsg(req, c.memory.GetView())

		return
	}

	// check if the node is leader
	if c.getCurrentLeader() != c.memory.GetNodeId() {
		// send the request to leader
		c.communication.Client().Request(c.getCurrentLeader(), msg)

		return
	}

	// find a sequence number for this request
	sequence := c.logs.InitRequest()

	// open a communication channel
	channel := make(chan *models.Packet, c.cfg.Total*2)
	c.requestsHandlersTable[sequence] = channel

	// update request metadata
	msg.SequenceNumber = int64(sequence)
	msg.Status = pbft.RequestStatus_REQUEST_STATUS_UNSPECIFIED

	// store it into datastore
	c.logs.SetRequest(sequence, msg)

	c.logger.Debug("new request got into the system", zap.Int("sequece", sequence), zap.Int64("time", msg.Transaction.GetTimestamp()))

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

// timerHandler creates a new timer and monitors the timer.
func (c *Consensus) timerHandler() {
	// stop the timer so that others can start it
	c.viewTimer.Stop()

	for {
		c.viewTimer.Notify()
		c.logger.Debug("timer expired")
	}
}
