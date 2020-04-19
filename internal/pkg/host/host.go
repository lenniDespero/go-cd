package host

import "fmt"

type Host struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Auth     string `mapstructure:"auth"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Key      string `mapstructure:"key"`
}

var Types = map[string]bool{
	"key":      true,
	"password": true,
}

type Error string

func (e Error) Error() string {
	return string(e)
}

//Errors
const (
	NoHostError        Error = "no host in config"
	NoAuthError        Error = "no auth type in config"
	WrongAuthTypeError Error = "wrong auth type in config"
	NoUserError        Error = "no user in config"
	NoPasswordError    Error = "no password in config"
	NoKeyError         Error = "no key in config"
)

func (h Host) CheckConfig() error {
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

func (h Host) GetConnectionString() string {
	port := "22"
	if h.Port != "" {
		port = h.Port
	}
	connectionString := fmt.Sprintf("%s:%s", h.Host, port)
	return connectionString
}
