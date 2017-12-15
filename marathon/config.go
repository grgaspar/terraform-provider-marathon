package marathon

import (
	"time"

	"github.com/gambol99/go-marathon"
)

type config struct {
	config                   marathon.Config
	Client                   marathon.Marathon
	DefaultDeploymentTimeout time.Duration
	DcosURL                  string
}

func (c *config) loadAndValidate() error {
	client, err := marathon.NewClient(c.config)
	c.Client = client
	return err
}
