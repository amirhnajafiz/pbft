package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/examples/data"
)

// Client has all RPCs to communicate with the gRPC servers.
type Client struct {
	nodeId string

	nodes   map[string]string
	clients map[string]string

	tlsConfig *tls.Config
	logger    *zap.Logger
}

// connect should be called in the beginning of each method to establish a connection.
func (c *Client) connect(address string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(credentials.NewTLS(c.tlsConfig)))
	if err != nil {
		return nil, fmt.Errorf("connection failed: %v", err)
	}

	return conn, nil
}

// LoadTLS get the path of keys and certificates and creates a TLS config for connections.
func (c *Client) LoadTLS(private, public, cas string) error {
	// load the client keys
	cert, err := tls.LoadX509KeyPair(data.Path(public), data.Path(private))
	if err != nil {
		return fmt.Errorf("failed to load certificates: %v", err)
	}

	// ca credentials
	ca := x509.NewCertPool()
	cac := data.Path(cas)
	cacBytes, err := os.ReadFile(cac)
	if err != nil {
		return fmt.Errorf("failed to read ca cert: %v", err)
	}
	if ok := ca.AppendCertsFromPEM(cacBytes); !ok {
		return fmt.Errorf("failed to append certs: %s", cac)
	}

	// set tls configs
	c.tlsConfig = &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
		ClientCAs:    ca,
	}

	return nil
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
