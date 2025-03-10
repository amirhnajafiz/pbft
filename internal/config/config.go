package config

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/f24-cse535/pbft/internal/config/node"
	"github.com/f24-cse535/pbft/internal/config/tls"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/tidwall/pretty"
)

// Prefix indicates environment variables prefix.
const Prefix = "pbft_"

// Config struct is a module that stores system configs.
type Config struct {
	CtlFiles []string `koanf:"ctl_files"` // a list of config files to run by ctl command
	LogLevel string   `koanf:"log_level"` // node logging level (debug, info, warn, error, panic, fatal)
	CSV      string   `koanf:"csv"`       // testcase file path

	Node node.Config `koanf:"node"` // node app configs
	TLS  tls.Config  `koanf:"tls"`  // node tls keys

	Clients []Pair `koanf:"clients"` // system clients
	IPTable []Pair `koanf:"iptable"` // system IP addresses
}

// New reads configuration with koanf, by loading a yaml config path into the Config struct.
func New(path string, print bool) Config {
	var instance Config

	k := koanf.New(".")

	// load default configuration from file
	if err := k.Load(structs.Provider(Default(), "koanf"), nil); err != nil {
		log.Fatalf("error loading default: %s", err)
	}

	// load configuration from file
	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		log.Printf("error loading config.yml: %s", err)
	}

	// load environment variables
	if err := k.Load(env.Provider(Prefix, ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(
			strings.TrimPrefix(s, Prefix)), "__", ".")
	}), nil); err != nil {
		log.Printf("error loading environment variables: %s", err)
	}

	if err := k.Unmarshal("", &instance); err != nil {
		log.Fatalf("error unmarshalling config: %s", err)
	}

	indent, err := json.MarshalIndent(instance, "", "\t")
	if err != nil {
		log.Fatalf("error marshaling config to json: %s", err)
	}

	indent = pretty.Color(indent, nil)
	tmpl := `
	================ Loaded Configuration ================
	%s
	=============================================
	`

	if print {
		log.Printf(tmpl, string(indent))
	}

	return instance
}
