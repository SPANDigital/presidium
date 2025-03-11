package cmd

import (
	"embed"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

//go:embed hugo-files/*
var hugoFiles embed.FS

type Service struct{}

func New() Service {
	return Service{}
}

func (s Service) Execute(args ...string) {
	// Write the embedded Hugo binary to a temporary file
	hugoBinaryPath := filepath.Join(os.TempDir(), "hugo")
	hugoBinary, err := hugoFiles.ReadFile("hugo-files/hugo")
	if err != nil {
		log.Fatalf("Failed to read embedded Hugo binary: %v", err)
	}
	if err := os.WriteFile(hugoBinaryPath, hugoBinary, 0755); err != nil {
		log.Fatalf("Failed to write Hugo binary to temporary file: %v", err)
	}

	cmd := exec.Command(hugoBinaryPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to run Hugo: %v", err)
	}
}

var (
	hugoCommand = &cobra.Command{
		Use:   "hugo",
		Short: "Runs hugo against your presidium site",
		Run: func(cmd *cobra.Command, args []string) {
			hugo := New()

			// Create a temporary directory to extract the embedded packages
			tempDir, err := os.MkdirTemp("", "hugo-files")
			if err != nil {
				log.Fatal(err)
			}
			defer os.RemoveAll(tempDir)

			// Extract the embedded packages to the temporary directory
			err = fs.WalkDir(hugoFiles, ".", func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() {
					return nil
				}
				data, err := hugoFiles.ReadFile(path)
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
