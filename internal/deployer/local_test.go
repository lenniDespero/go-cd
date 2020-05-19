package deployer

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/maraino/testify/require"

	"github.com/lenniDespero/go-cd/internal/pkg/config"
)

func prepareConfig(t *testing.T) *config.Config {
	c, err := config.ReadConfig("../../testdata/test.app.yml")
	require.Nil(t, err)
	return &c
}

func getDeployer(t *testing.T) *LocalDeployer {
	conf := prepareConfig(t)
	deployer := &LocalDeployer{conf: conf.Targets["devel"]}
	return deployer
}

func (l *LocalDeployer) removeLock(t *testing.T) {
	err := os.RemoveAll(filepath.Join(l.absPth, ".lock"))
	require.Nil(t, err)
}

func TestInitDeployer(t *testing.T) {
	conf := prepareConfig(t)
	_, err := NewLocalDeployer(conf.Targets["devel"])
	require.Nil(t, err)
}

func TestDeployLocal_PrepareFail(t *testing.T) {
	deployer := getDeployer(t)
	err := deployer.Prepare()
	require.Nil(t, err)
	err = deployer.Prepare()
	require.NotNil(t, err)
	deployer.removeLock(t)
}

func TestDeployLocal_PrepareFailPath(t *testing.T) {
	deployer := getDeployer(t)
	deployer.conf.Path = "/testdata/releases"
	err := deployer.Prepare()
	require.NotNil(t, err)
	deployer.removeLock(t)
}

func TestDeployLocal_UpdateSourceBadGit(t *testing.T) {
	deployer := getDeployer(t)
	err := deployer.Prepare()
	require.Nil(t, err)
	err = deployer.UpdateSource("fsuegjsehfgjhseb")
	require.NotNil(t, err)
	deployer.removeLock(t)
}

func TestDeployLocal_RunPipeWrongPath(t *testing.T) {
	conf := prepareConfig(t)
	deployer := getDeployer(t)
	err := deployer.Prepare()
	require.Nil(t, err)
	err = deployer.UpdateSource(conf.Git)
	require.Nil(t, err)
	deployer.timeNamePath += "syeugfsjgefgsef"
	err = deployer.RunPipe()
	require.NotNil(t, err)
	deployer.removeLock(t)
}

func TestDeployLocal_CleanUp(t *testing.T) {
	conf := prepareConfig(t)
	deployer := getDeployer(t)
	err := deployer.Prepare()
	require.Nil(t, err)
	err = deployer.UpdateSource(conf.Git)
	require.Nil(t, err)
	err = deployer.RunPipe()
	require.Nil(t, err)
	err = deployer.MakeLinks()
	require.Nil(t, err)
	err = deployer.CleanUp(conf.Count)
	require.Nil(t, err)

	files, err := ioutil.ReadDir(deployer.absPth)
	require.Nil(t, err)
	var folders []string
	for _, rec := range files {
		if rec.IsDir() {
			folders = append(folders, rec.Name())
		}
	}
	if len(folders) > conf.Count {
		t.Errorf("expected %d folders, get %d", conf.Count, len(folders))
	}
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.RemoveAll("../../testdata/releases")
	os.Exit(code)
}
