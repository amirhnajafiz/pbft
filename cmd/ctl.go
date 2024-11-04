package cmd

import (
	"strings"
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

			// get the instance name from path
			parts := strings.Split(path, "/")
			name := parts[len(parts)-2]

			// create a new node instance
			cfg := config.New(path, false)
			node := Node{
				Cfg:    cfg,
				Logger: c.Logger.Named("node." + name),
			}

			c.Logger.Info("instance started", zap.String("instance", name), zap.String("file", path))

			if err := node.Main(); err != nil {
				c.Logger.Error("failed to start instance", zap.String("instance", name), zap.Error(err))
			} else {
				c.Logger.Info("instance terminated", zap.String("instance", name), zap.String("file", key))
			}
		}(key, &wg)
	}

	// wait for all
	wg.Wait()

	return nil
}
