package cmd

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/maraino/testify/require"
)

func prepareCmd() Config {
	return Config{
		Command: "ln",
		Args:    []string{"-s", "pathFrom", "PathTo"},
	}
}

func TestCmd_CheckConfig(t *testing.T) {
	cmd := prepareCmd()
	err := cmd.CheckConfig()
	require.Nil(t, err)
	cmd.Command = ""
	err = cmd.CheckConfig()
	require.NotNil(t, err)
	require.Error(t, err, ErrNoCommand)
}

func TestCmd_ExecuteOnLocal(t *testing.T) {
	cmd := prepareCmd()
	dir, err := ioutil.TempDir("", "test-")
	require.Nil(t, err)
	defer func() {
		_ = os.RemoveAll(dir)
	}()
	err = os.Chdir(dir)
	require.Nil(t, err)
	err = cmd.ExecuteOnLocal()
	require.Nil(t, err)
	err = cmd.ExecuteOnLocal()
	require.NotNil(t, err)
	if err.Error() != "exit status 1" {
		t.Errorf("Unexpected error: %s, expected exit status 1", err.Error())
	}
}

func TestCmd_GetRemoteCommand(t *testing.T) {
	cmd := prepareCmd()
	str := cmd.GetRemoteCommand()
	s := []string{"ln"}
	s = append(s, cmd.Args...)
	real := strings.Join(s, " ")
	if str != real {
		t.Errorf("Wrong remote command string, get: %s, expected :%s", str, real)
	}
}
