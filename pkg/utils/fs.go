package utils

import (
	"github.com/SPANDigital/presidium-hugo/pkg/filesystem"
	"github.com/SPANDigital/presidium-hugo/pkg/log"
	"os"
)

func FileExists(path string) bool {
	if _, err := filesystem.AFS.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// WalkRename queues path rename during walk and processes them in reverse to prevent nil pointer exception
func WalkRename(path string, rename func(path string, info os.FileInfo) (*string, error)) error {
	type Rename struct {
		from string
		to   string
	}

	var renames []Rename
	err := filesystem.AFS.Walk(path, func(path string, info os.FileInfo, err error) error {
		newPath, err := rename(path, info)
		if err != nil {
			return err
		}

		if newPath != nil {
			renames = append(renames, Rename{
				from: path,
				to:   *newPath,
			})
		}
		return nil
	})

	for i := len(renames) - 1; i >= 0; i-- {
		to := renames[i]
		log.Infof("renaming directory %s to %s", to.from, to.to)
		if err = filesystem.AFS.Rename(to.from, to.to); err != nil {
			return err
		}
	}

	return err
}
