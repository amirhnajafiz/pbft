package client

import (
	"context"

	"github.com/f24-cse535/pbft/pkg/rpc/liveness"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"
	"google.golang.org/protobuf/types/known/emptypb"

	"go.uber.org/zap"
)

// Commit calls the Commit RPC on the target machine (nodes to nodes).
func (c *Client) Commit(target string, msg *pbft.CommitMsg) {
	address := c.nodes[target]
	msg.NodeId = c.nodeId

	// base connection
	conn, err := c.connect(address)
	if err != nil {
		c.logger.Debug("failed to connect", zap.String("address", address), zap.Error(err))
		return
	}
	defer conn.Close()

	// call commit RPC
	if _, err := pbft.NewPBFTClient(conn).Commit(context.Background(), msg); err != nil {
		c.logger.Debug("failed to call commit RPC", zap.String("address", address), zap.Error(err))
	}
}

// PrePrepare calls the PrePrepare RPC on the target machine (nodes to nodes).
func (c *Client) PrePrepare(target string, msg *pbft.PrePrepareMsg) {
	address := c.nodes[target]
	msg.NodeId = c.nodeId

	// base connection
	conn, err := c.connect(address)
	if err != nil {
		c.logger.Debug("failed to connect", zap.String("address", address), zap.Error(err))
		return
	}
	defer conn.Close()

	// call preprepare RPC
	if _, err := pbft.NewPBFTClient(conn).PrePrepare(context.Background(), msg); err != nil {
		c.logger.Debug("failed to call preprepare RPC", zap.String("address", address), zap.Error(err))
	}
}

// Prepare calls the Prepare RPC on the target machine (nodes to nodes).
func (c *Client) Prepare(target string, msg *pbft.PrepareMsg) {
	address := c.nodes[target]
	msg.NodeId = c.nodeId

	// base connection
	conn, err := c.connect(address)
	if err != nil {
		c.logger.Debug("failed to connect", zap.String("address", address), zap.Error(err))
		return
	}
	defer conn.Close()

	// call prepare RPC
	if _, err := pbft.NewPBFTClient(conn).Prepare(context.Background(), msg); err != nil {
		c.logger.Debug("failed to call prepare RPC", zap.String("address", address), zap.Error(err))
	}
}

// Reply calls the Reply RPC on the target machine (nodes to clients).
func (c *Client) Reply(target string, msg *pbft.ReplyMsg) {
	address := c.clients[target]
	msg.NodeId = c.nodeId

	// base connection
	conn, err := c.connect(address)
	if err != nil {
		c.logger.Debug("failed to connect", zap.String("address", address), zap.Error(err))
		return
	}
	defer conn.Close()

	// call reply RPC
	if _, err := pbft.NewPBFTClient(conn).Reply(context.Background(), msg); err != nil {
		c.logger.Debug("failed to call reply RPC", zap.String("address", address), zap.Error(err))
	}
}

// Request calls the Request RPC on the target machine (clients to nodes).
func (c *Client) Request(target string, msg *pbft.RequestMsg) {
	address := c.nodes[target]
	msg.NodeId = c.nodeId

	// base connection
	conn, err := c.connect(address)
	if err != nil {
		c.logger.Debug("failed to connect", zap.String("address", address), zap.Error(err))
		return
	}
	defer conn.Close()

	// call request RPC
	if _, err := pbft.NewPBFTClient(conn).Request(context.Background(), msg); err != nil {
		c.logger.Debug("failed to call request RPC", zap.String("address", address), zap.Error(err))
	}
}

// PrintDB gets a target datastore.
func (c *Client) PrintDB(target string) []*pbft.TransactionMsg {
	address := c.nodes[target]

	// base connection
	conn, err := c.connect(address)
	if err != nil {
		c.logger.Debug("failed to connect", zap.String("address", address), zap.Error(err))

		return nil
	}
	defer conn.Close()

	// call ping RPC
	_, err = liveness.NewLivenessClient(conn).Ping(context.Background(), &emptypb.Empty{})
	if err != nil {
		c.logger.Debug("failed to call ping RPC", zap.String("address", address), zap.Error(err))

		return nil
	}

	// server is ok
	return nil
}

// Ping is used to send a ping request to a server. If the server is available, it returns true.
func (c *Client) PrintLog(target string) bool {
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

// Ping is used to send a ping request to a server. If the server is available, it returns true.
func (c *Client) PrintStatus(target string) bool {
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

// Ping is used to send a ping request to a server. If the server is available, it returns true.
func (c *Client) PrintView(target string) bool {
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

// Ping is used to send a ping request to a server. If the server is available, it returns true.
func (c *Client) Transaction(target string) bool {
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
