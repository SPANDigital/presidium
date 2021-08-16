package wizard

import "github.com/SPANDigital/presidium-hugo/pkg/presidiumerr"

type (
	Item struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	Template int
	Theme    int
)

const (
	SpanTemplate Template = iota
	OnBoardingTemplate
	DesignTemplate
)

const (
	PresidiumTheme Theme = iota
)

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
