package link

import (
	"fmt"
	"os"
	"path/filepath"
)

//Link struct for config
type Link struct {
	From string `mapstructure:"from" json:"from"`
	To   string `mapstructure:"to" json:"to"`
}

//Error implementation for package
type Error string

//Error implementation for package
func (e Error) Error() string {
	return string(e)
}

//Errors
const (
	NoLinkFromError Error = "no link from"
	NoLinkToError   Error = "no link to"
)

//CheckConfig will check config for errors
func (link Link) CheckConfig() error {
	if link.From == "" {
		return NoLinkFromError
	}
	if link.To == "" {
		return NoLinkToError
	}
	return nil
}

//ExecuteOnLocal make link on local machine
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

//GetRemoteCommand prepare command to make link on remote machine
func (link Link) GetRemoteCommand() string {
	return fmt.Sprintf("ln -s %s %s", link.From, link.To)
}
