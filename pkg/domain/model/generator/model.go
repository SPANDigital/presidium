package generator

import (
	"path/filepath"

	"github.com/SPANDigital/presidium-hugo/pkg/presidiumerr"
	"github.com/google/uuid"
)

type (
	Template             int
	WhenSiteTargetExists int // What should happen if the generator targets an existing site
)

const (
	SpanTemplate Template = iota
	OnBoardingTemplate
	DesignTemplate
	BlogTemplate
)

const (
	AbortWhenTargetSiteExists WhenSiteTargetExists = iota // Abort with an error
	ReplaceTargetSiteIfExists                             // Replaces the content!
)

var (
	SupportedTemplates = []Template{
		SpanTemplate,
		OnBoardingTemplate,
		DesignTemplate,
		BlogTemplate,
	}
)

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
		BrandingModelUrl    string               // The Hugo model used for branding
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
		Brand:       t.BrandingModelUrl,
		Uuid:        uuid.NewString(),
	}
}

// TemplateParameters are fields which gets injected into the template to generate the final skeleton site
type TemplateParameters struct {
	Title       string `json:"title"`
	ProjectName string `json:"project_name"`
	Template    string `json:"template"`
	Brand       string `json:"brand"`
	Uuid        string `json:"uuid"`
}

func (t Template) Name() string {
	return [...]string{
		"SPAN Default Template",
		"SPAN On-boarding Template",
		"SPAN Design Template",
		"SPAN Blog Template",
	}[t]
}

func (t Template) Description() string {
	return [...]string{
		"SPAN's default template",
		"SPAN's on-boarding template",
		"SPAN's design template",
		"SPAN's blog template",
	}[t]
}

func (t Template) Code() string {
	return [...]string{
		"default",
		"onboarding",
		"design",
		"blog",
	}[t]
}

func GetTemplate(code string) (Template, error) {
	switch code {
	case SpanTemplate.Code():
		return SpanTemplate, nil
	case OnBoardingTemplate.Code():
		return OnBoardingTemplate, nil
	case DesignTemplate.Code():
		return DesignTemplate, nil
	case BlogTemplate.Code():
		return BlogTemplate, nil
	default:
		return 0, presidiumerr.GenericError{Code: presidiumerr.UnsupportedTemplate}
	}
}
