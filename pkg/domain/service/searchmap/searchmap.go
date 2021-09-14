package searchmap

import (
	"encoding/json"
	"github.com/SPANDigital/presidium-hugo/pkg/filesystem"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const (
	nameDirContent        = "content"
	nameDirSite           = "public"
	nameFileSearchMapJson = "searchmap.json"
	fileExtMarkdown       = ".md"
)

type Validation struct {
	MissingMarkdownFiles []string // All files found in the content Dir
	// -- Private state follows, please do not add public state after this line --
	projectDir        string
	contentDir        string // the site content dir
	siteDir           string // public directory of site
	searchMapJsonFile string // the path to the `searchmap.json` file
}

/*
New creates a new Validation and ensures that projectDir points to valid presidium project structure. To
start the validation just call Validation.Run() on it.
*/
func New(projectDir string) (*Validation, error) {

	fileSystem := filesystem.New()

	v := &Validation{
		MissingMarkdownFiles: []string{},
		projectDir:           projectDir,
		contentDir:           filepath.Join(projectDir, nameDirContent),
		siteDir:              filepath.Join(projectDir, nameDirSite),
		searchMapJsonFile:    filepath.Join(projectDir, nameDirSite, nameFileSearchMapJson),
	}

	if err := fileSystem.RequireDir(v.projectDir); err != nil {
		return nil, err
	}

	if err := fileSystem.RequireDir(v.siteDir); err != nil {
		return nil, err
	}

	if err := fileSystem.RequireRegularFile(v.searchMapJsonFile); err != nil {
		return nil, err
	}

	return v, nil
}

func (v *Validation) Run() error {

	searchMapContent, err := v.readSearchMapFile()

	if err != nil {
		return err
	}

	return filepath.Walk(v.contentDir, func(path string, info fs.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if info.Mode().IsRegular() {
			ext := filepath.Ext(info.Name())
			ext = strings.ToLower(ext)
			fileName := filepath.Base(info.Name())
			fileName = strings.ToLower(fileName)
			if ext == fileExtMarkdown && fileName != "_index.md" {
				if _, found := searchMapContent[path]; !found {
					v.MissingMarkdownFiles = append(v.MissingMarkdownFiles, path)
				}
			}
		}

		return nil
	})
}

func (v *Validation) readSearchMapFile() (map[string]entry, error) {

	file, err := os.Open(v.searchMapJsonFile)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var entries []entry

	if err := json.NewDecoder(file).Decode(&entries); err != nil {
		return nil, err
	}

	mapped := make(map[string]entry)

	for _, en := range entries {
		uri := en.Uri
		entry := en
		mapped[uri] = entry
	}

	return mapped, nil
}

func (v *Validation) Failed() bool {
	return len(v.MissingMarkdownFiles) > 0
}

type entry struct {
	Uri      string `json:"id"`
	Category string `json:"category"`
	Content  string `json:"content"`
}
