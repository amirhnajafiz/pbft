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

// IncView increases the view one unit.
func (m *Memory) IncView() {
	m.view = (m.view + 1) % m.totalNodes
}

// SetView updates the view value.
func (m *Memory) SetView(in int) {
	m.view = in
}

// SetBalances is used to set clients balances.
func (m *Memory) SetBalances(balances map[string]int) {
	m.balances = make(map[string]int)

	for key, value := range balances {
		m.balances[key] = value
	}
}

// SetNodes is used to set nodes and their index.
func (m *Memory) SetNodes(nodes map[string]int) {
	m.nodes = make(map[int]string)

	for key, value := range nodes {
		m.nodes[value] = key
	}
}

// SetBalance updates a client balance.
func (m *Memory) SetBalance(key string, value int) {
	m.balances[key] = value
}

// SetTimestamp sets the current request timestamp.
func (m *Memory) SetTimestamp(in int64) {
	m.currentTimestamp = in
}
