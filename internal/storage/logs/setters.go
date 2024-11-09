package logs

import (
	"fmt"

	"github.com/f24-cse535/pbft/pkg/models"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"
)

// InitRequest places a new request at the end of datastore.
func (l *Logs) InitRequest() int {
	l.lock.Lock()
	for {
		index := l.index
		l.index++

		if l.GetRequest(index) == nil {
			l.lock.Unlock()

			return index
		}
	}
}

// SetRequest adds a new request into the datastore.
func (l *Logs) SetRequest(index int, req *pbft.RequestMsg) {
	l.datastore[index] = &models.Log{
		Request: req,
	}
}

// SetPreprepare sets the preprepare message of a request.
func (l *Logs) SetPreprepare(index int, pp *pbft.PrePrepareMsg) {
	l.datastore[index].PrePrepare = pp
}

// AppendLog adds a new log entry to the logs.
func (l *Logs) AppendLog(prefix, log string) {
	l.logs = append(l.logs, fmt.Sprintf("[%s] %s", prefix, log))
}

// AppendViewChange gets all view change messages and stores them.
func (l *Logs) AppendViewChange(view int, msg *pbft.ViewChangeMsg) {
	if _, ok := l.viewChanges[view]; !ok {
		l.viewChanges[view] = &models.ViewLog{
			ViewChangeMsgs: make([]*pbft.ViewChangeMsg, 0),
		}
	}

	l.viewChanges[view].ViewChangeMsgs = append(l.viewChanges[view].ViewChangeMsgs, msg)
}

// AppendNewView adds a new view message to a view log.
func (l *Logs) AppendNewView(view int, msg *pbft.NewViewMsg) {
	l.viewChanges[view].NewViewMsg = msg
}

// Reset turns the values back to initial state.
func (l *Logs) Reset() {
	l.datastore = make(map[int]*models.Log)
	l.logs = make([]string, 0)
	l.index = 0
}

// SetRequestStatus accepts an index and status, and updates it if the new status is higher than what it is.
func (l *Logs) SetRequestStatus(index int, status pbft.RequestStatus) {
	if l.datastore[index].Request.GetStatus().Number()+1 == status.Number() {
		l.datastore[index].Request.Status = status
	}
}
