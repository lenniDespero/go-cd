package models

type Cmd struct {
	Command string   `mapstructure:"command"`
	Args    []string `mapstructure:"args"`
}

type CmdError string

func (e CmdError) Error() string {
	return string(e)
}

//Errors
const (
	NoCommandError CmdError = "no command"
)

func (command Cmd) checkArgsConfig() error {
	if command.Command == "" {
		return NoCommandError
	}
	return nil
}
