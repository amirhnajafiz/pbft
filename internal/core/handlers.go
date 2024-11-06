package core

import "github.com/f24-cse535/pbft/pkg/rpc/pbft"

// requestHandler gets a transaction from client, and follows PBFT procedure to send it.
func (c *Core) requestHandler(client string, trx *pbft.TransactionMsg) {
	// try to send the request to leader
	id := c.getCurrentLeader()
	c.cli.Request(id, &pbft.RequestMsg{
		Transaction: trx,
		Response:    &pbft.TransactionRsp{},
	})

	// wait for f+1 reply messages
	target := -1
	counts := make(map[int]int)
	replys := make(map[int]*pbft.ReplyMsg)

	for {
		// get reply message
		reply := <-c.handlers[client]

		// only get the replys with matching timestamp
		if reply.GetTimestamp() != trx.GetTimestamp() {
			continue
		}

		key := int(reply.SequenceNumber + reply.Timestamp)
		if _, ok := counts[key]; !ok {
			replys[key] = reply
			counts[key] = 0
		}

		counts[key]++

		for key, value := range counts {
			if c.replys >= value {
				target = key
			}
		}

		if target != -1 {
			break
		}
	}

	// update our view
	c.memory.SetView(int(replys[target].GetView()))

	// return the user response
	c.Done(client, replys[target].GetResponse())
}
