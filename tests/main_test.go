// +build integration

package main

import (
	"os"
	"testing"

	"github.com/cucumber/godog"

	"github.com/lenniDespero/go-cd/internal/pkg/config"
	"github.com/lenniDespero/go-cd/internal/pkg/target"
)

var instance *ConfTest

type ConfTest struct {
	configFile string
	config     target.Config
	exitCode   int
}

func TestMain(m *testing.M) {
	status := godog.RunWithOptions("integration", func(s *godog.Suite) {
		godog.SuiteContext(s)
		SSHDeployFeatureContext(s)
		DeployerOptionsFeatureContext(s)
	}, godog.Options{
		Format:      "pretty",
		Paths:       []string{"features"},
		Randomize:   0,
		Concurrency: 0,
	})

	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}

func getConfig() *ConfTest {
	if instance == nil {
		conf := ConfTest{configFile: "../testdata/test.app.yml"}
		c, err := config.ReadConfig(conf.configFile)
		if err != nil {
			os.Exit(1)
		}
		conf.config = c.Targets["tests"]
		instance = &conf
		return instance
	}
	return instance
}
