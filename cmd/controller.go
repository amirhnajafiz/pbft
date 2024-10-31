package cmd

import (
	"github.com/f24-cse535/pbft/internal/config"

	"go.uber.org/zap"
)

// Controller is used to communicate with our distributed system using gRPC calls.
type Controller struct {
	Cfg    config.Config
	Logger *zap.Logger
}

func (c Controller) Main() error {
	return nil
}
