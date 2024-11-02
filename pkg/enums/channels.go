package enums

type ChannelType int

// list of the channel types
const (
	ChCommits ChannelType = iota + 1
	ChPrePrepares
	ChPrepares
	ChRequests
	ChReplys
	ChPrePrepareds
	ChPrepareds
	ChTransactions
)

// ListNodeChannels returns a list of available channels for a node, it does not include Transactions channel.
func ListNodeChannels() []ChannelType {
	return []ChannelType{ChCommits, ChPrePrepares, ChPrepares, ChRequests, ChPrePrepareds, ChPrepareds}
}

// ListClientChannels returns a list of available channels for a client, it does not include Transactions channel.
func ListClientChannels() []ChannelType {
	return []ChannelType{ChReplys}
}
