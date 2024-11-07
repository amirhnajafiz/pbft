package cmd

import (
	"sync"

	"github.com/f24-cse535/pbft/internal/config"

	"go.uber.org/zap"
)

// CTl is an app for creating system's nodes.
type CTL struct {
	Cfg    config.Config
	Logger *zap.Logger
}

func (c CTL) Main() error {
	// init a waitgroup
	var wg sync.WaitGroup

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
