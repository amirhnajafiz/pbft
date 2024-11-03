package logs

import (
	"fmt"

	"github.com/f24-cse535/pbft/pkg/rpc/pbft"
)

// InitRequest places a new request at the end of datastore.
func (l *Logs) InitRequest() int {
	for {
		index := l.index
		l.index++
		if l.GetRequest(index) == nil {
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

// Reset turns the values back to initial state.
func (l *Logs) Reset() {
	l.datastore = make(map[int]*pbft.RequestMsg)
	l.logs = make([]string, 0)
	l.index = 0
}

// SetRequestStatus accepts an index and status, and updates it if the new status is higher than what it is.
func (l *Logs) SetRequestStatus(index int, status pbft.RequestStatus) {
	if l.datastore[index].GetStatus().Number() < status.Number() {
		l.datastore[index].Status = status
	}
}
