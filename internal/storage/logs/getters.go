package logs

import "github.com/f24-cse535/pbft/pkg/rpc/pbft"

// GetRequest is returns a request by its index.
func (l *Logs) GetRequest(index int) *pbft.RequestMsg {
	if value, ok := l.datastore[index]; ok {
		return value
	}

	return nil
}

// GetAllRequests returns an array of the node requests.
func (l *Logs) GetAllRequests() []*pbft.RequestMsg {
	list := make([]*pbft.RequestMsg, len(l.datastore))

	for key, value := range l.datastore {
		list[key] = value
	}

	return list
}

// GetViewChanges returns a list of stored view changes.
func (l *Logs) GetViewChanges(view int) []interface{} {
	if list, ok := l.viewChanges[view]; ok {
		return list
	}

	return nil
}

// GetLogs returns the node datalog.
func (l *Logs) GetLogs() []string {
	return l.logs
}
