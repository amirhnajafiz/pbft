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
		c.Logger.Debug(
			"request cannot get executed yet",
			zap.Int("sequence number", sequence),
		)
		return
	}

	// follow sequence until one is not committed, execute them
	index := sequence
	for {
		if tmp := c.Logs.GetLog(index); tmp != nil && tmp.Request.Status == pbft.RequestStatus_REQUEST_STATUS_C {
			// execute request
			c.executeRequest(tmp.Request)

			// update the log and set the status of prepare
			tmp.Request.Status = pbft.RequestStatus_REQUEST_STATUS_E
			c.Logs.SetLog(sequence, tmp.Request)

			// send the reply message
			c.Client.Reply(tmp.Request.GetClientId(), &pbft.ReplyMsg{
				SequenceNumber: tmp.Request.GetSequenceNumber(),
				View:           int64(c.Memory.GetView()),
				Timestamp:      tmp.Request.GetTransaction().GetTimestamp(),
				ClientId:       tmp.Request.GetClientId(),
				Response:       tmp.Request.GetResponse().GetText(),
			})

			c.Logger.Info(
				"message executed and reply sent",
				zap.String("client", tmp.Request.GetClientId()),
				zap.Int64("sequence number", tmp.Request.GetSequenceNumber()),
				zap.Int64("timestamp", tmp.Request.GetTransaction().GetTimestamp()),
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

	// get the message from our logs
	message := c.Logs.GetLog(int(msg.GetSequenceNumber()))
	if message == nil {
		c.Logger.Info(
			"message not found",
			zap.Int64("sequence number", msg.GetSequenceNumber()),
		)
		return
	}

	// update the log and set the status of prepare
	message.Request.Status = pbft.RequestStatus_REQUEST_STATUS_C
	c.Logs.SetLog(int(msg.GetSequenceNumber()), message.Request)

	// message committed
	c.Logger.Info(
		"message committed",
		zap.Int64("sequence number", message.Request.GetSequenceNumber()),
		zap.Int64("timestamp", message.Request.GetTransaction().GetTimestamp()),
	)

	// call execute handler
	c.handleExecute(int(msg.GetSequenceNumber()))
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

	// update the log and set the status of preprepared
	msg.GetRequest().Status = pbft.RequestStatus_REQUEST_STATUS_PP
	c.Logs.SetLog(int(msg.GetSequenceNumber()), msg.GetRequest())

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

	// get the message from our logs
	message := c.Logs.GetLog(int(msg.GetSequenceNumber()))
	if message == nil {
		c.Logger.Info(
			"message not found",
			zap.Int64("sequence number", msg.GetSequenceNumber()),
			zap.Int64("timestamp", msg.GetRequest().GetTransaction().GetTimestamp()),
		)
		return
	}

	// set the digest message
	msg.Digest = hashing.MD5(message.Request)

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
		zap.Int64("sequence number", message.Request.GetSequenceNumber()),
		zap.Int64("timestamp", message.Request.GetTransaction().GetTimestamp()),
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

	// get the message from our logs
	message := c.Logs.GetLog(int(msg.GetSequenceNumber()))
	if message == nil {
		c.Logger.Info(
			"message not found",
			zap.Int64("sequence number", msg.GetSequenceNumber()),
		)
		return
	}

	// get digest message
	digest := hashing.MD5(message.Request)

	// validate the message
	if !c.validatePrepareMsg(digest, msg) {
		c.Logger.Info(
			"prepare message is not valid",
			zap.Int64("sequence number", msg.GetSequenceNumber()),
		)
		return
	}

	// update the log and set the status of prepare
	message.Request.Status = pbft.RequestStatus_REQUEST_STATUS_P
	c.Logs.SetLog(int(msg.GetSequenceNumber()), message.Request)

	c.Logger.Info(
		"prepared a message",
		zap.Int64("sequence number", message.Request.GetSequenceNumber()),
		zap.Int64("timestamp", message.Request.GetTransaction().GetTimestamp()),
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

	// get the message from our logs
	message := c.Logs.GetLog(int(msg.GetSequenceNumber()))
	if message == nil {
		c.Logger.Info(
			"message not found",
			zap.Int64("sequence number", msg.GetSequenceNumber()),
		)
		return
	}

	// get digest message
	digest := hashing.MD5(message.Request)

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
		zap.Int64("sequence number", message.Request.GetSequenceNumber()),
		zap.Int64("timestamp", message.Request.GetTransaction().GetTimestamp()),
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
	seqn := c.Logs.InitLog()

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

		// store it into logs
		c.Logs.SetLog(sequence, message)

		c.Logger.Info(
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
		c.Logger.Info("received preprepared messages", zap.Int("messages", count))

		// broadcast to all using prepare
		go c.Client.BroadcastPrepare(&pbft.PrepareMsg{
			SequenceNumber: int64(sequence),
			View:           int64(c.Memory.GetView()),
			Digest:         hashing.MD5(message),
		})

		// wait for 2f+1 prepared messages
		count = c.waitForPrepareds(c.channels[sequence])
		c.Logger.Info("received prepared messages", zap.Int("messages", count))

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
