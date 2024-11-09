package models

import "github.com/f24-cse535/pbft/pkg/rpc/pbft"

// Log is a wrapper for request and prepare messages.
type Log struct {
	Request    *pbft.RequestMsg
	PrePrepare *pbft.PrePrepareMsg
}

// ViewLog is used as a wrapper for view messages.
type ViewLog struct {
	ViewChangeMsgs []*pbft.ViewChangeMsg
	NewViewMsg     *pbft.NewViewMsg
}
