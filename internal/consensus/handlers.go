package consensus

import (
	"github.com/f24-cse535/pbft/internal/utils/hashing"
	"github.com/f24-cse535/pbft/pkg/enum"
	"github.com/f24-cse535/pbft/pkg/models"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"

	"go.uber.org/zap"
)

// handleExecute gets a sequence and message and does the execution.
func (c *Consensus) handleExecute(sequence int) {
	// check if this sequence is executable
	if !c.canExecute(sequence) {
		c.Logger.Info(
			"request cannot get executed yet",
			zap.Int("sequence number", sequence),
		)
		return
	}

	// follow sequence until one is not committed, execute them
	index := sequence
	for {
		if msg := c.Logs.GetRequest(index); msg != nil && msg.GetStatus() == pbft.RequestStatus_REQUEST_STATUS_C {
			// execute request
			c.executeRequest(msg)

			// update the request and set the status of prepare
			c.Logs.SetRequestStatus(index, pbft.RequestStatus_REQUEST_STATUS_E)

			// send the reply message using helper functions
			c.helpSendReply(msg)

			c.Logger.Info(
				"request executed",
				zap.Int64("sequence number", msg.GetSequenceNumber()),
			)
		} else {
			break
		}

		index++
	}
}

// handleCommit gets a commit message, changes it status and calls the handleExecute.
func (c *Consensus) handleCommit(pkt interface{}) {
	// parse the input message
	msg := pkt.(*pbft.AckMsg)
	sequence := int(msg.GetSequenceNumber())

	// update the request and set the status of prepare
	c.Logs.SetRequestStatus(sequence, pbft.RequestStatus_REQUEST_STATUS_C)

	// message committed
	c.Logger.Debug(
		"message committed",
		zap.Int("sequence number", sequence),
	)

	// call execute handler
	c.handleExecute(sequence)
}

// handle preprepare accepts a preprepare message and validates it to call preprepared RPC.
func (c *Consensus) handlePrePrepare(pkt interface{}) {
	// parse the input message
	msg := pkt.(*pbft.PrePrepareMsg)

	// get the digest of input request
	digest := hashing.MD5(msg.GetRequest())

	// validate the message
	if !c.validatePrePrepareMsg(msg, digest) {
		c.Logger.Info(
			"preprepare message is not valid",
			zap.Int64("sequence number", msg.GetSequenceNumber()),
			zap.Int64("timestamp", msg.GetRequest().GetTransaction().GetTimestamp()),
		)
		return
	}

	// update the request and set the status of preprepared
	c.Logs.SetRequest(int(msg.GetSequenceNumber()), msg.GetRequest())
	c.Logs.SetRequestStatus(int(msg.GetSequenceNumber()), pbft.RequestStatus_REQUEST_STATUS_PP)

	c.Logger.Debug(
		"preprepared a message",
		zap.Int64("sequence number", msg.GetSequenceNumber()),
		zap.Int64("timestamp", msg.GetRequest().GetTransaction().GetTimestamp()),
	)

	// call preprepared message
	go c.Client.PrePrepared(msg.GetNodeId(), &pbft.AckMsg{
		View:           int64(c.Memory.GetView()),
		SequenceNumber: msg.GetSequenceNumber(),
		Digest:         msg.GetDigest(),
	})
}

// handlePrePrepared gets preprepared message and passes it to the correct request handler.
func (c *Consensus) handlePrePrepared(pkt interface{}) {
	// parse the input message
	msg := pkt.(*pbft.AckMsg)

	// get the message from our datastore
	message := c.Logs.GetRequest(int(msg.GetSequenceNumber()))
	if message == nil {
		c.Logger.Info(
			"request not found",
			zap.Int64("sequence number", msg.GetSequenceNumber()),
		)
		return
	}

	// get the digest of input request
	digest := hashing.MD5(message)

	// validate the message
	if !c.validateAckMessage(msg, digest) {
		c.Logger.Info(
			"preprepared message is not valid",
			zap.Int64("sequence number", msg.GetSequenceNumber()),
		)
		return
	}

	c.Logger.Debug(
		"preprepared received",
		zap.Int64("sequence number", message.GetSequenceNumber()),
		zap.Int64("timestamp", message.GetTransaction().GetTimestamp()),
	)

	// publish over correct request handler
	c.channels[int(msg.GetSequenceNumber())] <- &models.InterruptMsg{
		Type:    enum.IntrPrePrepared,
		Payload: msg,
	}
}

