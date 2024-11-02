package cmd

import (
	"github.com/f24-cse535/pbft/internal/config"
	"github.com/f24-cse535/pbft/internal/consensus"
	"github.com/f24-cse535/pbft/internal/grpc"
	"github.com/f24-cse535/pbft/internal/grpc/client"
	"github.com/f24-cse535/pbft/internal/storage/local"
	"github.com/f24-cse535/pbft/internal/storage/logs"

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

	// create a datalog instance
	datalog := logs.NewLogs()

	// create a new client.go
	cli := client.NewClient(
		n.Logger.Named("client.go"),
		n.Cfg.Node.NodeId,
		n.Cfg.IPTable.GetNodes(),
		n.Cfg.IPTable.GetClients(),
	)

	// create a new consensus module
	instance := consensus.Consensus{
		Client: cli,
		Memory: mem,
		Logs:   datalog,
		Logger: n.Logger.Named("consensus"),
	}
	instance.Start()

	// create a new gRPC bootstrap instance and execute the server by running the boot commands
	boot := grpc.Bootstrap{
		Port:      n.Cfg.Node.Port,
		Logger:    n.Logger.Named("grpc"),
		Consensus: &instance,
	}

	return boot.ListenAnsServer()
}
