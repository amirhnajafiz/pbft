package application

import (
	"fmt"

	"github.com/f24-cse535/pbft/pkg/rpc/pbft"
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
func (a *App) requestHandler(_ *pbft.TransactionMsg) string {
	return ""
}
