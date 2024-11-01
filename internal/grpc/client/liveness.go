package client

import (
	"context"

	"github.com/f24-cse535/pbft/pkg/rpc/liveness"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Ping is used to send a ping request to a server. If the server is available, it returns true.
func (c *Client) Ping(target string) bool {
	address := c.nodes[target]

	// base connection
	conn, err := c.connect(address)
	if err != nil {
		c.logger.Debug("failed to connect", zap.String("address", address), zap.Error(err))

		return false
	}
	defer conn.Close()

	// call ping RPC
	_, err = liveness.NewLivenessClient(conn).Ping(context.Background(), &emptypb.Empty{})
	if err != nil {
		c.logger.Debug("failed to call ping RPC", zap.String("address", address), zap.Error(err))

		return false
	}

	// server is ok
	return true
}

// ChangeState is used to modify the state of a gRPC server.
func (c *Client) ChangeState(target string, state, byzantine bool) {
	address := c.nodes[target]

	// base connection
	conn, err := c.connect(address)
	if err != nil {
		c.logger.Debug("failed to connect", zap.String("address", address), zap.Error(err))
	}
	defer conn.Close()

	// call change status RPC
	_, err = liveness.NewLivenessClient(conn).ChangeStatus(context.Background(), &liveness.StatusMsg{
		Status:    state,
		Byzantine: byzantine,
	})
	if err != nil {
		c.logger.Debug("failed to call ChangeState RPC", zap.String("address", address), zap.Error(err))

		return
	}
}
