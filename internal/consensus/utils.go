package consensus

import (
	"github.com/f24-cse535/pbft/internal/utils/hashing"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"
)

// getCurrentLeader returns the current leader id.
func (c *Consensus) getCurrentLeader() string {
	return c.Memory.GetNodeByIndex(c.Memory.GetView())
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

	// get the message from logs
	message := c.Logs.GetLog(int(msg.GetSequenceNumber()))
	if message == nil {
		return false
	}

	digest := hashing.MD5(message.Request)
	if digest != hashing.MD5(msg.GetRequest()) { // not the same digest
		return false
	}

	msg.Digest = digest

	return true
}
