package consensus

import (
	"github.com/f24-cse535/pbft/internal/config/node/bft"
	"github.com/f24-cse535/pbft/internal/consensus/modules"
	"github.com/f24-cse535/pbft/internal/grpc/client"
	"github.com/f24-cse535/pbft/internal/storage/local"
	"github.com/f24-cse535/pbft/internal/storage/logs"
	"github.com/f24-cse535/pbft/pkg/enum"
	"github.com/f24-cse535/pbft/pkg/models"

	"go.uber.org/zap"
)

// Consensus module is the main module of PBFT that manages packets.
type Consensus struct {
	logs          *logs.Logs
	cfg           *bft.Config
	memory        *local.Memory
	logger        *zap.Logger
	communication *modules.Communication
	waiter        *modules.Waiter

	consensusHandlersTable map[enum.PacketType]chan *models.Packet // a map of consensus handlers and their input channels
	requestsHandlersTable  map[int]chan *models.Packet             // a map of requests handlers and their input channels
}

// NewConsensus returns a consensus instance.
func NewConsensus(
	logs *logs.Logs,
	mem *local.Memory,
	logr *zap.Logger,
	cfg *bft.Config,
	cli *client.Client,
) *Consensus {
	// create a new consensus instance
	c := &Consensus{
		logs:   logs,
		memory: mem,
		logger: logr,
		cfg:    cfg,
	}

	// create consensus modules
	c.communication = modules.NewCommunicationModule(cli)
	c.waiter = modules.NewWaiter(cfg)

	// create consensus tables
	c.requestsHandlersTable = make(map[int]chan *models.Packet)
	c.consensusHandlersTable = map[enum.PacketType]chan *models.Packet{
		enum.PktPP:  make(chan *models.Packet, cfg.Total), // size of total
		enum.PktP:   make(chan *models.Packet, cfg.Total), // size of total
		enum.PktCmt: make(chan *models.Packet, cfg.Total), // size of total
	}

	// start handlers go-routines
	go c.preprepareHandler()
	go c.prepareHandler()
	go c.commitHandler()

	return c
}

// SignalToHandlers sends a packet from gRPC level to consensus handlers without waiting for a response.
func (c *Consensus) SignalToHandlers(pkt *models.Packet) {
	if ch, ok := c.consensusHandlersTable[pkt.Type]; ok {
		c.logger.Debug("received a signal to handler", zap.Any("handler", pkt.Type))

		ch <- pkt
	}
}

// SignalToReqHandlers sends a packet from gRPC level to request handlers without waiting for a response.
func (c *Consensus) SignalToReqHandlers(pkt *models.Packet) {
	c.logger.Debug("received a signal to request handler", zap.Any("handler", pkt.Type), zap.Int("sequence", pkt.Sequence))

	if ch, ok := c.requestsHandlersTable[pkt.Sequence]; ok {
		ch <- pkt // if the request handler exists, pass the packet to it
	} else if pkt.Type == enum.PktReq {
		go c.requestHandler(pkt) // if a new request is arrived create handler
	}
}
