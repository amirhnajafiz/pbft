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

// Node is a app in our transaction system.
type Node struct {
	Cfg    config.Config
	Logger *zap.Logger
}

func (n Node) Main() error {
	// create a memory instance
	mem := local.NewMemory(n.Cfg.Node.NodeId, n.Cfg.Node.BFT.Total)
	mem.SetBalances(n.Cfg.GetClients())
	mem.SetNodes(n.Cfg.GetNodesMeta())

	// create a datalog instance
	datalog := logs.NewLogs()

	// load tls configs
	tlsC, _ := n.Cfg.TLS.TLS()

	// create a new client.go
	cli := client.NewClient(
		n.Logger.Named("client.go"),
		tlsC,
		n.Cfg.Node.NodeId,
		n.Cfg.GetNodes(),
	)

	// create a new gRPC bootstrap instance and execute the server by running the boot commands
	boot := grpc.Bootstrap{
		Port:       n.Cfg.Node.Port,
		PrivateKey: n.Cfg.TLS.PrivateKey,
		PublicKey:  n.Cfg.TLS.PublicKey,
		CAKey:      n.Cfg.TLS.CaKey,
		Memory:     mem,
		Logs:       datalog,
		Logger:     n.Logger.Named("grpc"),
		Consensus:  consensus.NewConsensus(datalog, mem, n.Logger.Named("consensus"), &n.Cfg.Node.BFT, cli),
	}

	// start the gRPC server
	return boot.ListenAnsServer()
}
