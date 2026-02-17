package template

import (
	"bytes"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/model/generator"
	"github.com/SPANDigital/presidium-hugo/pkg/filesystem"
)

// templatesFS holds the embedded filesystem set from main via SetFS.
var templatesFS fs.FS

// SetFS sets the embedded filesystem used to read templates.
func SetFS(fsys fs.FS) {
	templatesFS = fsys
}

type Service struct {
	templates fs.FS
}

func New() Service {
	fsys := templatesFS
	if fsys == nil {
		// Fallback for tests: read templates from the module root filesystem.
		fsys = os.DirFS(findModuleRoot())
	}
	// Sub-FS into "templates" so callers can use template names directly (e.g. "default").
	sub, err := fs.Sub(fsys, "templates")
	if err != nil {
		panic("templates directory not found: " + err.Error())
	}
	return Service{
		templates: sub,
	}
}

// findModuleRoot walks up from the working directory to find the module root.
func findModuleRoot() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "."
		}
		dir = parent
	}
}

// GetListing returns a list of files under a given template directory.
func (s Service) GetListing(templateDir string) ([]string, error) {
	var listing []string
	err := fs.WalkDir(s.templates, templateDir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			listing = append(listing, p)
		}
		return nil
	})
	return listing, err
}

func (s Service) ProcessDirTemplates(templateDir string, outputDir string, model generator.TemplateParameters) error {
	return fs.WalkDir(s.templates, templateDir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		relativePath := strings.TrimPrefix(filepath.Dir(p), templateDir)
		outputPath := path.Join(outputDir, relativePath)
		return s.ProcessTemplate(outputPath, p, model)
	})
}

func (s Service) ProcessTemplate(dir, theTemplate string, model generator.TemplateParameters) error {
	filename := filepath.Base(theTemplate)
	// Strip .tmpl suffix so e.g. "go.mod.tmpl" becomes "go.mod"
	filename = strings.TrimSuffix(filename, ".tmpl")

	data, err := fs.ReadFile(s.templates, theTemplate)
	if err != nil {
		return err
	}
	templateString := string(data)

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
