package cmd

import (
	"fmt"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/validation"
	"github.com/SPANDigital/presidium-hugo/pkg/log"
	"github.com/spf13/cobra"
)

var (
	validateCommand = &cobra.Command{
		Use:   "validatelinks",
		Short: "Validates page links in a Presidium site",
		Run: func(cmd *cobra.Command, args []string) {
			validation, err := validation.New(args[0], 1, func(link validation.Link) {
				status := " ok"
				if !link.Valid {
					status = "err"
				}
				fmt.Printf("%s: [%s]", status, link.Uri)
				if len(link.Text) > 0 {
					fmt.Printf(" (%s)", link.Text)
				}
				if !link.Valid {
					fmt.Printf("\t%s", link.Message)
				}
				fmt.Println()
			})
			if err != nil {
				log.ErrorWithFields("Unable to validate links on the Presidium site", log.Fields{
					"url":   args[0],
					"error": err.Error(),
				})
				return
			}
			fmt.Printf("Validating links: %s\n", args[0])
			validation.Start()
		},
	}
)

func init() {
	rootCmd.AddCommand(validateCommand)
}
