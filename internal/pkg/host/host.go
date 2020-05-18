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

const (
	TypeKey      string = "key"
	TypePassword string = "password"
)

// Errors
var (
	ErrNoHost        = errors.New("no host in config")
	ErrNoAuth        = errors.New("no auth type in config")
	ErrWrongAuthType = errors.New("wrong auth type in config")
	ErrNoUser        = errors.New("no user in config")
	ErrNoPassword    = errors.New("no password in config")
	ErrNoKey         = errors.New("no key in config")
)

// CheckConfig will check config for errors
func (h Config) CheckConfig() error {
	if h.Host == "" {
		return ErrNoHost
	}
	if h.Auth == "" {
		return ErrNoAuth
	}
	if h.User == "" {
		return ErrNoUser
	}
	switch h.Auth {
	case TypePassword:
		if h.Password == "" {
			return ErrNoPassword
		}
	case TypeKey:
		if h.Key == "" {
			return ErrNoKey
		}
	default:
		if !Types[h.Auth] {
			return ErrWrongAuthType
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
