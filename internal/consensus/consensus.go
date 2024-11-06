package consensus

import (
	"github.com/f24-cse535/pbft/internal/config/node/bft"
	"github.com/f24-cse535/pbft/internal/grpc/client"
	"github.com/f24-cse535/pbft/internal/storage/local"
	"github.com/f24-cse535/pbft/internal/storage/logs"
	"github.com/f24-cse535/pbft/pkg/enum"
	"github.com/f24-cse535/pbft/pkg/models"

	"go.uber.org/zap"
)

// Consensus module is the main module of PBFT that manages packets.
type Consensus struct {
	Client *client.Client // client is used to make RPC calls
	Logs   *logs.Logs     // data log is used to store and retrive logs
	Memory *local.Memory  // memory is needed to update the node state
	Logger *zap.Logger    // logger is needed for tracing
	BFTCfg *bft.Config    // bft config is used inside consensus handlers

	consensusHandlersTable map[enum.PacketType]chan *models.Packet // a map of consensus handlers and their input channels
	requestsHandlersTable  map[int]chan *models.Packet             // a map of requests handlers and their input channels
}

// Init method, initializes the consensus fields.
func (c *Consensus) Init() {
	// create consensus tables
	c.requestsHandlersTable = make(map[int]chan *models.Packet, c.BFTCfg.Total*2) // size of 2*total nodes
	c.consensusHandlersTable = map[enum.PacketType]chan *models.Packet{
		enum.PktPP:  make(chan *models.Packet, c.BFTCfg.Total), // size of total
		enum.PktP:   make(chan *models.Packet, c.BFTCfg.Total), // size of total
		enum.PktCmt: make(chan *models.Packet, c.BFTCfg.Total), // size of total
	}

	// start handlers go-routines
	go c.preprepareHandler()
	go c.prepareHandler()
	go c.commitHandler()
}

// SignalToHandlers sends a packet from gRPC level to consensus handlers without waiting for a response.
func (c *Consensus) SignalToHandlers(pkt *models.Packet) {
	if ch, ok := c.consensusHandlersTable[pkt.Type]; ok {
		ch <- pkt
	}
}

// SignalToReqHandlers sends a packet from gRPC level to request handlers without waiting for a response.
func (c *Consensus) SignalToReqHandlers(pkt *models.Packet) {
	sequence := pkt.Headers["seq"].(int)

	// if the request handler exists, pass the packet to it
	if ch, ok := c.requestsHandlersTable[sequence]; ok {
		ch <- pkt
	} else if pkt.Type == enum.PktReq {
		// if a new request is arrived, create a new channel and handler
		c.requestsHandlersTable[sequence] = make(chan *models.Packet, c.BFTCfg.Total) // size of total
		c.requestsHandlersTable[sequence] <- pkt

		go c.requestHandler()
	}
}

// // SignalAndGo sends a packet from gRPC level to a new go-routine handler.
// func (c *Consensus) SignalAndGo(target enum.Interrupt, pkt interface{}) {
// 	if target == enum.IntrCommit || target == enum.IntrRequest {
// 		go c.interruptTable[target](pkt)
// 	}
// }

// // SignalAndWait sends a packet from gRPC level to handlers and sends a channel to get the response.
// func (c *Consensus) SignalAndWait(target enum.Interrupt, pkt interface{}) chan interface{} {
// 	if target == enum.IntrTransaction {
// 		// check to see if a transaction is in process or not
// 		if c.inTransactionChannel != nil {
// 			return nil
// 		}

// 		// initial transaction handler in and out channels
// 		c.inTransactionChannel = make(chan *models.InterruptMsg)
// 		c.outTransactionChannel = make(chan interface{})

// 		// call the proper transaction handler
// 		go c.interruptTable[enum.IntrTransaction](pkt)

// 		// return the out channel
// 		return c.outTransactionChannel
// 	}

// 	return nil
// }
