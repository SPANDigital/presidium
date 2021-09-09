package colors

import (
	"github.com/logrusorgru/aurora"
	"github.com/spf13/viper"
)

type StyleLabel = func(arg interface{}) aurora.Value

type labels struct {
	Underline StyleLabel
	Warning   StyleLabel
	Wanted    StyleLabel
	Unwanted  StyleLabel
	Info      StyleLabel
	Normal    StyleLabel
}

var Labels labels

func Setup() {
	Labels = makeLabels(viper.GetBool("enableColor"))
}

func makeLabels(enabled bool) labels {
	au := aurora.NewAurora(enabled)
	var labels = labels{
		Underline: au.Underline,
		Warning:   au.Red,
		Wanted:    au.Green,
		Unwanted:  au.Yellow,
		Info:      au.BrightBlue,
		Normal:    au.Black,
	}
	return labels
}
