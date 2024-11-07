package client

import (
	"context"

	"github.com/f24-cse535/pbft/pkg/rpc/app"

	"go.uber.org/zap"
)

// Reply calls the Reply RPC on the target machine (nodes to clients).
func (c *Client) Reply(target string, msg *app.ReplyMsg) {
	address := c.nodes[target]
	msg.NodeId = c.nodeId

	// base connection
	conn, err := c.connect(address)
	if err != nil {
		c.logger.Debug("failed to connect", zap.String("address", address), zap.Error(err))
		return
	}
	defer conn.Close()

	// call reply RPC
	if _, err := app.NewAppClient(conn).Reply(context.Background(), msg); err != nil {
		c.logger.Debug("failed to call Reply RPC", zap.String("address", address), zap.Error(err))
	}

	c.logger.Debug("reply sent", zap.String("to", target))
}
