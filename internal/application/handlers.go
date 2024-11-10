package application

import (
	"errors"
	"fmt"
	"time"

	"github.com/f24-cse535/pbft/pkg/enum"
	"github.com/f24-cse535/pbft/pkg/models"
	"github.com/f24-cse535/pbft/pkg/rpc/app"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"
)

// transactionHandler gets client transactions and call request handler.
func (a *App) transactionHandler(client string, ch chan *models.Transaction) {
	for {
		// get model transactions
		trx := <-ch

		// get a timestamp
		ts := a.memory.GetTimestamp()

		fmt.Printf("\t - %d processing (%s, %s, %d)\n", ts, trx.Sender, trx.Receiver, trx.Amount)

		// call requestHandler
		resp := a.requestHandler(
			client,
			&pbft.TransactionMsg{
				Sender:    trx.Sender,
				Reciever:  trx.Receiver,
				Amount:    trx.Amount,
				Timestamp: int64(ts),
			},
		)

		fmt.Printf("\t - %d done (%s, %s, %d): `%s`\n", ts, trx.Sender, trx.Receiver, trx.Amount, resp)
	}
}

// requestHandler gets a transaction and starts PBFT procedures.
func (a *App) requestHandler(client string, trx *pbft.TransactionMsg) string {
	defer func() {
		go a.emptyChannel(client) // empty the channel on exit
	}()

	// create a new pbft request
	req := &pbft.RequestMsg{
		Transaction: trx,
		ClientId:    a.memory.GetNodeId(),
		Response:    &pbft.TransactionRsp{},
	}

	var (
		resp *app.ReplyMsg
		err  error
	)

	currentLeader := a.getCurrentLeader() // get current leader id

	// create a new channel for this handler
	a.lock.Lock()
	a.handlers[client] = make(chan *app.ReplyMsg, a.cfg.Total)
	ch := a.handlers[client]
	a.lock.Unlock()

	// send the request to the leader
	a.cli.Request(currentLeader, req)

	// wait for the reply
	resp, err = a.replyHandler(ch, trx)
	if err != nil {
		flag := false

		for i := 0; i < 3; i++ {
			// if the number of live servers is less than 2f+1, then raise an error
			if count := a.broadcastRequest(req); count < a.cfg.Majority {
				return enum.RespNotEnoughServers
			}

			// wait for f+1 messages
			resp, err = a.replyHandler(ch, trx)
			if err == nil {
				flag = true
				break
			}
		}

		if !flag {
			return enum.RespSystemFailed
		}
	}

	a.memory.SetView(int(resp.GetView())) // update the view

	return resp.GetResponse()
}

// replyHandler waits for f+1 messages. it will raise an error if it get's timeout.
func (a *App) replyHandler(ch chan *app.ReplyMsg, trx *pbft.TransactionMsg) (*app.ReplyMsg, error) {
	// create a map
	hashMap := make(map[string]bool)
	replyMap := make(map[int]*app.ReplyMsg)
	countMap := make(map[int]int)

	// create a new timer
	timer := time.NewTimer(time.Duration(a.cfg.RequestTimeout) * time.Second)

	// reply handler loop
	for {
		select {
		case <-timer.C:
			return nil, errors.New("request timeout") // request timeout
		case reply := <-ch:
			if reply.GetTimestamp() != trx.GetTimestamp() {
				continue
			}

			// drop redundunt messages from servers
			node := reply.GetNodeId()
			if _, ok := hashMap[node]; ok {
				continue
			}

			hashMap[node] = true

			// a logic to count the number of unique replys
			key := int(reply.GetTimestamp() + reply.GetSequenceNumber())
			if _, ok := countMap[key]; !ok {
				countMap[key] = 0
				replyMap[key] = reply
			}

			countMap[key]++

			for key, value := range countMap {
				if value >= a.cfg.Responses {
					return replyMap[key], nil
				}
			}
		}
	}
}

// gRPCHandler dispatchs all reply messages to their handlers.
func (a *App) gRPCHandler() {
	for {
		// get the message from replys
		msg := <-a.replys

		// publish to the right handler
		a.lock.Lock()
		if ch, ok := a.handlers[msg.GetSender()]; ok {
			ch <- msg
		}
		a.lock.Unlock()
	}
}
