package local

// SetStatus gets a bool (in) to update status.
func (m *Memory) SetStatus(in bool) {
	m.status = in
}

// SetByzantine sets the byzantine behaviour.
func (m *Memory) SetByzantine(in bool) {
	m.byzantine = in
}

// Reset turns the values back to initial state.
func (m *Memory) Reset() {
	m.status = true
	m.byzantine = false
}
