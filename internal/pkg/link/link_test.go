package link

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func prepareLink() Link {
	return Link{
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
	if err != nil {
		t.Errorf("Unexpected error: %s, expected nil", err.Error())
	}
}

func TestLink_CheckConfigNoPathTo(t *testing.T) {
	link := prepareLink()
	link.To = ""
	err := link.CheckConfig()
	if err != nil {
		if err != NoLinkToError {
			t.Errorf("Unexpected error: %s, expected: %s", err.Error(), NoLinkToError.Error())
		}
	} else {
		t.Errorf("Expected error, get nil")
	}
}

func TestLink_CheckConfigNoPathFrom(t *testing.T) {
	link := prepareLink()
	link.From = ""
	err := link.CheckConfig()
	if err != nil {
		if err != NoLinkFromError {
			t.Errorf("Unexpected error: %s, expected: %s", err.Error(), NoLinkFromError.Error())
		}
	} else {
		t.Errorf("Expected error, get nil")
	}
}

func TestLink_ExecuteOnLocal(t *testing.T) {
	link := prepareLink()
	dir, err := ioutil.TempDir("", "test-")
	if err != nil {
		t.Errorf("Unexpected error: %s, expected nil", err.Error())
	}
	defer func() {
		_ = os.RemoveAll(dir)
	}()
	err = os.Chdir(dir)
	if err != nil {
		t.Errorf("Unexpected error: %s, expected nil", err.Error())
	}
	err = link.ExecuteOnLocal()
	if err != nil {
		t.Errorf("Unexpected error: %s, expected nil", err.Error())
	}
	err = link.ExecuteOnLocal()
	if err == nil {
		t.Errorf("Expected error, get nil")
	} else if os.IsExist(err) == false {
		t.Errorf("Unexpected error: %s, expected:%s", err.Error(), os.ErrExist.Error())
	}
}
