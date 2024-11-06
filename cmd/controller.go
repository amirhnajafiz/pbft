package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/f24-cse535/pbft/internal/config"
	goclient "github.com/f24-cse535/pbft/internal/grpc/client"
	"github.com/f24-cse535/pbft/internal/utils/lists"
	"github.com/f24-cse535/pbft/internal/utils/parser"

	"go.uber.org/zap"
)

// Controller is used to run testcases on the system.
type Controller struct {
	Cfg    config.Config
	Logger *zap.Logger
}

func (c Controller) Main() error {
	// create a new client.go
	cli := goclient.NewClient(
		c.Logger.Named("client.go"),
		"",
		c.Cfg.GetNodes(),
	)
	if err := cli.LoadTLS(c.Cfg.TLS.PrivateKey, c.Cfg.TLS.PublicKey, c.Cfg.TLS.CaKey); err != nil {
		return fmt.Errorf("failed to load TLS keys: %v", err)
	}

	// load the test-case file
	ts, err := parser.CSVInput(c.Cfg.Controller.CSV)
	if err != nil {
		return fmt.Errorf("failed to load the test file: %v", err)
	}

	// set the initial values of our parameters that are used in transaction sending
	index := 0
	timestamp := 10

	// print some metadata
	fmt.Printf("read %s with %d test sets.\n", c.Cfg.Controller.CSV, len(ts))
	for _, t := range ts {
		fmt.Printf(
			"transactions of %s is %d (Live Servers=%d, Byzantine Servers=%d)\n",
			t.Index, len(t.Transactions),
			len(t.LiveServers),
			len(t.ByzantineServers),
		)
	}

	// create a new node list and remove the target client from it
	nodes := c.Cfg.GetNodes()
	delete(nodes, c.Cfg.Controller.Client)

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
			return nil
		case "next":
			if index == len(ts) {
				fmt.Println("test-sets are over.")
			} else {
				testcase := ts[index]

				for key := range nodes {
					cli.Flush(key)
					cli.ChangeState(key, lists.IsInList(key, testcase.LiveServers), lists.IsInList(key, testcase.ByzantineServers))
				}

				fmt.Printf("executing test set %d\n", index)
				for _, trx := range testcase.Transactions {
					fmt.Printf("(timestamp: %d) (%s, %s, %s)\n", timestamp, trx.Sender, trx.Receiver, trx.Amount)

					amount, _ := strconv.Atoi(trx.Amount)
					fmt.Println(cli.Transaction(c.Cfg.Controller.Client, trx.Sender, trx.Receiver, amount, timestamp))

					timestamp++
				}

				for key := range nodes {
					cli.ChangeState(key, true, false)
				}

				index++
			}
		case "printlog":
			for _, item := range cli.PrintLog(parts[1]) {
				fmt.Println(item)
			}
		case "printdb":
			for _, item := range cli.PrintDB(parts[1]) {
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
				fmt.Printf("%s : %s\n", key, cli.PrintStatus(key, seq))
			}
		default:
			fmt.Printf("command `%s` not found.\n", parts[1])
		}
	}
}
