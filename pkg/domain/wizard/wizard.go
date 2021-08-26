package wizard

import "github.com/manifoldco/promptui"

const (
	TrueValue  = "y"
	FalseValue = "N"
)

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

func GetConfirmationFromUser(label string, defaultValue bool) (bool, error) {
	def := FalseValue
	if defaultValue {
		def = TrueValue
	}
	prompt := promptui.Prompt{
		Label:     label,
		Default:   def,
		IsConfirm: true,
	}
	confirm, err := prompt.Run()
	if err != nil || confirm != TrueValue {
		return false, nil
	} else {
		return true, nil
	}
}

func GetSelectTemplate() *promptui.SelectTemplates {
	return &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U00002714 {{ .Name | cyan }} - ({{ .Description }})",
		Inactive: "  {{ .Name | cyan }} - ({{ .Description }})",
		Selected: "\U00002714 {{ .Name | cyan }} - ({{ .Description }})",
	}
}
