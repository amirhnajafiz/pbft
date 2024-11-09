package cmd

import (
	"sync"

	"github.com/f24-cse535/pbft/internal/config"

	"go.uber.org/zap"
)

// Controller is a program for creating system's nodes.
type Controller struct {
	Cfg    config.Config
	Logger *zap.Logger
}

func (c Controller) Main() error {
	// init a waitgroup
	var wg sync.WaitGroup

	// loop over config files and run each node
	for _, key := range c.Cfg.CtlFiles {
		wg.Add(1)
		go func(path string, waitGroup *sync.WaitGroup) {
			// release waitgroup
			defer func() {
				wg.Done()
			}()

			// create a new node instance
			cfg := config.New(path, false)
			node := Node{
				Cfg:    cfg,
				Logger: c.Logger.Named("node." + cfg.Node.NodeId),
			}

			c.Logger.Info("instance started", zap.String("instance", cfg.Node.NodeId), zap.String("file", path))

			if err := node.Main(); err != nil {
				c.Logger.Error("failed to start instance", zap.String("instance", cfg.Node.NodeId), zap.Error(err))
			} else {
				c.Logger.Info("instance terminated", zap.String("instance", cfg.Node.NodeId), zap.String("file", key))
			}
		}(key, &wg)
	}

	// wait for all
	wg.Wait()

	return nil
}
