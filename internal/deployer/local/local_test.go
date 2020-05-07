package local

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/lenniDespero/go-cd/internal/pkg/config"
)

func prepareConfig(t *testing.T) *config.Config {
	c := config.Config{}
	err := config.ReadConfig(&c, "../../../testdata/test.app.yml")
	if err != nil {
		t.Errorf("Unexpected error : %s", err.Error())
	}
	return &c
}

func getDeployer(t *testing.T) (*DeployLocal, error) {
	conf := prepareConfig(t)
	deployer, err := InitDeployer(conf.Targets["devel"])
	if err != nil {
		return nil, err
	}
	return deployer, nil
}

func (l *DeployLocal) removeLock(t *testing.T) {
	err := os.RemoveAll(filepath.Join(l.absPth, ".lock"))
	if err != nil {
		t.Errorf("Unexpected error : %s", err.Error())
	}
}

func TestInitDeployer(t *testing.T) {
	_, err := getDeployer(t)
	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestDeployLocal_PrepareFail(t *testing.T) {
	deployer, err := getDeployer(t)
	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
	err = deployer.Prepare()
	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
	err = deployer.Prepare()
	if err == nil {
		t.Errorf("Expected err, get nil")
	}
	deployer.removeLock(t)
}

func TestDeployLocal_PrepareFailPath(t *testing.T) {
	deployer, err := getDeployer(t)
	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
	deployer.conf.Path = "/testdata/releases"
	err = deployer.Prepare()
	if err == nil {
		t.Errorf("Expected err, get nil")
	}
	deployer.removeLock(t)
}

func TestDeployLocal_UpdateSourceBadGit(t *testing.T) {
	conf := prepareConfig(t)
	deployer, err := getDeployer(t)
	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
	err = deployer.Prepare()
	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
	err = deployer.UpdateSource(conf.Git + "fsuegjsehfgjhseb")
	if err == nil {
		t.Errorf("Expected err, get nil")
	}
	deployer.removeLock(t)
}

func TestDeployLocal_RunPipeWrongPath(t *testing.T) {
	conf := prepareConfig(t)
	deployer, err := getDeployer(t)
	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
	err = deployer.Prepare()
	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
	err = deployer.UpdateSource(conf.Git)
	if err != nil {
		t.Errorf("Unexpected err : %s", err.Error())
	}
	deployer.timeName = "syeugfsjgefgsef"
	err = deployer.RunPipe()
	if err == nil {
		t.Errorf("Expected err, get nil")
	}
	deployer.removeLock(t)
}

func TestDeployLocal_CleanUp(t *testing.T) {
	conf := prepareConfig(t)
	deployer, err := getDeployer(t)
	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
	err = deployer.Prepare()
	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
	err = deployer.UpdateSource(conf.Git)
	if err != nil {
		t.Errorf("Unexpected err : %s", err.Error())
	}
	err = deployer.RunPipe()
	if err != nil {
		t.Errorf("Unexpected err : %s", err.Error())
	}
	err = deployer.MakeLinks()
	if err != nil {
		t.Errorf("Unexpected err : %s", err.Error())
	}
	err = deployer.CleanUp(conf.Count)
	if err != nil {
		t.Errorf("Unexpected err : %s", err.Error())
	}

	files, err := ioutil.ReadDir(deployer.absPth)
	if err != nil {
		t.Errorf("Unexpected err : %s", err.Error())
	}
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
	err := os.Mkdir("../../../testdata/releases", 0755)
	if err != nil {
		os.Exit(1)
	}
	defer func() {
		_ = os.RemoveAll("../../../testdata/releases")
	}()
	m.Run()
}
