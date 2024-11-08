package application

import (
	"crypto/tls"
	"fmt"
	"net"

	"github.com/f24-cse535/pbft/internal/config/node/bft"
	"github.com/f24-cse535/pbft/internal/grpc/client"
	"github.com/f24-cse535/pbft/internal/storage/local"
	"github.com/f24-cse535/pbft/pkg/models"
	"github.com/f24-cse535/pbft/pkg/rpc/app"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// App is the main module of the client's program.
type App struct {
	memory *local.Memory
	cli    *client.Client
	cfg    *bft.Config
	logger *zap.Logger

	clients  map[string]chan *models.Transaction // for each client there is go-routine that accepts using these channels
	handlers map[string]chan *app.ReplyMsg       // handlers is a channel to get gRPC messages
}

// NewApp returns a new app instance.
func NewApp(
	logr *zap.Logger,
	mem *local.Memory,
	cfg *bft.Config,
	cli *client.Client,
	clients map[string]int,
) *App {
	// create a new app instance
	a := &App{
		logger: logr,
		memory: mem,
		cfg:    cfg,
		cli:    cli,
	}

	// initial channels
	a.clients = make(map[string]chan *models.Transaction)
	a.handlers = make(map[string]chan *app.ReplyMsg)

	// for each client, run a transaction handler
	for key := range clients {
		a.clients[key] = make(chan *models.Transaction)
		go a.transactionHandler(key)
	}

	return a
}

// Client returns the app client for direct calls.
func (a *App) Client() *client.Client {
	return a.cli
}

// Transaction sends a new transaction to the transaction handler.
func (a *App) Transaction(trx *models.Transaction) {
	a.clients[trx.Sender] <- trx
}

// service starts the gRPC server.
func (a *App) Service(port int, tlsConfig *tls.Config) error {
	// on the local network, listen to a port
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to start the listener server: %v", err)
	}

	// create a new grpc instance
	server := grpc.NewServer(
		grpc.Creds(credentials.NewTLS(tlsConfig)), // set the server certificates
	)

	// register all gRPC services
	app.RegisterAppServer(server, &service{
		channels: a.handlers,
	})

	// starting the server
	if err := server.Serve(listener); err != nil {
		return fmt.Errorf("failed to start services: %v", err)
	}

	return nil
}
