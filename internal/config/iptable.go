package config

// IPTable is map table for converting users and nodes to their IP address.
type IPTable struct {
	Nodes   []Pair `koanf:"nodes"`
	Clients []Pair `koanf:"clients"`
}

// GetNodes returns a map of nodes and their IP addresses.
func (i IPTable) GetNodes() map[string]string {
	hashMap := make(map[string]string)

	for _, pair := range i.Nodes {
		hashMap[pair.Key] = pair.Value
	}

	return hashMap
}

// GetClients returns a map of clients and their IP addresses.
func (i IPTable) GetClients() map[string]string {
	hashMap := make(map[string]string)

	for _, pair := range i.Clients {
		hashMap[pair.Key] = pair.Value
	}

	return hashMap
}
