package cmd

import (
	"github.com/SPANDigital/presidium-hugo/pkg/domain/versioning"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCommand())
}

func versionCommand() *cobra.Command {
	enabled := false
	cmd := &cobra.Command{
		Use:   "versioning",
		Short: "managing versioning of presidium site",
		Run: func(cmd *cobra.Command, args []string) {
			if !enabled {
				v:= versioning.New(".")
				if !v.IsEnabled() {
					v.SetEnabled(true)
				}
			}
		},
	}
	cmd.Flags().BoolVar(&enabled, "enable", false, "enabled version of presidium site")
	cmd.AddCommand(activateNextVersionCommand(), syncLatestCommand())
	return cmd
}


func activateNextVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use: "next",
		Short: "creates a next version based on last version",
		Run: func(cmd *cobra.Command, args []string) {
			v := versioning.New(".")
			v.NextVersion()
		},
	}
}


func syncLatestCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Updates the current version based on your content as is.",
		Run: func(cmd *cobra.Command, args []string) {
			v := versioning.New(".")
			v.GrabLatest()
		},
	}
}


