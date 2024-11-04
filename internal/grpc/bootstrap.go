package grpc

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"

	"github.com/f24-cse535/pbft/internal/consensus"
	"github.com/f24-cse535/pbft/internal/grpc/services"
	"github.com/f24-cse535/pbft/pkg/rpc/liveness"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/examples/data"
)

// Bootstrap is a wrapper that holds requirements for the gRPC services.
type Bootstrap struct {
	Port       int
	PrivateKey string
	PublicKey  string
	CAC        string

	Consensus *consensus.Consensus // consensus module is the core module
	Logger    *zap.Logger          // logger is needed for tracing
}

// ListenAnsServer creates a new gRPC instance with all services.
func (b *Bootstrap) ListenAnsServer() error {
	// on the local network, listen to a port
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", b.Port))
	if err != nil {
		return fmt.Errorf("failed to start the core listener server: %v", err)
	}

	// load server's certificate
	cert, err := tls.LoadX509KeyPair(data.Path(b.PublicKey), data.Path(b.PrivateKey))
	if err != nil {
		return fmt.Errorf("failed to load certificates: %v", err)
	}

	// create the CA data
	ca := x509.NewCertPool()
	cac := data.Path(b.CAC)
	cacBytes, err := os.ReadFile(cac)
	if err != nil {
		return fmt.Errorf("failed to read ca cert: %v", err)
	}
	if ok := ca.AppendCertsFromPEM(cacBytes); !ok {
		return fmt.Errorf("failed to append certs: %s", cac)
	}

	// tls configs
	tlsConfig := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
		ClientCAs:    ca,
	}

	// create a new grpc instance
	server := grpc.NewServer(
		grpc.Creds(credentials.NewTLS(tlsConfig)),      // set the server certificates
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
	})

	// starting the server
	b.Logger.Info("gRPC server started", zap.Int("port", b.Port))
	if err := server.Serve(listener); err != nil {
		return fmt.Errorf("failed to start servers: %v", err)
	}

	return nil
}
