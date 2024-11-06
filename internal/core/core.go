package core

import (
	"github.com/f24-cse535/pbft/internal/grpc/client"
	"github.com/f24-cse535/pbft/internal/storage/local"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"
)

// Core is the client main module that get's users requests, and sends them to nodes.
type Core struct {
	memory *local.Memory
	cli    *client.Client
	replys int

	handlers map[string]chan *pbft.ReplyMsg
	clients  map[string]chan string
}

// NewTransaction creates a new request handler and returns a channel for user response.
func (c *Core) NewTransaction(client string, trx *pbft.TransactionMsg) chan string {
	c.handlers[client] = make(chan *pbft.ReplyMsg)
	c.clients[client] = make(chan string)

	go c.requestHandler(client, trx)

	return c.clients[client]
}

// NewReply sends a reply message to a transaction handler if it exists.
func (c *Core) NewReply(client string, rply *pbft.ReplyMsg) {
	if ch, ok := c.handlers[client]; ok {
		ch <- rply
	}
}

// Done sends a reply message to the client and closes the channels.
func (c *Core) Done(client, txt string) {
	c.clients[client] <- txt

	delete(c.handlers, client)
	delete(c.clients, client)
}

// NewCore returns an instance of the core struct.
func NewCore(mem *local.Memory, cli *client.Client, replys int) *Core {
	return &Core{
		memory:   mem,
		cli:      cli,
		replys:   replys,
		handlers: make(map[string]chan *pbft.ReplyMsg),
		clients:  make(map[string]chan string),
	}
}
