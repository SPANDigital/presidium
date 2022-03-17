package markdown

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"

	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/conversion/colors"
	"github.com/spf13/viper"
)

var excludes map[string]bool

type FrontMatter struct {
	Title  string `yaml:"title,omitempty"`
	Slug   string `yaml:"slug,omitempty"`
	URL    string `yaml:"url,omitempty"`
	Weight string `yaml:"weight,omitempty"`
	Author string `yaml:"author,omitempty"`
	Github string `yaml:"github,omitempty"`
	Status string `yaml:"status,omitempty"`
	Roles string `yaml:"roles,omitempty"`
}

// SetupExcludes initialize excludes from Viper
func SetupExcludes() {
	excludes = make(map[string]bool)
	excludes["url"] = !viper.GetBool("urlBasedOnFilename")
	excludes["weight"] = !viper.GetBool("weightBasedOnFilename")
}

// AddFrontMatter Add front matter keys and values to an existing markdown file
func AddFrontMatter(path string, fm FrontMatter) error {
	fmt.Println("Adding front matter", colors.Labels.Wanted(fm), "to", colors.Labels.Info(path))
	return ManipulateMarkdown(path, func(frontMatter []byte, w io.Writer) error {
		out, err := yaml.Marshal(fm)
		if err != nil {
			return err
		}

		_, err = io.WriteString(w, string(out))
		if err != nil {
			return err
		}

		return nil
	}, nil)
}
