package logs

import (
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"
)

// InitLog places a new log at the end of logs.
func (l *Logs) InitLog() int {
	for {
		index := l.index
		l.index++
		if l.GetLog(index) == nil {
			return index
		}
	}
}

// SetLog adds a new log into the logs.
func (l *Logs) SetLog(index int, req *pbft.RequestMsg) {
	l.logs[index] = req
}

// Reset turns the values back to initial state.
func (l *Logs) Reset() {
	l.datastore = make(map[int]*pbft.RequestMsg)
	l.logs = make(map[int]*pbft.RequestMsg)
	l.index = 0
}

// SetLogStatus accepts an index and status, and updates it if the new status is higher than what it is.
func (l *Logs) SetLogStatus(index int, status pbft.RequestStatus) {
	if l.logs[index].GetStatus() > status {
		l.logs[index].Status = status
	}
}
