package node

// Config parameters for running a node.
type Config struct {
	Port     int    `koanf:"port"`      // gRPC server port
	NodeId   string `koanf:"node_id"`   // a unique id for the node
	LogLevel string `koanf:"log_level"` // node logging level (debug, info, warn, error, panic, fatal)
}
