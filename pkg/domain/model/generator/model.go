package generator

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	"github.com/SPANDigital/presidium-hugo/pkg/presidiumerr"
	"github.com/google/uuid"
)

const templatesDir = "templates"

type (
	Template struct {
		code        string
		name        string
		description string
	}
	WhenSiteTargetExists int // What should happen if the generator targets an existing site
)

const (
	AbortWhenTargetSiteExists WhenSiteTargetExists = iota // Abort with an error
	ReplaceTargetSiteIfExists                             // Replaces the content!
)

// SupportedTemplates is populated by LoadTemplates at startup by scanning the
// immediate subdirectories of the embedded "templates" directory.
var SupportedTemplates []Template

// LoadTemplates discovers the available site templates by reading the
// immediate subdirectories of "templates" inside fsys. Each subdirectory
// becomes one entry in SupportedTemplates, ordered alphabetically by code.
func LoadTemplates(fsys fs.FS) error {
	entries, err := fs.ReadDir(fsys, templatesDir)
	if err != nil {
		return fmt.Errorf("reading %s directory: %w", templatesDir, err)
	}

	discovered := make([]Template, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		code := e.Name()
		label := humanize(code)
		discovered = append(discovered, Template{
			code:        code,
			name:        fmt.Sprintf("%s Template", label),
			description: fmt.Sprintf("Presidium %s template", strings.ToLower(label)),
		})
	}
	sort.Slice(discovered, func(i, j int) bool {
		return discovered[i].code < discovered[j].code
	})
	SupportedTemplates = discovered
	return nil
}

// humanize converts a directory name like "on-boarding" or "key_concepts"
// into a display label like "On Boarding" / "Key Concepts".
func humanize(code string) string {
	parts := strings.FieldsFunc(code, func(r rune) bool {
		return r == '-' || r == '_'
	})
	for i, p := range parts {
		if len(p) == 0 {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + p[1:]
	}
	return strings.Join(parts, " ")
}

type (
	ItemSelection struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	// InitialSiteTarget models the requirement for an initial Presidium site
	InitialSiteTarget struct {
		SiteTargetDirectory string               // Where the site must be generator to
		SiteName            string               // The name of the site
		SiteTitle           string               // The title for the site
		Template            Template             // Template to use
		WhenSiteExists      WhenSiteTargetExists // What should happen when the site already exists.
		Uuid                string               // Unique identifier for the site
	}
)

func (t *InitialSiteTarget) AssetsDir() string {
	return filepath.Join(t.SiteTargetDirectory, "static")
}

func (t *InitialSiteTarget) ContentDir() string {
	return filepath.Join(t.SiteTargetDirectory, "content")
}

func (t *InitialSiteTarget) GetTemplateParameters() TemplateParameters {

	or := func(s1, s2 string) string {
		if len(s1) > 0 {
			return s1
		} else {
			return s2
		}
	}

	_, projectName := filepath.Split(t.SiteTargetDirectory)

	return TemplateParameters{
		Title:       or(t.SiteTitle, t.SiteName),
		ProjectName: projectName,
		Template:    t.Template.Code(),
		Uuid:        uuid.NewString(),
	}
}

// TemplateParameters are fields which gets injected into the template to generate the final skeleton site
type TemplateParameters struct {
	Title       string `json:"title"`
	ProjectName string `json:"project_name"`
	Template    string `json:"template"`
	Uuid        string `json:"uuid"`
}

func (t Template) Name() string        { return t.name }
func (t Template) Description() string { return t.description }
func (t Template) Code() string        { return t.code }

func GetTemplate(code string) (Template, error) {
	for _, t := range SupportedTemplates {
		if t.code == code {
			return t, nil
		}
	}
	return Template{}, presidiumerr.GenericError{Code: presidiumerr.UnsupportedTemplate}
}
