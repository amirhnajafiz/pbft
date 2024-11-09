package logs

import (
	"math"

	"github.com/f24-cse535/pbft/pkg/models"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"
)

// GetRequest is returns a request by its index.
func (l *Logs) GetRequest(index int) *pbft.RequestMsg {
	if value, ok := l.datastore[index]; ok {
		return value.Request
	}

	return nil
}

// GetPreprepare returns the preprepare message of a request.
func (l *Logs) GetPreprepare(index int) *pbft.PrePrepareMsg {
	if value, ok := l.datastore[index]; ok {
		return value.PrePrepare
	}

	return nil
}

// GetAllRequests returns an array of the node requests.
func (l *Logs) GetAllRequests() map[int]*pbft.RequestMsg {
	list := make(map[int]*pbft.RequestMsg)

	for key, value := range l.datastore {
		list[key] = value.Request
	}

	return list
}

// GetViewChanges returns a list of stored view changes.
func (l *Logs) GetViewChanges(view int) []*pbft.ViewChangeMsg {
	if list, ok := l.viewChanges[view]; ok {
		return list.ViewChangeMsgs
	}

	return make([]*pbft.ViewChangeMsg, 0)
}

// GetAllViewChanges returns a map of views and their view change messages.
func (l *Logs) GetAllViewChanges() map[int]*models.ViewLog {
	return l.viewChanges
}

// GetLogs returns the node datalog.
func (l *Logs) GetLogs() []string {
	return l.logs
}

// GetSequenceNumber returns the minimun executed sequence number.
func (l *Logs) GetSequenceNumber(index int) int {
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

// GetPreprepares returns a list of preprepare messages from the given sequence.
func (l *Logs) GetPreprepares(from int, index int) []*pbft.PrePrepareMsg {
	list := make([]*pbft.PrePrepareMsg, 0)

	limit := int(math.Max(float64(from), float64(l.index)))

	for key, value := range l.datastore {
		if key >= limit && value.Request.GetStatus() > pbft.RequestStatus_REQUEST_STATUS_PP {
			list = append(list, value.PrePrepare)
		}
	}

	return list
}

// GetCheckpoints returns the checkpoints log.
func (l *Logs) GetCheckpoints() map[int][]*pbft.CheckpointMsg {
	return l.checkpoints
}
