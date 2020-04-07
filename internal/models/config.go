package models

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	ProjectName string            `mapstructure:"projectName"`
	Git         string            `mapstructure:"git"`
	Count       int               `mapstructure:"count"`
	Targets     map[string]Target `mapstructure:"targets"`
}

type ConfigError string

func (e ConfigError) Error() string {
	return string(e)
}

//Errors
const (
	NoNameError    ConfigError = "no project name in config"
	NoGitError     ConfigError = "no git in config"
	NoCountError   ConfigError = "no count releases in config"
	NoTargetsError ConfigError = "no targets in config"
)

func ReadConfig(C *Config) {
	viper.SetConfigName("app")
	viper.AddConfigPath("./cmd")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Couldn't read configuration file: %s", err.Error())
	}
	err := viper.Unmarshal(&C)
	if err != nil {
		log.Fatal(err)
	}
	err = C.checkConfig()
	if err != nil {
		log.Fatal(err)
	}
}

func (C Config) checkConfig() error {
	if C.ProjectName == "" {
		return NoNameError
	}
	if C.Git == "" {
		return NoGitError
	}
	if C.Count == 0 {
		return NoCountError
	}
	if len(C.Targets) == 0 {
		return NoTargetsError
	}
	for name, target := range C.Targets {
		err := target.checkTargetConfig()
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("in target %s ", name))
		}
	}
	log.Println("Config test complete: all ok")
	return nil
}
