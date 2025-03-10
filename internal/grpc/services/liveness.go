package services

import (
	"context"

	"github.com/f24-cse535/pbft/internal/storage/local"
	"github.com/f24-cse535/pbft/internal/storage/logs"
	"github.com/f24-cse535/pbft/pkg/rpc/liveness"

	"go.uber.org/zap"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// liveness server handles the running state of the gRPC server.
type Liveness struct {
	liveness.UnimplementedLivenessServer

	Memory *local.Memory
	Logs   *logs.Logs
	Logger *zap.Logger
}

// Ping RPC is used to check if a server is alive and can process or not.
func (l *Liveness) Ping(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	if !l.Memory.GetStatus() {
		return nil, status.Error(13, "service is not responding")
	}

	return &emptypb.Empty{}, nil
}

// ChangeStatus is used to update the liveness fields of the gRPC server.
func (l *Liveness) ChangeStatus(ctx context.Context, input *liveness.StatusMsg) (*emptypb.Empty, error) {
	l.Memory.SetStatus(input.GetStatus())
	l.Memory.SetByzantine(input.GetByzantine())

	return &emptypb.Empty{}, nil
}

// Flush is used to remove everything from the node's memory.
func (l *Liveness) Flush(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	l.Memory.Reset()
	l.Logs.Reset()

	return &emptypb.Empty{}, nil
}
