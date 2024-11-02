package client

import (
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client has all RPCs to communicate with the gRPC servers.
type Client struct {
	nodeId string

	nodes   map[string]string
	clients map[string]string

	logger *zap.Logger
}

// connect should be called in the beginning of each method to establish a connection.
func (c *Client) connect(address string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("connection failed: %v", err)
	}

	return conn, nil
}

// NewClient returns a new RPC client to make RPC to the gRPC server.
func NewClient(logr *zap.Logger, nodeId string, nodes map[string]string, clients map[string]string) *Client {
	return &Client{
		nodeId:  nodeId,
		nodes:   nodes,
		clients: clients,
		logger:  logr,
	}
}
