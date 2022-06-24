package cmd

import (
	"github.com/SPANDigital/presidium-hugo/pkg/domain/wizard/initwzd"
	"github.com/spf13/cobra"
)

var (
	initCommand = &cobra.Command{
		Use:   "init",
		Short: "Init a new Presidium site",
		Run: func(cmd *cobra.Command, args []string) {
			initWizard := initwzd.New()
			initWizard.Run()
		},
	}
)

func init() {
	rootCmd.AddCommand(initCommand)
}
