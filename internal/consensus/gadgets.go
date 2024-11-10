package consensus

import (
	"errors"
	"time"

	"github.com/f24-cse535/pbft/internal/consensus/modules"
	"github.com/f24-cse535/pbft/internal/utils/hashing"
	"github.com/f24-cse535/pbft/pkg/enum"
	"github.com/f24-cse535/pbft/pkg/models"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"

	"go.dedis.ch/kyber/v4/sign/tbls"
	"go.uber.org/zap"
)

// ackMsgReceivedGadget validates a preprepared or prepared ack message.
func (c *Consensus) ackMsgReceivedGadget(msg *pbft.AckMsg) *pbft.AckMsg {
	message := c.logs.GetRequest(int(msg.GetSequenceNumber()))
	if message == nil {
		return nil
	}

	digest := hashing.MD5HashRequestMsg(message) // get the digest of input request

	if !c.validateMsg(digest, msg.GetDigest(), msg.GetView()) {
		return nil
	}

	return msg
}

// msgProcessingGadget gets a preprepare message and sequence number to perform BFT protocol.
func (c *Consensus) msgProcessingGadget(sequence int, msg *pbft.PrePrepareMsg) {
	channel := make(chan *models.Packet, c.cfg.Total*2) // open a communication channel

	c.lock.Lock()
	c.requestsHandlersTable[sequence] = channel
	c.lock.Unlock()

	go c.communication.SendPreprepareMsg(msg) // send preprepare messages

	c.logs.SetRequestStatus(sequence, pbft.RequestStatus_REQUEST_STATUS_PP) // update our own status

	// byzantine leader stops after sending preprepare
	if c.memory.GetByzantine() {
		return
	}

	// wait for 2f+1 preprepared messages (count our own)
	if count := c.waiter.StartWaiting(channel, enum.PktPPed, c.ackMsgReceivedGadget); count+1 < c.cfg.Total {
		go c.communication.SendPrepareMsg(sequence, c.memory.GetView(), msg.GetDigest())

		c.logs.SetRequestStatus(sequence, pbft.RequestStatus_REQUEST_STATUS_P)

		c.waiter.StartWaiting(channel, enum.PktPed, c.ackMsgReceivedGadget)
	} else {
		c.logger.Info("optimized mode execution", zap.Int("count", count+1))
	}

	// broadcast to all using commit
	go c.communication.SendCommitMsg(sequence, c.memory.GetView(), msg.GetDigest())

	// delete our input channel as soon as possible
	c.lock.Lock()
	delete(c.requestsHandlersTable, sequence)
	c.lock.Unlock()

	// update our own status
	c.logs.SetRequestStatus(sequence, pbft.RequestStatus_REQUEST_STATUS_C)

	// send the sequence to the execute handler
	c.executionChannel <- sequence
}

// executionGadget gets a sequence number and performs the execution logic.
func (c *Consensus) executionGadget(sequence int) {
	if !c.canExecuteRequest(sequence) {
		return
	}

	// follow sequence until one is not committed, execute them
	index := sequence
	msg := c.logs.GetRequest(index)

	for {
		// don't reexecute a request
		if msg.GetStatus() != pbft.RequestStatus_REQUEST_STATUS_E {
			msg.GetResponse().Text = c.executeTransaction(msg.GetTransaction())
			c.logs.SetRequest(index, msg)
		}

		c.logs.SetRequestStatus(index, pbft.RequestStatus_REQUEST_STATUS_E) // update the request and set the status of prepare

		c.communication.SendReplyMsg(index, c.memory.GetView(), msg)

		index++
		if msg = c.logs.GetRequest(index); msg == nil || msg.GetStatus() != pbft.RequestStatus_REQUEST_STATUS_C {
			break
		}
	}

	// check if we can checkpoint on this node
	if target, ok, preprepares := c.canCheckpoint(); ok {
		checkpointMsg := pbft.CheckpointMsg{
			SequenceNumber:     int64(target),
			NodeId:             c.memory.GetNodeId(),
			PreprepareMessages: preprepares,
		}

		// send checkpoint message and notify our own handler
		c.communication.SendCheckpoint(&checkpointMsg)
		c.SignalToHandlers(models.NewPacket(&checkpointMsg, enum.PktCP, target))
	}
}

// enterViewChangeGadget stops the node for a new view change process.
func (c *Consensus) enterViewChangeGadget() {
	c.inViewChangeMode = true
	c.viewChangeGadgetChannel = make(chan *pbft.ViewChangeMsg)

	go func() {
		for {
			if err := c.viewChangeGadget(); err == nil {
				return
			}

			time.Sleep(100 * time.Millisecond)
		}
	}()
}

