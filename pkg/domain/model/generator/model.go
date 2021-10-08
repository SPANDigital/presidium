package generator

import (
	"github.com/SPANDigital/presidium-hugo/pkg/domain/wizard"
	"github.com/SPANDigital/presidium-hugo/pkg/presidiumerr"
)

type (
	Template             int
	Theme                int
	WhenSiteTargetExists int // What should happen if the generator targets an existing site
)

const (
	SpanTemplate Template = iota
	OnBoardingTemplate
	DesignTemplate
)

const (
	PresidiumTheme Theme = iota
)

const (
	AbortWhenTargetSiteExists WhenSiteTargetExists = iota // Abort with an error
	ReplaceTargetSiteIfExists                             // Replaces the content!
)

var (
	SupportedTemplates = []wizard.Template{
		wizard.SpanTemplate,
		wizard.OnBoardingTemplate,
		wizard.DesignTemplate,
	}
	SupportedThemes = []wizard.Theme{
		wizard.PresidiumTheme,
	}
)

// InitialSiteTarget models the requirement for an initial Presidium site
type InitialSiteTarget struct {
	SiteTargetDirectory string               // Where the site must be generator to
	SiteName            string               // The name of the site
	SiteTitle           string               // The title for the site
	BrandingModelUrl    string               // The Hugo model used for branding
	Theme               Theme                // Theme to use
	Template            Template             // Template to use
	WhenSiteExists      WhenSiteTargetExists // What should happen when the site already exists.
}

// TemplateParameters are fields which gets injected into the template to generate the final skeleton site
type TemplateParameters struct {
	Title       string `json:"title"`
	ProjectName string `json:"project_name"`
	Theme       string `json:"theme"`
	Template    string `json:"template"`
	Brand       string `json:"brand"`
}

func (t Template) Name() string {
	return [...]string{
		"SPAN Default Template",
		"SPAN On-boarding Template",
		"SPAN Design Template",
	}[t]
}

func (t Template) Description() string {
	return [...]string{
		"SPAN's default template",
		"SPAN's on-boarding template",
		"SPAN's design template",
	}[t]
}

func (t Template) Code() string {
	return [...]string{
		"default",
		"onboarding",
		"design",
	}[t]
}

func (t Theme) Name() string {
	return [...]string{
		"Presidium Theme",
	}[t]
}

func (t Theme) Description() string {
	return [...]string{
		"Presidium's default theme",
	}[t]
}

func (t Theme) Code() string {
	return [...]string{
		"presidium",
	}[t]
}

func (t Theme) ModulePath() string {
	return [...]string{
		"github.com/spandigital/presidium-theme-website",
	}[t]
}

func GetTemplate(name string) (Template, error) {
	switch name {
	case SpanTemplate.Code():
		return SpanTemplate, nil
	case OnBoardingTemplate.Code():
		return OnBoardingTemplate, nil
	case DesignTemplate.Code():
		return DesignTemplate, nil
	default:
		return 0, presidiumerr.GenericError{Code: presidiumerr.UnsupportedTemplate}
	}
}

func GetTheme(name string) (Theme, error) {
	switch name {
	case PresidiumTheme.Code():
		return PresidiumTheme, nil
	default:
		return 0, presidiumerr.GenericError{Code: presidiumerr.UnsupportedTheme}
	}
}
