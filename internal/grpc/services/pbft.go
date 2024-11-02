package services

import (
	"context"

	"github.com/f24-cse535/pbft/internal/consensus"
	"github.com/f24-cse535/pbft/pkg/enum"
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

// Commit RPC forwards a commit message into consensus.handleCommit
func (p *PBFT) Commit(ctx context.Context, msg *pbft.CommitMsg) (*emptypb.Empty, error) {
	p.Consensus.Signal(enum.IntrCommit, msg)

	return &emptypb.Empty{}, nil
}

// PrePrepare RPC forwards a preprepare message into consensus.handlePrePrepare
func (p *PBFT) PrePrepare(ctx context.Context, msg *pbft.PrePrepareMsg) (*emptypb.Empty, error) {
	p.Consensus.Signal(enum.IntrPrePrepare, msg)

	return &emptypb.Empty{}, nil
}

// PrePrepared RPC forwards a preprepared message into consensus.handlePrePrepared
func (p *PBFT) PrePrepared(ctx context.Context, msg *pbft.PrePreparedMsg) (*emptypb.Empty, error) {
	p.Consensus.Signal(enum.IntrPrePrepared, msg)

	return &emptypb.Empty{}, nil
}

// Prepare RPC forwards a prepare message into consensus.handlePrepare
func (p *PBFT) Prepare(ctx context.Context, msg *pbft.PrepareMsg) (*emptypb.Empty, error) {
	p.Consensus.Signal(enum.IntrPrepare, msg)

	return &emptypb.Empty{}, nil
}

// Prepared RPC forwards a prepared message into consensus.handlePrepared
func (p *PBFT) Prepared(ctx context.Context, msg *pbft.PreparedMsg) (*emptypb.Empty, error) {
	p.Consensus.Signal(enum.IntrPrepared, msg)

	return &emptypb.Empty{}, nil
}

// Reply RPC forwards a reply message into consensus.handleReply
func (p *PBFT) Reply(ctx context.Context, msg *pbft.ReplyMsg) (*emptypb.Empty, error) {
	p.Consensus.Signal(enum.IntrReply, msg)

	return &emptypb.Empty{}, nil
}

// Request RPC forwards a request message into consensus.handleRequest
func (p *PBFT) Request(ctx context.Context, msg *pbft.RequestMsg) (*emptypb.Empty, error) {
	p.Consensus.Signal(enum.IntrRequest, msg)

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

// Transaction RPC calls signal and wait on consensus and waits for a response.
func (p *PBFT) Transaction(ctx context.Context, msg *pbft.TransactionMsg) (*pbft.TransactionRsp, error) {
	resp := pbft.TransactionRsp{}

	// call signal and wait
	ch := p.Consensus.SignalAndWait(enum.IntrTransaction, msg)
	if ch != nil {
		text := <-ch
		resp.Text = text.(string)
	} else {
		// the channel is not returned, it means it is in progress
		resp.Text = "server is busy with processing another transaction"
	}

	return &resp, nil
}
