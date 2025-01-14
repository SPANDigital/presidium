package initwzd

import (
	"strings"

	"github.com/SPANDigital/presidium-hugo/pkg/config"
	. "github.com/SPANDigital/presidium-hugo/pkg/domain/model/generator"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/generator"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/wizard"
	"github.com/SPANDigital/presidium-hugo/pkg/log"
	"github.com/SPANDigital/presidium-hugo/pkg/presidiumerr"
	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
	"golang.org/x/mod/module"
)

type (
	initWizard struct{}
)

func New() wizard.Wizard {
	return &initWizard{}
}

func (i initWizard) Run() {

	err := askProjectName()
	if err != nil {
		return
	}

	err = askTitle()
	if err != nil {
		return
	}

	promptSupportedTemplates()

	err = askBrandRepo()
	if err != nil {
		log.Error(err)
		return
	}

	g := generator.New()

	siteModel := generateSiteModel()
	err = g.Run(siteModel)
	if err != nil {
		log.Error(err)
	}

}

func generateSiteModel() InitialSiteTarget {

	mustHaveTemplate := func() Template {
		templateName := viper.GetString(config.TemplateNameKey)
		template, err := GetTemplate(templateName)
		if err != nil {
			log.FatalWithFields(err, log.Fields{
				"template_name": templateName,
			})
		}
		return template
	}

	return InitialSiteTarget{
		SiteTargetDirectory: viper.GetString(config.ProjectNameKey),
		SiteName:            viper.GetString(config.ProjectNameKey),
		SiteTitle:           viper.GetString(config.TitleKey),
		BrandingModelUrl:    viper.GetString(config.BrandKey),
		Template:            mustHaveTemplate(),
		WhenSiteExists:      AbortWhenTargetSiteExists,
	}
}

func askTitle() error {
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

func askProjectName() error {
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

func askBrandRepo() error {

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
	items := make([]ItemSelection, 0)
	for _, item := range SupportedTemplates {
		items = append(items, ItemSelection{
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
	selected := SupportedTemplates[idx]
	viper.Set(config.TemplateNameKey, selected.Code())
}
