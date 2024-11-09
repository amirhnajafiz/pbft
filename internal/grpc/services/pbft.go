package services

import (
	"context"

	"github.com/f24-cse535/pbft/internal/consensus"
	"github.com/f24-cse535/pbft/internal/storage/local"
	"github.com/f24-cse535/pbft/internal/storage/logs"
	"github.com/f24-cse535/pbft/pkg/enum"
	"github.com/f24-cse535/pbft/pkg/models"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

// pbft server is responsible for handling the protocol states and communication within nodes.
type PBFT struct {
	pbft.UnimplementedPBFTServer

	Consensus *consensus.Consensus
	Memory    *local.Memory
	Logs      *logs.Logs
	Logger    *zap.Logger
}

// Commit RPC generates a packet for consensus' commit handler.
func (p *PBFT) Commit(ctx context.Context, msg *pbft.AckMsg) (*emptypb.Empty, error) {
	p.Logs.AppendLog("Commit", msg.String())
	p.Consensus.SignalToHandlers(models.NewPacket(msg, enum.PktCmt, int(msg.GetSequenceNumber())))

	return &emptypb.Empty{}, nil
}

// PrePrepare RPC generates a packet for consensus' preprepare handler.
func (p *PBFT) PrePrepare(ctx context.Context, msg *pbft.PrePrepareMsg) (*emptypb.Empty, error) {
	p.Logs.AppendLog("PrePrepare", msg.String())
	p.Consensus.SignalToHandlers(models.NewPacket(msg, enum.PktPP, int(msg.GetSequenceNumber())))

	return &emptypb.Empty{}, nil
}

// Prepare RPC generates a packet for consensus' prepare handler.
func (p *PBFT) Prepare(ctx context.Context, msg *pbft.AckMsg) (*emptypb.Empty, error) {
	p.Logs.AppendLog("Prepare", msg.String())
	p.Consensus.SignalToHandlers(models.NewPacket(msg, enum.PktP, int(msg.GetSequenceNumber())))

	return &emptypb.Empty{}, nil
}

// Request RPC generates a packet for consensus' request handler.
func (p *PBFT) Request(ctx context.Context, msg *pbft.RequestMsg) (*emptypb.Empty, error) {
	p.Logs.AppendLog("Request", msg.String())
	p.Consensus.SignalToReqHandlers(models.NewPacket(msg, enum.PktReq, 0))

	return &emptypb.Empty{}, nil
}

// PrePrepared RPC generates a packet for consensus' request handler.
func (p *PBFT) PrePrepared(ctx context.Context, msg *pbft.AckMsg) (*emptypb.Empty, error) {
	p.Logs.AppendLog("PrePrepared", msg.String())
	p.Consensus.SignalToReqHandlers(models.NewPacket(msg, enum.PktPPed, int(msg.GetSequenceNumber())))

	return &emptypb.Empty{}, nil
}

// Prepared RPC generates a packet for consensus' request handler.
func (p *PBFT) Prepared(ctx context.Context, msg *pbft.AckMsg) (*emptypb.Empty, error) {
	p.Logs.AppendLog("Prepared", msg.String())
	p.Consensus.SignalToReqHandlers(models.NewPacket(msg, enum.PktPed, int(msg.GetSequenceNumber())))

	return &emptypb.Empty{}, nil
}

// ViewChange RPC generates a packet for consensus' view change handler.
func (p *PBFT) ViewChange(ctx context.Context, msg *pbft.ViewChangeMsg) (*emptypb.Empty, error) {
	p.Logs.AppendLog("ViewChange", msg.String())
	p.Consensus.SignalToHandlers(models.NewPacket(msg, enum.PktVC, int(msg.GetSequenceNumber())))

	return &emptypb.Empty{}, nil
}

// NewView RPC generates a packet for consensus' new view handler.
func (p *PBFT) NewView(ctx context.Context, msg *pbft.NewViewMsg) (*emptypb.Empty, error) {
	p.Logs.AppendLog("NewView", msg.String())
	p.Consensus.SignalToHandlers(models.NewPacket(msg, enum.PktNV, 0))

	return &emptypb.Empty{}, nil
}

// Checkpoint RPC generates a packet for consensus' checkpoint handler.
func (p *PBFT) Checkpoint(ctx context.Context, msg *pbft.CheckpointMsg) (*emptypb.Empty, error) {
	p.Logs.AppendLog("Checkpoint", msg.String())
	p.Consensus.SignalToHandlers(models.NewPacket(msg, enum.PktCP, int(msg.GetSequenceNumber())))

	return &emptypb.Empty{}, nil
}

// PrintDB returns the current datastore of this node.
func (p *PBFT) PrintDB(_ *emptypb.Empty, stream pbft.PBFT_PrintDBServer) error {
	ds := p.Logs.GetAllRequests()

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
	logs := p.Logs.GetLogs()

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
	status := pbft.StatusRsp{
		Status: -1,
	}

	if value := p.Logs.GetRequest(int(msg.GetSequenceNumber())); value != nil {
		status.Status = value.GetStatus()
	}

	return &status, nil
}

// PrintView exports the view messages of this node.
func (p *PBFT) PrintView(_ *emptypb.Empty, stream pbft.PBFT_PrintViewServer) error {
	views := p.Logs.GetAllViewChanges()

	// publish views one by one
	for key, value := range views {
		if err := stream.Send(&pbft.ViewRsp{
			View:     int64(key),
			Messages: value.ViewChangeMsgs,
			NewView:  value.NewViewMsg,
		}); err != nil {
			return err
		}
	}

	return nil
}

// PrintCheckpoints exports the checkpoint messages of this node.
func (p *PBFT) PrintCheckpoints(_ *emptypb.Empty, stream pbft.PBFT_PrintCheckpointsServer) error {
	views := p.Logs.GetCheckpoints()

	// publish views one by one
	for key, value := range views {
		if err := stream.Send(&pbft.CheckpointRsp{
			SequenceNumber: int64(key),
			Checkpoints:    value,
		}); err != nil {
			return err
		}
	}

	return nil
}
