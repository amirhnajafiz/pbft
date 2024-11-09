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
	if msgView != int64(c.memory.GetView()) { // not the same view
		return false
	}

	if msgDigest != digest { // not the same digest
		return false
	}

	return true
}

// checkRequestExecution checks the timestamps to see if a request is executed or not.
func (c *Consensus) checkRequestExecution(ts int64) (*pbft.RequestMsg, bool) {
	for _, value := range c.logs.GetAllRequests() {
		if value.GetTransaction().GetTimestamp() == ts {
			return value, value.GetStatus() == pbft.RequestStatus_REQUEST_STATUS_E
		}
	}

	return nil, false
}

// canExecuteRequest gets a sequence number and checks if all requests before that are executed or not.
func (c *Consensus) canExecuteRequest(sequence int) bool {
	// loop from the first sequence and check the execution status
	for key, value := range c.logs.GetAllRequests() {
		if key < sequence && value.GetStatus() != pbft.RequestStatus_REQUEST_STATUS_E {
			return false
		}
	}

	return true
}

// executeRequest takes a request message and executes its transaction.
func (c *Consensus) executeRequest(msg *pbft.RequestMsg) {
	senderBalance := c.memory.GetBalance(msg.GetTransaction().GetSender())
	receiverBalance := c.memory.GetBalance(msg.GetTransaction().GetReciever())
	amount := msg.GetTransaction().GetAmount()

	msg.GetResponse().Text = enum.RespBadBalance

	// check if the balance is sufficient
	if amount <= int64(senderBalance) {
		c.memory.SetBalance(msg.GetTransaction().GetSender(), senderBalance-int(amount))
		c.memory.SetBalance(msg.GetTransaction().GetReciever(), receiverBalance+int(amount))

		msg.GetResponse().Text = enum.RespSuccess
	}
}

// shouldCheckpoint checks the number of executions to see if checkpoint is needed or not.
func (c *Consensus) shouldCheckpoint() int {
	requests := c.logs.GetPartialRequests(c.memory.GetLowWaterMark())

	count := 0
	target := 0

	for key, value := range requests {
		if value.GetStatus() == pbft.RequestStatus_REQUEST_STATUS_E {
			count++
			target = key
		}
	}

	if count >= c.cfg.Checkpoint {
		return target
	}

	return 0
}
