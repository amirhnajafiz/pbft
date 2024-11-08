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
func (l *Logs) GetViewChanges(view int) []*pbft.ViewChangeMsg {
	if list, ok := l.viewChanges[view]; ok {
		return list
	}

	return make([]*pbft.ViewChangeMsg, 0)
}

// GetLogs returns the node datalog.
func (l *Logs) GetLogs() []string {
	return l.logs
}

// GetSequenceNumber returns the minimun executed sequence number.
func (l *Logs) GetSequenceNumber() int {
	index := 0

	for {
		if req := l.GetRequest(index); req == nil {
			break
		} else {
			if req.GetStatus() == pbft.RequestStatus_REQUEST_STATUS_E {
				index++
			} else {
				break
			}
		}
	}

	return index
}
