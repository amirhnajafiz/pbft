package consensus

import (
	"github.com/f24-cse535/pbft/internal/grpc/client"
	"github.com/f24-cse535/pbft/internal/storage/local"
	"github.com/f24-cse535/pbft/internal/storage/logs"
	"github.com/f24-cse535/pbft/pkg/enum"
	"github.com/f24-cse535/pbft/pkg/models"

	"go.uber.org/zap"
)

// Consensus module is the core module that runs consensus protocols by getting the gRPC level packets.
type Consensus struct {
	Client *client.Client // client is used to make RPC calls
	Logs   *logs.Logs     // data log is used to store and retrive logs
	Memory *local.Memory  // memory is needed to update the node state
	Logger *zap.Logger    // logger is needed for tracing

	interrupts     chan *models.InterruptMsg            // interrupts is an internal channel for dispatching gRPC level packets
	interruptTable map[enum.Interrupt]func(interface{}) // interrupt table is a map for interrupts and their handlers

	channels map[int]chan *models.InterruptMsg // the channels is used for communication between consensus sub-processes

	// these channels will be used by the client node to handle the user's transactions
	inTransactionChannel  chan *models.InterruptMsg
	outTransactionChannel chan interface{}
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
	// the main two signal and wait procedures are request and transaction
	if target == enum.IntrTransaction {
		// check to see if a transaction is in process or not
		if c.inTransactionChannel != nil {
			return nil
		}

		// create input channels
		c.inTransactionChannel = make(chan *models.InterruptMsg)
		c.outTransactionChannel = make(chan interface{})

		// call the proper transaction handler and return the channel
		go c.interruptTable[target](pkt)

		return c.outTransactionChannel
	}

	return nil
}

// Start consensus captures all gRPC level interrupts and dispatchs them to handlers.
func (c *Consensus) Start() {
	// create the interrupts channel
	c.interrupts = make(chan *models.InterruptMsg)
	c.channels = make(map[int]chan *models.InterruptMsg)

	// create a map of interrupts and handlers (interrupt table)
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

	// consensus loop
	go func() {
		for {
			intr := <-c.interrupts // capture interrupts

			// call the proper interrupt handler
			c.interruptTable[intr.Type](intr.Payload)

			// note: in this table, we expect to not get request and transaction interrupts
			// cause they need a signal and wait procedure.
		}
	}()
}
