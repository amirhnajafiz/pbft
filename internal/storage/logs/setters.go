package logs

import (
	"fmt"

	"github.com/f24-cse535/pbft/pkg/models"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"
)

// InitRequest places a new request at the end of datastore.
func (l *Logs) InitRequest() int {
	l.datastoreLock.Lock()
	for {
		index := l.index
		l.index++

		if l.GetRequest(index) == nil {
			l.datastoreLock.Unlock()

			return index
		}
	}
}

// ResetRequest sets a position free.
func (l *Logs) ResetRequest(index int) {
	l.datastoreLock.Lock()
	delete(l.datastore, index)
	l.datastoreLock.Unlock()
}

// SetRequest adds a new request into the datastore.
func (l *Logs) SetRequest(index int, req *pbft.RequestMsg, pp *pbft.PrePrepareMsg) {
	l.datastoreLock.Lock()
	l.datastore[index] = &models.Log{
		Request:    req,
		PrePrepare: pp,
	}

	if index > l.lastProcessingSeq {
		l.lastProcessingSeq = index
	}
	l.datastoreLock.Unlock()
}

// AppendLog adds a new log entry to the logs.
func (l *Logs) AppendLog(prefix, log string) {
	l.logs = append(l.logs, fmt.Sprintf("[%s] %s", prefix, log))
}

// AppendViewChange gets all view change messages and stores them.
func (l *Logs) AppendViewChange(view int, msg *pbft.ViewChangeMsg) {
	l.viewChangesLock.Lock()
	if _, ok := l.viewChanges[view]; !ok {
		l.viewChanges[view] = &models.ViewLog{
			ViewChangeMsgs: make([]*pbft.ViewChangeMsg, 0),
		}
	}

	for _, tmp := range l.viewChanges[view].ViewChangeMsgs {
		if tmp.GetNodeId() == msg.GetNodeId() {
			l.viewChangesLock.Unlock()
			return
		}
	}

	l.viewChanges[view].ViewChangeMsgs = append(l.viewChanges[view].ViewChangeMsgs, msg)
	l.viewChangesLock.Unlock()
}

// AppendNewView adds a new view message to a view log.
func (l *Logs) AppendNewView(view int, msg *pbft.NewViewMsg) {
	l.viewChangesLock.Lock()
	l.viewChanges[view].NewViewMsg = msg
	l.viewChangesLock.Unlock()
}

// Reset turns the values back to initial state.
func (l *Logs) Reset() {
	l.datastore = make(map[int]*models.Log)
	l.logs = make([]string, 0)
	l.viewChanges = make(map[int]*models.ViewLog)
	l.checkpoints = make(map[int][]*pbft.CheckpointMsg)
	l.index = 0
	l.lastProcessingSeq = 0
}

// SetRequestStatus accepts an index and status, and updates it if the new status is higher than what it is.
func (l *Logs) SetRequestStatus(index int, status pbft.RequestStatus) {
	l.datastoreLock.Lock()
	if l.datastore[index].Request.GetStatus().Number()+1 == status.Number() {
		l.datastore[index].Request.Status = status
	}
	l.datastoreLock.Unlock()
}

// SetRequestStatusForce accepts an index and status, and updates it.
func (l *Logs) SetRequestStatusForce(index int, status pbft.RequestStatus) {
	l.datastoreLock.Lock()
	l.datastore[index].Request.Status = status
	l.datastoreLock.Unlock()
}

// AppendCheckpoint adds a new checkpoint log.
func (l *Logs) AppendCheckpoint(key int, list []*pbft.CheckpointMsg) {
	l.checkpointsLock.Lock()
	l.checkpoints[key] = list
	l.checkpointsLock.Unlock()
}

// SetLastCheckpoint updates the value of last checkpoint.
func (l *Logs) SetLastCheckpoint(in int) {
	l.datastoreLock.Lock()
	l.lastCheckpoint = in
	l.datastoreLock.Unlock()
}
