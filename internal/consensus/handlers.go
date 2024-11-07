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
			c.logger.Debug(
				"preprepare message is not valid",
				zap.Int("sequence number", raw.Sequence),
				zap.Int64("timestamp", msg.GetRequest().GetTransaction().GetTimestamp()),
			)
			continue
		}

		// update the request and set the status of preprepared
		c.logs.SetRequest(raw.Sequence, msg.GetRequest())
		c.logs.SetRequestStatus(raw.Sequence, pbft.RequestStatus_REQUEST_STATUS_PP)

		c.logger.Debug(
			"preprepared a message",
			zap.Int("sequence number", raw.Sequence),
			zap.Int64("timestamp", msg.GetRequest().GetTransaction().GetTimestamp()),
		)

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
			c.logger.Debug(
				"request not found",
				zap.Int("sequence number", raw.Sequence),
			)
			continue
		}

		digest := hashing.MD5(message) // get the digest of input request

		if !c.memory.GetByzantine() { // byzantine nodes don't prepare messages
			if !c.validateMsg(digest, msg.GetDigest(), msg.GetView()) {
				c.logger.Debug(
					"prepare message is not valid",
					zap.Int("sequence number", raw.Sequence),
				)
				continue
			}

			// update the request and set the status of prepare
			c.logs.SetRequestStatus(raw.Sequence, pbft.RequestStatus_REQUEST_STATUS_P)

			c.logger.Debug(
				"prepared a message",
				zap.Int("sequence number", raw.Sequence),
				zap.Int64("timestamp", message.GetTransaction().GetTimestamp()),
			)
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

		c.logger.Debug(
			"request committed",
			zap.Int("sequence number", raw.Sequence),
		)

		// execute the request
		c.newExecutionGadget(raw.Sequence)
	}
}

// requestHandler gets a request message and performs the request handling logic.
func (c *Consensus) requestHandler(pkt *models.Packet) {
	// parse the input message
	msg := pkt.Payload.(*pbft.RequestMsg)

	// check if we had a request with the given timestamp
	if req := c.checkRequestExecution(msg.GetTransaction().GetTimestamp()); req != nil {
		c.logger.Debug(
			"redundant request",
			zap.Int64("timestamp", req.GetTransaction().GetTimestamp()),
			zap.Int64("sequence number", req.GetSequenceNumber()),
		)

		c.communication.SendReplyMsg(req, c.memory.GetView())

		return
	}

	// check if the node is leader
	if c.getCurrentLeader() != c.memory.GetNodeId() {
		c.logger.Debug(
			"backup received a request",
			zap.String("current leader", c.getCurrentLeader()),
		)

		// send the request to leader
		c.communication.Client().Request(c.getCurrentLeader(), msg)

		return
	}

	// find a sequence number for this request
	sequence := c.logs.InitRequest()

	// open a communication channel
	channel := make(chan *models.Packet, c.cfg.Total*2)
	c.requestsHandlersTable[sequence] = channel

	c.logger.Debug(
		"new sequence",
		zap.Int("sequence number", sequence),
	)

	// update request metadata
	msg.SequenceNumber = int64(sequence)
	msg.Status = pbft.RequestStatus_REQUEST_STATUS_UNSPECIFIED

	// store it into datastore
	c.logs.SetRequest(sequence, msg)

	c.logger.Debug(
		"new request received",
		zap.Int("sequence number", sequence),
		zap.Int64("timestamp", msg.GetTransaction().GetTimestamp()),
	)

	// send preprepare messages
	go c.communication.SendPreprepareMsg(msg, c.memory.GetView())

	// update our own status
	c.logs.SetRequestStatus(sequence, pbft.RequestStatus_REQUEST_STATUS_PP)

	// wait for 2f+1 preprepared messages (count our own)
	count := c.waiter.NewPrePreparedWaiter(channel, c.newAckGadget)
	c.logger.Debug("received preprepared messages", zap.Int("messages", count+1))

	// broadcast to all using prepare
	go c.communication.SendPrepareMsg(msg, c.memory.GetView())

	// update our own status
	if !c.memory.GetByzantine() {
		c.logs.SetRequestStatus(sequence, pbft.RequestStatus_REQUEST_STATUS_P)
	}

	// wait for 2f+1 prepared messages (count our own)
	count = c.waiter.NewPreparedWaiter(channel, c.newAckGadget)
	c.logger.Debug("received prepared messages", zap.Int("messages", count+1))

	// broadcast to all using commit, make sure everyone get's it
	go c.communication.SendCommitMsg(msg, c.memory.GetView())

	// delete our input channel as soon as possible
	delete(c.requestsHandlersTable, sequence)

	// update our own status
	c.logs.SetRequestStatus(sequence, pbft.RequestStatus_REQUEST_STATUS_C)

	// execute our own requests
	c.newExecutionGadget(sequence)
}
