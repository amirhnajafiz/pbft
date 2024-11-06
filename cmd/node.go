package cmd

import (
	"github.com/f24-cse535/pbft/internal/config"
	"github.com/f24-cse535/pbft/internal/consensus"
	"github.com/f24-cse535/pbft/internal/core"
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
	mem := local.NewMemory(n.Cfg.Node.NodeId, n.Cfg.Node.BFT.Total)
	mem.SetBalances(n.Cfg.GetClients())
	mem.SetNodes(n.Cfg.GetNodesMeta())

	// create a datalog instance
	datalog := logs.NewLogs()

	// create a new client.go
	cli := client.NewClient(
		n.Logger.Named("client.go"),
		n.Cfg.Node.NodeId,
		n.Cfg.GetNodes(),
	)
	if err := cli.LoadTLS(n.Cfg.TLS.PrivateKey, n.Cfg.TLS.PublicKey, n.Cfg.TLS.CaKey); err != nil {
		return err
	}

	// create a new gRPC bootstrap instance and execute the server by running the boot commands
	boot := grpc.Bootstrap{
		Port:       n.Cfg.Node.Port,
		PrivateKey: n.Cfg.TLS.PrivateKey,
		PublicKey:  n.Cfg.TLS.PublicKey,
		CAKey:      n.Cfg.TLS.CaKey,
		Memory:     mem,
		Logs:       datalog,
		Logger:     n.Logger.Named("grpc"),
	}

	// for core nodes we use consensus, for clients we use core
	if n.Cfg.Node.CoreNode {
		boot.Consensus = consensus.NewConsensus(datalog, mem, n.Logger.Named("consensus"), &n.Cfg.Node.BFT, cli)
	} else {
		boot.Core = core.NewCore(mem, cli, n.Logger.Named("core"), n.Cfg.Node.BFT.Responses)
	}

	return boot.ListenAnsServer()
}
