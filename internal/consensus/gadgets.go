package consensus

import (
	"errors"
	"time"

	"github.com/f24-cse535/pbft/internal/consensus/modules"
	"github.com/f24-cse535/pbft/internal/utils/hashing"
	"github.com/f24-cse535/pbft/pkg/enum"
	"github.com/f24-cse535/pbft/pkg/models"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"

	"go.dedis.ch/kyber/v4/sign/bls"
	"go.dedis.ch/kyber/v4/sign/tbls"
	"go.uber.org/zap"
)

// newAckGadget validates a preprepared or prepared ack message.
func (c *Consensus) newAckGadget(msg *pbft.AckMsg) *pbft.AckMsg {
	// get the message from our datastore
	message := c.logs.GetRequest(int(msg.GetSequenceNumber()))
	if message == nil {
		return nil
	}

	// get the digest of input request
	digest := hashing.MD5Req(message)

	// validate the message
	if !c.validateMsg(digest, msg.GetDigest(), msg.GetView()) {
		return nil
	}

	return msg
}

func (c *Consensus) newProcessingGadet(sequence int, msg *pbft.PrePrepareMsg) {
	// open a communication channel
	channel := make(chan *models.Packet, c.cfg.Total*2)
	c.lock.Lock()
	c.requestsHandlersTable[sequence] = channel
	c.lock.Unlock()

	// send preprepare messages
	go c.communication.SendPreprepareMsg(msg, c.memory.GetView())

	// update our own status
	c.logs.SetRequestStatus(sequence, pbft.RequestStatus_REQUEST_STATUS_PP)

	// wait for 2f+1 preprepared messages (count our own)
	count := c.waiter.NewPrePreparedWaiter(channel, c.newAckGadget)

	// optimized mode
	if count+1 < c.cfg.Total {
		// broadcast to all using prepare
		go c.communication.SendPrepareMsg(msg.GetRequest(), c.memory.GetView())

		// update our own status
		if !c.memory.GetByzantine() {
			c.logs.SetRequestStatus(sequence, pbft.RequestStatus_REQUEST_STATUS_P)
		}

		// wait for 2f+1 prepared messages (count our own)
		c.waiter.NewPreparedWaiter(channel, c.newAckGadget)
	} else {
		c.logger.Info("optimized mode", zap.Int("count", count+1))
	}

	// broadcast to all using commit, make sure everyone get's it
	go c.communication.SendCommitMsg(msg.GetRequest(), c.memory.GetView())

	// delete our input channel as soon as possible
	c.lock.Lock()
	delete(c.requestsHandlersTable, sequence)
	c.lock.Unlock()

	// update our own status
	c.logs.SetRequestStatus(sequence, pbft.RequestStatus_REQUEST_STATUS_C)

	// send the sequence to the execute handler
	c.executionChannel <- sequence
}

// newExecutionGadget gets a sequence number and performs the execution logic.
func (c *Consensus) newExecutionGadget(sequence int) {
	if !c.memory.GetByzantine() && !c.canExecuteRequest(sequence) {
		c.logger.Debug("cannot execute this request yet", zap.Int("sequence", sequence))
		return
	}

	// follow sequence until one is not committed, execute them
	index := sequence
	msg := c.logs.GetRequest(index)

	// don't reexecute a request
	if msg.GetStatus() == pbft.RequestStatus_REQUEST_STATUS_E {
		return
	}

	for {
		c.executeRequest(msg) // execute request

		if !c.memory.GetByzantine() {
			c.logs.SetRequestStatus(index, pbft.RequestStatus_REQUEST_STATUS_E) // update the request and set the status of prepare
		}

		c.communication.SendReplyMsg(msg, c.memory.GetView()) // send the reply message using helper functions
		c.logger.Debug("request executed", zap.Int("sequence number", index))

		index++
		if msg = c.logs.GetRequest(index); msg == nil || msg.GetStatus() != pbft.RequestStatus_REQUEST_STATUS_C {
			break
		}
	}

	// check the number of executions
	if c.shouldCheckpoint() {
		c.communication.SendCheckpoint(&pbft.CheckpointMsg{
			SequenceNumber: int64(c.memory.GetLowWaterMark()),
		})
	}
}

// newViewChangeModeGadget stops the node for view change.
func (c *Consensus) newViewChangeModeGadget() {
	c.inViewChangeMode = true
	c.viewChangeGadgetChannel = make(chan *pbft.ViewChangeMsg)

	go func() {
		for {
			if err := c.newViewChangeGadget(); err == nil {
				return
			}

			time.Sleep(10 * time.Millisecond)
		}
	}()
}

