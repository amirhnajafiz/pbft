package client

import (
	"context"

	"github.com/f24-cse535/pbft/pkg/rpc/liveness"
	"google.golang.org/protobuf/types/known/emptypb"

	"go.uber.org/zap"
)

// Ping is used to send a ping request to a server. If the server is available, it returns true.
func (c *Client) Ping(address string) bool {
	// base connection
	conn, err := c.connect(address)
	if err != nil {
		c.logger.Debug("failed to connect", zap.String("address", address), zap.Error(err))

		return false
	}
	defer conn.Close()

	// call RPC of ping
	_, err = liveness.NewLivenessClient(conn).Ping(context.Background(), &emptypb.Empty{})
	if err != nil {
		c.logger.Debug("failed to call ping RPC", zap.String("address", address), zap.Error(err))

		return false
	}

	// server is ok
	return true
}

// ChangeState is used to modify the state of a gRPC server.
// if the state is true, then the server is alive, else the server will be blocked.
func (c *Client) ChangeState(address string, state, byzantine bool) {
	// base connection
	conn, err := c.connect(address)
	if err != nil {
		c.logger.Debug("failed to connect", zap.String("address", address), zap.Error(err))
	}
	defer conn.Close()

	// call RPC of change status
	_, err = liveness.NewLivenessClient(conn).ChangeStatus(context.Background(), &liveness.StatusMsg{
		Status:    state,
		Byzantine: byzantine,
	})
	if err != nil {
		c.logger.Debug("failed to call ChangeState RPC", zap.String("address", address), zap.Error(err))

		return
	}
}
