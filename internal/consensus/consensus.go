package consensus

import (
	"github.com/f24-cse535/pbft/internal/storage/local"

	"go.uber.org/zap"
)

// Consensus module is the core module that runs consensus protocols by getting the gRPC level packets.
type Consensus struct {
	Memory *local.Memory // memory is needed to update the node state
	Logger *zap.Logger   // logger is needed for tracing
}
