package tls

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"
)

// Config stores all tls communication keys.
type Config struct {
	PrivateKey string `koanf:"private_key"` // the app private key
	PublicKey  string `koanf:"public_key"`  // the app public key
	CaKey      string `koanf:"ca_key"`      // the certificate authority key
}

// TLS returns a crypto tls config type.
func (c *Config) TLS() (*tls.Config, error) {
	// load the clients keys
	prkBytes, err := os.ReadFile(c.PrivateKey)
	if err != nil {
		return nil, err
	}

	pukBytes, err := os.ReadFile(c.PublicKey)
	if err != nil {
		return nil, err
	}

	// load certificate
	cert, err := tls.X509KeyPair(pukBytes, prkBytes)
	if err != nil {
		return nil, err
	}

	// create the CA data
	ca := x509.NewCertPool()
	cacBytes, err := os.ReadFile(c.CaKey)
	if err != nil {
		return nil, err
	}
	if ok := ca.AppendCertsFromPEM(cacBytes); !ok {
		return nil, errors.New("failed to append certs")
	}

	return &tls.Config{
		ClientAuth:         tls.RequireAndVerifyClientCert,
		Certificates:       []tls.Certificate{cert},
		ClientCAs:          ca,
		InsecureSkipVerify: true,
	}, nil
}
