package local

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/lenniDespero/go-cd/internal/logger"

	"github.com/lenniDespero/go-cd/internal/copyer"

	"github.com/lenniDespero/go-cd/internal/pkg/pipe"
	"github.com/lenniDespero/go-cd/internal/pkg/target"

	"github.com/mitchellh/mapstructure"
)

type DeployLocal struct {
	conf     target.Target
	tmpdir   string
	timeName string
}

func InitDeployer(config target.Target) (*DeployLocal, error) {
	return &DeployLocal{conf: config}, nil
}

//Prepare Will check lock file
//And create lock file
func (l *DeployLocal) Prepare() error {
	logger.Debug("Prepare to deploy")
	path, err := filepath.Abs(filepath.Join(l.conf.Path, ".lock"))
	if err != nil {
		return err
	}
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			_, err := os.Create(path)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	if stat != nil {
		return errors.New("lock file already exists")
	}
	return nil
}

//UpdateSource will clone sources from git to tmp folder,
//Then copy files from temp folder to deploy folder
//Then remove tmp folder
func (l *DeployLocal) UpdateSource(gitPath string) error {
	logger.Debug("Download source from git")
	dir, err := ioutil.TempDir("", "deploy-")
	if err != nil {
		return err
	}
	l.tmpdir = dir
	defer func() {
		_ = os.RemoveAll(l.tmpdir)
	}()
	cmd := exec.Command("git", "clone", gitPath, l.tmpdir)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	now := strconv.FormatInt(time.Now().Unix(), 10)
	l.timeName = now
	path, err := filepath.Abs(l.conf.Path)
	if err != nil {
		return err
	}
	err = copyer.Copy(l.tmpdir, filepath.Join(path, l.timeName))
	if err != nil {
		_ = os.RemoveAll(filepath.Join(path, l.timeName))
		return err
	}

	return nil
}

func (l *DeployLocal) MakeLinks() error {
	logger.Debug("Make Links")
	path, err := filepath.Abs(l.conf.Path)
	if err != nil {
		return err
	}
	err = os.Symlink(filepath.Join(path, l.timeName), filepath.Join(path, l.timeName+"link"))
	if err != nil {
		return err
	}
	err = os.Rename(filepath.Join(path, l.timeName+"link"), filepath.Join(path, "current"))
	if err != nil {
		return err
	}

	return nil
}

func (l *DeployLocal) RunPipe() error {
	logger.Debug("Run pipes work")
	path, err := filepath.Abs(l.conf.Path)
	if err != nil {
		return err
	}
	err = os.Chdir(filepath.Join(path, l.timeName))
	if err != nil {
		return err
	}
	for _, p := range l.conf.Pipe {
		inter := pipe.Names[p.Type]
		for _, args := range p.Args {
			err := mapstructure.Decode(args, inter)
			if err != nil {
				return err
			}
			jsonInter, err := json.Marshal(inter)
			if err != nil {
				return err
			}
			pipeint := pipe.NamesInt[p.Type]
			if err := json.Unmarshal(jsonInter, &pipeint); err != nil {
				return err
			}
			err = pipeint.ExecuteOnLocal()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (l *DeployLocal) CleanUp(cnt int) error {
	logger.Debug("CleanUp work")
	path, err := filepath.Abs(l.conf.Path)
	if err != nil {
		return err
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	var folders []string
	for _, rec := range files {
		if rec.IsDir() {
			folders = append(folders, rec.Name())
		}
	}
	if len(folders) > cnt {
		logger.Debug("Clean folders")
		for _, folder := range folders[0:(len(folders) - cnt)] {
			_ = os.RemoveAll(filepath.Join(path, folder))
		}
	}
	if err != nil {
		return err
	}
	err = os.RemoveAll(filepath.Join(path, ".lock"))
	if err != nil {
		return err
	}
	return nil
}
