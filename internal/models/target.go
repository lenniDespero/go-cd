package models

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

type Target struct {
	Type string `mapstructure:"type"`
	Host Host   `mapstructure:"host"`
	Path string `mapstructure:"path"`
	Pipe []Pipe `mapstructure:"pipe"`
}

var TargetTypes = map[string]bool{
	"local": true,
	"host":  true,
}

type TargetError string

func (e TargetError) Error() string {
	return string(e)
}

//Errors
const (
	NoTypeError     TargetError = "no type in target"
	NotInTypesError TargetError = "not in legal types (local or host)"
	NoHostError     TargetError = "no host in target"
	NoPathError     TargetError = "no path to deploy in target"
	NoPipesError    TargetError = "no pipes in target"
)

func (target Target) checkTargetConfig() error {
	if target.Type == "" {
		return NoTypeError
	}
	if TargetTypes[target.Type] == false {
		return errors.Wrap(NotInTypesError, fmt.Sprintf("type %s ", target.Type))
	}
	if target.Type == "host" && reflect.DeepEqual(target.Host, Host{}) {
		return NoHostError
	}
	if target.Path == "" {
		return NoPathError
	}
	if len(target.Pipe) == 0 {
		return NoPipesError
	}
	for _, pipe := range target.Pipe {
		err := pipe.checkPipeConfig()
		if err != nil {
			return err
		}
		//inter := PipeNames[pipe.Type]
		//for _, args := range pipe.Args {
		//	err := mapstructure.Decode(args, inter)
		//	if err != nil {
		//		return errors.Wrap(NotInTypesError, fmt.Sprintf("error on decode pipe args in %s -> %s ", target.Type, pipe.Name))
		//	}
		//	//inter.checkArgsConfig()
		//	fmt.Printf("%T", inter)
		//}
	}
	return nil
}
