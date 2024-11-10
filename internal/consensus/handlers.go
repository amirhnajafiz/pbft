package consensus

import (
	"github.com/f24-cse535/pbft/internal/utils/hashing"
	"github.com/f24-cse535/pbft/pkg/enum"
	"github.com/f24-cse535/pbft/pkg/models"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"

	"go.uber.org/zap"
)

// preprepareHandler gets gRPC packets of type PP and handles them.
func (c *Consensus) preprepareHandler() {
	for {
		// get raw PP packets and cast it
		raw := <-c.consensusHandlersTable[enum.PktPP]
		msg := raw.Payload.(*pbft.PrePrepareMsg)

		// don't accept messages in view change mode
		if c.inViewChangeMode {
			continue
		}

		// start the timer
		c.viewTimer.AccumaccumulativeStart()

		// if the input request is not in watermark, don't accept it
		if raw.Sequence < c.logs.GetLastCheckpoint() || raw.Sequence > c.logs.GetHighWaterMark() {
			continue
		}

		digest := hashing.MD5HashRequestMsg(msg.GetRequest()) // get the digest of request
		if !c.validateMsg(digest, msg.GetDigest(), msg.GetView()) {
			continue
		}

		// update the request and set the status of preprepared
		c.logs.SetRequest(raw.Sequence, msg.GetRequest())
		c.logs.SetRequestStatus(raw.Sequence, pbft.RequestStatus_REQUEST_STATUS_PP)
		c.logs.SetPreprepare(raw.Sequence, msg)

		// call preprepared RPC to notify the sender
		c.communication.Client().PrePrepared(msg.GetNodeId(), &pbft.AckMsg{
			View:           int64(c.memory.GetView()),
			SequenceNumber: int64(raw.Sequence),
			Digest:         digest,
		})
	}
}

// prepareHandler gets gRPC packets of type P and handles them.
func (c *Consensus) prepareHandler() {
	for {
		// get raw P packets and cast them
		raw := <-c.consensusHandlersTable[enum.PktP]
		msg := raw.Payload.(*pbft.AckMsg)

		// don't accept messages in view change mode
		if c.inViewChangeMode {
			continue
		}

		// start the timer
		c.viewTimer.Start()

		// get the message from our datastore
		message := c.logs.GetRequest(raw.Sequence)
		if message == nil {
			continue
		}

		digest := hashing.MD5HashRequestMsg(message) // get the digest of input request

		if !c.memory.GetByzantine() { // byzantine nodes don't prepare messages
			if !c.validateMsg(digest, msg.GetDigest(), msg.GetView()) {
				continue
			}

			// update the request and set the status of prepare
			c.logs.SetRequestStatus(raw.Sequence, pbft.RequestStatus_REQUEST_STATUS_P)

			// call prepared RPC to notify the sender
			c.communication.Client().Prepared(msg.GetNodeId(), &pbft.AckMsg{
				View:           int64(c.memory.GetView()),
				SequenceNumber: int64(raw.Sequence),
				Digest:         digest,
			})
		}
	}
}

// commitHandler gets gRPC packets of type C and handles them.
func (c *Consensus) commitHandler() {
	for {
		// get raw C packets
		raw := <-c.consensusHandlersTable[enum.PktCmt]

		// don't accept messages in view change mode
		if c.inViewChangeMode {
			continue
		}

		// start the timer
		c.viewTimer.Start()

		// update the request and set the status of prepare
		c.logs.SetRequestStatus(raw.Sequence, pbft.RequestStatus_REQUEST_STATUS_C)

		// stop the timer
		c.viewTimer.Dismiss()

		// send the sequence to the execute handler
		c.executionChannel <- raw.Sequence
	}
}

// executeHandler gets sequence numbers from handlers and executes a request.
func (c *Consensus) executeHandler() {
	for {
		seq := <-c.executionChannel // get the sequence number
		c.executionGadget(seq)      // execute the request
	}
}

