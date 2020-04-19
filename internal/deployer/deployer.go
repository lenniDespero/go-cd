package deployer

import (
	"errors"

	"github.com/lenniDespero/go-cd/internal/pkg/config"

	"github.com/lenniDespero/go-cd/internal/deployer/host"

	"github.com/lenniDespero/go-cd/internal/deployer/local"
)

type DeployRunner struct {
	deployer DeployInterface
	git      string
	rCount   int
}

func (d DeployRunner) Run() error {
	err := d.deployer.Prepare()
	if err != nil {
		return err
	}
	defer func() {
		ferr := d.deployer.CleanUp(d.rCount)
		if ferr != nil {
			err = ferr
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

func GetDeployer(config config.Config, target string) (DeployRunner, error) {
	d := DeployRunner{}
	d.rCount = config.Count
	d.git = config.Git
	var err error
	if config.Targets[target].Type == "local" {
		d.deployer, err = local.InitDeployer(config.Targets[target])
		return d, err
	} else if config.Targets[target].Type == "host" {
		d.deployer, err = host.InitDeployer(config.Targets[target])
		return d, err
	}

	return d, errors.New("unknown deploy runner")
}
