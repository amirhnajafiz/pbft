package services

import (
	"context"

	"github.com/f24-cse535/pbft/internal/consensus"
	"github.com/f24-cse535/pbft/pkg/enum"
	"github.com/f24-cse535/pbft/pkg/models"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"

	"go.uber.org/zap"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// pbft server is responsible for handling the protocol states and communication within nodes.
type PBFT struct {
	pbft.UnimplementedPBFTServer

	Consensus *consensus.Consensus
	Logger    *zap.Logger
}

// Commit RPC generates a packet for consensus' commit handler.
func (p *PBFT) Commit(ctx context.Context, msg *pbft.AckMsg) (*emptypb.Empty, error) {
	p.Consensus.Logs.AppendLog("Commit", msg.String())
	p.Consensus.SignalToHandlers(models.NewPacket(msg, enum.PktCmt, int(msg.GetSequenceNumber())))

	return &emptypb.Empty{}, nil
}

// PrePrepare RPC generates a packet for consensus' preprepare handler.
func (p *PBFT) PrePrepare(ctx context.Context, msg *pbft.PrePrepareMsg) (*emptypb.Empty, error) {
	p.Consensus.Logs.AppendLog("PrePrepare", msg.String())
	p.Consensus.SignalToHandlers(models.NewPacket(msg, enum.PktPP, int(msg.GetSequenceNumber())))

	return &emptypb.Empty{}, nil
}

// Prepare RPC generates a packet for consensus' prepare handler.
func (p *PBFT) Prepare(ctx context.Context, msg *pbft.AckMsg) (*emptypb.Empty, error) {
	p.Consensus.Logs.AppendLog("Prepare", msg.String())
	p.Consensus.SignalToHandlers(models.NewPacket(msg, enum.PktP, int(msg.GetSequenceNumber())))

	return &emptypb.Empty{}, nil
}

// Request RPC generates a packet for consensus' request handler.
func (p *PBFT) Request(ctx context.Context, msg *pbft.RequestMsg) (*emptypb.Empty, error) {
	p.Consensus.Logs.AppendLog("Request", msg.String())
	p.Consensus.SignalToReqHandlers(models.NewPacket(msg, enum.PktReq, 0))

	return &emptypb.Empty{}, nil
}

// PrePrepared RPC generates a packet for consensus' request handler.
func (p *PBFT) PrePrepared(ctx context.Context, msg *pbft.AckMsg) (*emptypb.Empty, error) {
	p.Consensus.Logs.AppendLog("PrePrepared", msg.String())
	p.Consensus.SignalToReqHandlers(models.NewPacket(msg, enum.PktPPed, int(msg.GetSequenceNumber())))

	return &emptypb.Empty{}, nil
}

// Prepared RPC generates a packet for consensus' request handler.
func (p *PBFT) Prepared(ctx context.Context, msg *pbft.AckMsg) (*emptypb.Empty, error) {
	p.Consensus.Logs.AppendLog("Prepared", msg.String())
	p.Consensus.SignalToReqHandlers(models.NewPacket(msg, enum.PktPed, int(msg.GetSequenceNumber())))

	return &emptypb.Empty{}, nil
}

// Transaction RPC calls signal and wait on consensus and waits for a response.
func (p *PBFT) Transaction(ctx context.Context, msg *pbft.TransactionMsg) (*pbft.TransactionRsp, error) {
	return nil, nil
}

// Reply RPC forwards a reply message into consensus.handleReply
func (p *PBFT) Reply(ctx context.Context, msg *pbft.ReplyMsg) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

// PrintDB returns the current datastore of this node.
func (p *PBFT) PrintDB(_ *emptypb.Empty, stream pbft.PBFT_PrintDBServer) error {
	ds := p.Consensus.Logs.GetAllRequests()

	// publish requests one by one
	for _, block := range ds {
		if err := stream.Send(block); err != nil {
			return err
		}
	}

	return nil
}

// PrintLog returns the datalog of this node.
func (p *PBFT) PrintLog(_ *emptypb.Empty, stream pbft.PBFT_PrintLogServer) error {
	logs := p.Consensus.Logs.GetLogs()

	// publish logs one by one
	for _, block := range logs {
		if err := stream.Send(&pbft.LogRsp{
			Text: block,
		}); err != nil {
			return err
		}
	}

	return nil
}

// PrintStatus gets a sequence number and returns the status of its log.
func (p *PBFT) PrintStatus(ctx context.Context, msg *pbft.StatusMsg) (*pbft.StatusRsp, error) {
	if value := p.Consensus.Logs.GetRequest(int(msg.GetSequenceNumber())); value != nil {
		return &pbft.StatusRsp{
			Status: value.GetStatus(),
		}, nil
	}

	return nil, status.Error(5, "item not found")
}

func (p *PBFT) PrintView(ctx context.Context, msg *emptypb.Empty) (*pbft.ViewRsp, error) {
	return nil, nil
}
