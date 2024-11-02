package consensus

import (
	"github.com/f24-cse535/pbft/pkg/enums"
)

func (c *Consensus) handleCommit() {
	// check the message
	// update log
	// do execution
	for {
		<-c.channels[enums.ChCommits]
	}
}

func (c *Consensus) handlePrePrepare() {
	// check the message
	// update log
	// return with preprepared message
	for {
		<-c.channels[enums.ChPrePrepares]
	}
}

func (c *Consensus) handlePrepare() {
	// check the message
	// update log
	// return with accept message
	for {
		<-c.channels[enums.ChPrepares]
	}
}

func (c *Consensus) handleRequest() {
	// update the request meta-data
	// broadcast to all using preprepare
	// wait for 2f+1
	// broadcast to all using prepare
	// wait for 2f+1
	// broadcast to all using commit
	// execute message if possible
	for {
		select {
		case <-c.channels[enums.ChRequests]:
			return
		case <-c.channels[enums.ChPrePrepareds]:
			return
		case <-c.channels[enums.ChPrepareds]:
			return
		}
	}
}

func (c *Consensus) handleReply() {
	// update the memory
	// notify the transaction handler
	for {
		<-c.channels[enums.ChReplys]
	}
}

func (c *Consensus) handleTransaction(_ chan interface{}) {
	// get the current leader
	// send request
	// wait for f+1 matching reply or timeout request (+ timer)
	// on the timeout, reset yourself
	// on the f+1 reply, send over channel
	defer func() {
		// reset the channel
		c.channels[enums.ChTransactions] = nil
	}()

	for {
		<-c.channels[enums.ChTransactions]
	}
}
