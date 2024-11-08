package logs

import (
	"fmt"

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
	l.datastore[index] = req
}

// AppendLog adds a new log entry to the logs.
func (l *Logs) AppendLog(prefix, log string) {
	l.logs = append(l.logs, fmt.Sprintf("[%s] %s", prefix, log))
}

// AppendViewChange gets all view change messages and stores them.
func (l *Logs) AppendViewChange(view int, msg interface{}) {
	if _, ok := l.viewChanges[view]; !ok {
		l.viewChanges[view] = make([]interface{}, 0)
	}

	l.viewChanges[view] = append(l.viewChanges[view], msg)
}

// Reset turns the values back to initial state.
func (l *Logs) Reset() {
	l.datastore = make(map[int]*pbft.RequestMsg)
	l.logs = make([]string, 0)
	l.index = 0
}

// SetRequestStatus accepts an index and status, and updates it if the new status is higher than what it is.
func (l *Logs) SetRequestStatus(index int, status pbft.RequestStatus) {
	if l.datastore[index].GetStatus().Number()+1 == status.Number() {
		l.datastore[index].Status = status
	}
}
