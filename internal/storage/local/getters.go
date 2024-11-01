package local

// GetStatus is used to check the node availability.
func (m *Memory) GetStatus() bool {
	return m.status
}

// GetByzantine is used to check the byzantine behaviour of the node.
func (m *Memory) GetByzantine() bool {
	return m.byzantine
}
