package cmd

import (
	"fmt"

	"github.com/f24-cse535/pbft/internal/config"
	"github.com/f24-cse535/pbft/internal/utils/parser"
)

// Controller is used to communicate with our distributed system using gRPC calls.
type Controller struct {
	Cfg config.Config
}

func (c Controller) Main() error {
	// load the test-case file
	ts, err := parser.CSVInput(c.Cfg.Controller.CSV)
	if err != nil {
		return err
	}

	fmt.Printf("read %s with %d test sets.\n", c.Cfg.Controller.CSV, len(ts))
	for _, t := range ts {
		fmt.Printf("transactions of %s is %d (LS=%d, BY=%d)\n", t.Index, len(t.Transactions), len(t.LiveServers), len(t.ByzantineServers))
	}

	return nil
}
