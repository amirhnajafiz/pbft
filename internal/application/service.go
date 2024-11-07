package application

import (
	"context"

	"github.com/f24-cse535/pbft/pkg/rpc/app"

	"google.golang.org/protobuf/types/known/emptypb"
)

// service is the gRPC server of our client app.
type service struct {
	app.UnimplementedAppServer

	channel chan *app.ReplyMsg
}

// Reply RPC forwards the reply message to the app channel.
func (s *service) Reply(ctx context.Context, msg *app.ReplyMsg) (*emptypb.Empty, error) {
	s.channel <- msg

	return &emptypb.Empty{}, nil
}
