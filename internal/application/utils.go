package application

import "github.com/f24-cse535/pbft/pkg/rpc/pbft"

// getCurrentLeader returns the current leader id.
func (a *App) getCurrentLeader() string {
	return a.memory.GetNodeByIndex(a.memory.GetView())
}

// broadcastRequest sends a request to all server nodes.
func (a *App) broadcastRequest(req *pbft.RequestMsg) int {
	count := 0
	for key := range a.cli.GetSystemNodes() {
		if er := a.cli.Request(key, req); er == nil {
			count++
		}
	}

	return count
}

// emptyChannel gets all input messages from a channel and deletes it.
func (a *App) emptyChannel(client string) {
	ch := a.handlers[client]

	for {
		select {
		case <-ch:
			continue
		default:
			delete(a.handlers, client)
			return
		}
	}
}
