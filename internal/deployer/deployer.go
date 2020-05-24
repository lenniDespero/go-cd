package deployer

import (
	"errors"

	"github.com/lenniDespero/go-cd/internal/logger"
	"github.com/lenniDespero/go-cd/internal/pkg/config"
	"github.com/lenniDespero/go-cd/internal/pkg/target"
)

// DeployRunner stores git, count releases to keep and deployer
type Deployer struct {
	deployer DeployInterface
	git      string
	rCount   int
}

// Run deploy stages
func (d Deployer) Run() error {
	err := d.deployer.Prepare()
	if err != nil {
		return err
	}
	defer func() {
		ferr := d.deployer.CleanUp(d.rCount)
		if ferr != nil {
			logger.Warn(ferr.Error())
		}
	}()

	err = d.deployer.UpdateSource(d.git)
	if err != nil {
		return err
	}
	err = d.deployer.RunPipe()
	if err != nil {
		return err
	}
	err = d.deployer.MakeLinks()
	if err != nil {
		return err
	}

	return nil
}

// NewDeployer return Deployer
func NewDeployer(config config.Config, targetString string) (Deployer, error) {
	d := Deployer{rCount: config.Count, git: config.Git}

	var err error
	switch config.Targets[targetString].Type {
	case target.TypeLocal:
		d.deployer, err = NewLocalDeployer(config.Targets[targetString])
		return d, err
	case target.TypeHost:
		d.deployer, err = NewHosDeployer(config.Targets[targetString])
		return d, err
	default:
		return d, errors.New("unknown deploy runner")
	}
}
