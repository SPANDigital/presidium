package template

import (
	"bytes"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/model/generator"
	"github.com/SPANDigital/presidium-hugo/pkg/filesystem"
	"github.com/SPANDigital/presidium-hugo/templates"
)

type Service struct {
	templates fs.FS
}

func New() Service {
	return Service{
		templates: templates.FS,
	}
}

// GetListing returns a list of files by a given template
func (s Service) GetListing(templateDir string) ([]string, error) {
	listing := make([]string, 0)
	return listing, fs.WalkDir(s.templates, templateDir, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			listing = append(listing, filePath)
		}
		return nil
	})
}

func (s Service) ProcessDirTemplates(templateDir string, outputDir string, model generator.TemplateParameters) error {
	err := fs.WalkDir(s.templates, templateDir, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			relativePath := strings.TrimPrefix(filepath.Dir(filePath), templateDir)
			outputPath := path.Join(outputDir, relativePath)
			return s.ProcessTemplate(outputPath, filePath, model)
		}
		return nil
	})
	return err
}

func (s Service) ProcessTemplate(dir, theTemplate string, model generator.TemplateParameters) error {
	filename := filepath.Base(theTemplate)
	if filename == "go.mod.tpl" {
		filename = "go.mod"
	}
	templateBytes, err := fs.ReadFile(s.templates, theTemplate)
	if err != nil {
		return err
	}
	templateString := string(templateBytes)
	finalPath := path.Join(dir, filename)
	err = filesystem.AFS.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	f, err := filesystem.AFS.Create(finalPath)
	if err != nil {
		return err
	}
	var b bytes.Buffer
	t := template.Must(template.New(filepath.Base(theTemplate)).Funcs(sprig.HermeticTxtFuncMap()).Parse(templateString))
	err = t.Execute(&b, model)
	if err != nil {
		return err
	}
	_, err = f.Write(b.Bytes())
	if err != nil {
		return err
	}
	return nil
}
