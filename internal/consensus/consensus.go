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

// Consensus module is the core module that runs consensus protocols and handles the gRPC level packets.
type Consensus struct {
	Client *client.Client // client is used to make RPC calls
	Logs   *logs.Logs     // data log is used to store and retrive logs
	Memory *local.Memory  // memory is needed to update the node state
	Logger *zap.Logger    // logger is needed for tracing
	BFTCfg bft.Config     // bft config is used inside consensus handlers

	interrupts            chan *models.InterruptMsg            // interrupts is an internal channel for dispatching gRPC level packets
	interruptTable        map[enum.Interrupt]func(interface{}) // interrupt table is a map for interrupts and their handlers
	channels              map[int]chan *models.InterruptMsg    // the channels for communication between consensus handlers
	inTransactionChannel  chan *models.InterruptMsg            // input channel to send data to transaction handler
	outTransactionChannel chan interface{}                     // output channel to send data back to gRPC level
}

// Signal sends a packet from gRPC level to handlers without waiting for a response.
func (c *Consensus) Signal(target enum.Interrupt, pkt interface{}) {
	c.interrupts <- &models.InterruptMsg{
		Type:    target,
		Payload: pkt,
	}
}

// SignalAndWait sends a packet from gRPC level to handlers and sends a channel to get the response.
func (c *Consensus) SignalAndWait(target enum.Interrupt, pkt interface{}) chan interface{} {
	if target == enum.IntrTransaction {
		// check to see if a transaction is in process or not
		if c.inTransactionChannel != nil {
			return nil
		}

		// initial transaction handler in and out channels
		c.inTransactionChannel = make(chan *models.InterruptMsg)
		c.outTransactionChannel = make(chan interface{})

		// call the proper transaction handler
		go c.interruptTable[enum.IntrTransaction](pkt)

		// return the out channel
		return c.outTransactionChannel
	}

	return nil
}

// Start consensus registers a go-routine that captures all gRPC level interrupts and forwards them to handlers.
func (c *Consensus) Start() {
	// initial consensus fields
	c.interrupts = make(chan *models.InterruptMsg)
	c.channels = make(map[int]chan *models.InterruptMsg)
	c.interruptTable = map[enum.Interrupt]func(interface{}){
		enum.IntrCommit:      c.handleCommit,
		enum.IntrPrePrepare:  c.handlePrePrepare,
		enum.IntrPrePrepared: c.handlePrePrepared,
		enum.IntrPrepare:     c.handlePrepare,
		enum.IntrPrepared:    c.handlePrepared,
		enum.IntrReply:       c.handleReply,
		enum.IntrRequest:     c.handleRequest,
		enum.IntrTransaction: c.handleTransaction,
	}

	go func() {
		// consensus loop
		for {
			intr := <-c.interrupts                    // capture interrupts
			c.interruptTable[intr.Type](intr.Payload) // call the proper interrupt handler
		}
	}()
}
