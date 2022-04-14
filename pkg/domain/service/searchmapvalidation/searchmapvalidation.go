package searchmapvalidation

import (
	"encoding/json"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/model/validate"
	"github.com/SPANDigital/presidium-hugo/pkg/filesystem"
	"io/fs"
	"path/filepath"
	"strings"
)

const (
	nameDirContent          = "content"
	nameDirSite             = "public"
	nameFileSearchMapJson   = "searchmap.json"
	fileExtMarkdown         = ".md"
	folderMarkdownIndexFile = "_index.md"
)

type SearchMapValidation interface {
	FindUndeclaredFiles(projectDir string) (*validate.FilesReport, error)
}

func New() SearchMapValidation {
	return validation{
		fs: filesystem.New(),
	}
}

type validation struct {
	fs filesystem.FsUtil
}

func (v validation) FindUndeclaredFiles(projectDir string) (*validate.FilesReport, error) {
	projectDir, err := filepath.Abs(projectDir)
	if err != nil {
		return nil, err
	}

	contentDir := filepath.Join(projectDir, nameDirContent)
	siteDir := filepath.Join(projectDir, nameDirSite)
	fileSearchMap := filepath.Join(projectDir, nameDirSite, nameFileSearchMapJson)

	if err = v.fs.RequireDir(contentDir); err != nil {
		return nil, err
	}

	if err = v.fs.RequireDir(siteDir); err != nil {
		return nil, err
	}

	if err = v.fs.RequireRegularFile(fileSearchMap); err != nil {
		return nil, err
	}

	entries, err := v.readSearchMapFile(fileSearchMap)
	if err != nil {
		return nil, err
	}

	missing, err := v.findMissingFiles(contentDir, entries)
	if err != nil {
		return nil, err
	}

	return &validate.FilesReport{
		Files: missing,
		Found: missing != nil && len(missing) > 0,
	}, nil
}

func (v *validation) readSearchMapFile(filePath string) (map[string]entry, error) {
	file, err := filesystem.AFS.Open(filePath)
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

func (v validation) findMissingFiles(contentDir string, searchMapEntries map[string]entry) ([]string, error) {
	missing := make([]string, 0)
	err := filesystem.AFS.Walk(contentDir, func(path string, info fs.FileInfo, err error) error {
		if info.Mode().IsRegular() {
			ext := filepath.Ext(info.Name())
			ext = strings.ToLower(ext)
			fileName := filepath.Base(info.Name())
			fileName = strings.ToLower(fileName)
			if ext == fileExtMarkdown && fileName != folderMarkdownIndexFile {
				_, found := searchMapEntries[path]
				if !found {
					missing = append(missing, path)
				}
			}
		}
		return nil
	})
	return missing, err
}

type entry struct {
	Uri      string `json:"id"`
	Category string `json:"category"`
	Content  string `json:"content"`
}
