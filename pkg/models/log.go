package models

import "github.com/f24-cse535/pbft/pkg/rpc/pbft"

// Log is a wrapper for request and prepare messages.
type Log struct {
	Request    *pbft.RequestMsg
	PrePrepare *pbft.PrePrepareMsg
}
