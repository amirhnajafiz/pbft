package logs

import (
	"github.com/f24-cse535/pbft/pkg/models"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"
)

// GetRequestByTimestamp checks the requests for matching timestamp request.
func (l *Logs) GetRequestByTimestamp(ts int64) (int, *pbft.RequestMsg) {
	for key, value := range l.datastore {
		if value.Request.Transaction.GetTimestamp() == ts {
			return key, value.Request
		}
	}

	return 0, nil
}

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

// GetRequestsAfterCheckpoint returns the requests that are after the checkpoint.
func (l *Logs) GetRequestsAfterCheckpoint() map[int]*pbft.RequestMsg {
	list := make(map[int]*pbft.RequestMsg)

	for key, value := range l.datastore {
		if key >= l.lastCheckpoint {
			list[key] = value.Request
		}
	}

	return list
}

// GetPreprepares returns a list of preprepare messages from the last checkpoint.
func (l *Logs) GetPrepreparesAfterCheckpoint() []*pbft.PrePrepareMsg {
	list := make([]*pbft.PrePrepareMsg, 0)

	for key, value := range l.datastore {
		if key >= l.lastCheckpoint && value.Request.GetStatus() > pbft.RequestStatus_REQUEST_STATUS_PP {
			list = append(list, value.PrePrepare)
		}
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

// GetCheckpoints returns the checkpoints log.
func (l *Logs) GetCheckpoints() map[int][]*pbft.CheckpointMsg {
	return l.checkpoints
}

// GetLastCheckpoint returns the sequence number of the last checkpoint.
func (l *Logs) GetLastCheckpoint() int {
	return l.lastCheckpoint
}

// GetLastCheckpoint messages returns a list of checkpoints certificates.
func (l *Logs) GetLastCheckpointMsgs() []*pbft.CheckpointMsg {
	return l.checkpoints[l.lastCheckpoint]
}

// GetHighWaterMark returns the limit of the nodes water mark.
func (l *Logs) GetHighWaterMark() int {
	return l.lastCheckpoint + l.kvalue
}

// GetLastProcessingSeq returns the last processing sequence number.
func (l *Logs) GetLastProcessingSeq() int {
	return l.lastProcessingSeq
}
