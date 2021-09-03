package filesystem

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"

	"path/filepath"
)

type FileSystem interface {
	Copy(src, dest string, mode fs.FileMode) error
	CopyDir(src, dest string) error
	DeleteDir(path string) error
	Rename(old string, new string) error
	MakeDirs(path string) error
}

type fileSystem struct {
}

func (f fileSystem) MakeDirs(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func (f fileSystem) Rename(old string, new string) error {
	return os.Rename(old, new)
}

func (f fileSystem) DeleteDir(dir string) error {

	info, err := os.Stat(dir)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return errors.New(fmt.Sprintf("path is not a directory: %s", dir))
	}

	parentDirs, err := os.Open(dir)
	if err != nil {
		return err
	}

	defer parentDirs.Close()
	dirNames, err := parentDirs.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range dirNames {
		if err = os.RemoveAll(filepath.Join(dir, name)); err != nil {
			return err
		}
	}

	return nil
}

func New() FileSystem {
	return &fileSystem{}
}

func (f fileSystem) Copy(src, dest string, mode fs.FileMode) error {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dest, input, mode)
}

func (f fileSystem) CopyDir(src, dest string) error {
	// Create dest directory
	err := os.MkdirAll(dest, fs.ModePerm)
	if err != nil {
		return err
	}

	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		relPath := strings.TrimPrefix(path, src)
		relPath = strings.TrimPrefix(relPath, "/")
		relPath = filepath.Join(dest, relPath)
		if relPath == "" {
			return nil
		}

		if d.IsDir() {
			err := os.MkdirAll(relPath, fs.ModePerm)
			if err != nil {
				return err
			}
		} else {
			err := f.Copy(path, relPath, fs.ModePerm)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
