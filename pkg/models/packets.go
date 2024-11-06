package models

import "github.com/f24-cse535/pbft/pkg/enum"

// Packet is a wrapper for gRPC messages to consensus and core modules.
type Packet struct {
	Type     enum.PacketType
	Sequence int
	Payload  interface{}
}

// NewPacket gets a payload, sequence number, and a type; then it returns a packet instance.
func NewPacket(payload interface{}, pktType enum.PacketType, sequence int) *Packet {
	return &Packet{
		Sequence: sequence,
		Type:     pktType,
		Payload:  payload,
	}
}
