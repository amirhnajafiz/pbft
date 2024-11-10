package modules

import (
	"github.com/f24-cse535/pbft/internal/grpc/client"
	"github.com/f24-cse535/pbft/pkg/rpc/app"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"
)

// Communication is a module that uses client.go to make RPC calls for consensus module.
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

// SendReplyMsg gets a request message and uses client.go to send a reply message.
func (c *Communication) SendReplyMsg(sequence, view int, msg *pbft.RequestMsg) {
	c.cli.Reply(
		msg.GetClientId(),
		&app.ReplyMsg{
			SequenceNumber: int64(sequence),
			View:           int64(view),
			Timestamp:      msg.GetTransaction().GetTimestamp(),
			ClientId:       msg.GetClientId(),
			Response:       msg.GetResponse().GetText(),
			Sender:         msg.GetTransaction().GetSender(),
		},
	)
}

// SendPreprepareMsg gets a preprepare message and uses client.go to broadcast that message.
func (c *Communication) SendPreprepareMsg(msg *pbft.PrePrepareMsg) {
	for key := range c.cli.GetSystemNodes() {
		c.cli.PrePrepare(key, msg)
	}
}

// SendPrepareMsg gets parameters needed to broadcast a prepare message and sends it.
func (c *Communication) SendPrepareMsg(sequence, view int, digest string) {
	msg := pbft.AckMsg{
		SequenceNumber: int64(sequence),
		View:           int64(view),
		Digest:         digest,
	}

	for key := range c.cli.GetSystemNodes() {
		c.cli.Prepare(key, &msg)
	}
}

// SendCommitMsg gets a request message and uses client.go to broadcast a commit message.
func (c *Communication) SendCommitMsg(sequence, view int, digest string, optimized bool) {
	msg := pbft.AckMsg{
		SequenceNumber: int64(sequence),
		View:           int64(view),
		Digest:         digest,
		IsOptimized:    optimized,
	}

	for key := range c.cli.GetSystemNodes() {
		c.cli.Commit(key, &msg)
	}
}

// SendViewChangeMsg gets view and uses client.go to broadcase a view change message.
func (c *Communication) SendViewChangeMsg(msg *pbft.ViewChangeMsg) int {
	count := 1
	for key := range c.cli.GetSystemNodes() {
		if err := c.cli.ViewChange(key, msg); err == nil {
			count++
		}
	}

	return count
}

// SendNewViewMsg broadcasts a new view message to all nodes.
func (c *Communication) SendNewViewMsg(msg *pbft.NewViewMsg) {
	for key := range c.cli.GetSystemNodes() {
		c.cli.NewView(key, msg)
	}
}

// SendCheckpoint broadcasts a checkpoint message to all nodes.
func (c *Communication) SendCheckpoint(msg *pbft.CheckpointMsg) {
	for key := range c.cli.GetSystemNodes() {
		c.cli.Checkpoint(key, msg)
	}
}
