package cmd

import (
	"github.com/f24-cse535/pbft/internal/config"
	"github.com/f24-cse535/pbft/internal/consensus"
	"github.com/f24-cse535/pbft/internal/grpc"
	"github.com/f24-cse535/pbft/internal/grpc/client"
	"github.com/f24-cse535/pbft/internal/storage/local"
	"github.com/f24-cse535/pbft/internal/storage/logs"

	"go.dedis.ch/kyber/v4/pairing/bn256"
	"go.dedis.ch/kyber/v4/share"
	"go.uber.org/zap"
)

// Node is a app in our transaction system.
type Node struct {
	Cfg    config.Config
	Logger *zap.Logger
	Suite  *bn256.Suite
	Share  *share.PriShare
	Pub    *share.PubPoly
}

func (n Node) Main() error {
	// create a memory instance
	mem := local.NewMemory(n.Cfg.Node.NodeId, n.Cfg.Node.BFT.Total, n.Cfg.Node.BFT.KWatermark)
	mem.SetBalances(n.Cfg.GetClients())
	mem.SetNodes(n.Cfg.GetNodesMeta())

	// create a datalog instance
	datalog := logs.NewLogs()

	// load tls configs
	creds, err := n.Cfg.TLS.Creds()
	if err != nil {
		return err
	}

	// create a new client.go
	cli := client.NewClient(
		creds,
		n.Cfg.Node.NodeId,
		n.Cfg.GetNodes(),
	)

	// create a new gRPC bootstrap instance and execute the server by running the boot commands
	boot := grpc.Bootstrap{
		Memory: mem,
		Logs:   datalog,
		Logger: n.Logger.Named("grpc"),
		Consensus: consensus.NewConsensus(
			datalog,
			mem,
			n.Logger.Named("consensus"),
			&n.Cfg.Node.BFT,
			cli,
			n.Suite,
			n.Share,
			n.Pub,
		),
	}

	// start the gRPC server
	return boot.ListenAnsServer(n.Cfg.Node.Port, creds)
}
