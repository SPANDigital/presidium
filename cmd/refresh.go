package cmd

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/hugo"
	"github.com/SPANDigital/presidium-hugo/pkg/filesystem"
	"github.com/SPANDigital/presidium-hugo/pkg/log"
	"github.com/spf13/cobra"
)

var (
	noBuild bool

	refreshCommand = &cobra.Command{
		Use:   "refresh",
		Short: "Clean build artifacts and refresh Hugo modules",
		Long: `Clean all build artifacts (public, resources, themes, etc.) and refresh Hugo modules.
Optionally rebuilds the site after cleaning.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := runRefresh(); err != nil {
				log.Error(err)
				os.Exit(1)
			}
		},
	}
)

func init() {
	refreshCommand.Flags().BoolVar(&noBuild, "no-build", false, "Skip building after refresh")
	rootCmd.AddCommand(refreshCommand)
}

func runRefresh() error {
	fsUtil := filesystem.New()
	cwd, err := fsUtil.GetWorkingDir()
	if err != nil {
		return err
	}

	// Directories and files to clean
	artifacts := []string{"public", "resources", "themes", "go.sum", ".hugo_build.lock"}

	log.Info("Cleaning build artifacts...")
	for _, artifact := range artifacts {
		path := filepath.Join(cwd, artifact)
		if _, err := os.Stat(path); err == nil {
			log.InfoWithFields("removing", log.Fields{"path": artifact})
			if err := os.RemoveAll(path); err != nil {
				log.WarnWithFields("failed to remove", log.Fields{"path": artifact, "error": err})
			}
		}
	}
	log.Info("Build artifacts cleaned.")

	// Run hugo mod clean
	log.Info("Running hugo mod clean...")
	hugoService := hugo.New()
	if err := hugoService.Execute("mod", "clean"); err != nil {
		log.WarnWithFields("hugo mod clean failed", log.Fields{"error": err})
	}

	// Run go mod tidy
	log.Info("Running go mod tidy...")
	goTidy := exec.Command("go", "mod", "tidy")
	goTidy.Stdout = os.Stdout
	goTidy.Stderr = os.Stderr
	if err := goTidy.Run(); err != nil {
		log.WarnWithFields("go mod tidy failed", log.Fields{"error": err})
	}

	// Run hugo mod tidy
	log.Info("Running hugo mod tidy...")
	if err := hugoService.Execute("mod", "tidy"); err != nil {
		log.WarnWithFields("hugo mod tidy failed", log.Fields{"error": err})
	}

	// Run hugo mod get
	log.Info("Running hugo mod get...")
	if err := hugoService.Execute("mod", "get"); err != nil {
		log.WarnWithFields("hugo mod get failed", log.Fields{"error": err})
	}

	// Optionally build
	if !noBuild {
		log.Info("Building site...")
		if err := hugoService.Execute("--templateMetrics", "--ignoreCache", "--logLevel", "info"); err != nil {
			return err
		}
	}

	log.Info("Refresh complete!")
	return nil
}
