package config

const (
	ApplicationName = "presidium"
	DebugKey        = "debug"
	ConfigFileKey   = "config"
	LoggingLevelKey = "logging.level"

	ProjectNameKey  = "init.project_name"
	TemplateNameKey = "init.template_name"
	ThemeKey        = "init.theme"
	TitleKey        = "init.title"
	BrandKey        = "init.brand"
	MarkupKey       = "init.markup"
)

type InitConfig struct {
	ProjectName  string `json:"project_name"`
	TemplateName string `json:"template_name"`
	Theme        string `json:"theme"`
	Brand        string `json:"brand"`
	Markup       string `json:"markup"`
}

type GeneratorConfig struct {
	Title       string `json:"title"`        // the actual "name" of the site - may contain spaces
	ProjectName string `json:"project_name"` // the folder the site must be generated to - must never contain spaces
	Theme       string `json:"theme"`        // theme module code
	Template    string `json:"template"`     // template code
	Brand       string `json:"brand"`        // the url of the repo to the brand module
	Markup      string `json:"markup"`       // the markup style selector
}
