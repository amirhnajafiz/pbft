package config

import (
	"github.com/f24-cse535/pbft/internal/config/controller"
	"github.com/f24-cse535/pbft/internal/config/node"
)

// Default return default configuration.
func Default() Config {
	return Config{
		LogLevel: "debug",
		Controller: controller.Config{
			CSV: "testcase.csv",
		},
		Node: node.Config{
			Port:     80,
			NodeId:   "unique",
			LogLevel: "debug",
		},
	}
}
