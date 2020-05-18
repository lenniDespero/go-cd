package pipe

import (
	"errors"

	"github.com/lenniDespero/go-cd/internal/pkg"
	"github.com/lenniDespero/go-cd/internal/pkg/cmd"
	"github.com/lenniDespero/go-cd/internal/pkg/link"
)

// Config pipe struct for config
type Config struct {
	Name string        `mapstructure:"name"`
	Type string        `mapstructure:"type"`
	Args []interface{} `mapstructure:"args"`
}

// Types of pipe stages
var Types = map[string]bool{
	"links":   true,
	"command": true,
}

// Names of pipe stages
var Names = map[string]interface{}{"links": &link.Config{}, "command": &cmd.Config{}}

// NamesInt return ArgsInterface for pipe stage
var NamesInt = map[string]pkg.ArgsInterface{"links": &link.Config{}, "command": &cmd.Config{}}

// Errors
var (
	ErrNoPipeType = errors.New("no type in pipe")
	ErrNotInPipes = errors.New("not in legal types (link or command)")
	ErrNoPipeName = errors.New("no pipe name in pipes")
	ErrNoPipeArgs = errors.New("no args in pipe")
)

// CheckConfig will check config for errors
func (pipe Config) CheckConfig() error {
	if pipe.Name == "" {
		return ErrNoPipeName
	}
	if pipe.Type == "" {
		return ErrNoPipeType
	}
	if !Types[pipe.Type] {
		return ErrNotInPipes
	}
	if len(pipe.Args) == 0 {
		return ErrNoPipeArgs
	}
	return nil
}
