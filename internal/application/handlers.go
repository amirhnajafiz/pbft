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

	for {
		// get model transactions
		trx := <-a.clients[client]

		// call requestHandler
		fmt.Println(a.requestHandler(&pbft.TransactionMsg{
			Sender:    trx.Sender,
			Reciever:  trx.Receiver,
			Amount:    trx.Amount,
			Timestamp: int64(ts),
		}))

		// increase timestamp
		ts++
	}
}

// requestHandler gets a transaction and starts PBFT procedures.
func (a *App) requestHandler(trx *pbft.TransactionMsg) string {
	// create a pbft request
	req := &pbft.RequestMsg{
		Transaction: trx,
		ClientId:    a.memory.GetNodeId(),
		Response:    &pbft.TransactionRsp{},
	}

	currentLeader := a.getCurrentLeader() // get current leader id

	// send the request
	if err := a.cli.Request(currentLeader, req); err != nil {
		a.logger.Debug("failed to send the request", zap.String("leader", currentLeader))
	}

	// handle the reply
	resp, err := a.replyHandler(trx)
	if err != nil {
		a.logger.Debug("request blocked or got timeout", zap.Error(err))
	}

	// update the view
	a.memory.SetView(int(resp.GetView()))

	return resp.GetResponse()
}

// replyHandler waits for f+1 messages. it will raise an error if it get's timeout.
func (a *App) replyHandler(trx *pbft.TransactionMsg) (*app.ReplyMsg, error) {
	// create a map
	replyMap := make(map[int]*app.ReplyMsg)
	countMap := make(map[int]int)

	// create a new timer
	timer := time.NewTimer(time.Duration(a.cfg.RequestTimeout) * time.Millisecond)

	// get the response channel
	ch := a.handlers[trx.GetSender()]

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
