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

type FsUtil interface {
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
	RequireRegularFile(path string) error // TODO: need to add unit test for this
	RequireDir(dir string) error          // TODO: need to add unit test for this
}

func SetFileSystem(fs afero.Fs) {
	AFS = &afero.Afero{Fs: fs}
	FS = fs
}

type fsUtil struct{}

func New() FsUtil {
	return &fsUtil{}
}

func (f fsUtil) RequireDir(dir string) error {
	info, err := AFS.Stat(dir)
	if err != nil {
		return err
	}

	if !info.Mode().IsDir() {
		return fmt.Errorf("expected directory here: %s", dir)
	}

	return nil
}

func (f fsUtil) RequireRegularFile(path string) error {
	info, err := AFS.Stat(path)
	if err != nil {
		return err
	}

	if !info.Mode().IsRegular() {
		return fmt.Errorf("expected regular file here: %s", path)
	}

	return nil
}

func (f fsUtil) DeleteDir(dir string) error {
	return FS.RemoveAll(dir)
}

func (f fsUtil) DirExists(dir string) bool {
	info, err := FS.Stat(dir)
	if err == nil {
		return info.IsDir()
	}
	return false
}

func (f fsUtil) GetWorkingDir() (string, error) {
	return os.Getwd()
}

func (f fsUtil) CopyWithOptions(src, dest string, options copy.Options) error {
	return copy.Copy(src, dest, options)
}

func (f fsUtil) AbsolutePath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}

	return filepath.Abs(path)
}

func (f fsUtil) MakeDirs(path string) error {
	return FS.MkdirAll(path, os.ModePerm)
}

func (f fsUtil) Rename(old string, new string) error {
	return FS.Rename(old, new)
}

// EmptyDir removes all content leaving an empty directory
func (f fsUtil) EmptyDir(dir string) error {
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

func (f fsUtil) Copy(src, dest string, mode fs.FileMode) error {
	input, err := AFS.ReadFile(src)
	if err != nil {
		return err
	}
	return AFS.WriteFile(dest, input, mode)
}

func (f fsUtil) CopyDir(src, dest string) error {
	// Create dest directory
	err := FS.MkdirAll(dest, fs.ModePerm)
	if err != nil {
		return err
	}

	return AFS.Walk(src, func(path string, info fs.FileInfo, err error) error {
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
