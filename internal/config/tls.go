package config

// TLS stores all secure communication keys.
type TLS struct {
	PrivateKey string `koanf:"private_key"` // the app private key
	PublicKey  string `koanf:"public_key"`  // the app public key
	CaKey      string `koanf:"ca_key"`      // the certificate authority key
}
