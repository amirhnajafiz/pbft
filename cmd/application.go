package cmd

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/f24-cse535/pbft/internal/application"
	"github.com/f24-cse535/pbft/internal/config"
	"github.com/f24-cse535/pbft/internal/grpc/client"
	"github.com/f24-cse535/pbft/internal/storage/local"
	"github.com/f24-cse535/pbft/internal/utils/lists"
	"github.com/f24-cse535/pbft/internal/utils/parser"

	"go.uber.org/zap"
)

// Application is the client program in our system.
type Application struct {
	Cfg    config.Config
	Logger *zap.Logger
}

func (a Application) Main() error {
	// create a memory instance
	mem := local.NewMemory(a.Cfg.Node.NodeId, a.Cfg.Node.BFT.Total)
	mem.SetNodes(a.Cfg.GetNodesMeta())

	// load tls configs
	creds, err := a.Cfg.TLS.Creds()
	if err != nil {
		return err
	}

	// create a new client.go instance
	cli := client.NewClient(creds, a.Cfg.Node.NodeId, a.Cfg.GetNodes())

	// create a new app instance
	app := application.NewApp(
		a.Logger.Named("app"),
		mem,
		&a.Cfg.Node.BFT,
		cli,
		a.Cfg.GetClients(),
	)

	// terminal goes into a for loop to get user inputs
	go a.terminal(app)

	// start the gRPC server to get requests' replys
	return app.Service(a.Cfg.Node.Port, creds)
}

// terminal method reads a csv test-case file, and waits for input commands.
func (a Application) terminal(app *application.App) {
	// load the test-case file
	testCases, err := parser.CSVInput(a.Cfg.CSV)
	if err != nil {
		panic(err)
	}

	// set the tests index to zero
	testCaseIndex := 0

	// print some metadata about the test-case file
	fmt.Printf("input test-case file %s with %d test sets.\n", a.Cfg.CSV, len(testCases))
	for _, testSet := range testCases {
		fmt.Printf(
			"\t- number transactions in set %s is %d (LiveServers=%d, ByzantineServers=%d)\n",
			testSet.Index, len(testSet.Transactions),
			len(testSet.LiveServers),
			len(testSet.ByzantineServers),
		)
	}

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

		// switch case on the first part of the input
		switch parts[0] {
		case "exit":
			os.Exit(0)
		case "next":
			if testCaseIndex < len(testCases) {
				currentTestSet := testCases[testCaseIndex]

				// reset the app
				app.Reset()

				// flush and reset our nodes status
				for key := range a.Cfg.GetNodes() {
					app.Client().Flush(key)
					app.Client().ChangeState(key, lists.IsInList(key, currentTestSet.LiveServers), lists.IsInList(key, currentTestSet.ByzantineServers))
				}

				fmt.Printf("- executing set %s\n", currentTestSet.Index)

				// loop over transactions and execute them
				for _, trx := range currentTestSet.Transactions {
					app.Transaction(trx) // see application.Transaction method
				}

				testCaseIndex++
			} else {
				fmt.Println("end of the test sets!")
			}
		case "unblock":
			for key := range a.Cfg.GetNodes() { // reset only the nodes status (without flushing)
				app.Client().ChangeState(key, true, false)
			}
		case "printlog":
			for _, item := range app.Client().PrintLog(parts[1]) {
				fmt.Printf("- %s\n", item)
			}
		case "printdb":
			for _, item := range app.Client().PrintDB(parts[1]) {
				fmt.Printf(
					"\t- %d : %d (%s, %s, %d) : %s\n",
					item.GetSequenceNumber(),
					item.GetRequest().GetTransaction().GetTimestamp(),
					item.GetRequest().GetTransaction().GetSender(),
					item.GetRequest().GetTransaction().GetReciever(),
					item.GetRequest().GetTransaction().GetAmount(),
					item.GetRequest().GetResponse().GetText(),
				)
			}
		case "printstatus":
			seq, _ := strconv.Atoi(parts[1]) // get the sequence number

			for key := range a.Cfg.GetNodes() {
				fmt.Printf("\t- %s : %s\n", key, app.Client().PrintStatus(key, seq))
			}
		case "printview":
			for key := range a.Cfg.GetNodes() {
				fmt.Printf("- node %s\n", key)

				for _, item := range app.Client().PrintView(key) {
					fmt.Printf("\t- new view: %d\n", item.GetNewviewMessage().GetView())

					fmt.Printf("\t\t- preprepare messages:\n")
					for _, msg := range item.GetNewviewMessage().GetPreprepareMessages() {
						fmt.Printf("\t\t\t- seq=%d timestamp=%d digest=%s\n", msg.GetSequenceNumber(), msg.Request.Transaction.GetTimestamp(), msg.GetDigest())
					}

					fmt.Printf("\t\tviewchange messages:\n")
					for _, msg := range item.GetViewchangeMessages() {
						fmt.Printf("\t\t\t- sender=%s sequence=%d\n", msg.GetNodeId(), msg.GetSequenceNumber())
						for _, msg := range msg.GetPreprepareMessages() {
							fmt.Printf("\t\t\t\t- seq=%d timestamp=%d digest=%s\n", msg.GetSequenceNumber(), msg.Request.Transaction.GetTimestamp(), msg.GetDigest())
						}
					}

					fmt.Printf("\t\t- treshold signature message: %s\n", item.GetNewviewMessage().GetViewchangeMessage())
					for _, sh := range item.GetNewviewMessage().GetShares() {
						fmt.Printf("\t\t\t- %s\n", base64.StdEncoding.EncodeToString(sh))
					}
				}
			}
		case "printcheckpoint":
			for key := range a.Cfg.GetNodes() {
				fmt.Printf("- node %s\n", key)
				for _, item := range app.Client().PrintCheckpoints(key) {
					fmt.Printf("\t- sequence=%d\n", item.GetSequenceNumber())
					for _, checkpoint := range item.GetCheckpointMessages() {
						fmt.Printf("\t\t- sender=%s sequence=%d\n", checkpoint.GetNodeId(), checkpoint.GetSequenceNumber())
						for _, msg := range checkpoint.GetPreprepareMessages() {
							fmt.Printf("\t\t\t- seq=%d timestamp=%d digest=%s\n", msg.GetSequenceNumber(), msg.GetRequest().GetTransaction().GetTimestamp(), msg.GetDigest())
						}
					}
				}
			}
		default:
			fmt.Printf("command `%s` not found.\n", parts[0])
		}
	}
}
