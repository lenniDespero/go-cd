package copyer

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestCopy(t *testing.T) {
	dir, err := ioutil.TempDir("", "test-deploy-")
	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
	defer func() {
		_ = os.RemoveAll(dir)
	}()
	err = Copy(dir, dir+"-new")
	defer func() {
		_ = os.RemoveAll(dir + "-new")
	}()
	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestCopyRecursive(t *testing.T) {
	dir, err := ioutil.TempDir("", "test-deploy-")
	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
	defer func() {
		_ = os.RemoveAll(dir)
	}()
	err = Copy(dir, filepath.Join(dir, "new"))

	if err == nil {
		t.Errorf("Expcted err: file name too long, nill given")
	}
}

func TestCopyWithLink(t *testing.T) {
	dir, err := ioutil.TempDir("", "test-deploy-")
	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
	defer func() {
		_ = os.RemoveAll(dir)
	}()
	err = os.Symlink(filepath.Join(dir, "PathFrom"), filepath.Join(dir, "link"))
	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
	err = Copy(dir, dir+"-new")
	defer func() {
		_ = os.RemoveAll(dir + "-new")
	}()
	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestCopyWrongPath(t *testing.T) {
	dir, err := ioutil.TempDir("", "test-deploy-")
	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
	defer func() {
		_ = os.RemoveAll(dir)
	}()
	err = Copy(dir+"fake", dir)
	if !os.IsNotExist(err) {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestCopyFiles(t *testing.T) {
	dir, err := ioutil.TempDir("", "test-deploy-")
	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
	defer func() {
		_ = os.RemoveAll(dir)
	}()
	err = ioutil.WriteFile(filepath.Join(dir, "file"), []byte{}, 0500)
	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
	err = Copy(dir, dir+"new")
	defer func() {
		_ = os.RemoveAll(dir + "new")
	}()
	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}
