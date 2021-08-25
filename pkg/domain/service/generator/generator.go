package generator

import (
	"github.com/SPANDigital/presidium-hugo/pkg/config"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/template"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/wizard"
	"github.com/SPANDigital/presidium-hugo/pkg/log"
	"github.com/spf13/viper"
	"io/fs"
	"os"
	"path"
)

type Generator struct {
}

func New() Generator {
	return Generator{}
}

// GenerateWithConfig generates a presidium site with the config given as parameter
func (g Generator) GenerateWithConfig(c Config) error {
	theTemplate, err := wizard.GetTemplate(c.Template)
	if err != nil {
		return err
	}
	tplSvc := template.New()
	err = tplSvc.ProcessDirTemplates(theTemplate.Code(), c)
	if err != nil {
		log.Error(err)
		return err
	}
	err = os.Mkdir(path.Join(c.ProjectName, "static"), fs.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// Generate uses viper to populate the default config and execute GenerateWithConfig
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
		Template:    viper.GetString(config.TemplateNameKey),
	}

	return g.GenerateWithConfig(c)
}
