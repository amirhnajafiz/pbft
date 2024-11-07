package application

import (
	"crypto/tls"
	"fmt"
	"net"

	"github.com/f24-cse535/pbft/internal/grpc/client"
	"github.com/f24-cse535/pbft/pkg/models"
	"github.com/f24-cse535/pbft/pkg/rpc/app"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// App is the main module of the client's program.
type App struct {
	Cli *client.Client

	clients  map[string]chan *models.Transaction
	handlers map[string]chan *app.ReplyMsg
}

// Start a go-routine to get and dispatch reply messages to handlers.
func (a *App) Start(clients map[string]int) {
	a.clients = make(map[string]chan *models.Transaction)
	a.handlers = make(map[string]chan *app.ReplyMsg)

	for key := range clients {
		a.clients[key] = make(chan *models.Transaction)
		a.handlers[key] = make(chan *app.ReplyMsg)

		go a.transactionHandler(key)
	}
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
