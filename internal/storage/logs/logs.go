package logs

import (
	"sync"

	"github.com/f24-cse535/pbft/pkg/rpc/pbft"
)

// Logs is a memory type that stores the node's logs and datastore.
type Logs struct {
	logs        []string
	datastore   map[int]*pbft.RequestMsg
	viewChanges map[int][]interface{}

	lock sync.Mutex

	index int
}

// NewLogs returns a new logs instance.
func NewLogs() *Logs {
	return &Logs{
		lock:        sync.Mutex{},
		logs:        make([]string, 0),
		datastore:   make(map[int]*pbft.RequestMsg),
		viewChanges: make(map[int][]interface{}),
		index:       0,
	}
}
