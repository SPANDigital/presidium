package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	model "github.com/SPANDigital/presidium-hugo/pkg/domain/model/validate"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/validate"
	"github.com/SPANDigital/presidium-hugo/pkg/log"
	"github.com/spf13/cobra"
)

var (
	pageLinksCommand = &cobra.Command{
		Use:   "pagelinks",
		Short: "Validates page links in a Presidium site",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			if !filepath.IsAbs(path) {
				cwd, err := os.Getwd()
				if err != nil {
					log.Fatal("failed to get working directory: ", err.Error())
				}
				path = filepath.Join(cwd, path)
			}

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

			printLinks(path, report, model.Broken)
			printLinks(path, report, model.Warning)
			printLinks(path, report, model.External)
		},
	}
)

func printLinks(path string, report model.Report, status model.Status) {

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

		uri := strings.TrimPrefix(link.Uri, path)
		location := strings.TrimPrefix(link.Location, path)

		fmt.Printf("%s: %s\nlabel: [%s]\noutput file: %s\nsource file: %s\nmessage:%s\n========================\n", status, uri, link.Label, location, link.DataId, message)
	}
}
