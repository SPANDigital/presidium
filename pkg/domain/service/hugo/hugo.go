package hugo

import (
	"fmt"
	"os"

	"github.com/SPANDigital/presidium-hugo/pkg/configtranslation"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/themes"
	"github.com/SPANDigital/presidium-hugo/pkg/filesystem"
	"github.com/SPANDigital/presidium-hugo/pkg/log"
	"github.com/gohugoio/hugo/commands"
)

const (
	moduleStylingBase = "github.com/spandigital/presidium-styling-base"
	moduleLayoutsBase = "github.com/spandigital/presidium-layouts-base"
)

type Service struct {
}

func New() Service {
	return Service{}
}

func (s Service) Execute(args ...string) error {
	// Initialize themes service
	themesService, err := themes.New()
	if err != nil {
		log.Warn(fmt.Sprintf("Failed to initialize themes service, falling back to remote theme fetch: %v", err))
		return commands.Execute(args)
	}

	// Extract embedded themes to temp directory
	tmpDir, replacements, err := themesService.Extract()
	if err != nil {
		log.Warn(fmt.Sprintf("Failed to extract embedded themes, falling back to remote theme fetch: %v", err))
		return commands.Execute(args)
	}

	// Capture the previous value (if any) so we can restore it after execution
	prevReplacements, hadPrevReplacements := os.LookupEnv("HUGO_MODULE_REPLACEMENTS")

	// Ensure cleanup happens regardless of how Hugo execution ends
	defer func() {
		_ = filesystem.AFS.RemoveAll(tmpDir)
		if hadPrevReplacements {
			os.Setenv("HUGO_MODULE_REPLACEMENTS", prevReplacements)
		} else {
			os.Unsetenv("HUGO_MODULE_REPLACEMENTS")
		}
	}()

	// Set environment variable to point Hugo to local themes
	os.Setenv("HUGO_MODULE_REPLACEMENTS", replacements)

	if err := validateModuleImportOrder(configPath()); err != nil {
		return err
	}

	// Execute Hugo with local themes - return the error to preserve Hugo's exit behavior
	return commands.Execute(args)
}

// configPath returns the path to the Hugo config file in the working directory.
// Returns an empty string if neither config.yaml nor config.yml exists.
func configPath() string {
	for _, name := range []string{"config.yaml", "config.yml"} {
		if _, err := os.Stat(name); err == nil {
			return name
		}
	}
	return ""
}

// validateModuleImportOrder checks that presidium-styling-base appears before
// presidium-layouts-base in the Hugo module imports. Returns nil when either
// module is absent — no opinion on configs that do not use both.
// Config read errors are silently ignored; Hugo will report them downstream.
func validateModuleImportOrder(configFile string) error {
	cfg, err := configtranslation.ReadHugoConfig(configFile)
	if err != nil {
		return nil
	}

	stylingIdx := -1
	layoutsIdx := -1
	for i, imp := range cfg.Module.Imports {
		switch imp.Path {
		case moduleStylingBase:
			stylingIdx = i
		case moduleLayoutsBase:
			layoutsIdx = i
		}
	}

	if stylingIdx == -1 || layoutsIdx == -1 {
		return nil
	}

	if layoutsIdx < stylingIdx {
		return fmt.Errorf(
			"invalid module import order: %q (index %d) must come before %q (index %d).\n"+
				"Fix config.yaml: list %s before %s.",
			moduleStylingBase, stylingIdx,
			moduleLayoutsBase, layoutsIdx,
			moduleStylingBase, moduleLayoutsBase,
		)
	}

	return nil
}
