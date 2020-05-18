package target

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"

	"github.com/lenniDespero/go-cd/internal/pkg/host"
	"github.com/lenniDespero/go-cd/internal/pkg/pipe"
)

// Config target struct for config
type Config struct {
	Type string        `mapstructure:"type"`
	Host host.Config   `mapstructure:"host"`
	Path string        `mapstructure:"path"`
	Pipe []pipe.Config `mapstructure:"pipe"`
}

// Types of targets
var Types = map[string]bool{
	"local": true,
	"host":  true,
}

const (
	TypeLocal string = "local"
	TypeHost  string = "host"
)

// Errors
var (
	ErrNoType     = errors.New("no type in target")
	ErrNotInTypes = errors.New("not in legal types (local or host)")
	ErrNoHost     = errors.New("no host in target")
	ErrNoPath     = errors.New("no path to deploy in target")
	ErrNoPipes    = errors.New("no pipes in target")
)

// CheckConfig will check config for errors
func (target Config) CheckConfig() error {
	if target.Type == "" {
		return ErrNoType
	}
	if !Types[target.Type] {
		return errors.Wrap(ErrNotInTypes, fmt.Sprintf("type %s ", target.Type))
	}
	if target.Type == "host" {
		if reflect.DeepEqual(target.Host, host.Config{}) {
			return ErrNoHost
		}
		err := target.Host.CheckConfig()
		if err != nil {
			return errors.Wrap(err, "error in target.host")
		}
	}
	if target.Path == "" {
		return ErrNoPath
	}
	if len(target.Pipe) == 0 {
		return ErrNoPipes
	}
	for _, p := range target.Pipe {
		err := p.CheckConfig()
		if err != nil {
			return err
		}
		inter, ok := pipe.Names[p.Type]
		if !ok {
			return ErrNotInTypes
		}
		for _, args := range p.Args {
			err := mapstructure.Decode(args, inter)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("error on decode pipe args in %s -> %s ", target.Type, p.Name))
			}
			jsonInter, err := json.Marshal(inter)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("error on marshal in %s -> %s ", target.Type, p.Name))
			}
			pipeint, ok := pipe.NamesInt[p.Type]
			if !ok {
				return errors.New(fmt.Sprintf("unexpected error on %s -> %s", target.Type, p.Name))
			}
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
