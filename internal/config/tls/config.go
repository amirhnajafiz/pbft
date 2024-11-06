package tls

// Config stores all tls communication keys.
type Config struct {
	PrivateKey string `koanf:"private_key"` // the app private key
	PublicKey  string `koanf:"public_key"`  // the app public key
	CaKey      string `koanf:"ca_key"`      // the certificate authority key
}
