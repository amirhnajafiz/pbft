package application

// getCurrentLeader returns the current leader id.
func (a *App) getCurrentLeader() string {
	return a.memory.GetNodeByIndex(a.memory.GetView())
}
