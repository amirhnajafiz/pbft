package services

import (
	"context"

	"github.com/f24-cse535/pbft/internal/consensus"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

// pbft server is responsible for handling the protocol states and communication within nodes.
type PBFT struct {
	pbft.UnimplementedPBFTServer

	Consensus *consensus.Consensus
	Logger    *zap.Logger
}

func (p *PBFT) Commit(context.Context, *pbft.CommitMsg) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (p *PBFT) PrePrepare(context.Context, *pbft.PrePrepareMsg) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (p *PBFT) Prepare(context.Context, *pbft.PrepareMsg) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (p *PBFT) Reply(context.Context, *pbft.ReplyMsg) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (p *PBFT) Request(context.Context, *pbft.RequestMsg) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
