package consensus

import (
	"github.com/f24-cse535/pbft/internal/grpc/client"
	"github.com/f24-cse535/pbft/internal/storage/local"
	"github.com/f24-cse535/pbft/internal/storage/logs"
	"github.com/f24-cse535/pbft/pkg/enums"

	"go.uber.org/zap"
)

// Consensus module is the core module that runs consensus protocols by getting the gRPC level packets.
type Consensus struct {
	Client *client.Client // client is used to make RPC calls
	Logs   *logs.Logs     // data log is used to store and retrive logs
	Memory *local.Memory  // memory is needed to update the node state
	Logger *zap.Logger    // logger is needed for tracing

	channels map[enums.ChannelType]chan interface{}
}

// Signal sends a packet from gRPC level to handlers.
func (c *Consensus) Signal(target enums.ChannelType, pkt interface{}) {
	c.channels[target] <- pkt
}

// Start will initialize all channels and all handlers.
func (c *Consensus) Start() {
	// loop over all channels and create them
	for _, ct := range enums.ListChannelTypes() {
		c.channels[ct] = make(chan interface{})
	}

	// start all handlers in go-routines
	go c.handleCommit()
	go c.handlePrePrepare()
	go c.handlePrepare()
	go c.handleRequest()
	go c.handleReply()
}
