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
)

type InitConfig struct {
	ProjectName  string `json:"project_name"`
	TemplateName string `json:"template_name"`
	Theme        string `json:"theme"`
	Brand        string `json:"brand"`
}
