package initwzd

import (
	"github.com/SPANDigital/presidium-hugo/pkg/config"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/generator"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/wizard"
	"github.com/SPANDigital/presidium-hugo/pkg/log"
	"github.com/SPANDigital/presidium-hugo/pkg/presidiumerr"
	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
	"golang.org/x/mod/module"
	"strings"
)

var (
	supportedTemplates = []wizard.Template{
		wizard.SpanTemplate,
		wizard.OnBoardingTemplate,
		wizard.DesignTemplate,
	}
	supportedThemes = []wizard.Theme{
		wizard.PresidiumTheme,
	}
)

type initWizard struct {
}

func New() wizard.Wizard {
	return &initWizard{}
}

func (i initWizard) Run() {
	err := getProjectName()
	if err != nil {
		return
	}
	err = getTitle()
	if err != nil {
		return
	}
	promptSupportedTemplates()
	promptSupportedThemes()
	err = getBrandRepo()
	if err != nil {
		log.Error(err)
		return
	}
	g := generator.New()
	err = g.Generate()
	if err != nil {
		log.Error(err)
		return
	}
}

func getTitle() error {
	validate := func(input string) error {
		if len(input) <= 0 {
			return presidiumerr.GenericError{Code: presidiumerr.InvalidTitle}
		}
		return nil
	}
	title, err := wizard.GetInputString("Title (The name of your documentation site)", "", validate)
	if err != nil {
		return err
	}
	viper.Set(config.TitleKey, title)
	return nil
}

func getProjectName() error {
	validate := func(input string) error {
		if len(input) <= 0 || strings.Contains(input, " ") {
			return presidiumerr.GenericError{Code: presidiumerr.InvalidProjectName}
		}
		return nil
	}
	projectName, err := wizard.GetInputString("Project Name (May NOT contain spaces)", "", validate)
	if err != nil {
		return err
	}
	viper.Set(config.ProjectNameKey, projectName)
	return nil
}

func getBrandRepo() error {

	isBrand, err := wizard.GetConfirmationFromUser("Do you want to add a brand?", false)
	if err != nil {
		log.Error(err)
		return err
	}

	if isBrand {
		validate := func(input string) error {
			return module.CheckPath(input)
		}
		repoURL, err := wizard.GetInputString("Provide your go module for branding", "", validate)
		if err != nil {
			return err
		}
		viper.Set(config.BrandKey, repoURL)
	}
	return nil
}

func promptSupportedTemplates() {
	items := make([]wizard.Item, 0)
	for _, item := range supportedTemplates {
		items = append(items, wizard.Item{
			Name:        item.Name(),
			Description: item.Description(),
		})
	}
	prompt := promptui.Select{
		Label:     "Select a template",
		Items:     items,
		Templates: wizard.GetSelectTemplate(),
	}

	idx, _, err := prompt.Run()
	if err != nil {
		log.FatalWithFields("error selecting template", log.Fields{"error": err})
	}
	selected := supportedTemplates[idx]
	viper.Set(config.TemplateNameKey, selected.Code())
}

func promptSupportedThemes() {
	items := make([]wizard.Item, 0)
	for _, item := range supportedThemes {
		items = append(items, wizard.Item{
			Name:        item.Name(),
			Description: item.Description(),
		})
	}
	prompt := promptui.Select{
		Label:     "Select a theme",
		Items:     items,
		Templates: wizard.GetSelectTemplate(),
	}

	idx, _, err := prompt.Run()
	if err != nil {
		log.FatalWithFields("error selecting theme", log.Fields{"error": err})
	}
	selected := supportedThemes[idx]
	viper.Set(config.ThemeKey, selected.Code())
}
