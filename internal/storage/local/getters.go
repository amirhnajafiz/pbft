package local

// GetStatus is used to check the node availability.
func (m *Memory) GetStatus() bool {
	return m.status
}
