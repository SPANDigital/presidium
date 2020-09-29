package markdown

import "io/ioutil"

type Markdown struct {
	FrontMatter map[string]string
	Content string
}

func Parse(path string) (*Markdown, error) {
	b, err := ioutil.ReadFile(path) // just pass the file name
	if err != nil {
		return nil, err
	}
	matches := MarkdownRe.FindSubmatch(b)
	if matches != nil {
		allFmMatches := FrontmatterRe.FindAllSubmatch(matches[2], -1)
		fm := make(map[string]string, len(allFmMatches))
		for _, fmMatches := range allFmMatches {
			fm[string(fmMatches[1])] = string(fmMatches[2])

		}
		return &Markdown{
			FrontMatter: fm,
			Content: string(matches[4]),
		}, nil
	}
	return nil, nil
}
