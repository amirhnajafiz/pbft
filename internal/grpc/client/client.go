package client

import (
	"crypto/tls"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Client has all RPCs to communicate with the gRPC servers.
type Client struct {
	nodeId string
	nodes  map[string]string
	creds  *tls.Config
}

// connect should be called in the beginning of each method to establish a connection.
func (c *Client) connect(address string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(credentials.NewTLS(c.creds)))
	if err != nil {
		return nil, fmt.Errorf("connection failed: %v", err)
	}

	return conn, nil
}

// NewClient returns a new RPC client to make RPC to the gRPC server.
func NewClient(creds *tls.Config, nodeId string, nodes map[string]string) *Client {
	delete(nodes, nodeId)

	return &Client{
		creds:  creds,
		nodeId: nodeId,
		nodes:  nodes,
	}
}
