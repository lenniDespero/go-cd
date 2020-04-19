package pipe

import (
	"github.com/lenniDespero/go-cd/internal/pkg"
	"github.com/lenniDespero/go-cd/internal/pkg/cmd"
	"github.com/lenniDespero/go-cd/internal/pkg/link"
)

//Pipe struct for config
type Pipe struct {
	Name string        `mapstructure:"name"`
	Type string        `mapstructure:"type"`
	Args []interface{} `mapstructure:"args"`
}

//Types of pipe stages
var Types = map[string]bool{
	"links":   true,
	"command": true,
}

//Names of pipe stages
var Names = map[string]interface{}{"links": &link.Link{}, "command": &cmd.Cmd{}}

//NamesInt return ArgsInterface for pipe stage
var NamesInt = map[string]pkg.ArgsInterface{"links": &link.Link{}, "command": &cmd.Cmd{}}

//Error implementation for package
type Error string

//Error implementation for package
func (e Error) Error() string {
	return string(e)
}

//Errors
const (
	NoPipeTypeError Error = "no type in pipe"
	NotInPipesError Error = "not in legal types (link or command)"
	NoPipeName      Error = "no pipe name in pipes"
	NoPipeArgs      Error = "no args in pipe"
)

//CheckConfig will check config for errors
func (pipe Pipe) CheckConfig() error {
	if pipe.Name == "" {
		return NoPipeName
	}
	if pipe.Type == "" {
		return NoPipeTypeError
	}
	if !Types[pipe.Type] {
		return NotInPipesError
	}
	if len(pipe.Args) == 0 {
		return NoPipeArgs
	}
	return nil
}
