package main

import (
	"os"

	"github.com/f24-cse535/pbft/cmd"
	"github.com/f24-cse535/pbft/internal/config"
	"github.com/f24-cse535/pbft/pkg/logger"

	"go.uber.org/zap"
)

// list of current system commands.
const (
	ControllerCmd = "controller"
	NodeCmd       = "node"
	CtlCmd        = "ctl"
)

func main() {
	argv := os.Args
	if len(argv) < 3 {
		panic("you did not provide enough arguments to run! (./main <command> <config-path>)")
	}

	// load configs into a config struct
	cfg := config.New(argv[2], true)

	// create a new zap logger instance
	logr := logger.NewLogger(cfg.LogLevel)

	// create cmd instances and pass needed parameters
	commands := map[string]cmd.CMD{
		ControllerCmd: cmd.Controller{
			Cfg:    cfg,
			Logger: logr.Named("controller"),
		},
		NodeCmd: cmd.Node{
			Cfg:    cfg,
			Logger: logr.Named("node." + cfg.Node.NodeId),
		},
		CtlCmd: cmd.CTL{
			Cfg:    cfg,
			Logger: logr.Named("ctl"),
		},
	}

	// command is the first argument variable
	if callback, ok := commands[argv[1]]; ok {
		if err := callback.Main(); err != nil {
			logr.Panic("failed to run the command", zap.Error(err), zap.String("command", argv[1]))
		}
	} else {
		panic("your input command must be the first argument variable")
	}
}
