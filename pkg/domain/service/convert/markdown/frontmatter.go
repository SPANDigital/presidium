package markdown

import (
	"fmt"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/convert/colors"
	"github.com/spf13/viper"
	"io"
)

var excludes map[string]bool

// initialize excludes from Viper
func SetupExcludes() {
	excludes = make(map[string]bool)
	excludes["url"] = !viper.GetBool("urlBasedOnFilename")
	excludes["weight"] = !viper.GetBool("weightBasedOnFilename")
	excludes["slug"] = !viper.GetBool("slugBasedOnFilename")
}

// Add front matter keys and values to an existing markdown file
func AddFrontMatter(path string, params map[string]interface{}) error {

	fmt.Println("Adding front matter", colors.Labels.Wanted(params), "to", colors.Labels.Info(path))

	return ManipulateMarkdown(path, func(frontMatter []byte, w io.Writer) error {
		_, err := w.Write(frontMatter)
		if err != nil {
			return err
		}
		for key, value := range params {
			if !excludes[key] {
				_, err := io.WriteString(w, fmt.Sprintf("%s: %s\n", key, value))
				if err != nil {
					return err
				}
			}
		}
		return nil
	}, nil)
	return nil
}
