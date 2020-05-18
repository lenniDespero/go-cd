package config

import (
	"errors"
	"testing"

	"github.com/maraino/testify/require"

	"github.com/lenniDespero/go-cd/internal/pkg/host"
	"github.com/lenniDespero/go-cd/internal/pkg/pipe"
	"github.com/lenniDespero/go-cd/internal/pkg/target"
)

func prepareConfig() Config {
	return Config{
		ProjectName: "someName",
		Git:         "https://github.com/lenniDespero/go-cd.git",
		Count:       3,
		Targets: map[string]target.Config{
			"dev": {
				Type: "local",
				Host: host.Config{},
				Path: "newPath",
				Pipe: []pipe.Config{
					{Name: "some links", Type: "links", Args: []interface{}{map[string]interface{}{"from": "pathFrom", "to": "pathTo"}}},
					{Name: "some Cmd", Type: "command", Args: []interface{}{map[string]interface{}{"command": "ls", "args": []string{"ls"}}}},
				},
			},
		},
	}
}

func check(t *testing.T, config Config, e error) {
	err := config.CheckConfig()
	require.NotNil(t, err)
	require.Error(t, err, e)
	//if err != nil {
	//	if err != e {
	//		t.Errorf("Unexpected error: %s, expected: %s", err.Error(), e.Error())
	//	}
	//}
}

func TestConfig_CheckConfig(t *testing.T) {
	config := prepareConfig()
	err := config.CheckConfig()
	require.Nil(t, err)
}

func TestConfig_CheckConfigNoGit(t *testing.T) {
	config := prepareConfig()
	config.Git = ""
	check(t, config, ErrNoGit)
}

func TestConfig_CheckConfigNoName(t *testing.T) {
	config := prepareConfig()
	config.ProjectName = ""
	check(t, config, ErrNoName)
}

func TestConfig_CheckConfigNoCount(t *testing.T) {
	config := prepareConfig()
	config.Count = 0
	check(t, config, ErrNoCount)
}

func TestConfig_CheckConfigNoTargets(t *testing.T) {
	config := prepareConfig()
	config.Targets = map[string]target.Config{}
	check(t, config, ErrNoTargets)
}

func TestReadConfigErr(t *testing.T) {
	newErr := errors.New("can't read configuration file: Unsupported Config Type \"\"")
	_, err := ReadConfig("errorPath")
	if err.Error() != newErr.Error() {
		t.Errorf("Unexpected error : %s", err.Error())
	}
	newErr = errors.New("can't read configuration file: open errorPath.yml: no such file or directory")
	_, err = ReadConfig("errorPath.yml")
	if err.Error() != newErr.Error() {
		t.Errorf("Unexpected error : %s", err.Error())
	}
}

func TestReadConfig(t *testing.T) {
	_, err := ReadConfig("../../../testdata/test.app.yml")
	require.Nil(t, err)
}

func TestReadConfigEmpty(t *testing.T) {
	_, err := ReadConfig("../../../testdata/test.bad.empty.yml")
	require.NotNil(t, err)
}

func TestReadConfigBad(t *testing.T) {
	newErr := errors.New("can't read configuration file: While parsing config: yaml: line 10: did not find expected '-' indicator")
	_, err := ReadConfig("../../../testdata/test.bad.yml")
	require.Error(t, err, newErr)
}
