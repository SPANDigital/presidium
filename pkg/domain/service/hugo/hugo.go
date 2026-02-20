package hugo

import (
	"fmt"
	"os"

	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/themes"
	"github.com/SPANDigital/presidium-hugo/pkg/log"
	"github.com/gohugoio/hugo/commands"
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

	// Ensure cleanup happens regardless of how Hugo execution ends
	defer func() {
		os.RemoveAll(tmpDir)
		os.Unsetenv("HUGO_MODULE_REPLACEMENTS")
	}()

	// Set environment variable to point Hugo to local themes
	os.Setenv("HUGO_MODULE_REPLACEMENTS", replacements)

	// Execute Hugo with local themes - return the error to preserve Hugo's exit behavior
	return commands.Execute(args)
}
