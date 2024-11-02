package enums

type ChannelType int

// list of the channel types
const (
	ChCommits ChannelType = iota + 1
	ChPrePrepares
	ChPrepares
	ChRequests
	ChReplys
	ChTransactions
)

// ListChannelTypes returns a list of available channels, it does not include Transactions channel.
func ListChannelTypes() []ChannelType {
	return []ChannelType{ChCommits, ChPrePrepares, ChPrepares, ChRequests, ChReplys}
}
