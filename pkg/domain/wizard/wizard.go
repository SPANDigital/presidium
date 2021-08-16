package wizard

import "github.com/manifoldco/promptui"

type Wizard interface {
	Run()
}

type SelectItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func GetInputString(label, defaultValue string, validateFn promptui.ValidateFunc) (string, error) {
	prompt := promptui.Prompt{
		Label:    label,
		Validate: validateFn,
		Default:  defaultValue,
	}
	result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return result, err
}

func GetSelectTemplate() *promptui.SelectTemplates {
	return &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U00002714 {{ .Name | cyan }} - ({{ .Description }})",
		Inactive: "  {{ .Name | cyan }} - ({{ .Description }})",
		Selected: "\U00002714 {{ .Name | cyan }} - ({{ .Description }})",
	}
}
