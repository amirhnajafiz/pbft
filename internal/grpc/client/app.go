package client

import (
	"context"

	"github.com/f24-cse535/pbft/pkg/rpc/app"
)

// Reply calls the Reply RPC on the target machine (nodes to clients).
func (c *Client) Reply(target string, msg *app.ReplyMsg) error {
	address := c.nodes[target]
	msg.NodeId = c.nodeId

	// base connection
	conn, err := c.connect(address)
	if err != nil {
		return err
	}
	defer conn.Close()

	// call reply RPC
	if _, err := app.NewAppClient(conn).Reply(context.Background(), msg); err != nil {
		return err
	}

	return nil
}
