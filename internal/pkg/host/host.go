package host

import (
	"errors"
	"fmt"
)

// Config host structure for remote machines
type Config struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Auth     string `mapstructure:"auth"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Key      string `mapstructure:"key"`
}

// Types of host auth
var Types = map[string]bool{
	"key":      true,
	"password": true,
}

// Errors
var (
	NoHostError        = errors.New("no host in config")
	NoAuthError        = errors.New("no auth type in config")
	WrongAuthTypeError = errors.New("wrong auth type in config")
	NoUserError        = errors.New("no user in config")
	NoPasswordError    = errors.New("no password in config")
	NoKeyError         = errors.New("no key in config")
)

// CheckConfig will check config for errors
func (h Config) CheckConfig() error {
	if h.Host == "" {
		return NoHostError
	}
	if h.Auth == "" {
		return NoAuthError
	}
	if h.User == "" {
		return NoUserError
	}
	switch h.Auth {
	case "password":
		if h.Password == "" {
			return NoPasswordError
		}
	case "key":
		if h.Key == "" {
			return NoKeyError
		}
	default:
		if !Types[h.Auth] {
			return WrongAuthTypeError
		}
	}
	return nil
}

// GetConnectionString prepare connection string to host
func (h Config) GetConnectionString() string {
	port := "22"
	if h.Port != "" {
		port = h.Port
	}
	connectionString := fmt.Sprintf("%s:%s", h.Host, port)
	return connectionString
}
