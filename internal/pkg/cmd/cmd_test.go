package cmd

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func prepareCmd() Cmd {
	return Cmd{
		Command: "ln",
		Args:    []string{"-s", "pathFrom", "PathTo"},
	}
}

func TestCmd_CheckConfig(t *testing.T) {
	cmd := prepareCmd()
	err := cmd.CheckConfig()
	if err != nil {
		t.Errorf("Unexpected error: %s, expected nil", err.Error())
	}
	cmd.Command = ""
	err = cmd.CheckConfig()
	if err != nil {
		if err != NoCommandError {
			t.Errorf("Unexpected error: %s, expected: %s", err.Error(), NoCommandError.Error())
		}
	} else {
		t.Errorf("Expected error, get nil")
	}
}

func TestCmd_ExecuteOnLocal(t *testing.T) {
	cmd := prepareCmd()
	dir, err := ioutil.TempDir("", "test-")
	if err != nil {
		t.Errorf("Unexpected error: %s, expected nil", err.Error())
	}
	defer func() {
		_ = os.RemoveAll(dir)
	}()
	err = os.Chdir(dir)
	if err != nil {
		t.Errorf("Unexpected error: %s, expected nil", err.Error())
	}
	err = cmd.ExecuteOnLocal()
	if err != nil {
		t.Errorf("Unexpected error: %s, expected nil", err.Error())
	}
	err = cmd.ExecuteOnLocal()
	if err == nil {
		t.Errorf("Expected error, get nil")
	} else if err.Error() != "exit status 1" {
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
