package cmd

import (
	"fmt"
	model "github.com/SPANDigital/presidium-hugo/pkg/domain/model/validate"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/validate"
	"github.com/spf13/cobra"
)

var (
	validateCommand = &cobra.Command{
		Use:   "validate",
		Short: "Validates page links in a Presidium site",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			report, err := validate.New(path).Validate()
			if err != nil {
				fmt.Printf("error validating : %s\n", path)
				fmt.Printf("%v\n", err.Error())
				return
			}

			fmt.Printf("\n")
			fmt.Printf("VALIDATION PATH: %s\n", path)
			fmt.Printf("\n")
			fmt.Printf("        total: %v\n", report.TotalLinks)
			fmt.Printf("  valid links: %v\n", report.Valid)
			fmt.Printf("       broken: %v\n", report.Broken)
			fmt.Printf("     external: %v\n", report.External)
			fmt.Printf("     warnings: %v\n", report.Warning)
			fmt.Printf("\n")

			printLinks(report, model.Broken)
			printLinks(report, model.Warning)
			printLinks(report, model.External)
		},
	}
)

func printLinks(report model.Report, status model.Status) {

	links, found := report.Data[status]

	if !found {
		return
	}

	fmt.Printf("%s\n", status)
	fmt.Printf("----------------------\n")

	for _, link := range links {
		message := ""
		if len(link.Message) > 0 {
			message = fmt.Sprintf(" %s", link.Message)
		}
		fmt.Printf("%s: %s [%s]%s\n", status, link.Uri, link.Label, message)
	}
}

func init() {
	rootCmd.AddCommand(validateCommand)
}
