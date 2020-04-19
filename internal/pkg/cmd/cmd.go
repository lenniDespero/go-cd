package cmd

import (
	"os"
	"os/exec"
	"strings"
)

type Cmd struct {
	Command string   `mapstructure:"command" json:"command"`
	Args    []string `mapstructure:"args" json:"args"`
}

type Error string

func (e Error) Error() string {
	return string(e)
}

//Errors
const (
	NoCommandError Error = "no command"
)

func (command Cmd) CheckConfig() error {
	if command.Command == "" {
		return NoCommandError
	}
	return nil
}

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

func (command Cmd) GetRemoteCommand() string {
	comString := []string{command.Command}
	comString = append(comString, command.Args...)
	return strings.Join(comString, " ")
}
