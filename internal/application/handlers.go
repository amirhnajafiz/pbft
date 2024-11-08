package application

import (
	"errors"
	"fmt"
	"time"

	"github.com/f24-cse535/pbft/pkg/rpc/app"
	"github.com/f24-cse535/pbft/pkg/rpc/pbft"
	"go.uber.org/zap"
)

// transactionHandler gets client transactions and call request handler.
func (a *App) transactionHandler(client string) {
	// timestamp variable
	ts := 10

	a.logger.Debug("transaction handler started", zap.String("client", client))

	for {
		// get model transactions
		trx := <-a.clients[client]

		fmt.Printf("processing (%s, %s, %d)\n", trx.Sender, trx.Receiver, trx.Amount)

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

		fmt.Printf("done (%s, %s, %d) at %d : `%s`\n", trx.Sender, trx.Receiver, trx.Amount, ts, resp)

		// increase timestamp
		ts++
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

	currentLeader := a.getCurrentLeader() // get current leader id

	a.handlers[client] = make(chan *app.ReplyMsg, a.cfg.Total) // get our channel
	ch := a.handlers[client]

	// send the request
	if err := a.cli.Request(currentLeader, req); err != nil {
		a.logger.Debug("failed to send the request", zap.String("leader", currentLeader))

		// send the request to all servers
		count := 0
		for key := range a.cli.GetSystemNodes() {
			if er := a.cli.Request(key, req); er != nil {
				a.logger.Debug("failed to resend the request", zap.String("target", key))
			} else {
				count++
			}
		}

		// if the number of live servers is less than 2f+1, then raise an error
		if count < a.cfg.Majority {
			a.logger.Debug("majority of servers are unavailable", zap.Int("servers", count))
			return "not enough servers are available"
		}
	}

	// handle the reply
	resp, err := a.replyHandler(ch, trx)
	if err != nil {
		a.logger.Debug("request blocked or got timeout", zap.Error(err))

		// send the request to all servers
		count := 0
		for key := range a.cli.GetSystemNodes() {
			if er := a.cli.Request(key, req); er != nil {
				a.logger.Debug("failed to resend the request", zap.String("target", key))
			} else {
				count++
			}
		}

		// if the number of live servers is less than 2f+1, then raise an error
		if count < a.cfg.Majority {
			a.logger.Debug("majority of servers are unavailable", zap.Int("servers", count))
			return "not enough servers are available"
		}

		// again, wait for replys
		resp, err = a.replyHandler(ch, trx)
		if err != nil {
			a.logger.Debug("servers are not responding", zap.Error(err))

			return "servers are not responding"
		}
	}

	// empty the channel
	go func(key string) {
		for {
			select {
			case <-ch:
				continue
			default:
				delete(a.clients, key)
				return
			}
		}
	}(client)

	// update the view
	a.memory.SetView(int(resp.GetView()))

	return resp.GetResponse()
}

// replyHandler waits for f+1 messages. it will raise an error if it get's timeout.
func (a *App) replyHandler(ch chan *app.ReplyMsg, trx *pbft.TransactionMsg) (*app.ReplyMsg, error) {
	// create a map
	replyMap := make(map[int]*app.ReplyMsg)
	countMap := make(map[int]int)

	// create a new timer
	timer := time.NewTimer(time.Duration(a.cfg.RequestTimeout) * time.Millisecond)

	// reply handler loop
	for {
		select {
		case <-timer.C:
			return nil, errors.New("request timeout") // request timeout
		case reply := <-ch:
			if reply.GetTimestamp() != trx.GetTimestamp() {
				continue
			}

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
