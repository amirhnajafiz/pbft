package consensus

import (
	"errors"
	"time"

	"github.com/f24-cse535/pbft/internal/consensus/modules"
	"github.com/f24-cse535/pbft/internal/utils/hashing"
	"github.com/f24-cse535/pbft/pkg/enum"
	"github.com/f24-cse535/pbft/pkg/models"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"

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
	c.requestsHandlersTable[sequence] = channel

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
	delete(c.requestsHandlersTable, sequence)

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
			return
		}
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
	seq := c.logs.GetSequenceNumber()

	// create a new view change msg
	message := pbft.ViewChangeMsg{
		NodeId:         c.memory.GetNodeId(),
		View:           int64(view),
		SequenceNumber: int64(seq),
		Preprepares:    c.logs.GetPreprepares(seq),
	}

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
	minSequence := c.logs.GetSequenceNumber()
	maxSequence := c.logs.GetSequenceNumber()

	// loop in all messages
	for _, msg := range messages {
		sequence := int(msg.GetSequenceNumber())

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

	// create a new view message
	newViewMsg := pbft.NewViewMsg{
		View:        int64(view),
		NodeId:      c.memory.GetNodeId(),
		Preprepares: requests,
		Messages:    messages,
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
