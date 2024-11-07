package client

import (
	"context"

	"github.com/f24-cse535/pbft/pkg/rpc/liveness"

	"google.golang.org/protobuf/types/known/emptypb"
)

// Ping is used to send a ping request to a server. If the server is available, it returns true.
func (c *Client) Ping(target string) bool {
	address := c.nodes[target]

	// base connection
	conn, err := c.connect(address)
	if err != nil {
		return false
	}
	defer conn.Close()

	// call ping RPC
	_, err = liveness.NewLivenessClient(conn).Ping(context.Background(), &emptypb.Empty{})

	return err == nil
}

// ChangeState is used to modify the state of a gRPC server.
func (c *Client) ChangeState(target string, state, byzantine bool) error {
	address := c.nodes[target]

	// base connection
	conn, err := c.connect(address)
	if err != nil {
		return err
	}
	defer conn.Close()

	// call change status RPC
	_, err = liveness.NewLivenessClient(conn).ChangeStatus(context.Background(), &liveness.StatusMsg{
		Status:    state,
		Byzantine: byzantine,
	})
	if err != nil {
		return err
	}

	return nil
}

// Flush calls the Flush RPC on a target node.
func (c *Client) Flush(target string) error {
	address := c.nodes[target]

	// base connection
	conn, err := c.connect(address)
	if err != nil {
		return err
	}
	defer conn.Close()

	// call flush RPC
	_, err = liveness.NewLivenessClient(conn).Flush(context.Background(), &emptypb.Empty{})
	if err != nil {
		return err
	}

	return err
}
