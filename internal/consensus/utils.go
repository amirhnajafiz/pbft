package consensus

import (
	"github.com/f24-cse535/pbft/pkg/enum"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"
)

// getCurrentLeader returns the current leader id.
func (c *Consensus) getCurrentLeader() string {
	return c.memory.GetNodeByIndex(c.memory.GetView())
}

// validateMsg is used to check if the view and digest of a message is valid.
func (c *Consensus) validateMsg(digest, msgDigest string, msgView int64) bool {
	if msgView != int64(c.memory.GetView()) { // different views
		return false
	}

	if msgDigest != digest { // mismatch digest
		return false
	}

	return true
}

// canExecuteRequest gets a sequence number and checks if all requests before that sequence are executed or not.
func (c *Consensus) canExecuteRequest(sequence int) bool {
	for key, value := range c.logs.GetRequestsAfterCheckpoint() {
		if key < sequence && value.GetStatus() != pbft.RequestStatus_REQUEST_STATUS_E {
			return false
		}
	}

	return true
}

// executeTransaction takes a transaction message and performs its operation.
func (c *Consensus) executeTransaction(trx *pbft.TransactionMsg) string {
	senderBalance := c.memory.GetBalance(trx.GetSender())
	receiverBalance := c.memory.GetBalance(trx.GetReciever())
	amount := trx.GetAmount()

	if amount <= int64(senderBalance) {
		c.memory.SetBalance(trx.GetSender(), senderBalance-int(amount))
		c.memory.SetBalance(trx.GetReciever(), receiverBalance+int(amount))

		return enum.RespSuccessTransaction
	}

	return enum.RespNotEnoughBalance
}

// shouldCheckpoint checks the number of executions to see if checkpoint is needed or not.
func (c *Consensus) canCheckpoint() (int, bool, []*pbft.PrePrepareMsg) {
	list := make([]*pbft.PrePrepareMsg, 0)
	index := c.logs.GetLastCheckpoint()
	count := 0

	for {
		if req := c.logs.GetRequest(index); req == nil {
			break
		} else {
			if req.GetStatus() == pbft.RequestStatus_REQUEST_STATUS_E {
				list = append(list, c.logs.GetPreprepare(index))
				count++
				index++
			} else {
				break
			}
		}
	}

	return index, count >= c.cfg.Checkpoint, list
}
