package consensus

import (
	"github.com/f24-cse535/pbft/pkg/enum"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"
)

// getCurrentLeader returns the current leader id.
func (c *Consensus) getCurrentLeader() string {
	return c.memory.GetNodeByIndex(c.memory.GetView())
}

// checkRequestExecution checks the timestamps to see if a request is executed or not.
func (c *Consensus) checkRequestExecution(ts int64) *pbft.RequestMsg {
	for _, key := range c.logs.GetAllRequests() {
		if key.GetTransaction().GetTimestamp() == ts && key.GetStatus() == pbft.RequestStatus_REQUEST_STATUS_E {
			return key
		}
	}

	return nil
}

// canExecuteRequest gets a sequence number and checks if all requests before that are executed or not.
func (c *Consensus) canExecuteRequest(sequence int) bool {
	// loop from the first sequence and check the execution status
	for i := 0; i < sequence; i++ {
		if tmp := c.logs.GetRequest(i); tmp != nil && tmp.GetStatus() != pbft.RequestStatus_REQUEST_STATUS_E {
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
