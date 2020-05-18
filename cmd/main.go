package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/pflag"

	"github.com/lenniDespero/go-cd/internal/deployer"
	"github.com/lenniDespero/go-cd/internal/logger"
	"github.com/lenniDespero/go-cd/internal/pkg/config"
)

var (
	configPath   = flag.String("config", "deploy.yml", "path to configuration flag")
	deployTarget = flag.String("target", "", "target to deploy")
	configTest   = flag.Bool("test", false, "only test config")
)

func main() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	flag.Parse()

	c, err := config.ReadConfig(*configPath)
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Info("config test complete: all ok")
	if *configTest {
		logger.Info("config test mode, exit")
		os.Exit(0)
	}

	logger.Notice(fmt.Sprintf("deployment script for project '%s'", c.ProjectName))
	if *deployTarget == "" {
		logger.Fatal("no target to deploy")
	}
	_, ok := c.Targets[*deployTarget]
	if !ok {
		logger.Fatal(fmt.Sprintf("target '%s' not found in config", *deployTarget))
	}

	jobber, err := deployer.NewDeployer(c, *deployTarget)
	if err != nil {
		logger.Fatal(fmt.Sprintf("deployer init error %s", err))
	}

	if err = jobber.Run(); err != nil {
		logger.Fatal(err.Error())
	}
	logger.Info("all done")
}
