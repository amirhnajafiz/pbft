package models

import "github.com/f24-cse535/pbft/pkg/enum"

// InterruptMsg is a wrapper for gRPC messages to consensus.
type InterruptMsg struct {
	Type    enum.Interrupt
	Payload interface{}
}