// requestHandler gets a request message and performs the request handling logic.
func (c *Consensus) requestHandler(pkt *models.Packet) {
	// parse the input message
	msg := pkt.Payload.(*pbft.RequestMsg)

	// don't accept messages in view change mode
	if c.inViewChangeMode {
		return
	}

	// check if we have a request with the given timestamp
	if seq, req := c.logs.GetRequestByTimestamp(msg.GetTransaction().GetTimestamp()); req != nil {
		if req.Status == pbft.RequestStatus_REQUEST_STATUS_E {
			c.communication.SendReplyMsg(seq, c.memory.GetView(), req)
		}

		return
	}

	// check if the node is leader
	if c.getCurrentLeader() != c.memory.GetNodeId() {
		// send the request to leader
		c.communication.Client().Request(c.getCurrentLeader(), msg)

		// start the timer
		c.viewTimer.Start()

		return
	}

	// find a sequence number for this request
	sequence := c.logs.InitRequest()

	// update request metadata
	msg.Status = pbft.RequestStatus_REQUEST_STATUS_UNSPECIFIED

	// create a preprepare message
	ppMessage := pbft.PrePrepareMsg{
		Request:        msg,
		View:           int64(c.memory.GetView()),
		NodeId:         c.memory.GetNodeId(),
		SequenceNumber: int64(sequence),
		Digest:         hashing.MD5HashRequestMsg(msg),
	}

	// store messages into datastore
	c.logs.SetRequest(sequence, msg)
	c.logs.SetPreprepare(sequence, &ppMessage)

	c.logger.Debug("new request got into the system", zap.Int("sequece", sequence), zap.Int64("time", msg.Transaction.GetTimestamp()))

	// run a processing gadget
	c.msgProcessingGadget(sequence, &ppMessage)
}

// timerHandler creates a new timer and monitors the timer.
func (c *Consensus) timerHandler() {
	// stop the timer so that others can start it
	c.viewTimer.Stop()

	for {
		<-c.viewTimer.Notify()
		c.logger.Debug("timer expired")

		// start view change
		if !c.inViewChangeMode {
			c.enterViewChangeGadget()
			c.viewTimer.Stop()
		}
	}
}

// viewChangeHandler waits for f+1 view changes, then it will start view change procedure.
func (c *Consensus) viewChangeHandler() {
	for {
		// capture view change messages
		raw := <-c.consensusHandlersTable[enum.PktVC]
		msg := raw.Payload.(*pbft.ViewChangeMsg)

		// check the destination
		if !c.inViewChangeMode {
			c.logs.AppendViewChange(int(msg.GetView()), msg)

			if msg.GetView() == int64(c.memory.GetView()) {
				continue
			}

			if len(c.logs.GetViewChanges(int(msg.GetView()))) >= c.cfg.Responses {
				c.enterViewChangeGadget()
			}
		} else {
			if c.viewChangeGadgetChannel != nil {
				c.viewChangeGadgetChannel <- msg
			}
		}
	}
}

// checkpointHandler captures checkpoint messages.
func (c *Consensus) checkpointHandler() {
	// a temporary list of checkpoints.
	checkpoints := make(map[int][]*pbft.CheckpointMsg)

	for {
		// capture raw checkpoint messages
		raw := <-c.consensusHandlersTable[enum.PktCP]
		msg := raw.Payload.(*pbft.CheckpointMsg)

		// checkpoint message is old drop it
		if raw.Sequence < c.logs.GetLastCheckpoint() {
			continue
		}

		if _, ok := checkpoints[int(msg.GetSequenceNumber())]; !ok {
			checkpoints[int(msg.GetSequenceNumber())] = make([]*pbft.CheckpointMsg, 0)
		}
		checkpoints[raw.Sequence] = append(checkpoints[raw.Sequence], msg)

		// call checkpoint gadget
		for _, key := range c.checkpointGadget(checkpoints) {
			delete(checkpoints, key)
		}
	}
}
