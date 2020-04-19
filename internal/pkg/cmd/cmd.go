package cmd

import (
	"os"
	"os/exec"
	"strings"
)

//Cmd part of config
type Cmd struct {
	Command string   `mapstructure:"command" json:"command"`
	Args    []string `mapstructure:"args" json:"args"`
}

//Error type for package cmd
type Error string

//Error interface implementation
func (e Error) Error() string {
	return string(e)
}

//Errors
const (
	NoCommandError Error = "no command"
)

//CheckConfig will check config for errors
func (command Cmd) CheckConfig() error {
	if command.Command == "" {
		return NoCommandError
	}
	return nil
}

//ExecuteOnLocal run cmd on local target
func (command Cmd) ExecuteOnLocal() error {
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
func (command Cmd) GetRemoteCommand() string {
	comString := []string{command.Command}
	comString = append(comString, command.Args...)
	return strings.Join(comString, " ")
}
