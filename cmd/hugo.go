package cmd

import (
	"github.com/gohugoio/hugo/commands"
	"github.com/spf13/cobra"
)

type Service struct{}

func New() Service {
	return Service{}
}
func (s Service) Execute(args ...string) {
	commands.Execute(args)
}

var (
	// hugoCommand wraps hugo into Presidium.  This allows you to run hugo
	// in Presidium, and makes it easier to debug etc.
	hugoCommand = &cobra.Command{
		Use:   "hugo",
		Short: "Runs hugo against your presidium site",
		Run: func(cmd *cobra.Command, args []string) {
			hugo := New()
			hugo.Execute(args...)
		},
	}
)

func init() {
	rootCmd.AddCommand(hugoCommand)
}
