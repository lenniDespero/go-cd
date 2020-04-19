package host

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/lenniDespero/go-cd/internal/logger"

	"github.com/lenniDespero/go-cd/internal/pkg/pipe"
	"github.com/mitchellh/mapstructure"

	ssh_client "github.com/lenniDespero/go-cd/internal/sshclient"

	"github.com/lenniDespero/go-cd/internal/pkg/target"
)

type DeployHost struct {
	conf     target.Target
	client   *ssh_client.Client
	timeName string
}

func InitDeployer(config target.Target) (*DeployHost, error) {
	host := config.Host
	switch host.Auth {
	case "password":
		client, err := ssh_client.DialWithPasswd(host.GetConnectionString(), host.User, host.Password)
		if err != nil {
			return &DeployHost{}, err
		}
		return &DeployHost{conf: config, client: client}, nil
	case "key":
		client, err := ssh_client.DialWithKey(host.GetConnectionString(), host.User, host.Key)
		if err != nil {
			return &DeployHost{}, err
		}
		return &DeployHost{conf: config, client: client}, nil
	default:
		return &DeployHost{}, fmt.Errorf("unknown type host auth: %s", host.Auth)
	}
}

func (h DeployHost) Prepare() error {
	logger.Debug("Prepare to deploy")
	script := h.client.Cmd("ls " + h.conf.Path + ".lock")
	err := script.Run()
	if err != nil {
		logger.Info(`No ".lock" file, make ".lock"`)
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

func (h *DeployHost) UpdateSource(git string) error {
	logger.Debug("Download source from git")
	now := strconv.FormatInt(time.Now().Unix(), 10)
	h.timeName = now
	script := h.client.Cmd(fmt.Sprintf("git clone %s %s", git, filepath.Join(h.conf.Path, h.timeName)))
	script.SetStdio(os.Stdout, os.Stderr)
	err := script.Run()
	if err != nil {
		defer func() {
			script := h.client.Cmd(fmt.Sprintf("rm -rf %s ", filepath.Join(h.conf.Path, h.timeName)))
			script.SetStdio(os.Stdout, os.Stderr)
			_ = script.Run()
		}()
		return err
	}
	return nil
}

func (h *DeployHost) RunPipe() error {
	logger.Debug("Run pipes work")
	cmds := []string{}
	cmds = append(cmds, fmt.Sprintf("cd %s", filepath.Join(h.conf.Path, h.timeName)))
	for _, p := range h.conf.Pipe {
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
			comm := pipeint.GetRemoteCommand()
			cmds = append(cmds, fmt.Sprintf(`echo "%s"`, p.Name), comm)
		}
	}
	err := h.client.Shell().Start(cmds)
	if err != nil {
		return err
	}
	return nil
}

func (h *DeployHost) MakeLinks() error {
	logger.Debug("Make Links")
	script := h.client.Cmd(fmt.Sprintf("ln -s -f %s %s", filepath.Join(h.conf.Path, h.timeName), filepath.Join(h.conf.Path, "current")))
	script.SetStdio(os.Stdout, os.Stderr)
	err := script.Run()
	if err != nil {
		return err
	}
	return nil
}

func (h *DeployHost) CleanUp(count int) error {
	logger.Debug("CleanUp work")
	script := h.client.Cmd(fmt.Sprintf("ls -d %s*", h.conf.Path))
	out, err := script.Output()
	if err != nil {
		log.Println(err)
		return err
	}
	folders := strings.Split(string(out), "\n")
	//link to current +1
	cnt := count + 1
	if len(folders) > cnt {
		logger.Debug("Clean folders")
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
