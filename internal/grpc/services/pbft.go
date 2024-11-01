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

func (p *PBFT) Commit(ctx context.Context, msg *pbft.CommitMsg) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (p *PBFT) PrePrepare(ctx context.Context, msg *pbft.PrePrepareMsg) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (p *PBFT) Prepare(ctx context.Context, msg *pbft.PrepareMsg) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (p *PBFT) Reply(ctx context.Context, msg *pbft.ReplyMsg) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (p *PBFT) Request(ctx context.Context, msg *pbft.RequestMsg) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (p *PBFT) PrintDB(_ *emptypb.Empty, stream pbft.PBFT_PrintDBServer) error {
	return nil
}

func (p *PBFT) PrintLog(_ *emptypb.Empty, stream pbft.PBFT_PrintLogServer) error {
	return nil
}

func (p *PBFT) PrintStatus(ctx context.Context, msg *pbft.StatusMsg) (*pbft.StatusRsp, error) {
	return nil, nil
}

func (p *PBFT) PrintView(ctx context.Context, msg *emptypb.Empty) (*pbft.ViewRsp, error) {
	return nil, nil
}

func (p *PBFT) Transaction(ctx context.Context, msg *pbft.TransactionMsg) (*pbft.TransactionRsp, error) {
	return nil, nil
}
