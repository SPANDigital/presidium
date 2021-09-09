package filesystem

import (
	"errors"
	"fmt"
	"github.com/otiai10/copy"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type FileSystem interface {
	Copy(src, dest string, mode fs.FileMode) error
	CopyWithOptions(src, dest string, options copy.Options) error
	CopyDir(src, dest string) error
	EmptyDir(path string) error
	Rename(old string, new string) error
	MakeDirs(path string) error
	AbsolutePath(path string) (string, error) // TODO: Need to add unit test for this
	GetWorkingDir() (string, error)
	DirExists(dir string) bool
	DeleteDir(dir string) error
}

type fileSystem struct {
}

func (f fileSystem) DeleteDir(dir string) error {
	return os.RemoveAll(dir)
}

func (f fileSystem) DirExists(dir string) bool {
	info, err := os.Stat(dir)
	return err != nil && info.IsDir()
}

func (f fileSystem) GetWorkingDir() (string, error) {
	return os.Getwd()
}

func (f fileSystem) CopyWithOptions(src, dest string, options copy.Options) error {
	return copy.Copy(src, dest, options)
}

func (f fileSystem) AbsolutePath(path string) (string, error) {

	if filepath.IsAbs(path) {
		return path, nil
	}

	return filepath.Abs(path)
}

func (f fileSystem) MakeDirs(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func (f fileSystem) Rename(old string, new string) error {
	return os.Rename(old, new)
}

// EmptyDir removes all content leaving an empty directory
func (f fileSystem) EmptyDir(dir string) error {

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
