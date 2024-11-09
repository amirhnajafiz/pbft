package config

import (
	"github.com/f24-cse535/pbft/internal/config/node"
	"github.com/f24-cse535/pbft/internal/config/node/bft"
	"github.com/f24-cse535/pbft/internal/config/tls"
)

// Default return default configuration.
func Default() Config {
	return Config{
		CtlFiles: make([]string, 0),
		LogLevel: "debug",
		CSV:      "",
		Node: node.Config{
			Port:   80,
			NodeId: "unique",
			BFT: bft.Config{
				Total:             0,
				Majority:          0,
				Responses:         0,
				RequestTimeout:    0,
				MajorityTimeout:   0,
				ViewTimeout:       0,
				ViewChangeTimeout: 0,
				NewViewTimeout:    0,
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
