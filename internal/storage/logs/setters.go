package logs

import "github.com/f24-cse535/pbft/pkg/rpc/pbft"

// InitLog places a new log at the end of logs.
func (l *Logs) InitLog(req *pbft.RequestMsg) int {
	index := l.index
	l.index++

	l.logs[index] = req

	return index
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
