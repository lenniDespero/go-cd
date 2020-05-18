package deployer

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/lenniDespero/go-cd/internal/pkg/host"

	"github.com/mitchellh/mapstructure"

	"github.com/lenniDespero/go-cd/internal/logger"
	"github.com/lenniDespero/go-cd/internal/pkg/pipe"
	"github.com/lenniDespero/go-cd/internal/pkg/target"
	"github.com/lenniDespero/go-cd/internal/ssh"
)

// HostDeployer struct for host deploy runner
type HostDeployer struct {
	conf         target.Config
	client       *ssh.Client
	timeNamePath string
}

// NewHosDeployer prepare HostDeployer
func NewHosDeployer(config target.Config) (*HostDeployer, error) {
	hostConf := config.Host
	switch hostConf.Auth {
	case host.TypePassword:
		client, err := ssh.DialWithPasswd(hostConf.GetConnectionString(), hostConf.User, hostConf.Password)
		if err != nil {
			return nil, err
		}
		return &HostDeployer{conf: config, client: client}, nil
	case host.TypeKey:
		client, err := ssh.DialWithKey(hostConf.GetConnectionString(), hostConf.User, hostConf.Key)
		if err != nil {
			return nil, err
		}
		return &HostDeployer{conf: config, client: client}, nil
	default:
		return nil, fmt.Errorf("unknown type host auth: %s", hostConf.Auth)
	}
}

// Prepare work for deploy
func (h HostDeployer) Prepare() error {
	logger.Debug("prepare to deploy")
	script := h.client.Cmd("ls " + h.conf.Path + ".lock")

	err := script.Run()
	if err != nil {
		logger.Info(`no ".lock" file, make ".lock"`)
		script := h.client.Cmd("touch " + h.conf.Path + ".lock")
		script.SetStdio(os.Stdout, os.Stderr)
		err := script.Run()
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New(`".lock" file already exists`)
}

// UpdateSource will download source code
func (h *HostDeployer) UpdateSource(git string) error {
	logger.Debug("download source from git")
	now := strconv.FormatInt(time.Now().Unix(), 10)
	timeName := now
	h.timeNamePath = filepath.Join(h.conf.Path, timeName)

	script := h.client.Cmd(fmt.Sprintf("git clone %s %s", git, h.timeNamePath))
	script.SetStdio(os.Stdout, os.Stderr)
	err := script.Run()
	if err != nil {
		defer func() {
			script := h.client.Cmd(fmt.Sprintf("rm -rf %s ", h.timeNamePath))
			script.SetStdio(os.Stdout, os.Stderr)
			_ = script.Run()
		}()
		return err
	}
	return nil
}

// RunPipe execute pipe
func (h *HostDeployer) RunPipe() error {
	logger.Debug("run pipes work")
	cmds := []string{
		fmt.Sprintf("cd %s", h.timeNamePath),
	}

	for _, p := range h.conf.Pipe {
		inter, ok := pipe.Names[p.Type]
		if !ok {
			return fmt.Errorf("unresolved %s type", p.Type)
		}
		for _, args := range p.Args {
			err := mapstructure.Decode(args, inter)
			if err != nil {
				return err
			}

			jsonInter, err := json.Marshal(inter)
			if err != nil {
				return err
			}

			pipeint, ok := pipe.NamesInt[p.Type]
			if !ok {
				return fmt.Errorf("unresolved %s type", p.Type)
			}

			if err := json.Unmarshal(jsonInter, &pipeint); err != nil {
				return err
			}

			comm := pipeint.GetRemoteCommand()
			cmds = append(cmds, fmt.Sprintf(`echo "%s"`, p.Name), comm)
		}
	}

	err := h.client.Shell().Start(cmds)
	if err != nil {
		defer func() {
			script := h.client.Cmd(fmt.Sprintf("rm -rf %s ", h.timeNamePath))
			script.SetStdio(os.Stdout, os.Stderr)
			_ = script.Run()
		}()
		return err
	}
	return nil
}

// MakeLinks make link to current version
func (h *HostDeployer) MakeLinks() error {
	logger.Debug("make Links")
	script := h.client.Cmd(fmt.Sprintf("chmod %d %s", 775, h.timeNamePath))
	script.SetStdio(os.Stdout, os.Stderr)

	err := script.Run()
	if err != nil {
		defer func() {
			script := h.client.Cmd(fmt.Sprintf("rm -rf %s ", h.timeNamePath))
			script.SetStdio(os.Stdout, os.Stderr)
			_ = script.Run()
		}()
		return err
	}

	script = h.client.Cmd(fmt.Sprintf("ln -s -f %s %s", h.timeNamePath, filepath.Join(h.conf.Path, "current")))
	script.SetStdio(os.Stdout, os.Stderr)

	err = script.Run()
	if err != nil {
		defer func() {
			script := h.client.Cmd(fmt.Sprintf("rm -rf %s ", h.timeNamePath))
			script.SetStdio(os.Stdout, os.Stderr)
			_ = script.Run()
		}()
		return err
	}
	return nil
}

// CleanUp after work
func (h *HostDeployer) CleanUp(count int) error {
	logger.Debug("cleanUp work")
	script := h.client.Cmd(fmt.Sprintf("ls -d %s*", h.conf.Path))
	out, err := script.Output()
	if err != nil {
		return err
	}

	folders := strings.Split(string(out), "\n")
	//link to current +1
	cnt := count + 1
	if len(folders) > cnt {
		logger.Debug("clean folders")
		for _, folder := range folders[0:(len(folders) - cnt)] {
			if folder != "" {
				script := h.client.Cmd(fmt.Sprintf("rm -rf %s", folder))
				script.SetStdio(os.Stdout, os.Stderr)
				err = script.Run()
				if err != nil {
					return err
				}
			}
		}
	}

	script = h.client.Cmd("rm " + h.conf.Path + ".lock")
	script.SetStdio(os.Stdout, os.Stderr)

	err = script.Run()
	if err != nil {
		return err
	}
	err = h.client.Close()
	if err != nil {
		return err
	}
	return nil
}
