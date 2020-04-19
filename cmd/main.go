package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"

	"github.com/lenniDespero/go-cd/internal/logger"

	"github.com/lenniDespero/go-cd/internal/pkg/config"
	"github.com/lenniDespero/go-cd/internal/pkg/target"

	"github.com/lenniDespero/go-cd/internal/deployer"

	"github.com/spf13/pflag"
)

func main() {
	var configPath = flag.String("config", "config/app.yml", "path to configuration flag")
	var deployTarget = flag.String("target", "", "target to deploy")
	var configTest = flag.Bool("test", false, "only test config")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	flag.Parse()
	var C config.Config
	err := config.ReadConfig(&C, *configPath)
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Info("Config test complete: all ok")
	if *configTest {
		logger.Info("Config test mode. Exit")
		os.Exit(0)
	}
	logger.Notice(fmt.Sprintf("Deployment script for project '%s'", C.ProjectName))
	if *deployTarget == "" {
		logger.Fatal("No target to deploy")
	}
	if reflect.DeepEqual(C.Targets[*deployTarget], target.Target{}) {
		logger.Fatal(fmt.Sprintf("Target '%s' not found in config.", *deployTarget))
	}
	jobber, err := deployer.GetDeployer(C, *deployTarget)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Deployer init error %s", err))
	}
	err = jobber.Run()
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Info("All done!")
}
