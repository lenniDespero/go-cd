package copyer

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Copy(src, dest string) error {
	info, err := os.Lstat(src)
	if err != nil {
		return err
	}
	return copy(src, dest, info)
}

func copy(src, dest string, info os.FileInfo) error {
	if info.Mode()&os.ModeSymlink != 0 {
		return linkCopy(src, dest)
	}
	if info.IsDir() {
		return dirCopy(src, dest, info)
	}
	return fileCopy(src, dest, info)
}

func fileCopy(src, dest string, info os.FileInfo) (err error) {
	err = os.MkdirAll(filepath.Dir(dest), os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer func() {
		dErr := f.Close()
		if err == nil {
			err = dErr
		}
	}()

	if err = os.Chmod(f.Name(), info.Mode()); err != nil {
		return err
	}

	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		dErr := s.Close()
		if err == nil {
			err = dErr
		}
	}()

	_, err = io.Copy(f, s)
	return err
}

func dirCopy(srcdir string, destdir string, info os.FileInfo) (err error) {
	originalMode := info.Mode()
	if err := os.MkdirAll(destdir, os.FileMode(0755)); err != nil {
		return err
	}
	defer chmod(destdir, originalMode, &err)

	contents, err := ioutil.ReadDir(srcdir)
	if err != nil {
		return err
	}

	for _, content := range contents {
		cSrc := filepath.Join(srcdir, content.Name())
		cDest := filepath.Join(destdir, content.Name())
		err = copy(cSrc, cDest, content)
		if err != nil {
			return err
		}
	}
	return nil
}

func linkCopy(src, dest string) error {
	src, err := os.Readlink(src)
	if err != nil {
		return err
	}
	return os.Symlink(src, dest)
}

func chmod(dir string, mode os.FileMode, reported *error) {
	if err := os.Chmod(dir, mode); *reported == nil {
		*reported = err
	}
}
