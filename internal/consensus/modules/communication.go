package modules

import (
	"github.com/f24-cse535/pbft/internal/grpc/client"
	"github.com/f24-cse535/pbft/internal/utils/hashing"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"
)

// Communication is a module that uses client.go to send messages.
type Communication struct {
	cli *client.Client
}

// NewCommunicationModule returns a new communication module instance.
func NewCommunicationModule(cli *client.Client) *Communication {
	return &Communication{
		cli: cli,
	}
}

// Client returns the client.go instance for direct calls.
func (c *Communication) Client() *client.Client {
	return c.cli
}

// sendReplyMsg gets a request message and uses client.go to send a reply message.
func (c *Communication) SendReplyMsg(msg *pbft.RequestMsg, view int) {
	c.cli.Reply(
		msg.GetClientId(),
		&pbft.ReplyMsg{
			SequenceNumber: msg.GetSequenceNumber(),
			View:           int64(view),
			Timestamp:      msg.GetTransaction().GetTimestamp(),
			ClientId:       msg.GetClientId(),
			Response:       msg.GetResponse().GetText(),
		},
	)
}

// sendPreprepareMsg gets a request message and uses client.go to broadcast a preprepare message.
func (c *Communication) SendPreprepareMsg(msg *pbft.RequestMsg, view int) {
	c.cli.BroadcastPrePrepare(&pbft.PrePrepareMsg{
		Request:        msg,
		SequenceNumber: msg.GetSequenceNumber(),
		View:           int64(view),
		Digest:         hashing.MD5(msg),
	})
}

// sendPrepareMsg gets a request message and uses client.go to broadcast a prepare message.
func (c *Communication) SendPrepareMsg(msg *pbft.RequestMsg, view int) {
	c.cli.BroadcastPrepare(&pbft.AckMsg{
		SequenceNumber: msg.GetSequenceNumber(),
		View:           int64(view),
		Digest:         hashing.MD5(msg),
	})
}

// sendCommitMsg gets a request message and uses client.go to broadcast a commit message.
func (c *Communication) SendCommitMsg(msg *pbft.RequestMsg, view int) {
	c.cli.BroadcastCommit(&pbft.AckMsg{
		SequenceNumber: msg.GetSequenceNumber(),
		View:           int64(view),
		Digest:         hashing.MD5(msg),
	})
}
