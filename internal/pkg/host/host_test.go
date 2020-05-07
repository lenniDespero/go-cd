package host

import (
	"strings"
	"testing"
)

func prepareHost() Host {
	return Host{
		Host:     "someHost",
		Port:     "123",
		User:     "username",
		Password: "pass",
		Auth:     "password",
		Key:      ".somekey.key",
	}
}

func check(t *testing.T, host Host, e Error) {
	err := host.CheckConfig()
	if err != nil {
		if err != e {
			t.Errorf("Unexpected error: %s, expected: %s", err.Error(), e.Error())
		}
	} else {
		t.Errorf("Expected error, get nil")
	}
}

func TestHost_GetConnectionStringWithPort(t *testing.T) {
	host := Host{
		Host: "someHostname",
		Port: "1234",
	}
	str := host.GetConnectionString()
	real := strings.Join([]string{host.Host, host.Port}, ":")
	if str != real {
		t.Errorf("Wrong connection string, get: %s, expected :%s", str, real)
	}
}

func TestHost_GetConnectionStringWithoutPort(t *testing.T) {
	host := Host{
		Host: "someHostname",
	}
	str := host.GetConnectionString()
	real := strings.Join([]string{host.Host, "22"}, ":")
	if str != real {
		t.Errorf("Wrong connection string, get: %s, expected :%s", str, real)
	}
}

func TestHost_CheckConfig(t *testing.T) {
	host := prepareHost()
	err := host.CheckConfig()
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
}

func TestHost_CheckConfigNoHost(t *testing.T) {
	host := prepareHost()
	host.Host = ""
	check(t, host, NoHostError)
}

func TestHost_CheckConfigNoAuth(t *testing.T) {
	host := prepareHost()
	host.Auth = ""
	check(t, host, NoAuthError)
}

func TestHost_CheckConfigNoPass(t *testing.T) {
	host := prepareHost()
	host.Password = ""
	check(t, host, NoPasswordError)
}

func TestHost_CheckConfigNoKey(t *testing.T) {
	host := prepareHost()
	host.Key = ""
	host.Auth = "key"
	check(t, host, NoKeyError)
}

func TestHost_CheckConfigWrongType(t *testing.T) {
	host := prepareHost()
	host.Auth = "strange"
	check(t, host, WrongAuthTypeError)
}

func TestHost_CheckConfigNoUser(t *testing.T) {
	host := prepareHost()
	host.User = ""
	check(t, host, NoUserError)
}
