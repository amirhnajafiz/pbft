package grpc

import (
	"crypto/tls"
	"fmt"
	"net"

	"github.com/f24-cse535/pbft/internal/consensus"
	"github.com/f24-cse535/pbft/internal/grpc/services"
	"github.com/f24-cse535/pbft/internal/storage/local"
	"github.com/f24-cse535/pbft/internal/storage/logs"
	"github.com/f24-cse535/pbft/pkg/rpc/liveness"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Bootstrap is a wrapper that holds requirements for the gRPC services.
type Bootstrap struct {
	Consensus *consensus.Consensus // consensus module is the node main module
	Memory    *local.Memory        // memory is needed for liveness
	Logs      *logs.Logs
	Logger    *zap.Logger // logger is needed for tracing
}

// ListenAnsServer creates a new gRPC instance with all services.
func (b *Bootstrap) ListenAnsServer(port int, creds *tls.Config) error {
	// on the local network, listen to a port
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to start the listener server: %v", err)
	}

	// create a new grpc instance
	server := grpc.NewServer(
		grpc.Creds(credentials.NewTLS(creds)),          // set the server certificates
		grpc.UnaryInterceptor(b.allUnaryInterceptor),   // set an unary interceptor
		grpc.StreamInterceptor(b.allStreamInterceptor), // set a stream interceptor
	)

	// register all gRPC services
	liveness.RegisterLivenessServer(server, &services.Liveness{
		Memory: b.Memory,
		Logs:   b.Logs,
		Logger: b.Logger.Named("liveness"),
	})
	pbft.RegisterPBFTServer(server, &services.PBFT{
		Consensus: b.Consensus,
		Memory:    b.Memory,
		Logs:      b.Logs,
		Logger:    b.Logger.Named("pbft"),
	})

	// starting the server
	b.Logger.Info("gRPC server started", zap.Int("port", port))
	if err := server.Serve(listener); err != nil {
		return fmt.Errorf("failed to start services: %v", err)
	}

	return nil
}
