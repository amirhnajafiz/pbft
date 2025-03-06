package main

import (
	"os"

	"github.com/f24-cse535/pbft/cmd"
	"github.com/f24-cse535/pbft/internal/config"
	"github.com/f24-cse535/pbft/pkg/logger"

	"go.uber.org/zap"
)

// list of current system commands.
// - app: creates a new client application
// - node: creates a single node instance
// - controller: sets up a process that builds node instances
const (
	AppCmd        = "app"
	ControllerCmd = "controller"
	NodeCmd       = "node"
)

func main() {
	argv := os.Args
	if len(argv) < 3 {
		panic("you did not provide enough arguments to run! (./main <command> <config-path>)")
	}

	// load the config file into a config struct
	cfg := config.New(argv[2], true)

	// create a new zap logger instance
	logr := logger.NewLogger(cfg.LogLevel)

	// create cmd instances with config struct and zap logger as input
	commands := map[string]cmd.CMD{
		ControllerCmd: cmd.Controller{
			Cfg:    cfg,
			Logger: logr.Named("controller"),
		},
		NodeCmd: cmd.Node{
			Cfg:    cfg,
			Logger: logr.Named("node"),
		},
		AppCmd: cmd.Application{
			Cfg:    cfg,
			Logger: logr.Named("app"),
		},
	}

	// program uses first argument variable as the command
	if callback, ok := commands[argv[1]]; ok {
		if err := callback.Main(); err != nil { // follow cmd.Main methods
			logr.Panic("failed to run the command", zap.Error(err), zap.String("command", argv[1]))
		}
	} else {
		panic("failed to find the input command")
	}
}
