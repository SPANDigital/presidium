package cmd

import (
	"embed"
	"fmt"
	"github.com/gohugoio/hugo/commands"
	"github.com/spf13/cobra"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

//go:embed hugo-packages/*
var hugoPackages embed.FS

type Service struct{}

func New() Service {
	return Service{}
}

func (s Service) Execute(args ...string) {
	commands.Execute(args)
}

var (
	hugoCommand = &cobra.Command{
		Use:   "hugo",
		Short: "Runs hugo against your presidium site",
		Run: func(cmd *cobra.Command, args []string) {
			hugo := New()
			// hugo.Execute(args...)

			//Verify embedded packages

			// This is failing to read the go.mod file
			//data, err := hugoPackages.ReadFile("hugo-packages/github.com/spandigital/presidium-layouts-base@v0.2.1/go.mod")
			//if err != nil {
			//	log.Fatalf("Error reading embedded file: %v", err)
			//}
			//fmt.Println("File content:", string(data))

			dirs, err := hugoPackages.ReadDir("hugo-packages/github.com")
			if err != nil {
				log.Fatal(err)
			}
			for _, dir := range dirs {
				if dir.IsDir() {
					fmt.Println("Directory name:", dir.Name())
				}
			}

			// Create a temporary directory to extract the embedded packages
			tempDir, err := os.MkdirTemp("", "hugo-packages")
			if err != nil {
				log.Fatal(err)
			}
			defer os.RemoveAll(tempDir)

			// Extract the embedded packages to the temporary directory
			err = fs.WalkDir(hugoPackages, ".", func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() {
					return nil
				}
				data, err := hugoPackages.ReadFile(path)
				if err != nil {
					return err
				}
				destPath := filepath.Join(tempDir, path)
				os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
				return os.WriteFile(destPath, data, os.ModePerm)
			})
			if err != nil {
				log.Fatal(err)
			}

			// Set the HUGO_MODULES_REPO environment variable to the temporary directory
			os.Setenv("HUGO_MODULES_REPO", tempDir)

			// Run Hugo with the embedded packages
			hugo.Execute(args...)
		},
	}
)

func init() {
	rootCmd.AddCommand(hugoCommand)
}
