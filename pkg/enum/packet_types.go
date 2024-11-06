package enum

// PacketType is a new type for identifying packets.
type PacketType int

// List of the packet types.
const (
	PktTrx PacketType = iota + 1
	PktReq
	PktPP
	PktPPed
	PktP
	PktPed
	PktCmt
	PktRpl
)
