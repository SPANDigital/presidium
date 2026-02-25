package cmd

import (
	"os"

	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/hugo"
	"github.com/SPANDigital/presidium-hugo/pkg/log"
	"github.com/spf13/cobra"
)

var (
	// hugoCommand wraps hugo into Presidium.  This allows you to run hugo
	// in Presidium, and makes it easier to debug etc.
	// All arguments and flags are passed through to Hugo unchanged.
	hugoCommand = &cobra.Command{
		Use:                "hugo",
		Short:              "Runs hugo with full access to all Hugo commands and flags",
		DisableFlagParsing: true, // Don't parse flags - let Hugo handle them
		Run: func(cmd *cobra.Command, args []string) {
			hugoService := hugo.New()
			err := hugoService.Execute(args...)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(hugoCommand)
}
