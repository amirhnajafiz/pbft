package services

import (
	"context"

	"github.com/f24-cse535/pbft/internal/consensus"
	"github.com/f24-cse535/pbft/pkg/rpc/controller"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

// controller server is handling RPC calls of controller app.
type Controller struct {
	controller.UnimplementedControllerServer

	Consensus *consensus.Consensus
	Logger    *zap.Logger
}

func (c *Controller) PrintDB(_ *emptypb.Empty, stream controller.Controller_PrintDBServer) error {
	return nil
}

func (c *Controller) PrintLog(_ *emptypb.Empty, stream controller.Controller_PrintLogServer) error {
	return nil
}

func (c *Controller) PrintStatus(ctx context.Context, msg *controller.StatusMsg) (*controller.StatusRsp, error) {
	return nil, nil
}

func (c *Controller) PrintView(ctx context.Context, msg *emptypb.Empty) (*controller.ViewRsp, error) {
	return nil, nil
}

func (c *Controller) Transaction(ctx context.Context, msg *controller.TransactionMsg) (*controller.TransactionRsp, error) {
	return nil, nil
}
