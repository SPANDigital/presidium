package generator

import (
	"github.com/SPANDigital/presidium-hugo/pkg/config"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/template"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/wizard"
	"github.com/SPANDigital/presidium-hugo/pkg/log"
	"github.com/spf13/viper"
)

type Generator struct {
}

func New() Generator {
	return Generator{}
}

func (g Generator) Generate() error {
	themeKey := viper.GetString(config.ThemeKey)
	theme, err := wizard.GetTheme(themeKey)
	if err != nil {
		return err
	}
	c := Config{
		Title:       viper.GetString(config.TitleKey),
		ProjectName: viper.GetString(config.ProjectNameKey),
		Theme:       theme.ModulePath(),
	}

	templateKey := viper.GetString(config.TemplateNameKey)
	theTemplate, err := wizard.GetTemplate(templateKey)
	if err != nil {
		log.Error(err)
		return err
	}

	tplSvc := template.New()
	err = tplSvc.ProcessDirTemplates(theTemplate.Code(), c)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
