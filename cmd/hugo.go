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

	// serverCommand provides a direct way to run 'hugo server'
	serverCommand = &cobra.Command{
		Use:   "server",
		Short: "Start the Hugo development server",
		Long:  "Start the Hugo development server with live reload and other development features",
		Run: func(cmd *cobra.Command, args []string) {
			hugoService := hugo.New()
			err := hugoService.Execute(append([]string{"server"}, args...)...)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
		},
	}

	// newCommand provides a direct way to run 'hugo new'
	newCommand = &cobra.Command{
		Use:   "new",
		Short: "Create new content for your Hugo site",
		Long:  "Create new content for your Hugo site, such as posts, pages, etc.",
		Run: func(cmd *cobra.Command, args []string) {
			hugoService := hugo.New()
			err := hugoService.Execute(append([]string{"new"}, args...)...)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
		},
	}

	// versionCommand provides a direct way to run 'hugo version'
	hugoVersionCommand = &cobra.Command{
		Use:   "hugo-version",
		Short: "Print the version number of Hugo",
		Long:  "Print the version number of Hugo that Presidium is using",
		Run: func(cmd *cobra.Command, args []string) {
			hugoService := hugo.New()
			err := hugoService.Execute(append([]string{"version"}, args...)...)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(hugoCommand)
	rootCmd.AddCommand(serverCommand)
	rootCmd.AddCommand(newCommand)
	rootCmd.AddCommand(hugoVersionCommand)
}
