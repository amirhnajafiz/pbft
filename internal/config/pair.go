package config

// Pair is a tiny module that is used to read key-value pairs from the input config file.
type Pair struct {
	Key   string `koanf:"key"`
	Value string `koanf:"value"`
}
