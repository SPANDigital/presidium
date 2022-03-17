package markdown

import (
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

type Markdown struct {
	FrontMatter FrontMatter
	Content     string
}

var af = afero.NewOsFs()

func Parse(path string) (*Markdown, error) {
	b, err := afero.ReadFile(af, path)
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
