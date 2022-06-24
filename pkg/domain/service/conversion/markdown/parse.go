package markdown

import (
	"github.com/SPANDigital/presidium-hugo/pkg/filesystem"
	"gopkg.in/yaml.v2"
)

type Markdown struct {
	FrontMatter FrontMatter
	Content     string
}

func Parse(path string) (*Markdown, error) {
	b, err := filesystem.AFS.ReadFile(path)
	if err != nil {
		return nil, err
	}

	matches := MarkdownRe.FindSubmatch(b)
	if matches != nil {
		var fm FrontMatter
		err = yaml.Unmarshal(matches[2], &fm)
		if err != nil {
			return nil, err
		}

		return &Markdown{
			FrontMatter: fm,
			Content:     string(matches[4]),
		}, nil
	}
	return &Markdown{}, nil
}
