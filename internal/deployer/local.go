package deployer

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"

	"github.com/lenniDespero/go-cd/internal/copyer"
	"github.com/lenniDespero/go-cd/internal/logger"
	"github.com/lenniDespero/go-cd/internal/pkg/pipe"
	"github.com/lenniDespero/go-cd/internal/pkg/target"
)

// LocalDeployer local deploy runner struct
type LocalDeployer struct {
	conf         target.Config
	tmpdir       string
	timeName     string
	timeNamePath string
	absPth       string
}

// NewLocalDeployer prepare local deploy runner
func NewLocalDeployer(config target.Config) (DeployInterface, error) {
	return &LocalDeployer{conf: config}, nil
}

// Prepare Will check lock file
// And create lock file
func (l *LocalDeployer) Prepare() error {
	logger.Debug("prepare to deploy")
	path, err := filepath.Abs(l.conf.Path)
	if err != nil {
		return err
	}
	l.absPth = path
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		err = os.MkdirAll(path, 0775)
		if err != nil {
			return err
		}
	}

	path, err = filepath.Abs(filepath.Join(l.absPth, ".lock"))
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

// UpdateSource will clone sources from git to tmp folder,
// Then copy files from temp folder to deploy folder
// Then remove tmp folder
func (l *LocalDeployer) UpdateSource(gitPath string) error {
	logger.Debug("download source from git")
	dir, err := ioutil.TempDir("", "deploy-")
	if err != nil {
		return err
	}
	l.tmpdir = dir
	defer os.RemoveAll(l.tmpdir)
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
	l.timeNamePath = filepath.Join(l.absPth, l.timeName)
	err = copyer.Copy(l.tmpdir, l.timeNamePath)
	if err != nil {
		_ = os.RemoveAll(l.timeNamePath)
		return err
	}
	err = os.Chmod(l.timeNamePath, 0775)
	if err != nil {
		_ = os.RemoveAll(l.timeNamePath)
		return err
	}
	return nil
}

// MakeLinks will make links to current release
func (l *LocalDeployer) MakeLinks() error {
	logger.Debug("make Links")
	err := os.Chdir(filepath.Join(l.absPth))
	if err != nil {
		return err
	}
	err = os.Symlink(l.timeName, l.timeName+"link")
	if err != nil {
		return err
	}
	err = os.Rename(l.timeName+"link", "current")
	if err != nil {
		return err
	}

	return nil
}

// RunPipe execute pipe stages
func (l *LocalDeployer) RunPipe() error {
	logger.Debug("run pipes work")
	err := os.Chdir(l.timeNamePath)
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

// CleanUp work after work
func (l *LocalDeployer) CleanUp(cnt int) error {
	logger.Debug("cleanUp work")
	files, err := ioutil.ReadDir(l.absPth)
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
		logger.Debug("clean folders")
		for _, folder := range folders[0:(len(folders) - cnt)] {
			_ = os.RemoveAll(filepath.Join(l.absPth, folder))
		}
	}
	if err != nil {
		return err
	}
	err = os.RemoveAll(filepath.Join(l.absPth, ".lock"))
	if err != nil {
		return err
	}
	return nil
}
