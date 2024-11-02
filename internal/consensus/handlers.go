package consensus

func (c *Consensus) handleCommit(pkt interface{}) {
	// check the message
	// update log
	// do execution
}

func (c *Consensus) handlePrePrepare(pkt interface{}) {
	// check the message
	// update log
	// return with preprepared message
}

func (c *Consensus) handlePrePrepared(pkt interface{}) {

}

func (c *Consensus) handlePrepare(pkt interface{}) {
	// check the message
	// update log
	// return with accept message
}

func (c *Consensus) handlePrepared(pkt interface{}) {

}

func (c *Consensus) handleReply(pkt interface{}) {
	// update the memory
	// notify the transaction handler
}

func (c *Consensus) handleRequest(pkt interface{}) {
	// update the request meta-data
	// broadcast to all using preprepare
	// wait for 2f+1
	// broadcast to all using prepare
	// wait for 2f+1
	// broadcast to all using commit
	// execute message if possible
}

func (c *Consensus) handleTransaction(pkt interface{}) {
	// get the current leader
	// send request
	// wait for f+1 matching reply or timeout request (+ timer)
	// on the timeout, reset yourself
	// on the f+1 reply, send over channel
}
