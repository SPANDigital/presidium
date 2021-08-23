package cmd

import (
	"fmt"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/validation"
	"github.com/spf13/cobra"
)

var (
	validateCommand = &cobra.Command{
		Use:   "validate",
		Short: "Validates page links in a Presidium site",
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			report, err := validation.New(path).Validate()
			if err != nil {
				fmt.Printf("error validating : %s\n", path)
				fmt.Printf("%v\n", err.Error())
				return
			}

			margin := 2
			width := 0
			width = longestMessage(width, report.Warnings)
			width = longestMessage(width, report.Valid)
			width = longestMessage(width, report.Broken)
			width = longestMessage(width, report.External)

			printHeader("Validation Report", width, margin)
			printLine(width, margin)
			fmt.Printf("\n")

			broken := len(report.Broken)
			valid := len(report.Valid)
			warnings := len(report.Warnings)
			external := len(report.External)
			total := broken + valid + warnings + external
			fmt.Printf("VALIDATION PATH: %s\n", path)
			fmt.Printf("\n")
			fmt.Printf("     total: %v\n", total)
			fmt.Printf("    breken: %v\n", broken)
			fmt.Printf("  external: %v\n", external)
			fmt.Printf("  warnings: %v\n", warnings)

			fmt.Printf("\n")
			printLine(width, margin)

		},
	}
)

func init() {
	rootCmd.AddCommand(validateCommand)
}

func longestMessage(longest int, links []validation.Link) int {

	for _, link := range links {
		if len(link.Message) > longest {
			longest = len(link.Message)
		}
	}

	return longest
}

func printHeader(header string, width int, margin int) {
}

func printLine(width int, margin int) {
}