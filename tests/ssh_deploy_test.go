// +build integration

package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/lenniDespero/go-cd/internal/ssh"

	"github.com/cucumber/godog"
)

func (c *ConfTest) iRunDeployer() error {
	args := []string{"--config", c.configFile, "-target", "tests"}
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

func (c *ConfTest) exitCodeWillBeZero() error {
	if c.exitCode != 0 {
		return fmt.Errorf("unexpected exit code %d", c.exitCode)
	}
	return nil
}

func (c *ConfTest) iHaveFoldersOnRemoteServer() error {
	client, err := ssh.DialWithPasswd(c.config.Host.GetConnectionString(), c.config.Host.User, c.config.Host.Password)
	if err != nil {
		return err
	}
	script := client.Cmd("ls " + c.config.Path + "current")
	err = script.Run()
	if err != nil {
		return err
	}
	return nil
}

func SSHDeployFeatureContext(s *godog.Suite) {
	conf := getConfig()
	s.Step(`^I run Deployer$`, conf.iRunDeployer)
	s.Step(`^Exit code will be zero$`, conf.exitCodeWillBeZero)
	s.Step(`^I have folders on remote server$`, conf.iHaveFoldersOnRemoteServer)
}
