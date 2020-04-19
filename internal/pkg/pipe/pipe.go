package pipe

import (
	"github.com/lenniDespero/go-cd/internal/pkg"
	"github.com/lenniDespero/go-cd/internal/pkg/cmd"
	"github.com/lenniDespero/go-cd/internal/pkg/link"
)

type Pipe struct {
	Name string        `mapstructure:"name"`
	Type string        `mapstructure:"type"`
	Args []interface{} `mapstructure:"args"`
}

var Types = map[string]bool{
	"links":   true,
	"command": true,
}
var Names = map[string]interface{}{"links": &link.Link{}, "command": &cmd.Cmd{}}
var NamesInt = map[string]pkg.ArgsInterface{"links": &link.Link{}, "command": &cmd.Cmd{}}

type Error string

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
