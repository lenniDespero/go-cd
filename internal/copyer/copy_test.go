package copyer

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/maraino/testify/require"
)

func TestCopy(t *testing.T) {
	dir, err := ioutil.TempDir("", "test-deploy-")
	require.Nil(t, err)
	defer func() {
		_ = os.RemoveAll(dir)
	}()
	err = Copy(dir, dir+"-new")
	defer func() {
		_ = os.RemoveAll(dir + "-new")
	}()
	require.Nil(t, err)
}

func TestCopyRecursive(t *testing.T) {
	dir, err := ioutil.TempDir("", "test-deploy-")
	require.Nil(t, err)
	defer func() {
		_ = os.RemoveAll(dir)
	}()
	err = Copy(dir, filepath.Join(dir, "new"))

	require.NotNil(t, err)
}

func TestCopyWithLink(t *testing.T) {
	dir, err := ioutil.TempDir("", "test-deploy-")
	require.Nil(t, err)
	defer func() {
		_ = os.RemoveAll(dir)
	}()
	err = os.Symlink(filepath.Join(dir, "PathFrom"), filepath.Join(dir, "link"))
	require.Nil(t, err)
	err = Copy(dir, dir+"-new")
	defer func() {
		_ = os.RemoveAll(dir + "-new")
	}()
	require.Nil(t, err)
}

func TestCopyWrongPath(t *testing.T) {
	dir, err := ioutil.TempDir("", "test-deploy-")
	require.Nil(t, err)
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
	require.Nil(t, err)
	defer func() {
		_ = os.RemoveAll(dir)
	}()
	err = ioutil.WriteFile(filepath.Join(dir, "file"), []byte{}, 0500)
	require.Nil(t, err)
	err = Copy(dir, dir+"new")
	defer func() {
		_ = os.RemoveAll(dir + "new")
	}()
	require.Nil(t, err)
}
