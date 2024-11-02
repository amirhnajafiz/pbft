package enums

type ChannelType int

const (
	ChCommits ChannelType = iota + 1
	ChPrePrepares
	ChPrepares
	ChRequests
	ChReplys
)

func ListChannelTypes() []ChannelType {
	return []ChannelType{ChCommits, ChPrePrepares, ChPrepares, ChRequests, ChReplys}
}
