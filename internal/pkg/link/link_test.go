package link

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/maraino/testify/require"
)

func prepareLink() Config {
	return Config{
		From: "somePathFrom",
		To:   "SomePathTo",
	}
}

func TestLink_GetRemoteCommand(t *testing.T) {
	link := prepareLink()
	str := link.GetRemoteCommand()
	real := strings.Join([]string{"ln -s", link.From, link.To}, " ")
	if str != real {
		t.Errorf("Wrong connection string, get: %s, expected :%s", str, real)
	}
}

func TestLink_CheckConfig(t *testing.T) {
	link := prepareLink()
	err := link.CheckConfig()
	require.Nil(t, err)
}

func TestLink_CheckConfigNoPathTo(t *testing.T) {
	link := prepareLink()
	link.To = ""
	err := link.CheckConfig()
	require.Error(t, err, ErrNoLinkTo)
}

func TestLink_CheckConfigNoPathFrom(t *testing.T) {
	link := prepareLink()
	link.From = ""
	err := link.CheckConfig()
	require.Error(t, err, ErrNoLinkFrom)
}

func TestLink_ExecuteOnLocal(t *testing.T) {
	link := prepareLink()
	dir, err := ioutil.TempDir("", "test-")
	require.Nil(t, err)
	defer func() {
		_ = os.RemoveAll(dir)
	}()
	err = os.Chdir(dir)
	require.Nil(t, err)
	err = link.ExecuteOnLocal()
	require.Nil(t, err)
	err = link.ExecuteOnLocal()
	require.NotNil(t, err)
	if os.IsExist(err) == false {
		t.Errorf("Unexpected error: %s, expected:%s", err.Error(), os.ErrExist.Error())
	}
}
