package client

import (
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"
	"go.uber.org/zap"
)

// BroadcastRequest sends a request message to all nodes.
func (c *Client) BroadcastRequest(msg *pbft.TransactionMsg) {
	for key := range c.nodes {
		c.Request(key, &pbft.RequestMsg{
			Transaction: msg,
			Response:    &pbft.TransactionRsp{},
		})
	}
}

// BroadcastPrePrepare sends a preprepare message to all nodes.
func (c *Client) BroadcastPrePrepare(msg *pbft.PrePrepareMsg) {
	for key := range c.nodes {
		c.logger.Debug("preprepare is sent", zap.String("to", key))
		c.PrePrepare(key, msg)
	}
}

// BroadcastPrepare sends a prepare message to all nodes.
func (c *Client) BroadcastPrepare(msg *pbft.AckMsg) {
	for key := range c.nodes {
		c.logger.Debug("prepare is sent", zap.String("to", key))
		c.Prepare(key, msg)
	}
}

// BroadcastCommit sends a commit message to all nodes.
func (c *Client) BroadcastCommit(msg *pbft.AckMsg) {
	for key := range c.nodes {
		c.logger.Debug("commit is sent", zap.String("to", key))
		c.Commit(key, msg)
	}
}
