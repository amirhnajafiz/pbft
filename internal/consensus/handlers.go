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

			// send the reply message
			c.Client.Reply(msg.GetClientId(), &pbft.ReplyMsg{
				SequenceNumber: msg.GetSequenceNumber(),
				View:           int64(c.Memory.GetView()),
				Timestamp:      msg.GetTransaction().GetTimestamp(),
				ClientId:       msg.GetClientId(),
				Response:       msg.GetResponse().GetText(),
			})

			c.Logger.Info(
				"request executed and reply sent",
				zap.String("client", msg.GetClientId()),
				zap.Int64("sequence number", msg.GetSequenceNumber()),
				zap.Int64("timestamp", msg.GetTransaction().GetTimestamp()),
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
	msg := pkt.(*pbft.CommitMsg)
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

	// validate the message
	if !c.validatePrePrepareMsg(msg) {
		c.Logger.Info(
			"preprepare message is not valid",
			zap.Int64("sequence number", msg.GetSequenceNumber()),
			zap.Int64("timestamp", msg.GetRequest().GetTransaction().GetTimestamp()),
		)
		return
	}

	// update the request and set the status of preprepared
	c.Logs.SetRequestStatus(int(msg.GetSequenceNumber()), pbft.RequestStatus_REQUEST_STATUS_PP)

	c.Logger.Debug(
		"preprepared a message",
		zap.Int64("sequence number", msg.GetSequenceNumber()),
		zap.Int64("timestamp", msg.GetRequest().GetTransaction().GetTimestamp()),
	)

	// call preprepared message
	go c.Client.PrePrepared(msg.GetNodeId(), &pbft.PrePreparedMsg{
		Request:        msg.GetRequest(),
		View:           int64(c.Memory.GetView()),
		SequenceNumber: msg.GetSequenceNumber(),
	})
}

// handlePrePrepared gets preprepared message and passes it to the correct request handler.
func (c *Consensus) handlePrePrepared(pkt interface{}) {
	// parse the input message
	msg := pkt.(*pbft.PrePreparedMsg)

	// get the message from our datastore
	message := c.Logs.GetRequest(int(msg.GetSequenceNumber()))
	if message == nil {
		c.Logger.Info(
			"request not found",
			zap.Int64("sequence number", msg.GetSequenceNumber()),
			zap.Int64("timestamp", msg.GetRequest().GetTransaction().GetTimestamp()),
		)
		return
	}

	// set the digest message
	msg.Digest = hashing.MD5(message)

	// validate the message
	if !c.validatePrePreparedMsg(msg) {
		c.Logger.Info(
			"preprepared message is not valid",
			zap.Int64("sequence number", msg.GetSequenceNumber()),
			zap.Int64("timestamp", msg.GetRequest().GetTransaction().GetTimestamp()),
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
	msg := pkt.(*pbft.PrepareMsg)

	// get the message from our datastore
	message := c.Logs.GetRequest(int(msg.GetSequenceNumber()))
	if message == nil {
		c.Logger.Info(
			"request not found",
			zap.Int64("sequence number", msg.GetSequenceNumber()),
		)
		return
	}

	// get digest message
	digest := hashing.MD5(message)

	// validate the message
	if !c.validatePrepareMsg(digest, msg) {
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

	// call prepared message
	go c.Client.Prepared(msg.GetNodeId(), &pbft.PreparedMsg{
		View:           int64(c.Memory.GetView()),
		SequenceNumber: msg.GetSequenceNumber(),
		Digest:         digest,
	})
}

// handlePrepared gets a prepared message and sends it to the gith request handler.
func (c *Consensus) handlePrepared(pkt interface{}) {
	// parse the input message
	msg := pkt.(*pbft.PreparedMsg)

	// get the message from our datastore
	message := c.Logs.GetRequest(int(msg.GetSequenceNumber()))
	if message == nil {
		c.Logger.Info(
			"request not found",
			zap.Int64("sequence number", msg.GetSequenceNumber()),
		)
		return
	}

	// get digest message
	digest := hashing.MD5(message)

	// validate the message
	if !c.validatePreparedMsg(digest, msg) {
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

	// store the log place
	seqn := c.Logs.InitRequest()

	// create our channel for input messages
	c.channels[seqn] = make(chan *models.InterruptMsg)

	// need to create a go-routine to not block the request
	go func(sequence int, message *pbft.RequestMsg) {
		defer func() {
			// reset our channel
			delete(c.channels, sequence)
		}()

		// update request metadata
		message.SequenceNumber = int64(sequence)
		message.Status = pbft.RequestStatus_REQUEST_STATUS_UNSPECIFIED

		// store it into datastore
		c.Logs.SetRequest(sequence, message)

		c.Logger.Debug(
			"new request received",
			zap.Int64("sequence number", message.GetSequenceNumber()),
			zap.Int64("timestamp", message.GetTransaction().GetTimestamp()),
		)

		// broadcast to all using preprepare
		go c.Client.BroadcastPrePrepare(&pbft.PrePrepareMsg{
			Request:        message,
			SequenceNumber: int64(sequence),
			View:           int64(c.Memory.GetView()),
			Digest:         hashing.MD5(message),
		})

		// wait for 2f+1 preprepared messages
		count := c.waitForPrePrepareds(c.channels[sequence])
		c.Logger.Debug("received preprepared messages", zap.Int("messages", count))

		// broadcast to all using prepare
		go c.Client.BroadcastPrepare(&pbft.PrepareMsg{
			SequenceNumber: int64(sequence),
			View:           int64(c.Memory.GetView()),
			Digest:         hashing.MD5(message),
		})

		// wait for 2f+1 prepared messages
		count = c.waitForPrepareds(c.channels[sequence])
		c.Logger.Debug("received prepared messages", zap.Int("messages", count))

		// broadcast to all using commit
		go c.Client.BroadcastCommit(&pbft.CommitMsg{
			SequenceNumber: int64(sequence),
			View:           int64(c.Memory.GetView()),
			Digest:         hashing.MD5(message),
		})
	}(seqn, msg)
}

// handle transaction accepts a new transaction and calls
// the request RPC on the current leader.
func (c *Consensus) handleTransaction(pkt interface{}) {
	defer func() {
		// reset the channel when transaction is done
		c.inTransactionChannel = nil
	}()

	// parse the input message
	msg := pkt.(*pbft.TransactionMsg)

	// get the current leader
	id := c.getCurrentLeader()

	c.Logger.Debug("current leader", zap.String("id", id))

	// send the transaction request to the leader
	c.Client.Request(id, &pbft.RequestMsg{
		Transaction: msg,
		Response:    &pbft.TransactionRsp{},
	})

	c.Logger.Info(
		"request is sent",
		zap.Int64("timestamp", msg.GetTimestamp()),
		zap.String("sender", msg.GetSender()),
		zap.String("receiver", msg.GetReciever()),
		zap.Int64("amount", msg.GetAmount()),
	)

	// wait for f+1 matching reply or timeout request
	resp := c.waitReplys(c.inTransactionChannel)

	c.Logger.Info(
		"received reply message",
		zap.Int64("timestamp", resp.GetTimestamp()),
		zap.Int64("sequence number", resp.GetSequenceNumber()),
	)

	// return the response in channel
	c.outTransactionChannel <- resp.GetResponse()
}
