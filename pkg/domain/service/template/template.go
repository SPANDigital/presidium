package template

import (
	"bytes"
	"github.com/Masterminds/sprig"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/model/generator"
	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/packr/v2"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

type Service struct {
	templates packd.Box
}

func New() Service {
	box := packr.New("templatesBox", "../../../../templates")
	return Service{
		templates: box,
	}
}

// GetListing returns a list of files by a given template
func (s Service) GetListing(templateDir string) ([]string, error) {
	listing := make([]string, 0)
	return listing, s.templates.WalkPrefix(templateDir, func(templateName string, file packd.File) error {
		listing = append(listing, templateName)
		return nil
	})
}

func (s Service) ProcessDirTemplates(templateDir string, outputDir string, model generator.TemplateParameters) error {
	err := s.templates.WalkPrefix(templateDir, func(templateName string, file packd.File) error {
		relativePath := strings.TrimPrefix(filepath.Dir(templateName), templateDir)
		outputPath := path.Join(outputDir, relativePath)
		return s.ProcessTemplate(outputPath, templateName, model)
	})
	return err
}

func (s Service) ProcessTemplate(dir, theTemplate string, model generator.TemplateParameters) error {
	filename := filepath.Base(theTemplate)
	templateString, err := s.templates.FindString(theTemplate)
	if err != nil {
		return err
	}
	finalPath := path.Join(dir, filename)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	f, err := os.Create(finalPath)
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
