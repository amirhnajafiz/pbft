package services

import (
	"context"

	"github.com/f24-cse535/pbft/internal/consensus"
	"github.com/f24-cse535/pbft/pkg/rpc/liveness"

	"go.uber.org/zap"
	"google.golang.org/grpc/status"
)

// liveness server handles the running state of the gRPC server.
type Liveness struct {
	liveness.UnimplementedLivenessServer

	Consensus *consensus.Consensus
	Logger    *zap.Logger
}

// Ping RPC is used to check if a server is alive and can process or not.
func (l *Liveness) Ping(ctx context.Context, input *liveness.LivePingMessage) (*liveness.LivePingMessage, error) {
	if !l.Consensus.Memory.GetStatus() {
		return nil, status.Error(13, "service is not responding")
	}

	return &liveness.LivePingMessage{}, nil
}

// ChangeStatus is used to update the liveness of the gRPC server.
func (l *Liveness) ChangeStatus(ctx context.Context, input *liveness.LiveChangeStatusMessage) (*liveness.LiveChangeStatusMessage, error) {
	l.Consensus.Memory.SetStatus(input.GetStatus())

	return &liveness.LiveChangeStatusMessage{Status: l.Consensus.Memory.GetStatus()}, nil
}
