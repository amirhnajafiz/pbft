package node

// Config parameters for running a node.
type Config struct {
	NodeId string `koanf:"node_id"` // a unique id for the node
}
