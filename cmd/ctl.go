package cmd

import (
	"log"
	"os/exec"

	"github.com/f24-cse535/pbft/internal/config"
)

type CTL struct {
	Cfg config.Config
}

func (c CTL) Main() error {
	for _, key := range c.Cfg.CtlFiles {
		go func(config string) {
			if err := exec.Command("./main", "node", config).Run(); err != nil {
				log.Println(err)
			}
		}(key)
	}
}
