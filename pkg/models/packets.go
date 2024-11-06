package models

import "github.com/f24-cse535/pbft/pkg/enum"

// Packet is a wrapper for gRPC messages to consensus and core modules.
type Packet struct {
	Type    enum.PacketType
	Headers map[string]interface{}
	Payload interface{}
}

// AddHeader adds a new key-value pair to the headers.
func (p Packet) AddHeader(key string, value interface{}) *Packet {
	p.Headers[key] = value
	return &p
}

// NewPacket gets a payload and a type, then it returns a packet instance.
func NewPacket(payload interface{}, pktType enum.PacketType) *Packet {
	return &Packet{
		Headers: make(map[string]interface{}),
		Type:    pktType,
		Payload: payload,
	}
}
