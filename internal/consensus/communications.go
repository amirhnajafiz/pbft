package consensus

import (
	"github.com/f24-cse535/pbft/internal/utils/hashing"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"
)

// sendReplyMsg gets a request message and uses client.go to send a reply message.
func (c *Consensus) sendReplyMsg(msg *pbft.RequestMsg) {
	c.Client.Reply(
		msg.GetClientId(),
		&pbft.ReplyMsg{
			SequenceNumber: msg.GetSequenceNumber(),
			View:           int64(c.Memory.GetView()),
			Timestamp:      msg.GetTransaction().GetTimestamp(),
			ClientId:       msg.GetClientId(),
			Response:       msg.GetResponse().GetText(),
		},
	)
}

// sendPreprepareMsg gets a request message and uses client.go to broadcast a preprepare message.
func (c *Consensus) sendPreprepareMsg(msg *pbft.RequestMsg) {
	c.Client.BroadcastPrePrepare(&pbft.PrePrepareMsg{
		Request:        msg,
		SequenceNumber: msg.GetSequenceNumber(),
		View:           int64(c.Memory.GetView()),
		Digest:         hashing.MD5(msg),
	})
}

// sendPrepareMsg gets a request message and uses client.go to broadcast a prepare message.
func (c *Consensus) sendPrepareMsg(msg *pbft.RequestMsg) {
	c.Client.BroadcastPrepare(&pbft.AckMsg{
		SequenceNumber: msg.GetSequenceNumber(),
		View:           int64(c.Memory.GetView()),
		Digest:         hashing.MD5(msg),
	})
}

// sendCommitMsg gets a request message and uses client.go to broadcast a commit message.
func (c *Consensus) sendCommitMsg(msg *pbft.RequestMsg) {
	c.Client.BroadcastCommit(&pbft.AckMsg{
		SequenceNumber: msg.GetSequenceNumber(),
		View:           int64(c.Memory.GetView()),
		Digest:         hashing.MD5(msg),
	})
}
