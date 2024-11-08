package client

import "github.com/f24-cse535/pbft/pkg/rpc/pbft"

// BroadcastPrePrepare sends a preprepare message to all nodes.
func (c *Client) BroadcastPrePrepare(msg *pbft.PrePrepareMsg) {
	for key := range c.nodes {
		c.PrePrepare(key, msg)
	}
}

// BroadcastPrepare sends a prepare message to all nodes.
func (c *Client) BroadcastPrepare(msg *pbft.AckMsg) {
	for key := range c.nodes {
		c.Prepare(key, msg)
	}
}

// BroadcastCommit sends a commit message to all nodes.
func (c *Client) BroadcastCommit(msg *pbft.AckMsg) {
	for key := range c.nodes {
		c.Commit(key, msg)
	}
}

// BroadcastViewChange sends a view change message to all nodes.
func (c *Client) BroadcastViewChange(msg *pbft.ViewChangeMsg) {
	for key := range c.nodes {
		c.ViewChange(key, msg)
	}
}
