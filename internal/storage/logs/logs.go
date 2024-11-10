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
	viewChanges map[int]*models.ViewLog
	checkpoints map[int][]*pbft.CheckpointMsg

	lock              sync.Mutex
	index             int
	lastCheckpoint    int
	kvalue            int
	lastProcessingSeq int
}

// NewLogs returns a new logs instance.
func NewLogs(kvalue int) *Logs {
	return &Logs{
		lock:              sync.Mutex{},
		logs:              make([]string, 0),
		datastore:         make(map[int]*models.Log),
		viewChanges:       make(map[int]*models.ViewLog),
		checkpoints:       make(map[int][]*pbft.CheckpointMsg),
		index:             0,
		lastCheckpoint:    0,
		kvalue:            kvalue,
		lastProcessingSeq: 0,
	}
}
