package consensus

import (
	"github.com/f24-cse535/pbft/internal/grpc/client"
	"github.com/f24-cse535/pbft/internal/storage/local"
	"github.com/f24-cse535/pbft/internal/storage/logs"

	"go.uber.org/zap"
)

// Consensus module is the core module that runs consensus protocols by getting the gRPC level packets.
type Consensus struct {
	Client *client.Client // client is used to make RPC calls
	Logs   *logs.Logs     // data log is used to store and retrive logs
	Memory *local.Memory  // memory is needed to update the node state
	Logger *zap.Logger    // logger is needed for tracing
}
