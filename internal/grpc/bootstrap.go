package grpc

import (
	"fmt"
	"net"

	"github.com/f24-cse535/pbft/internal/consensus"
	"github.com/f24-cse535/pbft/internal/grpc/services"
	"github.com/f24-cse535/pbft/internal/monitoring/metrics"
	"github.com/f24-cse535/pbft/pkg/rpc/liveness"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Bootstrap is a wrapper that holds requirements for the gRPC services.
type Bootstrap struct {
	Port int

	Consensus *consensus.Consensus // consensus module is the core module
	Logger    *zap.Logger          // logger is needed for tracing
	Metrics   *metrics.Metrics     // metrics is needed for monitoring
}

// ListenAnsServer creates a new gRPC instance with all services.
func (b *Bootstrap) ListenAnsServer() error {
	// on the local network, listen to a port
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", b.Port))
	if err != nil {
		return fmt.Errorf("failed to start the listener server: %v", err)
	}

	// create a new grpc instance
	server := grpc.NewServer(
		grpc.UnaryInterceptor(b.allUnaryInterceptor),   // set an unary interceptor
		grpc.StreamInterceptor(b.allStreamInterceptor), // set a stream interceptor
	)

	// register all gRPC services
	liveness.RegisterLivenessServer(server, &services.Liveness{
		Consensus: b.Consensus,
		Logger:    b.Logger.Named("liveness"),
	})
	pbft.RegisterPBFTServer(server, &services.PBFT{
		Consensus: b.Consensus,
		Logger:    b.Logger.Named("pbft"),
		Metrics:   b.Metrics,
	})

	// starting the server
	b.Logger.Info("grpc server started", zap.Int("port", b.Port))
	if err := server.Serve(listener); err != nil {
		return fmt.Errorf("failed to start the server: %v", err)
	}

	return nil
}