// handlePrepare gets a prepare message and updates a log status.
func (c *Consensus) handlePrepare(pkt interface{}) {
	// parse the input message
	msg := pkt.(*pbft.AckMsg)

	// get the message from our datastore
	message := c.Logs.GetRequest(int(msg.GetSequenceNumber()))
	if message == nil {
		c.Logger.Info(
			"request not found",
			zap.Int64("sequence number", msg.GetSequenceNumber()),
		)
		return
	}

	// get the digest of input request
	digest := hashing.MD5(message)

	if !c.Memory.GetByzantine() { // byzantine nodes don't prepare messages
		// validate the message
		if !c.validateAckMessage(msg, digest) {
			c.Logger.Info(
				"prepare message is not valid",
				zap.Int64("sequence number", msg.GetSequenceNumber()),
			)
			return
		}

		// update the request and set the status of prepare
		c.Logs.SetRequestStatus(int(msg.GetSequenceNumber()), pbft.RequestStatus_REQUEST_STATUS_P)

		c.Logger.Debug(
			"prepared a message",
			zap.Int64("sequence number", message.GetSequenceNumber()),
			zap.Int64("timestamp", message.GetTransaction().GetTimestamp()),
		)
	}

	// call prepared message
	go c.Client.Prepared(msg.GetNodeId(), &pbft.AckMsg{
		View:           int64(c.Memory.GetView()),
		SequenceNumber: msg.GetSequenceNumber(),
		Digest:         digest,
	})
}

// handlePrepared gets a prepared message and sends it to the gith request handler.
func (c *Consensus) handlePrepared(pkt interface{}) {
	// parse the input message
	msg := pkt.(*pbft.AckMsg)

	// get the message from our datastore
	message := c.Logs.GetRequest(int(msg.GetSequenceNumber()))
	if message == nil {
		c.Logger.Info(
			"request not found",
			zap.Int64("sequence number", msg.GetSequenceNumber()),
		)
		return
	}

	// get the digest of input request
	digest := hashing.MD5(message)

	// validate the message
	if !c.validateAckMessage(msg, digest) {
		c.Logger.Info(
			"prepared message is not valid",
			zap.Int64("sequence number", msg.GetSequenceNumber()),
		)
		return
	}

	c.Logger.Debug(
		"prepared received",
		zap.Int64("sequence number", message.GetSequenceNumber()),
		zap.Int64("timestamp", message.GetTransaction().GetTimestamp()),
	)

	// publish over correct request handler
	c.channels[int(msg.GetSequenceNumber())] <- &models.InterruptMsg{
		Type:    enum.IntrPrepared,
		Payload: msg,
	}
}

// handleReply gets a reply message and passes it to the right request handler.
func (c *Consensus) handleReply(pkt interface{}) {
	// parse the input message
	msg := pkt.(*pbft.ReplyMsg)

	// ignore the messages that are not for this node
	if msg.GetClientId() != c.Memory.GetNodeId() {
		c.Logger.Debug(
			"reply dropped",
			zap.String("to", msg.GetClientId()),
			zap.String("whoami", c.Memory.GetNodeId()),
		)
		return
	}

	// publish over correct request handler
	c.inTransactionChannel <- &models.InterruptMsg{
		Type:    enum.IntrReply,
		Payload: msg,
	}
}

// handle request accepts a new request and creates a go-routine
// to collect all preprepared and prepared messages.
func (c *Consensus) handleRequest(pkt interface{}) {
	// parse the input message
	msg := pkt.(*pbft.RequestMsg)

	// check if we had a request with the given timestamp
	if rsp := c.isRequestExecuted(msg.GetTransaction().GetTimestamp()); rsp != nil {
		c.helpSendReply(rsp) // send the reply message using helper function
		return
	}

	// store the log place
	seqn := c.Logs.InitRequest()

	// create our channel for input messages
	c.channels[seqn] = make(chan *models.InterruptMsg)

	// need to create a go-routine to not block the request and call helper functions
	go c.helpProcessTransaction(seqn, msg)
}

// handle transaction checks a new transaction to call request RPC.
func (c *Consensus) handleTransaction(pkt interface{}) {
	defer func() {
		// reset the channel when transaction is done
		c.inTransactionChannel = nil
	}()

	// parse the input message
	msg := pkt.(*pbft.TransactionMsg)

	// send the request using helper functions
	c.helpSendRequest(msg)

	c.Logger.Debug(
		"request is sent",
		zap.Int64("timestamp", msg.GetTimestamp()),
		zap.String("sender", msg.GetSender()),
		zap.String("receiver", msg.GetReciever()),
		zap.Int64("amount", msg.GetAmount()),
	)

	// get the response by calling helper functions
	resp := c.helpReceiveResponse(msg)

	c.Logger.Debug(
		"received reply message",
		zap.Int64("timestamp", resp.GetTimestamp()),
		zap.Int64("sequence number", resp.GetSequenceNumber()),
	)

	// reset the view
	c.Memory.SetView(int(resp.GetView()))

	// return the response in channel
	c.outTransactionChannel <- resp.GetResponse()
}
