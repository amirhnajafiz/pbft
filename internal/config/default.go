package config

import (
	"github.com/f24-cse535/pbft/internal/config/controller"
	"github.com/f24-cse535/pbft/internal/config/node"
	"github.com/f24-cse535/pbft/internal/config/node/bft"
)

// Default return default configuration.
func Default() Config {
	return Config{
		LogLevel: "debug",
		Controller: controller.Config{
			CSV: "tests/case.csv",
		},
		Node: node.Config{
			Port:   80,
			NodeId: "unique",
			BFT: bft.Config{
				Total:          0,
				Majority:       0,
				Responses:      0,
				RequestTimeout: 0,
			},
		},
		IPTable: IPTable{
			Nodes:   make([]Pair, 0),
			Clients: make([]Pair, 0),
		},
	}
}
