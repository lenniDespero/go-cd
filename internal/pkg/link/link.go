package link

import (
	"fmt"
	"os"
	"path/filepath"
)

type Link struct {
	From string `mapstructure:"from" json:"from"`
	To   string `mapstructure:"to" json:"to"`
}
type Error string

func (e Error) Error() string {
	return string(e)
}

//Errors
const (
	NoLinkFromError Error = "no link from"
	NoLinkToError   Error = "no link to"
)

func (link Link) CheckConfig() error {
	if link.From == "" {
		return NoLinkFromError
	}
	if link.To == "" {
		return NoLinkToError
	}
	return nil
}

func (link Link) ExecuteOnLocal() error {
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

func (link Link) GetRemoteCommand() string {
	return fmt.Sprintf("ln -s %s %s", link.From, link.To)
}
