package config

// GetNodes returns a map of nodes and their IP addresses.
func (c Config) GetNodes() map[string]string {
	hashMap := make(map[string]string)

	for _, pair := range c.IPTable {
		hashMap[pair.Key] = pair.Value
	}

	return hashMap
}

// GetNodesMeta returns a map of nodes and their metadata.
func (c Config) GetNodesMeta() map[string]int {
	hashMap := make(map[string]int)

	for _, pair := range c.IPTable {
		hashMap[pair.Key] = pair.Metadata
	}

	return hashMap
}

// GetClients return a map of clients and their balances.
func (c Config) GetClients() map[string]int {
	hashMap := make(map[string]int)

	for _, pair := range c.Clients {
		hashMap[pair.Key] = pair.Metadata
	}

	return hashMap
}
