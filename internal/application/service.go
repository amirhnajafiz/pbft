package application

import (
	"context"

	"github.com/f24-cse535/pbft/pkg/rpc/app"

	"google.golang.org/protobuf/types/known/emptypb"
)

// service is the gRPC server of our client app.
type service struct {
	app.UnimplementedAppServer

	channels chan *app.ReplyMsg
}

// Reply RPC forwards the reply message to the apps request handlers.
func (s *service) Reply(ctx context.Context, msg *app.ReplyMsg) (*emptypb.Empty, error) {
	s.channels <- msg

	return &emptypb.Empty{}, nil
}
