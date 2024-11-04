package bft

// Config stores BFT protocl parameters.
type Config struct {
	Total           int `koanf:"total"`            // 3f+1
	Majority        int `koanf:"majority"`         // 2f+1
	Responses       int `koanf:"responses"`        // f+1
	RequestTimeout  int `koanf:"request_timeout"`  // in milliseconds
	MajorityTimeout int `koanf:"majority_timeout"` // in milliseconds
}
