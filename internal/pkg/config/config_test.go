package config

import (
	"errors"
	"testing"

	"github.com/lenniDespero/go-cd/internal/pkg/host"
	"github.com/lenniDespero/go-cd/internal/pkg/pipe"
	"github.com/lenniDespero/go-cd/internal/pkg/target"
)

func prepareConfig() Config {
	return Config{
		ProjectName: "someName",
		Git:         "https://github.com/lenniDespero/go-cd.git",
		Count:       3,
		Targets: map[string]target.Target{
			"dev": {
				Type: "local",
				Host: host.Host{},
				Path: "newPath",
				Pipe: []pipe.Pipe{
					{Name: "some links", Type: "links", Args: []interface{}{map[string]interface{}{"from": "pathFrom", "to": "pathTo"}}},
					{Name: "some Cmd", Type: "command", Args: []interface{}{map[string]interface{}{"command": "ls", "args": []string{"ls"}}}},
				},
			},
		},
	}
}

func check(t *testing.T, config Config, e error) {
	err := config.CheckConfig()
	if err != nil {
		if err != e {
			t.Errorf("Unexpected error: %s, expected: %s", err.Error(), e.Error())
		}
	} else {
		t.Errorf("Expected error, get nil")
	}
}

func TestConfig_CheckConfig(t *testing.T) {
	config := prepareConfig()
	err := config.CheckConfig()
	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
	}
}

func TestConfig_CheckConfigNoGit(t *testing.T) {
	config := prepareConfig()
	config.Git = ""
	check(t, config, NoGitError)
}

func TestConfig_CheckConfigNoName(t *testing.T) {
	config := prepareConfig()
	config.ProjectName = ""
	check(t, config, NoNameError)
}

func TestConfig_CheckConfigNoCount(t *testing.T) {
	config := prepareConfig()
	config.Count = 0
	check(t, config, NoCountError)
}

func TestConfig_CheckConfigNoTargets(t *testing.T) {
	config := prepareConfig()
	config.Targets = map[string]target.Target{}
	check(t, config, NoTargetsError)
}

func TestReadConfigErr(t *testing.T) {
	c := Config{}
	newErr := errors.New("can't read configuration file: Unsupported Config Type \"\"")
	err := ReadConfig(&c, "errorPath")
	if err.Error() != newErr.Error() {
		t.Errorf("Unexpected error : %s", err.Error())
	}
	newErr = errors.New("can't read configuration file: open errorPath.yml: no such file or directory")
	err = ReadConfig(&c, "errorPath.yml")
	if err.Error() != newErr.Error() {
		t.Errorf("Unexpected error : %s", err.Error())
	}
}

func TestReadConfig(t *testing.T) {
	c := Config{}
	err := ReadConfig(&c, "../../../testdata/test.app.yml")
	if err != nil {
		t.Errorf("Unexpected error : %s", err.Error())
	}
}

func TestReadConfigEmpty(t *testing.T) {
	c := Config{}
	err := ReadConfig(&c, "../../../testdata/test.bad.empty.yml")
	if err == nil {
		t.Errorf("Expected err, get nil")
	}
}

func TestReadConfigBad(t *testing.T) {
	c := Config{}
	newErr := errors.New("can't read configuration file: While parsing config: yaml: line 10: did not find expected '-' indicator")
	err := ReadConfig(&c, "../../../testdata/test.bad.yml")
	if err != nil {
		if err.Error() != newErr.Error() {
			t.Errorf("Unexpected error : %s", err.Error())
		}
	}
}
