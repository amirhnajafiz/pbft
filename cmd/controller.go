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

// Controller is used to communicate with our distributed system using gRPC calls.
type Controller struct {
	Cfg    config.Config
	Logger *zap.Logger

	client *goclient.Client
	index  int
}

func (c Controller) Main() error {
	// create a new client.go
	c.client = goclient.NewClient(
		c.Logger.Named("client.go"),
		"",
		c.Cfg.IPTable.GetNodes(),
		c.Cfg.IPTable.GetClients(),
	)
	if err := c.client.LoadTLS(c.Cfg.PrivateKey, c.Cfg.PublicKey, c.Cfg.CAC); err != nil {
		return err
	}

	// load the test-case file
	ts, err := parser.CSVInput(c.Cfg.Controller.CSV)
	if err != nil {
		return err
	}

	// reset the index
	c.index = 0
	timestamp := 10

	// print some metadata
	fmt.Printf("read %s with %d test sets.\n", c.Cfg.Controller.CSV, len(ts))
	for _, t := range ts {
		fmt.Printf("transactions of %s is %d (Live Servers=%d, Byzantine Servers=%d)\n", t.Index, len(t.Transactions), len(t.LiveServers), len(t.ByzantineServers))
	}

	// in a for loop, read user commands
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")

		input, _ := reader.ReadString('\n') // read input until newline
		input = strings.TrimSpace(input)

		// no input
		if len(input) == 0 {
			continue
		}

		// break by space and split into parts
		parts := strings.Split(input, " ")

		// switch on the first part
		switch parts[0] {
		case "exit":
			return nil
		case "next":
			if c.index == len(ts) {
				fmt.Println("test-sets are over.")
			} else {
				fmt.Printf("executing set %d\n", c.index)

				c.resetNodesStatus()
				c.updateNodesStatus(ts[c.index].LiveServers, ts[c.index].ByzantineServers)

				for _, trx := range ts[c.index].Transactions {
					fmt.Printf("execute (timestamp: %d) (%s, %s, %s)\n", timestamp, trx.Sender, trx.Receiver, trx.Amount)

					amount, _ := strconv.Atoi(trx.Amount)
					fmt.Println(c.client.Transaction(trx.Sender, trx.Receiver, amount, timestamp))

					timestamp++
				}

				c.index++
			}
		case "printlog":
			for _, item := range c.client.PrintLog(parts[1]) {
				fmt.Println(item)
			}
		case "printdb":
			for _, item := range c.client.PrintDB(parts[1]) {
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
			for key := range c.Cfg.IPTable.GetNodes() {
				fmt.Printf("%s : %s\n", key, c.client.PrintStatus(key, seq))
			}
		default:
			fmt.Printf("command `%s` not found.\n", parts[1])
		}
	}
}

// updateNodesStatus before running each transaction.
func (c Controller) updateNodesStatus(liveservers []string, byzantines []string) {
	for key := range c.Cfg.IPTable.GetNodes() {
		c.client.ChangeState(key, lists.IsInList(key, liveservers), lists.IsInList(key, byzantines))
	}
}

// resetNodesStatus by calling flush and reset status.
func (c Controller) resetNodesStatus() {
	for key := range c.Cfg.IPTable.GetNodes() {
		c.client.ChangeState(key, true, false)
		c.client.Flush(key)
	}
}
