package filesystem

import (
	"errors"
	"fmt"
	"github.com/otiai10/copy"
	"github.com/spf13/afero"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var (
	FS     afero.Fs
	FSUtil *afero.Afero
)

func init() {
	FS = afero.NewOsFs()
	FSUtil = &afero.Afero{Fs: FS}
}

type FileSystem interface {
	Copy(src, dest string, mode fs.FileMode) error
	CopyWithOptions(src, dest string, options copy.Options) error
	CopyDir(src, dest string) error
	EmptyDir(path string) error
	Rename(old string, new string) error
	MakeDirs(path string) error
	AbsolutePath(path string) (string, error)
	GetWorkingDir() (string, error)
	DirExists(dir string) bool
	DeleteDir(dir string) error
}

type fileSystem struct {
}

func (f fileSystem) DeleteDir(dir string) error {
	return FS.RemoveAll(dir)
}

func (f fileSystem) DirExists(dir string) bool {
	info, err := FS.Stat(dir)
	if err == nil {
		return info.IsDir()
	}
	return false
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
	return FS.MkdirAll(path, os.ModePerm)
}

func (f fileSystem) Rename(old string, new string) error {
	return FS.Rename(old, new)
}

// EmptyDir removes all content leaving an empty directory
func (f fileSystem) EmptyDir(dir string) error {

	info, err := FS.Stat(dir)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return errors.New(fmt.Sprintf("path is not a directory: %s", dir))
	}

	parentDirs, err := FS.Open(dir)
	if err != nil {
		return err
	}

	defer parentDirs.Close()
	dirNames, err := parentDirs.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range dirNames {
		if err = FS.RemoveAll(filepath.Join(dir, name)); err != nil {
			return err
		}
	}

	return nil
}

func New() FileSystem {
	return &fileSystem{}
}

func (f fileSystem) Copy(src, dest string, mode fs.FileMode) error {
	input, err := FSUtil.ReadFile(src)
	if err != nil {
		return err
	}
	return FSUtil.WriteFile(dest, input, mode)
}

func (f fileSystem) CopyDir(src, dest string) error {
	// Create dest directory
	err := FS.MkdirAll(dest, fs.ModePerm)
	if err != nil {
		return err
	}

	return FSUtil.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath := strings.TrimPrefix(path, src)
		relPath = strings.TrimPrefix(relPath, "/")
		relPath = filepath.Join(dest, relPath)
		if relPath == "" {
			return nil
		}

		if info.IsDir() {
			err := FS.MkdirAll(relPath, fs.ModePerm)
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
