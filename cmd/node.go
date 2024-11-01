package cmd

import (
	"github.com/f24-cse535/pbft/internal/config"
	"github.com/f24-cse535/pbft/internal/consensus"
	"github.com/f24-cse535/pbft/internal/grpc"
	"github.com/f24-cse535/pbft/internal/storage/local"

	"go.uber.org/zap"
)

// Each node of our transaction system runs using this main function.
type Node struct {
	Cfg    config.Config
	Logger *zap.Logger
}

func (n Node) Main() error {
	// create a local storage (aka memory)
	mem := local.NewMemory()

	// create a new consensus module
	instance := consensus.Consensus{
		Memory: mem,
		Logger: n.Logger.Named("consensus"),
	}

	// create a new gRPC bootstrap instance and execute the server by running the boot commands
	boot := grpc.Bootstrap{
		Port:      n.Cfg.Node.Port,
		Logger:    n.Logger.Named("grpc"),
		Consensus: &instance,
	}

	return boot.ListenAnsServer()
}
