package consensus

import (
	"time"

	"github.com/f24-cse535/pbft/pkg/rpc/pbft"

	"go.uber.org/zap"
)

// promiseRequest in a loop trys to send your request to a leader.
func (c *Consensus) promiseRequest(msg *pbft.TransactionMsg) {
	for {
		// get the current leader id
		id := c.getCurrentLeader()

		c.Logger.Debug("current leader", zap.String("id", id))

		// send the transaction request to the leader
		if err := c.Client.Request(id, &pbft.RequestMsg{
			Transaction: msg,
			Response:    &pbft.TransactionRsp{},
		}); err != nil {
			// if the leader failed, increament the view, and wait 1 second before resending your request
			c.Logger.Debug("leader failed", zap.Error(err))
			c.Memory.IncView()

			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}
}

// promiseReceive in a loop trys to get a response for user request.
func (c *Consensus) promiseReceive(msg *pbft.TransactionMsg) *pbft.ReplyMsg {
	for {
		// wait for f+1 matching reply or timeout request
		resp := c.waitReplys(c.inTransactionChannel)
		if resp != nil {
			return resp
		}

		c.Logger.Debug("request timeout")

		// send message to all nodes
		c.Client.BroadcastRequest(msg)

		// sleep one second before resending the request
		time.Sleep(1 * time.Second)
	}
}

// promiseClear makes sure that nothing is inside a dead channel.
func (c *Consensus) promiseClear(sequence int) {
	channel := c.channels[sequence]
	timer := time.NewTimer(60 * time.Second)

	for {
		select {
		case <-channel:
			continue
		case <-timer.C:
			delete(c.channels, sequence)
			return
		}
	}
}
