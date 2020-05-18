package config

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/lenniDespero/go-cd/internal/pkg/target"
)

// Config - base config of deploy app
type Config struct {
	ProjectName string                   `mapstructure:"projectName"`
	Git         string                   `mapstructure:"git"`
	Count       int                      `mapstructure:"count"`
	Targets     map[string]target.Config `mapstructure:"targets"`
}

// Errors
var (
	ErrNoName    = errors.New("no project name in config")
	ErrNoGit     = errors.New("no git in config")
	ErrNoCount   = errors.New("no count releases in config")
	ErrNoTargets = errors.New("no targets in config")
)

// ReadConfig from file and check for errors
func ReadConfig(configPath string) (Config, error) {
	c := Config{}
	viper.SetConfigFile(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		return c, errors.Wrap(err, "can't read configuration file")
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	err = viper.Unmarshal(&c)
	if err != nil {
		return c, err
	}
	return c, c.CheckConfig()
}

// CheckConfig will check config for errors
func (c Config) CheckConfig() error {
	if c.ProjectName == "" {
		return ErrNoName
	}
	if c.Git == "" {
		return ErrNoGit
	}
	if c.Count < 1 {
		return ErrNoCount
	}
	if len(c.Targets) == 0 {
		return ErrNoTargets
	}
	for name, target := range c.Targets {
		err := target.CheckConfig()
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("in target %s ", name))
		}
	}
	return nil
}
