package cmd

import (
	"sync"

	"github.com/f24-cse535/pbft/internal/config"

	"go.dedis.ch/kyber/v4/pairing/bn256"
	"go.dedis.ch/kyber/v4/share"
	"go.uber.org/zap"
)

// Controller is a process for creating system's nodes.
type Controller struct {
	Cfg    config.Config
	Logger *zap.Logger
}

func (c Controller) Main() error {
	// init a waitgroup
	var wg sync.WaitGroup

	// create threshold signature keys
	suite := bn256.NewSuite()
	n := c.Cfg.Node.BFT.Total
	t := c.Cfg.Node.BFT.Majority

	// generate the shared secret and ploynomial
	secret := suite.G1().Scalar().Pick(suite.RandomStream())
	priPoly := share.NewPriPoly(suite.G2(), t, secret, suite.RandomStream())
	pubPoly := priPoly.Commit(suite.G2().Point().Base())

	// loop over config files and run each node
	for index, key := range c.Cfg.CtlFiles {
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
				suite:  suite,
				share:  priPoly.Shares(n)[index],
				pub:    pubPoly,
			}

			c.Logger.Info("instance started", zap.String("instance", cfg.Node.NodeId), zap.String("file", path))

			if err := node.Main(); err != nil {
				c.Logger.Error("failed to start instance", zap.String("instance", cfg.Node.NodeId), zap.Error(err))
			} else {
				c.Logger.Info("instance terminated", zap.String("instance", cfg.Node.NodeId), zap.String("file", key))
			}
		}(key, &wg)
	}

	// wait for all instances to terminate
	wg.Wait()

	return nil
}
