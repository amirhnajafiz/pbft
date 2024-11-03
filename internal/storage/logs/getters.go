package logs

import "github.com/f24-cse535/pbft/pkg/rpc/pbft"

// GetLog is returns a log by its index.
func (l *Logs) GetLog(index int) *pbft.RequestMsg {
	if value, ok := l.logs[index]; ok {
		return value
	}

	return nil
}

func (l *Logs) GetAllLogs() []*pbft.RequestMsg {
	return nil
}
