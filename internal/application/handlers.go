package application

import (
	"errors"
	"fmt"
	"time"

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
	// create a pbft request
	req := &pbft.RequestMsg{
		Transaction:    trx,
		ClientId:       a.memory.GetNodeId(),
		Response:       &pbft.TransactionRsp{},
		SequenceNumber: -1,
	}

	var (
		resp *app.ReplyMsg
		err  error
	)

	currentLeader := a.getCurrentLeader() // get current leader id

	a.handlers[client] = make(chan *app.ReplyMsg, a.cfg.Total) // get our channel
	ch := a.handlers[client]

	// send the request
	if err := a.cli.Request(currentLeader, req); err != nil {
		// if the number of live servers is less than 2f+1, then raise an error
		if count := a.broadcastRequest(req); count < a.cfg.Majority {
			return "not enough servers are available"
		}
	}

	// follow the reply
	for {
		// handle the reply
		resp, err = a.replyHandler(ch, trx)
		if err == nil {
			break
		}

		// if the number of live servers is less than 2f+1, then raise an error
		if count := a.broadcastRequest(req); count < a.cfg.Majority {
			return "not enough servers are available"
		}
	}

	// empty the channel
	go a.emptyChannel(client)

	// update the view
	a.memory.SetView(int(resp.GetView()))

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
		if ch, ok := a.handlers[msg.GetSender()]; ok {
			ch <- msg
		}
	}
}
