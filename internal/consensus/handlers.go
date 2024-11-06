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
		// get raw PP packets
		raw := <-c.consensusHandlersTable[enum.PktPP]

		// parse the input message
		msg := raw.Payload.(*pbft.PrePrepareMsg)

		// get the digest of request
		digest := hashing.MD5(msg.GetRequest())

		// validate the message
		if !c.validatePrePrepareMsg(msg, digest) {
			c.Logger.Debug(
				"preprepare message is not valid",
				zap.Int("sequence number", raw.Sequence),
				zap.Int64("timestamp", msg.GetRequest().GetTransaction().GetTimestamp()),
			)
			continue
		}

		// update the request and set the status of preprepared
		c.Logs.SetRequest(raw.Sequence, msg.GetRequest())
		c.Logs.SetRequestStatus(raw.Sequence, pbft.RequestStatus_REQUEST_STATUS_PP)

		c.Logger.Debug(
			"preprepared a message",
			zap.Int("sequence number", raw.Sequence),
			zap.Int64("timestamp", msg.GetRequest().GetTransaction().GetTimestamp()),
		)

		// call preprepared RPC to notify the sender
		c.Client.PrePrepared(msg.GetNodeId(), &pbft.AckMsg{
			View:           int64(c.Memory.GetView()),
			SequenceNumber: int64(raw.Sequence),
			Digest:         msg.GetDigest(),
		})
	}
}

// prepareHandler gets gRPC packets of type P and handles them.
func (c *Consensus) prepareHandler() {
	for {
		// get raw P packets
		raw := <-c.consensusHandlersTable[enum.PktP]

		// parse the input message
		msg := raw.Payload.(*pbft.AckMsg)

		// get the message from our datastore
		message := c.Logs.GetRequest(raw.Sequence)
		if message == nil {
			c.Logger.Debug(
				"request not found",
				zap.Int("sequence number", raw.Sequence),
			)
			continue
		}

		// get the digest of input request
		digest := hashing.MD5(message)

		if !c.Memory.GetByzantine() { // byzantine nodes don't prepare messages
			if !c.validateAckMessage(msg, digest) {
				c.Logger.Debug(
					"prepare message is not valid",
					zap.Int("sequence number", raw.Sequence),
				)
				continue
			}

			// update the request and set the status of prepare
			c.Logs.SetRequestStatus(raw.Sequence, pbft.RequestStatus_REQUEST_STATUS_P)

			c.Logger.Debug(
				"prepared a message",
				zap.Int("sequence number", raw.Sequence),
				zap.Int64("timestamp", message.GetTransaction().GetTimestamp()),
			)
		}

		// call prepared RPC to notify the sender
		c.Client.Prepared(msg.GetNodeId(), &pbft.AckMsg{
			View:           int64(c.Memory.GetView()),
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
		c.Logs.SetRequestStatus(raw.Sequence, pbft.RequestStatus_REQUEST_STATUS_C)

		c.Logger.Debug(
			"request committed",
			zap.Int("sequence number", raw.Sequence),
		)

		// execute the request
		c.newExecutionGadget(raw.Sequence)
	}
}

// requestHandler gets a request message and performs the request handling logic.
func (c *Consensus) requestHandler(pkt interface{}) {
	// parse the input message
	msg := pkt.(*pbft.RequestMsg)

	// check if we had a request with the given timestamp
	if req := c.checkRequestExecution(msg.GetTransaction().GetTimestamp()); req != nil {
		c.Logger.Debug(
			"redundant request",
			zap.Int64("timestamp", req.GetTransaction().GetTimestamp()),
			zap.Int64("sequence number", req.GetSequenceNumber()),
		)

		c.Communication.SendReplyMsg(req, c.Memory.GetView())

		return
	}

	// check if the node is leader
	if c.getCurrentLeader() != c.Memory.GetNodeId() {
		c.Logger.Debug(
			"backup received a request",
			zap.String("current leader", c.getCurrentLeader()),
		)

		return // drop the request if not leader
	}

	// find a sequence number for this request
	sequence := c.Logs.InitRequest()

	// open a communication channel
	channel := make(chan *models.Packet)
	c.requestsHandlersTable[sequence] = channel

	c.Logger.Debug(
		"new sequence",
		zap.Int("sequence number", sequence),
	)

	// update request metadata
	msg.SequenceNumber = int64(sequence)
	msg.Status = pbft.RequestStatus_REQUEST_STATUS_UNSPECIFIED

	// store it into datastore
	c.Logs.SetRequest(sequence, msg)

	c.Logger.Debug(
		"new request received",
		zap.Int("sequence number", sequence),
		zap.Int64("timestamp", msg.GetTransaction().GetTimestamp()),
	)

	// send preprepare messages
	go c.Communication.SendPreprepareMsg(msg, c.Memory.GetView())

	// update our own status
	c.Logs.SetRequestStatus(sequence, pbft.RequestStatus_REQUEST_STATUS_PP)

	// wait for 2f+1 preprepared messages (count our own)
	count := c.Waiter.NewPrePreparedWaiter(channel, c.newAckGadget)
	c.Logger.Debug("received preprepared messages", zap.Int("messages", count+1))

	// broadcast to all using prepare
	go c.Communication.SendPrepareMsg(msg, c.Memory.GetView())

	// update our own status
	if !c.Memory.GetByzantine() {
		c.Logs.SetRequestStatus(sequence, pbft.RequestStatus_REQUEST_STATUS_P)
	}

	// wait for 2f+1 prepared messages (count our own)
	count = c.Waiter.NewPreparedWaiter(channel, c.newAckGadget)
	c.Logger.Debug("received prepared messages", zap.Int("messages", count+1))

	// broadcast to all using commit, make sure everyone get's it
	go c.Communication.SendCommitMsg(msg, c.Memory.GetView())

	// update our own status
	c.Logs.SetRequestStatus(sequence, pbft.RequestStatus_REQUEST_STATUS_C)

	// execute our own requests
	c.newExecutionGadget(sequence)
}
