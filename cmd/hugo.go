package cmd

import (
	"archive/zip"
	"bytes"
	"embed"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed hugo-files/hugo
var hugoFiles embed.FS

//go:embed hugo-files/hugo-files.zip
var hugoFilesZip embed.FS

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

			// Extract the embedded zip file to the temporary directory
			zipFile, err := hugoFilesZip.Open("hugo-files/hugo-files.zip")
			if err != nil {
				log.Fatalf("Failed to open embedded zip file: %v", err)
			}
			defer zipFile.Close()

			zipData, err := io.ReadAll(zipFile)
			if err != nil {
				log.Fatalf("Failed to read embedded zip file: %v", err)
			}

			zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
			if err != nil {
				log.Fatalf("Failed to create zip reader: %v", err)
			}

			if err := unzip(zipReader, tempDir); err != nil {
				log.Fatalf("Failed to unzip embedded file: %v", err)
			}

			//// Walk the temporary directory and debug the files to screen
			//filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
			//	if err != nil {
			//		return err
			//	}
			//	log.Println(path)
			//	return nil
			//})

			// Set the GOPROXY environment variable to off
			os.Setenv("GOPROXY", "off")

			// Set the HUGO_MODULES_REPO environment variable to the temporary directory
			os.Setenv("GOPATH", tempDir)

			fmt.Printf("The GOPATH is", tempDir)

			// Run Hugo with the embedded packages
			hugo.Execute(args...)
		},
	}
)

func init() {
	rootCmd.AddCommand(hugoCommand)
}

// unzip extracts a zip file to the specified destination directory
func unzip(zipReader *zip.Reader, dest string) error {
	for _, f := range zipReader.File {
		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}
