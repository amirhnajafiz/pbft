package local

// Memory is a local storage that is used for each node. It keeps the state of each node.
type Memory struct {
	status     bool   // status is for node's availability
	byzantine  bool   // byzantine is for node's behavior
	nodeId     string // nodes name
	view       int    // systems view
	totalNodes int    // number of total-nodes

	balances map[string]int // balances is the holder for clients and their balance value
	nodes    map[int]string // nodes is map used for tracking leaders
}

// NewMemory returns an instance of the memory struct.
func NewMemory(nodeId string, totalNodes int) *Memory {
	return &Memory{
		status:     true,  // the init status of node is true
		byzantine:  false, // the init behavior node is non-byzantine
		nodeId:     nodeId,
		view:       0,
		totalNodes: totalNodes,
	}
}
