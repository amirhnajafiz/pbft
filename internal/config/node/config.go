package node

import "github.com/f24-cse535/pbft/internal/config/node/bft"

// Config parameters for running a node.
type Config struct {
	Port   int        `koanf:"port"`    // gRPC server port
	NodeId string     `koanf:"node_id"` // a unique id for the node
	BFT    bft.Config `koanf:"bft"`     // node PBFT config values
}
