package consensus

import (
	"github.com/f24-cse535/pbft/pkg/enum"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"
)

// getCurrentLeader returns the current leader id.
func (c *Consensus) getCurrentLeader() string {
	return c.Memory.GetNodeByIndex(c.Memory.GetView())
}

// isRequestExecuted checks the timestamps to see if a request is executed or not.
func (c *Consensus) isRequestExecuted(ts int64) *pbft.RequestMsg {
	for _, key := range c.Logs.GetAllRequests() {
		if key.GetTransaction().GetTimestamp() == ts && key.GetStatus() == pbft.RequestStatus_REQUEST_STATUS_E {
			return key
		}
	}

	return nil
}

// canExecute gets a sequence number and checks if all requests before that are executed or not.
func (c *Consensus) canExecute(sequence int) bool {
	// loop from the first sequence and check the execution status
	for i := 0; i < sequence; i++ {
		if tmp := c.Logs.GetRequest(i); tmp != nil && tmp.GetStatus() != pbft.RequestStatus_REQUEST_STATUS_E {
			return false
		}
	}

	return true
}

// executeRequest takes a request message and executes its transaction.
func (c *Consensus) executeRequest(msg *pbft.RequestMsg) {
	// fetch all balance values
	senderBalance := c.Memory.GetBalance(msg.GetTransaction().GetSender())
	receiverBalance := c.Memory.GetBalance(msg.GetTransaction().GetReciever())
	amount := msg.GetTransaction().GetAmount()

	msg.GetResponse().Text = enum.RespBadBalance

	// check if the balance is sufficient
	if amount <= int64(senderBalance) {
		c.Memory.SetBalance(msg.GetTransaction().GetSender(), senderBalance-int(amount))
		c.Memory.SetBalance(msg.GetTransaction().GetReciever(), receiverBalance+int(amount))

		msg.GetResponse().Text = enum.RespSuccess
	}
}

// validatePrePrepareMsg checks the view and digest of a preprepare message.
func (c *Consensus) validatePrePrepareMsg(msg *pbft.PrePrepareMsg, digest string) bool {
	if msg.GetView() != int64(c.Memory.GetView()) { // not the same view
		return false
	}

	if msg.GetDigest() != digest { // not the same digest
		return false
	}

	return true
}

// validateAckMessage checks the view and digest of a ack message.
func (c *Consensus) validateAckMessage(msg *pbft.AckMsg, digest string) bool {
	if msg.GetView() != int64(c.Memory.GetView()) { // not the same view
		return false
	}

	if msg.GetDigest() != digest { // not the same digest
		return false
	}

	return true
}
