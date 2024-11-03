package consensus

import (
	"time"

	"github.com/f24-cse535/pbft/pkg/rpc/pbft"

	"go.uber.org/zap"
)

// helpSendRequest in a loop trys to send your request to a leader.
func (c *Consensus) helpSendRequest(msg *pbft.TransactionMsg) {
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
