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
				zap.Int64("sequence number", msg.GetSequenceNumber()),
				zap.Int64("timestamp", msg.GetRequest().GetTransaction().GetTimestamp()),
			)
			continue
		}

		// update the request and set the status of preprepared
		c.Logs.SetRequest(int(msg.GetSequenceNumber()), msg.GetRequest())
		c.Logs.SetRequestStatus(int(msg.GetSequenceNumber()), pbft.RequestStatus_REQUEST_STATUS_PP)

		c.Logger.Debug(
			"preprepared a message",
			zap.Int64("sequence number", msg.GetSequenceNumber()),
			zap.Int64("timestamp", msg.GetRequest().GetTransaction().GetTimestamp()),
		)

		// call preprepared RPC to notify the sender
		c.Client.PrePrepared(msg.GetNodeId(), &pbft.AckMsg{
			View:           int64(c.Memory.GetView()),
			SequenceNumber: msg.GetSequenceNumber(),
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
		message := c.Logs.GetRequest(int(msg.GetSequenceNumber()))
		if message == nil {
			c.Logger.Debug(
				"request not found",
				zap.Int64("sequence number", msg.GetSequenceNumber()),
			)
			continue
		}

		// get the digest of input request
		digest := hashing.MD5(message)

		if !c.Memory.GetByzantine() { // byzantine nodes don't prepare messages
			// validate the message
			if !c.validateAckMessage(msg, digest) {
				c.Logger.Debug(
					"prepare message is not valid",
					zap.Int64("sequence number", msg.GetSequenceNumber()),
				)
				continue
			}

			// update the request and set the status of prepare
			c.Logs.SetRequestStatus(int(msg.GetSequenceNumber()), pbft.RequestStatus_REQUEST_STATUS_P)

			c.Logger.Debug(
				"prepared a message",
				zap.Int64("sequence number", message.GetSequenceNumber()),
				zap.Int64("timestamp", message.GetTransaction().GetTimestamp()),
			)
		}

		// call prepared RPC to notify the sender
		c.Client.Prepared(msg.GetNodeId(), &pbft.AckMsg{
			View:           int64(c.Memory.GetView()),
			SequenceNumber: msg.GetSequenceNumber(),
			Digest:         digest,
		})
	}
}

// commitHandler gets gRPC packets of type C and handles them.
func (c *Consensus) commitHandler() {
	for {
		// get raw C packets
		raw := <-c.consensusHandlersTable[enum.PktCmt]

		// parse the input message
		msg := raw.Payload.(*pbft.AckMsg)

		// update the request and set the status of prepare
		c.Logs.SetRequestStatus(int(msg.GetSequenceNumber()), pbft.RequestStatus_REQUEST_STATUS_C)

		c.Logger.Debug(
			"request committed",
			zap.Int64("sequence number", msg.GetSequenceNumber()),
		)

		// execute the request
		c.newExecutionGadget(int(msg.GetSequenceNumber()))
	}
}

func (c *Consensus) requestHandler() {
	for {

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
