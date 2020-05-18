package cmd

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

// Config cmd part of config
type Config struct {
	Command string   `mapstructure:"command" json:"command"`
	Args    []string `mapstructure:"args" json:"args"`
}

//Errors
var (
	ErrNoCommand = errors.New("no command")
)

//CheckConfig will check config for errors
func (command Config) CheckConfig() error {
	if command.Command == "" {
		return ErrNoCommand
	}
	return nil
}

//ExecuteOnLocal run cmd on local target
func (command Config) ExecuteOnLocal() error {
	cmdString := command.Command
	var args []string
	args = append(args, command.Args...)
	cmd := exec.Command(cmdString, args...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

//GetRemoteCommand prepare cmd string for remote machine
func (command Config) GetRemoteCommand() string {
	comString := []string{command.Command}
	comString = append(comString, command.Args...)
	return strings.Join(comString, " ")
}
