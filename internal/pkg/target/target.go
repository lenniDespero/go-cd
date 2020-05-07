package target

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/lenniDespero/go-cd/internal/pkg/host"
	"github.com/lenniDespero/go-cd/internal/pkg/pipe"

	"github.com/mitchellh/mapstructure"

	"github.com/pkg/errors"
)

//Target struct for config
type Target struct {
	Type string      `mapstructure:"type"`
	Host host.Host   `mapstructure:"host"`
	Path string      `mapstructure:"path"`
	Pipe []pipe.Pipe `mapstructure:"pipe"`
}

//Types of targets
var Types = map[string]bool{
	"local": true,
	"host":  true,
}

//Error implementation for package
type Error string

//Error implementation for package
func (e Error) Error() string {
	return string(e)
}

//Errors
const (
	NoTypeError     Error = "no type in target"
	NotInTypesError Error = "not in legal types (local or host)"
	NoHostError     Error = "no host in target"
	NoPathError     Error = "no path to deploy in target"
	NoPipesError    Error = "no pipes in target"
)

//CheckConfig will check config for errors
func (target Target) CheckConfig() error {
	if target.Type == "" {
		return NoTypeError
	}
	if !Types[target.Type] {
		return errors.Wrap(NotInTypesError, fmt.Sprintf("type %s ", target.Type))
	}
	if target.Type == "host" {
		if reflect.DeepEqual(target.Host, host.Host{}) {
			return NoHostError
		}
		err := target.Host.CheckConfig()
		if err != nil {
			return errors.Wrap(err, "error in target.host")
		}
	}
	if target.Path == "" {
		return NoPathError
	}
	if len(target.Pipe) == 0 {
		return NoPipesError
	}
	for _, p := range target.Pipe {
		err := p.CheckConfig()
		if err != nil {
			return err
		}
		inter := pipe.Names[p.Type]
		for _, args := range p.Args {
			err := mapstructure.Decode(args, inter)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("error on decode pipe args in %s -> %s ", target.Type, p.Name))
			}
			jsonInter, err := json.Marshal(inter)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("error on marshal in %s -> %s ", target.Type, p.Name))
			}
			pipeint := pipe.NamesInt[p.Type]
			if err := json.Unmarshal(jsonInter, &pipeint); err != nil {
				return errors.Wrap(err, fmt.Sprintf("error on unmarshal in %s -> %s ", target.Type, p.Name))
			}
			err = pipeint.CheckConfig()
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("error on check pipe args in %s -> %s ", target.Type, p.Name))
			}
		}
	}
	return nil
}
