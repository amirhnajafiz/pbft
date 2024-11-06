package consensus

import "github.com/f24-cse535/pbft/pkg/rpc/pbft"

// validatePrePrepareMsg checks the view and digest of a preprepare message.
func (c *Consensus) validatePrePrepareMsg(msg *pbft.PrePrepareMsg, digest string) bool {
	if msg.GetView() != int64(c.memory.GetView()) { // not the same view
		return false
	}

	if msg.GetDigest() != digest { // not the same digest
		return false
	}

	return true
}

// validateAckMessage checks the view and digest of a ack message.
func (c *Consensus) validateAckMessage(msg *pbft.AckMsg, digest string) bool {
	if msg.GetView() != int64(c.memory.GetView()) { // not the same view
		return false
	}

	if msg.GetDigest() != digest { // not the same digest
		return false
	}

	return true
}
