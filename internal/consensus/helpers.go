package consensus

import (
	"time"

	"github.com/f24-cse535/pbft/internal/utils/hashing"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"

	"go.uber.org/zap"
)

// helpSendReply sends a reply message to the client by parsing the request message.
func (c *Consensus) helpSendReply(msg *pbft.RequestMsg) {
	c.Client.Reply(msg.GetClientId(), &pbft.ReplyMsg{
		SequenceNumber: msg.GetSequenceNumber(),
		View:           int64(c.Memory.GetView()),
		Timestamp:      msg.GetTransaction().GetTimestamp(),
		ClientId:       msg.GetClientId(),
		Response:       msg.GetResponse().GetText(),
	})
}

// helpSendRequest in a loop trys to send your request to a leader.
func (c *Consensus) helpSendRequest(msg *pbft.TransactionMsg) {
	for {
		// get the current leader id
		id := c.getCurrentLeader()

		c.Logger.Debug("current leader", zap.String("id", id))

		// send the transaction request to the leader
		if err := c.Client.Request(id, &pbft.RequestMsg{
			Transaction: msg,
			Response:    &pbft.TransactionRsp{},
		}); err != nil {
			// if the leader failed, increament the view, and wait 1 second before resending your request
			c.Logger.Debug("leader failed", zap.Error(err))
			c.Memory.IncView()

			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}
}

// helpReceiveResponse in a loop trys to get a response for user request.
func (c *Consensus) helpReceiveResponse(msg *pbft.TransactionMsg) *pbft.ReplyMsg {
	for {
		// wait for f+1 matching reply or timeout request
		resp := c.waitReplys(c.inTransactionChannel)
		if resp != nil {
			return resp
		}

		c.Logger.Debug("request timeout")

		// send message to all nodes
		c.Client.BroadcastRequest(msg)

		// sleep one second before resending the request
		time.Sleep(1 * time.Second)
	}
}

// helpProcessTransaction follows the pbft protocol.
func (c *Consensus) helpProcessTransaction(sequence int, message *pbft.RequestMsg) {
	defer func() {
		delete(c.channels, sequence) // reset our channel
	}()

	// update request metadata
	message.SequenceNumber = int64(sequence)
	message.Status = pbft.RequestStatus_REQUEST_STATUS_UNSPECIFIED

	// store it into datastore
	c.Logs.SetRequest(sequence, message)

	c.Logger.Debug(
		"new request received",
		zap.Int64("sequence number", message.GetSequenceNumber()),
		zap.Int64("timestamp", message.GetTransaction().GetTimestamp()),
	)

	// broadcast to all using preprepare
	go c.Client.BroadcastPrePrepare(&pbft.PrePrepareMsg{
		Request:        message,
		SequenceNumber: int64(sequence),
		View:           int64(c.Memory.GetView()),
		Digest:         hashing.MD5(message),
	})

	// wait for 2f+1 preprepared messages
	count := c.waitForPrePrepareds(c.channels[sequence])
	c.Logger.Debug("received preprepared messages", zap.Int("messages", count))

	// broadcast to all using prepare
	go c.Client.BroadcastPrepare(&pbft.PrepareMsg{
		SequenceNumber: int64(sequence),
		View:           int64(c.Memory.GetView()),
		Digest:         hashing.MD5(message),
	})

	// wait for 2f+1 prepared messages
	count = c.waitForPrepareds(c.channels[sequence])
	c.Logger.Debug("received prepared messages", zap.Int("messages", count))

	// broadcast to all using commit
	go c.Client.BroadcastCommit(&pbft.CommitMsg{
		SequenceNumber: int64(sequence),
		View:           int64(c.Memory.GetView()),
		Digest:         hashing.MD5(message),
	})
}
