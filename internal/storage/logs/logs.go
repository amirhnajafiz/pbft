package logs

import (
	"sync"

	"github.com/f24-cse535/pbft/pkg/models"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"
)

// Logs is a memory type that stores the node's logs and datastore.
type Logs struct {
	logs        []string
	datastore   map[int]*models.Log
	viewChanges map[int][]*pbft.ViewChangeMsg

	lock sync.Mutex

	index int
}

// NewLogs returns a new logs instance.
func NewLogs() *Logs {
	return &Logs{
		lock:        sync.Mutex{},
		logs:        make([]string, 0),
		datastore:   make(map[int]*models.Log),
		viewChanges: make(map[int][]*pbft.ViewChangeMsg),
		index:       0,
	}
}
