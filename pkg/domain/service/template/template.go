package template

import (
	"bytes"
	"github.com/Masterminds/sprig"
	"github.com/SPANDigital/presidium-hugo/pkg/config"
	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/packr/v2"
	"github.com/spf13/viper"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

type Service struct {
	templates *packr.Box
}

func New() Service {
	box := packr.New("templatesBox", "../../../../templates")
	return Service{
		templates: box,
	}
}

func (s Service) ProcessDirTemplates(dir string, obj interface{}) error {
	err := s.templates.WalkPrefix(dir, func(templateName string, file packd.File) error {
		relativePath := strings.TrimPrefix(filepath.Dir(templateName), dir)
		outputPath := path.Join(viper.GetString(config.ProjectNameKey), relativePath)
		return s.ProcessTemplate(outputPath, templateName, obj)
	})
	return err
}

func (s Service) ProcessTemplate(dir, tpl string, obj interface{}) error {
	filename := filepath.Base(tpl)
	tplStr, err := s.templates.FindString(tpl)
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
	t := template.Must(template.New(filepath.Base(tpl)).Funcs(sprig.HermeticTxtFuncMap()).Parse(tplStr))
	err = t.Execute(&b, obj)
	if err != nil {
		return err
	}
	_, err = f.Write(b.Bytes())
	if err != nil {
		return err
	}
	return nil
}
