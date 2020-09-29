package colors

import (
	"github.com/logrusorgru/aurora"
	"github.com/spf13/viper"
)

type labels struct {
	Underline func(arg interface{}) aurora.Value
	Warning   func(arg interface{}) aurora.Value
	Wanted    func(arg interface{}) aurora.Value
	Unwanted  func(arg interface{}) aurora.Value
	Info      func(arg interface{}) aurora.Value
}

var Labels labels

func Setup() {
	Labels = makeLabels(viper.GetBool("enableColor"))
}

func makeLabels(enabled bool) labels {
	au := aurora.NewAurora(enabled)
	var labels = labels{
		Underline: au.Underline,
		Warning:  au.Red,
		Wanted:   au.Green,
		Unwanted: au.Yellow,
		Info:     au.BrightBlue,
	}
	return labels
}
