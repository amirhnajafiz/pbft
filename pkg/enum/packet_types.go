package enum

// PacketType is a new type for identifying packets.
type PacketType int

// List of the packet types.
const (
	PktReq PacketType = iota + 1
	PktPP
	PktPPed
	PktP
	PktPed
	PktCmt
	PktVC
)
