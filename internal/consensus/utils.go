package consensus

// GetCurrentLeader returns the current leader id.
func (c *Consensus) GetCurrentLeader() string {
	return c.Memory.GetNodeByIndex(c.Memory.GetView())
}
