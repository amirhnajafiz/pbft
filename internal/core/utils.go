package core

// getCurrentLeader returns the current leader id.
func (c *Core) getCurrentLeader() string {
	return c.memory.GetNodeByIndex(c.memory.GetView())
}
