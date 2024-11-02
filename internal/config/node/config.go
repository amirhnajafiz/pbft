package node

import "github.com/f24-cse535/pbft/internal/config/node/bft"

// Config parameters for running a node.
type Config struct {
	Port     int        `koanf:"port"`      // gRPC server port
	NodeId   string     `koanf:"node_id"`   // a unique id for the node
	LogLevel string     `koanf:"log_level"` // node logging level (debug, info, warn, error, panic, fatal)
	IsClient bool       `koanf:"is_client"` // if set true, the node will behaive as client
	BFT      bft.Config `koanf:"bfp"`       // node PBFT config values
}
