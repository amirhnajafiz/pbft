package cmd

import (
	"github.com/f24-cse535/pbft/internal/config"

	"go.uber.org/zap"
)

// Each node of our transaction system runs using this main function.
type Node struct {
	Cfg    config.Config
	Logger *zap.Logger
}

func (n Node) Main() error {
	return nil
}
