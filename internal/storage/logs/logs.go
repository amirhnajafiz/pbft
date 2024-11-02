package logs

import (
	"github.com/f24-cse535/pbft/pkg/models"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"
)

// Logs is a memory type that stores the node's logs and datastore.
type Logs struct {
	logs      map[int]*models.Log
	datastore map[int]*pbft.RequestMsg

	index int
}

// NewLogs returns a new logs instance.
func NewLogs() *Logs {
	return &Logs{
		logs:      make(map[int]*models.Log),
		datastore: make(map[int]*pbft.RequestMsg),
		index:     0,
	}
}
