package consensus

import (
	"github.com/f24-cse535/pbft/internal/utils/hashing"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"

	"go.uber.org/zap"
)

// getCurrentLeader returns the current leader id.
func (c *Consensus) getCurrentLeader() string {
	return c.Memory.GetNodeByIndex(c.Memory.GetView())
}

// canExecute gets a sequence number and checks if all logs before that
// are executed or not.
func (c *Consensus) canExecute(sequence int) bool {
	for i := 0; i < sequence; i++ {
		if tmp := c.Logs.GetLog(i); tmp != nil && tmp.Request.Status != pbft.RequestStatus_REQUEST_STATUS_E {
			return false
		}
	}

	return true
}

// executeRequest takes a request message and executes its transaction.
func (c *Consensus) executeRequest(sequence int, msg *pbft.RequestMsg) {
	senderBalance := c.Memory.GetBalance(msg.GetTransaction().GetSender())
	receiverBalance := c.Memory.GetBalance(msg.GetTransaction().GetReciever())
	amount := msg.GetTransaction().GetAmount()

	if amount > int64(senderBalance) {
		msg.GetResponse().Text = "not enough balance"
	} else {
		c.Memory.SetBalance(msg.GetTransaction().GetSender(), senderBalance-int(amount))
		c.Memory.SetBalance(msg.GetTransaction().GetReciever(), receiverBalance+int(amount))
		msg.GetResponse().Text = "submitted"
	}

	// update the log and set the status of prepare
	msg.Status = pbft.RequestStatus_REQUEST_STATUS_E
	c.Logs.SetLog(sequence, msg)

	// send the reply message
	c.Client.Reply(msg.GetClientId(), &pbft.ReplyMsg{
		SequenceNumber: msg.GetSequenceNumber(),
		View:           int64(c.Memory.GetView()),
		Timestamp:      msg.GetTransaction().GetTimestamp(),
		ClientId:       msg.GetClientId(),
		Response:       msg.GetResponse().GetText(),
	})

	c.Logger.Info("reply sent",
		zap.String("client", msg.GetClientId()),
		zap.Int64("sequence number", msg.GetSequenceNumber()),
		zap.Int64("timestamp", msg.GetTransaction().GetTimestamp()),
	)
}

// checkLogsForPossibleExecution executes all available requests that are after this sequence.
func (c *Consensus) checkLogsForPossibleExecution(sequence int) {
	index := sequence + 1
	for {
		if tmp := c.Logs.GetLog(index); tmp != nil && tmp.Request.Status == pbft.RequestStatus_REQUEST_STATUS_C {
			c.executeRequest(index, tmp.Request)
		} else {
			break
		}

		index++
	}
}

// validatePrePrepareMsg checks the view and digest of a preprepare message.
func (c *Consensus) validatePrePrepareMsg(msg *pbft.PrePrepareMsg) bool {
	if msg.GetView() != int64(c.Memory.GetView()) { // not the same view
		return false
	}

	if msg.GetDigest() != hashing.MD5(msg.GetRequest()) { // not the same digest
		return false
	}

	return true
}

// validatePrePreparedMsg checks the view and digest of a preprepared message.
func (c *Consensus) validatePrePreparedMsg(msg *pbft.PrePreparedMsg) bool {
	if msg.GetView() != int64(c.Memory.GetView()) { // not the same view
		return false
	}

	if msg.GetDigest() != hashing.MD5(msg.GetRequest()) { // not the same digest
		return false
	}

	return true
}

// validatePrepareMsg checks the view and digest of a prepare message.
func (c *Consensus) validatePrepareMsg(digest string, msg *pbft.PrepareMsg) bool {
	if msg.GetView() != int64(c.Memory.GetView()) { // not the same view
		return false
	}

	if msg.GetDigest() != digest { // not the same digest
		return false
	}

	return true
}

// validatePreparedMsg checks the view and digest of a prepared message.
func (c *Consensus) validatePreparedMsg(digest string, msg *pbft.PreparedMsg) bool {
	if msg.GetView() != int64(c.Memory.GetView()) { // not the same view
		return false
	}

	if msg.GetDigest() != digest { // not the same digest
		return false
	}

	return true
}
