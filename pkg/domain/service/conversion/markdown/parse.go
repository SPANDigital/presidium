package markdown

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Markdown struct {
	FrontMatter FrontMatter
	Content     string
}

func Parse(path string) (*Markdown, error) {
	b, err := ioutil.ReadFile(path) // just pass the file name
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
	return nil, nil
}
