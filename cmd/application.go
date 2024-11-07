package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/f24-cse535/pbft/internal/application"
	"github.com/f24-cse535/pbft/internal/config"
	"github.com/f24-cse535/pbft/internal/grpc/client"
	"github.com/f24-cse535/pbft/internal/utils/parser"

	"go.uber.org/zap"
)

// Application is the client program in our system.
type Application struct {
	Cfg    config.Config
	Logger *zap.Logger
}

func (a Application) Main() error {
	// load tls configs
	creds, err := a.Cfg.TLS.TLS()
	if err != nil {
		return err
	}

	// create a new client.go
	cli := client.NewClient(
		creds,
		a.Cfg.Node.NodeId,
		a.Cfg.GetNodes(),
	)

	// create a new app instance
	app := application.NewApp(cli, config.Default().GetClients())

	// start getting user inputs
	go a.terminal(app)

	// start the gRPC server
	return app.Service(a.Cfg.Node.Port, creds)
}

func (a Application) terminal(app *application.App) {
	// load the test-case file
	ts, err := parser.CSVInput(a.Cfg.Controller.CSV)
	if err != nil {
		panic(err)
	}

	// set the tests index to zero
	index := 0

	// print some metadata
	fmt.Printf("read %s with %d test sets.\n", a.Cfg.Controller.CSV, len(ts))
	for _, t := range ts {
		fmt.Printf(
			"transactions of %s is %d (Live Servers=%d, Byzantine Servers=%d)\n",
			t.Index, len(t.Transactions),
			len(t.LiveServers),
			len(t.ByzantineServers),
		)
	}

	// create a new node list and remove the target client from it
	nodes := a.Cfg.GetNodes()
	delete(nodes, a.Cfg.Controller.Client)

	// in a for loop, read user commands
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")

		input, _ := reader.ReadString('\n') // read input until newline
		input = strings.TrimSpace(input)

		if len(input) == 0 {
			continue
		}

		parts := strings.Split(input, " ") // break by space and split into parts

		switch parts[0] {
		case "exit":
			os.Exit(0)
		case "next":
			if index < len(ts) {
				for _, trx := range ts[index].Transactions {
					app.Transaction(trx)
				}

				index++
			}
		case "printlog":
			for _, item := range app.Client().PrintLog(parts[1]) {
				fmt.Println(item)
			}
		case "printdb":
			for _, item := range app.Client().PrintDB(parts[1]) {
				fmt.Printf(
					"%d : %d (%s, %s, %d) : %s\n",
					item.GetSequenceNumber(),
					item.GetTransaction().GetTimestamp(),
					item.GetTransaction().GetSender(),
					item.GetTransaction().GetReciever(),
					item.GetTransaction().GetAmount(),
					item.GetResponse().GetText(),
				)
			}
		case "printstatus":
			seq, _ := strconv.Atoi(parts[1])
			for key := range nodes {
				fmt.Printf("%s : %s\n", key, app.Client().PrintStatus(key, seq))
			}
		default:
			fmt.Printf("command `%s` not found.\n", parts[1])
		}
	}
}
