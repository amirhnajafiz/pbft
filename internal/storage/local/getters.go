package local

// GetStatus is used to check the node availability.
func (m *Memory) GetStatus() bool {
	return m.status
}

// GetByzantine is used to check the byzantine behaviour of the node.
func (m *Memory) GetByzantine() bool {
	return m.byzantine
}

// GetNodeId returns the node id.
func (m *Memory) GetNodeId() string {
	return m.nodeId
}

// GetView returns the current view of this node.
func (m *Memory) GetView() int {
	return m.view
}

// GetNodeByIndex returns the node name by index.
func (m *Memory) GetNodeByIndex(index int) string {
	return m.nodes[index]
}

// GetBalance returns the balance value of a client.
func (m *Memory) GetBalance(key string) int {
	return m.balances[key]
}

// GetTimestamp returns the current timestamp.
func (m *Memory) GetTimestamp() int64 {
	return m.currentTimestamp
}
