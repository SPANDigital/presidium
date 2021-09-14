package cmd

import (
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/searchmap"
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
	failureExitCode          = 3
	doNotExitWithError       = false
	validateSearchMapCommand = &cobra.Command{
		Use:   "validate-searchmap",
		Short: "Validates the searchmap.json file in presidium site.",
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

			validation, err := searchmap.New(projectDir)
			if err != nil {
				log.Fatal(err)
			}

			err = validation.Run()
			if err != nil {
				log.Fatal(err)
			}

			if validation.Failed() {
				println("[ERROR]")
				println("The following markdown files are not included in the searchmap.json file:")
				println("-------------------------------------------------------------------------")
				for _, missingMarkdownFile := range validation.MissingMarkdownFiles {
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
	validateSearchMapCommand.Flags().IntVar(&failureExitCode, flagFailureExitCode, failureExitCode, usageFailureExitCode)
	validateSearchMapCommand.Flags().BoolVar(&doNotExitWithError, flagDoNotExitWithErrorCode, doNotExitWithError, usageDoNotExitErrorWithCode)
	rootCmd.AddCommand(validateSearchMapCommand)
}
