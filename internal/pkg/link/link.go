package link

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Config link struct for config
type Config struct {
	From string `mapstructure:"from" json:"from"`
	To   string `mapstructure:"to" json:"to"`
}

// Errors
var (
	ErrNoLinkFrom = errors.New("no link from")
	ErrNoLinkTo   = errors.New("no link to")
)

// CheckConfig will check config for errors
func (link Config) CheckConfig() error {
	if link.From == "" {
		return ErrNoLinkFrom
	}
	if link.To == "" {
		return ErrNoLinkTo
	}
	return nil
}

// ExecuteOnLocal make link on local machine
func (link Config) ExecuteOnLocal() error {
	from, err := filepath.Abs(link.From)
	if err != nil {
		return err
	}
	to, err := filepath.Abs(link.To)
	if err != nil {
		return err
	}

	err = os.Symlink(from, to)
	if err != nil {
		return err
	}
	return nil
}

// GetRemoteCommand prepare command to make link on remote machine
func (link Config) GetRemoteCommand() string {
	return fmt.Sprintf("ln -s %s %s", link.From, link.To)
}
