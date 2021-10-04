package cmd

import "github.com/spf13/cobra"

var (
	reportCommand = &cobra.Command{
		Use:   "report [command]",
		Short: "Provides some validation checks",
	}
)

func init() {
	rootCmd.AddCommand(reportCommand)
	reportCommand.AddCommand(pageLinksCommand)
}
