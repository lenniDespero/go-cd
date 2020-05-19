//+build integration

package main

import (
	"errors"
	"os"
	"os/exec"

	"github.com/cucumber/godog"
)

func (c *ConfTest) iRunDeployerForTest() error {
	args := []string{"--config", c.configFile, "-test"}
	cmd := exec.Command("../bin/deployer", args...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			c.exitCode = exitError.ExitCode()
		}
	}
	return nil
}

func (c *ConfTest) iRunDeployerWithoutTarget() error {
	args := []string{"--config", c.configFile}
	cmd := exec.Command("../bin/deployer", args...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			c.exitCode = exitError.ExitCode()
		}
	}
	return nil
}

func (c *ConfTest) exitCodeWillNotZero() error {
	if c.exitCode == 0 {
		return errors.New("Expected error exit code")
	}
	c.exitCode = 0
	return nil
}

func (c *ConfTest) iRunDeployerWithWrongTarget() error {
	args := []string{"--config", c.configFile, "-target", "strange"}
	cmd := exec.Command("../bin/deployer", args...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			c.exitCode = exitError.ExitCode()
		}
	}
	return nil
}

func DeployerOptionsFeatureContext(s *godog.Suite) {
	c := getConfig()
	s.Step(`^I run Deployer for test$`, c.iRunDeployerForTest)
	s.Step(`^Exit code will be zero$`, c.exitCodeWillBeZero)
	s.Step(`^I run Deployer without target$`, c.iRunDeployerWithoutTarget)
	s.Step(`^Exit code will not zero$`, c.exitCodeWillNotZero)
	s.Step(`^I run Deployer with wrong target$`, c.iRunDeployerWithWrongTarget)
	s.Step(`^Exit code will not zero$`, c.exitCodeWillNotZero)
}
