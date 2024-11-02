package consensus

import (
	"github.com/f24-cse535/pbft/pkg/enums"
)

func (c *Consensus) handleCommit() {
	for {
		<-c.channels[enums.ChCommits]
	}
}

func (c *Consensus) handlePrePrepare() {
	for {
		<-c.channels[enums.ChPrePrepares]
	}
}

func (c *Consensus) handlePrepare() {
	for {
		<-c.channels[enums.ChPrepares]
	}
}

func (c *Consensus) handleRequest() {
	for {
		<-c.channels[enums.ChRequests]
	}
}

func (c *Consensus) handleReply() {
	for {
		<-c.channels[enums.ChReplys]
	}
}
