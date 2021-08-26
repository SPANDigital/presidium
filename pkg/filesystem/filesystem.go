package filesystem

import (
	"io/fs"
	"io/ioutil"
	"os"
	"strings"

	"path/filepath"
)

type FileSystem interface {
	Copy(src, dest string, mode fs.FileMode) error
	CopyDir(src, dest string) error
}

type fileSystem struct {
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
