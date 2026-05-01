package template

import (
	"bytes"
	"fmt"
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

// verbatimPathPrefixes lists template subdirectories whose contents are copied
// to the generated site as-is, without running them through text/template.
// These paths typically contain Hugo template syntax ({{ ... }}) which
// collides with Go's template syntax and would otherwise fail to parse.
//
// TODO: when requirements-specific layouts are extracted into a Hugo module
// (e.g. github.com/spandigital/presidium-layouts-requirements), "layouts/"
// can be removed from this list — at that point no template should ship
// inlined Hugo layouts.
var verbatimPathPrefixes = []string{
	"layouts/",
	"static/",
	"data/",
	"scripts/",
	"assets/",
	"discovery/",
	"resources/", // Hugo build cache (e.g. _gen/) — must not be templated
	"public/",    // Hugo output cache shipped with some templates
}

// SetFS sets the embedded filesystem used to read templates.
func SetFS(fsys fs.FS) {
	templatesFS = fsys
}

type Service struct {
	templates fs.FS
}

func New() (Service, error) {
	fsys := templatesFS
	if fsys == nil {
		// Fallback for tests: read templates from the module root filesystem.
		root, err := findModuleRoot()
		if err != nil {
			return Service{}, fmt.Errorf("locating module root: %w", err)
		}
		fsys = os.DirFS(root)
	}
	// Sub-FS into "templates" so callers can use template names directly (e.g. "default").
	sub, err := fs.Sub(fsys, "templates")
	if err != nil {
		return Service{}, fmt.Errorf("templates directory not found: %w", err)
	}
	return Service{
		templates: sub,
	}, nil
}

// findModuleRoot walks up from the working directory to find the module root.
func findModuleRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("unable to determine working directory: %w", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ".", nil
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
		if isVerbatimPath(p, templateDir) {
			return s.CopyVerbatim(outputPath, p)
		}
		return s.ProcessTemplate(outputPath, p, model)
	})
}

// isVerbatimPath reports whether the file at p (relative to the embedded
// templates filesystem) lives under one of the verbatimPathPrefixes for the
// given templateDir, and therefore must be copied without text/template
// processing.
func isVerbatimPath(p, templateDir string) bool {
	rel := strings.TrimPrefix(p, templateDir)
	rel = strings.TrimPrefix(rel, "/")
	for _, prefix := range verbatimPathPrefixes {
		if strings.HasPrefix(rel, prefix) {
			return true
		}
	}
	return false
}

// CopyVerbatim writes the source template file to dir/filename without
// running it through text/template. The .tmpl suffix is still stripped so
// that authors can opt out of accidental shell/IDE interference by naming a
// file e.g. raw-payload.json.tmpl.
func (s Service) CopyVerbatim(dir, theTemplate string) error {
	filename := filepath.Base(theTemplate)
	filename = strings.TrimSuffix(filename, ".tmpl")

	data, err := fs.ReadFile(s.templates, theTemplate)
	if err != nil {
		return err
	}

	finalPath := path.Join(dir, filename)
	if err := filesystem.AFS.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	f, err := filesystem.AFS.Create(finalPath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
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
