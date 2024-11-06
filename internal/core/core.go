package core

import "github.com/f24-cse535/pbft/pkg/models"

// Core is the client main module that get's users requests, and sends them to nodes.
type Core struct {
	clients map[string]chan *models.Packet
}

// NewCore returns an instance of the core struct.
func NewCore() *Core {
	return &Core{
		clients: make(map[string]chan *models.Packet),
	}
}
