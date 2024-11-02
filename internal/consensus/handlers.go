package consensus

import (
	"github.com/f24-cse535/pbft/pkg/models"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"
	"go.uber.org/zap"
)

func (c *Consensus) handleCommit(pkt interface{}) {
	// check the message
	// update log
	// do execution
}

func (c *Consensus) handlePrePrepare(pkt interface{}) {
	// check the message
	// update log
	// return with preprepared message
}

func (c *Consensus) handlePrePrepared(pkt interface{}) {

}

func (c *Consensus) handlePrepare(pkt interface{}) {
	// check the message
	// update log
	// return with accept message
}

func (c *Consensus) handlePrepared(pkt interface{}) {

}

func (c *Consensus) handleReply(pkt interface{}) {
	// update the memory
	// notify the transaction handler
}

func (c *Consensus) handleRequest(pkt interface{}) {
	// parse the input message
	msg := pkt.(*pbft.RequestMsg)

	// store the log
	seqn := c.Logs.InitLog(msg)

	// create our channel for input messages
	c.channels[seqn] = make(chan *models.InterruptMsg)

	// need to create a go-routine to not block the request
	go func(sequence int, message *pbft.RequestMsg) {
		defer func() {
			// reset our channel
			delete(c.channels, sequence)
		}()

		c.Logger.Info("new request", zap.Int("sequence number", sequence))

		// update the request meta-data
		// broadcast to all using preprepare
		// wait for 2f+1
		// broadcast to all using prepare
		// wait for 2f+1
		// broadcast to all using commit
		// execute message if possible
		// send the reply
	}(seqn, msg)
}

func (c *Consensus) handleTransaction(pkt interface{}) {
	defer func() {
		// reset the channel when transaction is done
		c.inTransactionChannel = nil
	}()

	// parse the input message
	msg := pkt.(*pbft.TransactionMsg)

	// get the current leader
	id := c.GetCurrentLeader()

	c.Logger.Debug("current leader", zap.String("id", id))

	// send the transaction request to the leader
	c.Client.Request(id, &pbft.RequestMsg{
		Transaction: msg,
	})

	c.Logger.Info(
		"request is send",
		zap.Int("timestamp", int(msg.GetTimestamp())),
		zap.String("sender", msg.GetSender()),
		zap.String("receiver", msg.GetReciever()),
		zap.Int64("amount", msg.GetAmount()),
	)

	// wait for f+1 matching reply or timeout request (+ timer)
	// on the timeout, reset yourself
	// on the f+1 reply, send over channel

	c.outTransactionChannel <- "received"
}
