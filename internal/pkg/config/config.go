package config

import "C"
import (
	"fmt"
	"strings"

	"github.com/lenniDespero/go-cd/internal/pkg/target"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	ProjectName string                   `mapstructure:"projectName"`
	Git         string                   `mapstructure:"git"`
	Count       int                      `mapstructure:"count"`
	Targets     map[string]target.Target `mapstructure:"targets"`
}

type Error string

func (e Error) Error() string {
	return string(e)
}

//Errors
const (
	NoNameError    Error = "no project name in config"
	NoGitError     Error = "no git in config"
	NoCountError   Error = "no count releases in config"
	NoTargetsError Error = "no targets in config"
)

func ReadConfig(c *Config, configPath string) error {
	viper.SetConfigFile(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, "can't read configuration file")
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	err = viper.Unmarshal(&c)
	if err != nil {
		return err
	}
	err = c.CheckConfig()
	if err != nil {
		return err
	}
	return nil
}

func (c Config) CheckConfig() error {
	if c.ProjectName == "" {
		return NoNameError
	}
	if c.Git == "" {
		return NoGitError
	}
	if c.Count == 0 {
		return NoCountError
	}
	if len(c.Targets) == 0 {
		return NoTargetsError
	}
	for name, target := range c.Targets {
		err := target.CheckConfig()
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("in target %s ", name))
		}
	}
	return nil
}
