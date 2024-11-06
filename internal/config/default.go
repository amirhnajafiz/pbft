package config

import (
	"github.com/f24-cse535/pbft/internal/config/controller"
	"github.com/f24-cse535/pbft/internal/config/node"
	"github.com/f24-cse535/pbft/internal/config/node/bft"
	"github.com/f24-cse535/pbft/internal/config/tls"
)

// Default return default configuration.
func Default() Config {
	return Config{
		LogLevel: "debug",
		Controller: controller.Config{
			CSV:    "tests/case.csv",
			Client: "",
		},
		Node: node.Config{
			Port:     80,
			NodeId:   "unique",
			CoreNode: false,
			BFT: bft.Config{
				Total:           0,
				Majority:        0,
				Responses:       0,
				RequestTimeout:  0,
				MajorityTimeout: 0,
				ViewTimeout:     0,
			},
		},
		IPTable: make([]Pair, 0),
		Clients: make([]Pair, 0),
		TLS: tls.Config{
			PrivateKey: "",
			PublicKey:  "",
			CaKey:      "",
		},
	}
}