// newViewChangeGadget gets a count number of view-change messages and starts view change procedure.
func (c *Consensus) newViewChangeGadget() error {
	// change the view to stop processing requests
	c.memory.IncView()

	view := c.memory.GetView()
	lwm := c.memory.GetLowWaterMark()
	seq := c.logs.GetSequenceNumber(lwm)

	// create a new view change msg
	message := pbft.ViewChangeMsg{
		NodeId:         c.memory.GetNodeId(),
		View:           int64(view),
		SequenceNumber: int64(seq),
		Preprepares:    c.logs.GetPreprepares(seq, lwm),
		Checkpoints:    c.logs.GetCheckpoints()[lwm],
	}

	// sign the message
	sig, err := tbls.Sign(c.suite, c.tss, []byte(hashing.MD5View(&message)))
	if err != nil {
		c.logger.Error("failed to sign the view message", zap.Error(err))
	}

	// set the signature
	message.Signature = sig

	// append own view change msg
	c.logs.AppendViewChange(view, &message)

	// send a view change message
	if count := c.communication.SendViewChangeMsg(&message); count < c.cfg.Majority {
		c.logger.Info("not enough available servers to start view change", zap.Int("live servers", count))
		return nil
	}

	// create a new timer
	timer := modules.NewTimer(c.cfg.ViewChangeTimeout, time.Millisecond)
	timer.Start(false)

	// wait for 2f+1 messages
	flag := true
	for {
		if !flag {
			break
		}

		select {
		case msg := <-c.viewChangeGadgetChannel:
			c.logs.AppendViewChange(view, msg)

			if len(c.logs.GetViewChanges(view)) >= c.cfg.Majority {
				flag = false
				timer.Stop(false)
			}
		case <-timer.Notify(): // if the timer expired, return and reset the view timer
			return errors.New("view change failed")
		}
	}

	// if the node is the leader, run a new leader gadget
	if c.getCurrentLeader() == c.memory.GetNodeId() {
		c.logger.Debug("new leader", zap.String("id", c.memory.GetNodeId()))
		if !c.memory.GetByzantine() {
			c.newLeaderGadget()
		}
	} else { // if the node is primary, it needs new-view message
		timer = modules.NewTimer(c.cfg.ViewChangeTimeout, time.Millisecond)
		flag = true

		for {
			if !flag {
				break
			}

			select {
			case raw := <-c.consensusHandlersTable[enum.PktNV]:
				msg := raw.Payload.(*pbft.NewViewMsg)

				// update the view
				c.memory.SetView(int(msg.GetView()))

				// set the message for view change
				c.logs.AppendNewView(int(msg.GetView()), msg)

				// update the log
				for _, msg := range msg.GetPreprepares() {
					c.logs.SetRequest(int(msg.GetSequenceNumber()), msg.GetRequest())
				}

				flag = false
			case <-timer.Notify(): // if the timer expired, return and reset the view timer
				return errors.New("view change failed")
			}
		}
	}

	// close our channel
	c.inViewChangeMode = false
	c.viewChangeGadgetChannel = nil

	return nil
}

// newLeaderGadget performs the procedure of new leader.
func (c *Consensus) newLeaderGadget() {
	// get current view
	view := c.memory.GetView()

	// get all view change messages from other nodes
	messages := c.logs.GetViewChanges(view)

	// create a log map to get requests
	logsMap := make(map[int]*pbft.PrePrepareMsg)

	// set the min and max
	lwm := c.memory.GetLowWaterMark()
	minSequence := c.logs.GetSequenceNumber(lwm)
	maxSequence := c.logs.GetSequenceNumber(lwm)

	var (
		message   *pbft.ViewChangeMsg
		sigShares [][]byte
	)

	// loop in all messages
	for _, msg := range messages {
		sequence := int(msg.GetSequenceNumber())

		message = msg
		sigShares = append(sigShares, []byte(msg.GetSignature()))

		// loop over preprepares to insert them inside a logs map
		for _, pp := range msg.GetPreprepares() {
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
	digest := hashing.MD5View(message)

	// recover the full signature
	sig, err := tbls.Recover(c.suite, c.pub, []byte(digest), sigShares[:c.cfg.Majority], c.cfg.Majority, c.cfg.Total)
	if err != nil {
		c.logger.Error("failed to recover signature", zap.Error(err))
	}

	// verify the record signature
	public := c.pub.Commit()
	if err := bls.Verify(c.suite, public, []byte(digest), sig); err != nil {
		c.logger.Error("failed to verify the signature", zap.Error(err))
	}

	// create a new view message
	newViewMsg := pbft.NewViewMsg{
		View:        int64(view),
		NodeId:      c.memory.GetNodeId(),
		Preprepares: requests,
		Messages:    messages,
		Message:     digest,
		Shares:      sigShares,
	}

	// save the entry
	c.logs.AppendNewView(view, &newViewMsg)

	// send new view
	c.communication.SendNewViewMsg(&newViewMsg)

	// start the protocol for every request
	for _, req := range requests {
		c.newProcessingGadet(int(req.GetSequenceNumber()), req)
	}
}
