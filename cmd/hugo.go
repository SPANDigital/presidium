package cmd

import (
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/hugo"
	"github.com/spf13/cobra"
)

var (
	hugoCommand = &cobra.Command{
		Use:   "hugo",
		Short: "Runs hugo against your presidium site",
		Run: func(cmd *cobra.Command, args []string) {
			hugo := hugo.New()
			hugo.Execute(args...)
		},
	}
)

func init() {
	rootCmd.AddCommand(hugoCommand)
}
