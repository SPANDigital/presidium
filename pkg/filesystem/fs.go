package filesystem

import "github.com/spf13/afero"

var (
	FS  afero.Fs
	AFS *afero.Afero
)

func init() {
	SetFileSystem(afero.NewOsFs())
}
