package enum

// Interrupt is a new type for separating messages between gRPC level and consensus.
type Interrupt int

// List of the system's interrutps.
const (
	IntrTransaction Interrupt = iota + 1
	IntrRequest
	IntrPrePrepare
	IntrPrePrepared
	IntrPrepare
	IntrPrepared
	IntrCommit
	IntrReply
)
