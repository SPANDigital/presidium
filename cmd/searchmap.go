package cmd

import (
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/searchmapvalidation"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
)

const (
	flagFailureExitCode         string = "failureExitCode"
	flagDoNotExitWithErrorCode  string = "doNotExitErrorWithCode"
	usageFailureExitCode        string = "exit code to use when validation failed."
	usageDoNotExitErrorWithCode string = "exit normal if validation failed."
)

var (
	failureExitCode    = 3
	doNotExitWithError = false
	searchMapCommand   = &cobra.Command{
		Use:   "searchmap",
		Short: "Validates the searchmap.json file in a presidium site.",
		Long:  "Ensures the `searchmap.json` file references all valid pages.",
		Run: func(cmd *cobra.Command, args []string) {

			projectDir := "."
			if len(args) > 0 {
				projectDir = args[0]
			}

			projectDir, err := filepath.Abs(projectDir)
			if err != nil {
				log.Fatal(err)
			}

			s := searchmapvalidation.New()

			undeclaredFiles, err := s.FindUndeclaredFiles(projectDir)

			if err != nil {
				log.Fatal(err)
			}

			if undeclaredFiles.Found {
				println("[ERROR]")
				println("The following markdown files are not included in the searchmap.json file:")
				println("-------------------------------------------------------------------------")
				for _, missingMarkdownFile := range undeclaredFiles.Files {
					println(missingMarkdownFile)
				}
				if doNotExitWithError {
					return
				}
				os.Exit(failureExitCode)
			} else {
				println("[OK]")
			}
		},
	}
)

func init() {
	searchMapCommand.Flags().IntVar(&failureExitCode, flagFailureExitCode, failureExitCode, usageFailureExitCode)
	searchMapCommand.Flags().BoolVar(&doNotExitWithError, flagDoNotExitWithErrorCode, doNotExitWithError, usageDoNotExitErrorWithCode)
}
