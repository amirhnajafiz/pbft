package bft

// Config stores BFT protocl parameters.
type Config struct {
	Total     int `koanf:"total"`
	Majority  int `koanf:"majority"`
	Responses int `koanf:"responses"`
}
