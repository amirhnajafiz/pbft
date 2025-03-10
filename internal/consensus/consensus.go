package consensus

import (
	"sync"
	"time"

	"github.com/f24-cse535/pbft/internal/config/node/bft"
	"github.com/f24-cse535/pbft/internal/consensus/modules"
	"github.com/f24-cse535/pbft/internal/grpc/client"
	"github.com/f24-cse535/pbft/internal/storage/local"
	"github.com/f24-cse535/pbft/internal/storage/logs"
	"github.com/f24-cse535/pbft/pkg/enum"
	"github.com/f24-cse535/pbft/pkg/models"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"

	"go.dedis.ch/kyber/v4/pairing/bn256"
	"go.dedis.ch/kyber/v4/share"
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
	viewTimer     *modules.Timer

	// lock is used when an entity trys to update consensus fields
	lock sync.Mutex

	// threshold signature keys
	tss   *share.PriShare
	suite *bn256.Suite
	pub   *share.PubPoly

	inViewChangeMode        bool                     // a flag for in view change mode
	viewChangeGadgetChannel chan *pbft.ViewChangeMsg // view change gadget channel forwards the view change messages

	executionChannel chan int // execution channel is the execution handler input channel

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
	suite *bn256.Suite,
	share *share.PriShare,
	pub *share.PubPoly,
) *Consensus {
	// create a new consensus instance
	c := &Consensus{
		logs:             logs,
		memory:           mem,
		logger:           logr,
		cfg:              cfg,
		inViewChangeMode: false,
		lock:             sync.Mutex{},
		suite:            suite,
		tss:              share,
		pub:              pub,
	}

	// create consensus modules
	c.communication = modules.NewCommunicationModule(cli)
	c.waiter = modules.NewWaiter(cfg)
	c.viewTimer = modules.NewTimer(c.cfg.ViewTimeout, time.Second)

	// create consensus tables
	c.requestsHandlersTable = make(map[int]chan *models.Packet)
	c.consensusHandlersTable = map[enum.PacketType]chan *models.Packet{
		enum.PktPP:  make(chan *models.Packet, cfg.Total), // size of total
		enum.PktP:   make(chan *models.Packet, cfg.Total), // size of total
		enum.PktCmt: make(chan *models.Packet, cfg.Total), // size of total
		enum.PktVC:  make(chan *models.Packet, cfg.Total), // size of total
		enum.PktNV:  make(chan *models.Packet, cfg.Total), // size of total
		enum.PktCP:  make(chan *models.Packet, cfg.Total), // size of total
	}

	// create side channels
	c.executionChannel = make(chan int)

	// start handlers in go-routines
	go c.preprepareHandler()
	go c.prepareHandler()
	go c.commitHandler()
	go c.executeHandler()
	go c.timerHandler()
	go c.viewChangeHandler()
	go c.checkpointHandler()

	return c
}

// SignalToHandlers sends a packet from gRPC level to consensus handlers without waiting for a response.
func (c *Consensus) SignalToHandlers(pkt *models.Packet) {
	if ch, ok := c.consensusHandlersTable[pkt.Type]; ok {
		ch <- pkt
	}
}

// SignalToReqHandlers sends a packet from gRPC level to request handlers without waiting for a response.
func (c *Consensus) SignalToReqHandlers(pkt *models.Packet) {
	if pkt.Type == enum.PktReq {
		go c.requestHandler(pkt)
	} else if ch, ok := c.requestsHandlersTable[pkt.Sequence]; ok {
		ch <- pkt // if the request handler exists, pass the packet to it
	}
}
