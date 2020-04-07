package models

type Pipe struct {
	Name string        `mapstructure:"name"`
	Type string        `mapstructure:"type"`
	Args []interface{} `mapstructure:"args"`
}

var Links Link
var Cmds Cmd

var PypeTypes = map[string]bool{
	"links":   true,
	"command": true,
}
var PipeNames = map[string]interface{}{"links": &Link{}, "command": &Cmd{}}
var PipeNamesInt = map[string]ArgsInterface{"links": Links, "command": &Cmds}

type PipeError string

func (e PipeError) Error() string {
	return string(e)
}

//Errors
const (
	NoPipeTypeError PipeError = "no type in pipe"
	NotInPipesError PipeError = "not in legal types (link or command)"
	NoPipeName      PipeError = "no pipe name in pipes"
	NoPipeArgs      PipeError = "no args in pipe"
)

func (pipe Pipe) checkPipeConfig() error {
	if pipe.Name == "" {
		return NoPipeName
	}
	if pipe.Type == "" {
		return NoPipeTypeError
	}
	if PypeTypes[pipe.Type] == false {
		return NotInPipesError
	}
	if len(pipe.Args) == 0 {
		return NoPipeArgs
	}
	return nil
}
