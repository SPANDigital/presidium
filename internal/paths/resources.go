package paths

import (
	"fmt"
	"github.com/SPANDigital/presidium-hugo/internal/colors"
	"github.com/SPANDigital/presidium-hugo/internal/settings"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var numberedPathRe = regexp.MustCompile(`(\d+\-)(.*)`)
var resourceList []string

func GatherResources(path string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && !strings.HasSuffix(path, ".md"){
			var resource = finalResourcePath(path)
			fmt.Println("Found resource",colors.Labels.Info(resource))
			resourceList = append(resourceList, resource)
		}
		return nil
	})
}

func finalResourcePart(part string) string {
	matches := numberedPathRe.FindStringSubmatch(part)
	if matches != nil {
		return matches[2]
	}
	return part
}

func finalResourcePath(path string) string {
	parts := strings.Split(strings.TrimPrefix(path, settings.Flags.StagingDir), string(filepath.Separator))
	for idx, part := range parts {
		parts[idx] = finalResourcePart(part)
	}
	return filepath.Join(parts...)
}
