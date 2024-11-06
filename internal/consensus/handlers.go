package consensus

import (
	"github.com/f24-cse535/pbft/internal/utils/hashing"
	"github.com/f24-cse535/pbft/pkg/enum"
	"github.com/f24-cse535/pbft/pkg/models"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"

	"go.uber.org/zap"
)

func (c *Consensus) preprepareHandler() {
	for {

	}
}

func (c *Consensus) prepareHandler() {
	for {

	}
}

func (c *Consensus) commitHandler() {
	for {

	}
}

func (c *Consensus) requestHandler() {
	for {

	}
}

// handleExecute gets a sequence and message and does the execution.
func (c *Consensus) handleExecute(sequence int) {
	// check if this sequence is executable
	if !c.Memory.GetByzantine() && !c.canExecute(sequence) {
		c.Logger.Info(
			"request cannot get executed yet",
			zap.Int("sequence number", sequence),
		)
		return
	}

	// follow sequence until one is not committed, execute them
	index := sequence
	msg := c.Logs.GetRequest(index)

	for {
		c.executeRequest(msg)  // execute request
		c.viewTimerRm <- index // notify the timer

		// update the request and set the status of prepare
		c.Logs.SetRequestStatus(index, pbft.RequestStatus_REQUEST_STATUS_E)

		c.promiseReply(msg) // send the reply message using helper functions

		c.Logger.Info(
			"request executed",
			zap.Int64("sequence number", msg.GetSequenceNumber()),
		)

		index++

		if msg = c.Logs.GetRequest(index); msg == nil || msg.GetStatus() != pbft.RequestStatus_REQUEST_STATUS_C {
			break
		}
	}
}

// handleCommit gets a commit message, changes it status and calls the handleExecute.
func (c *Consensus) handleCommit(pkt interface{}) {
	// parse the input message
	msg := pkt.(*pbft.AckMsg)
	sequence := int(msg.GetSequenceNumber())
	c.viewTimerIn <- sequence

	// update the request and set the status of prepare
	c.Logs.SetRequestStatus(sequence, pbft.RequestStatus_REQUEST_STATUS_C)

	// message committed
	c.Logger.Debug(
		"message committed",
		zap.Int("sequence number", sequence),
	)

	c.handleExecute(sequence) // call execute handler
}

// handle preprepare accepts a preprepare message and validates it to call preprepared RPC.
func (c *Consensus) handlePrePrepare(pkt interface{}) {
	// parse the input message
	msg := pkt.(*pbft.PrePrepareMsg)

	c.Logger.Debug("inside preprepare handler", zap.String("node", c.Memory.GetNodeId()))

	// get the digest of input request
	digest := hashing.MD5(msg.GetRequest())

	c.Logger.Debug("inside preprepare handler after digest", zap.String("node", c.Memory.GetNodeId()))

	// validate the message
	if !c.validatePrePrepareMsg(msg, digest) {
		c.Logger.Info(
			"preprepare message is not valid",
			zap.Int64("sequence number", msg.GetSequenceNumber()),
			zap.Int64("timestamp", msg.GetRequest().GetTransaction().GetTimestamp()),
		)
		return
	}

	c.Logger.Debug("inside preprepare handler after validate", zap.String("node", c.Memory.GetNodeId()))

	// update the request and set the status of preprepared
	c.Logs.SetRequest(int(msg.GetSequenceNumber()), msg.GetRequest())
	c.Logs.SetRequestStatus(int(msg.GetSequenceNumber()), pbft.RequestStatus_REQUEST_STATUS_PP)

	c.Logger.Debug(
		"preprepared a message",
		zap.Int64("sequence number", msg.GetSequenceNumber()),
		zap.Int64("timestamp", msg.GetRequest().GetTransaction().GetTimestamp()),
	)

	c.viewTimerIn <- int(msg.GetSequenceNumber())

	// call preprepared message
	c.Client.PrePrepared(msg.GetNodeId(), &pbft.AckMsg{
		View:           int64(c.Memory.GetView()),
		SequenceNumber: msg.GetSequenceNumber(),
		Digest:         msg.GetDigest(),
	})
}

// handlePrepare gets a prepare message and updates a log status.
func (c *Consensus) handlePrepare(pkt interface{}) {
	// parse the input message
	msg := pkt.(*pbft.AckMsg)
	c.viewTimerIn <- int(msg.GetSequenceNumber())

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
	c.Client.Prepared(msg.GetNodeId(), &pbft.AckMsg{
		View:           int64(c.Memory.GetView()),
		SequenceNumber: msg.GetSequenceNumber(),
		Digest:         digest,
	})
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

	if msg.GetTimestamp() != c.Memory.GetTimestamp() {
		c.Logger.Debug(
			"old reply dropped",
			zap.Int64("timestamp", msg.GetTimestamp()),
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
		c.Logger.Debug(
			"redundant request",
			zap.Int64("timestamp", rsp.GetTransaction().GetTimestamp()),
			zap.Int64("sequence number", rsp.GetSequenceNumber()),
		)

		c.promiseReply(rsp) // send the reply message using helper function
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

	// store the log place
	seqn := c.Logs.InitRequest()

	c.Logger.Debug(
		"new sequence",
		zap.Int("sequence number", seqn),
	)

	// create our channel for input messages
	c.channels[seqn] = make(chan *models.InterruptMsg)

	//call helper functions to process the transaction
	c.promiseProcess(seqn, msg)
}

// handle transaction checks a new transaction to call request RPC.
func (c *Consensus) handleTransaction(pkt interface{}) {
	defer func() {
		c.Memory.SetTimestamp(0)

		// reset the channel when transaction is done
		c.inTransactionChannel = nil
	}()

	// parse the input message
	msg := pkt.(*pbft.TransactionMsg)

	c.Logger.Debug("new transaction", zap.Int64("timestamp", msg.GetTimestamp()))

	// send the request using helper functions
	c.promiseRequest(msg)

	// after the request is sent, update the current timestamp
	c.Memory.SetTimestamp(msg.GetTimestamp())

	// get the response by calling helper functions
	resp := c.promiseReceive(msg)

	c.Logger.Debug("received reply", zap.Int64("sequence number", resp.GetSequenceNumber()))

	// reset the view
	c.Memory.SetView(int(resp.GetView()))

	// return the response in channel
	c.outTransactionChannel <- resp.GetResponse()
}
