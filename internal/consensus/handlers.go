package consensus

import (
	"github.com/f24-cse535/pbft/internal/utils/hashing"
	"github.com/f24-cse535/pbft/pkg/enum"
	"github.com/f24-cse535/pbft/pkg/models"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"

	"go.uber.org/zap"
)

func (c *Consensus) handleCommit(pkt interface{}) {
	// check the message
	// update log
	// do execution
}

// handle preprepare accepts a preprepare message and validates it to call preprepared RPC.
func (c *Consensus) handlePrePrepare(pkt interface{}) {
	// parse the input message
	msg := pkt.(*pbft.PrePrepareMsg)

	// validate the message
	if !c.validatePrePrepareMsg(msg) {
		c.Logger.Info("preprepared message is not valid")
		return
	}

	// update the log and set the status of preprepared
	msg.GetRequest().Status = pbft.RequestStatus_REQUEST_STATUS_PP
	c.Logs.SetLog(int(msg.GetSequenceNumber()), msg.GetRequest())

	c.Logger.Info(
		"preprepared a message",
		zap.Int("sequence number", int(msg.GetSequenceNumber())),
		zap.Int("timestamp", int(msg.GetRequest().GetTransaction().GetTimestamp())),
	)

	// call preprepared message
	c.Client.PrePrepared(msg.GetNodeId(), &pbft.PrePreparedMsg{
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
			zap.Int("sequence number", int(msg.GetSequenceNumber())),
		)
		return
	}

	// set the digest message
	msg.Digest = hashing.MD5(message.Request)

	// validate the message
	if !c.validatePrePreparedMsg(msg) {
		c.Logger.Info(
			"preprepared message is not valid",
			zap.Int("sequence number", int(msg.GetSequenceNumber())),
		)
		return
	}

	c.Logger.Debug(
		"preprepared received",
		zap.Int("sequence number", int(msg.GetSequenceNumber())),
		zap.Int("timestamp", int(msg.GetRequest().GetTransaction().GetTimestamp())),
	)

	// publish over correct request handler
	c.channels[int(msg.GetSequenceNumber())] <- &models.InterruptMsg{
		Type:    enum.IntrPrePrepared,
		Payload: msg,
	}
}

func (c *Consensus) handlePrepare(pkt interface{}) {
	// check the message
	// update log
	// return with accept message
}

func (c *Consensus) handlePrepared(pkt interface{}) {

}

// handleReply gets a reply message and passes it to the right request handler.
func (c *Consensus) handleReply(pkt interface{}) {
	// parse the input message
	msg := pkt.(*pbft.ReplyMsg)

	// ignore the messages that are not for this node
	if msg.GetClientId() != c.Memory.GetNodeId() {
		return
	}

	// publish over correct request handler
	c.channels[int(msg.GetSequenceNumber())] <- &models.InterruptMsg{
		Type:    enum.IntrPrePrepared,
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

		// store it into logs
		c.Logs.SetLog(sequence, msg)

		c.Logger.Info(
			"new request received",
			zap.Int("sequence number", int(message.GetSequenceNumber())),
			zap.Int("timestamp", int(message.GetTransaction().GetTimestamp())),
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
		// wait for 2f+1
		// broadcast to all using commit
		// execute message if possible

		c.Logger.Info("reply sent",
			zap.String("client", message.GetClientId()),
			zap.Int("sequence number", int(message.GetSequenceNumber())),
			zap.Int("timestamp", int(message.GetTransaction().GetTimestamp())),
		)

		// send the reply
		c.Client.Reply(message.GetClientId(), &pbft.ReplyMsg{
			SequenceNumber: message.GetSequenceNumber(),
			View:           int64(c.Memory.GetView()),
			Timestamp:      message.GetTransaction().GetTimestamp(),
			ClientId:       message.GetClientId(),
			Response:       "received",
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
	})

	c.Logger.Info(
		"request is sent",
		zap.Int("timestamp", int(msg.GetTimestamp())),
		zap.String("sender", msg.GetSender()),
		zap.String("receiver", msg.GetReciever()),
		zap.Int64("amount", msg.GetAmount()),
	)

	// wait for f+1 matching reply or timeout request
	resp := c.waitReplys(c.inTransactionChannel)

	// return the response in channel
	c.outTransactionChannel <- resp.GetResponse()
}