// viewChangeGadget starts a view change procedure.
func (c *Consensus) viewChangeGadget() error {
	// get the current view and increase it
	view := c.memory.GetView()
	view++

	// create a view change message
	message := pbft.ViewChangeMsg{
		NodeId:             c.memory.GetNodeId(),
		View:               int64(view),
		SequenceNumber:     int64(c.logs.GetLastCheckpoint()),
		CheckpointMessages: c.logs.GetLastCheckpointMsgs(),
		PreprepareMessages: c.logs.GetPrepreparesAfterCheckpoint(),
	}

	// sign the message by using its digest
	sig, err := tbls.Sign(c.suite, c.tss, []byte(hashing.MD5HashViewMsg(&message)))
	if err != nil {
		c.logger.Error("failed to sign the view message", zap.Error(err))
	}

	// set the signature message
	message.Signature = sig

	// append own view change msg
	c.logs.AppendViewChange(view, &message)

	// send a view change message
	if count := c.communication.SendViewChangeMsg(&message); count < c.cfg.Majority {
		c.logger.Info("not enough available servers to start view change", zap.Int("live servers", count))
		return nil
	}

	// create a new timer
	flag := true
	timer := modules.NewTimer(c.cfg.ViewChangeTimeout, time.Millisecond)
	timer.Start()

	// wait for 2f+1 messages
	for {
		if !flag {
			break
		}

		select {
		case msg := <-c.viewChangeGadgetChannel:
			if msg.GetView() != int64(view) {
				continue
			}

			c.logs.AppendViewChange(view, msg)

			if len(c.logs.GetViewChanges(view)) >= c.cfg.Majority {
				flag = false
				timer.Stop()
			}
		case <-timer.Notify(): // if the timer expired, return and reset the view timer
			return errors.New("view change timeout")
		}
	}

	// update the node view
	c.memory.SetView(view)

	// if the node is the leader, run a new leader gadget
	if c.getCurrentLeader() == c.memory.GetNodeId() {
		c.logger.Info("new leader elected", zap.String("id", c.memory.GetNodeId()))

		if !c.memory.GetByzantine() {
			c.newViewLeaderGadget(view)
		}
	} else { // if the node is backup, it needs new-view message
		if err := c.newViewBackupGadget(view); err != nil {
			return err
		}
	}

	// return from view change mode
	c.inViewChangeMode = false
	c.viewChangeGadgetChannel = nil

	return nil
}

// newViewBackupGadget waits for a new view message from new leader.
func (c *Consensus) newViewBackupGadget(view int) error {
	var (
		timer = modules.NewTimer(100*c.cfg.ViewChangeTimeout, time.Millisecond)
		msg   *pbft.NewViewMsg
	)

	// leader should response before timeout
	select {
	case raw := <-c.consensusHandlersTable[enum.PktNV]:
		msg = raw.Payload.(*pbft.NewViewMsg)
		timer.Stop()
	case <-timer.Notify(): // if the timer expired, return and reset the view timer
		return errors.New("leader didn't send new view")
	}

	c.logs.AppendNewView(view, msg) // set the message for view change

	return nil
}

// newLeaderGadget performs the procedure of new leader.
func (c *Consensus) newViewLeaderGadget(view int) {
	// get all view change messages from other nodes
	messages := c.logs.GetViewChanges(view)

	// create a log map to get requests
	logsMap := make(map[int]*pbft.PrePrepareMsg)

	// set the min and max
	minSequence := c.logs.GetLastCheckpoint()
	maxSequence := c.logs.GetLastCheckpoint()

	var (
		message   *pbft.ViewChangeMsg
		sigShares [][]byte
	)

	// loop in all view change messages
	for _, msg := range messages {
		sequence := int(msg.GetSequenceNumber())

		message = msg
		sigShares = append(sigShares, []byte(msg.GetSignature()))

		// loop over preprepares to insert them inside a logs map
		for _, pp := range msg.GetPreprepareMessages() {
			logsMap[int(pp.GetSequenceNumber())] = pp
		}

		// update the minimum and maximum bounds
		if sequence <= minSequence {
			minSequence = sequence
		}
		if sequence >= maxSequence {
			maxSequence = sequence
		}
	}

	// create an array to store sequences
	requests := make([]*pbft.PrePrepareMsg, 0)

	// collect all requets that are prepared
	for i := minSequence; i <= maxSequence; i++ {
		if item := c.logs.GetPreprepare(i); item != nil {
			requests = append(requests, item)
		} else if item, ok := logsMap[i]; ok {
			requests = append(requests, item)
		}
	}

	// get view digest
	digest := hashing.MD5HashViewMsg(message)

	// create a new view message
	newViewMsg := pbft.NewViewMsg{
		View:               int64(view),
		NodeId:             c.memory.GetNodeId(),
		ViewchangeMessage:  digest,
		PreprepareMessages: requests,
		Shares:             sigShares,
	}

	// save the entry
	c.logs.AppendNewView(view, &newViewMsg)

	// send new view
	c.communication.SendNewViewMsg(&newViewMsg)

	// start the protocol for every request
	go func() {
		for _, req := range requests {
			c.msgProcessingGadget(int(req.GetSequenceNumber()), req)
		}
	}()
}

// checkpointGadget gets a list of checkpoints and updates its logs.
func (c *Consensus) checkpointGadget(checkpoints map[int][]*pbft.CheckpointMsg) []int {
	// keys will be returned to the handler to delete used checkpoints.
	keys := make([]int, 0)

	// check if 2f+1 matching
	for key, value := range checkpoints {
		if len(value) >= c.cfg.Majority {
			c.logs.AppendCheckpoint(key, value)
			c.logs.SetLastCheckpoint(key)

			keys = append(keys, key)
		}
	}

	return keys
}
