package markdown

import (
	"fmt"
	. "github.com/SPANDigital/presidium-hugo/pkg/filesystem"
	"io"
	"os"
)

// WriteFrontMatterFunc are for callbacks which allow you to customize
// the front matter of markdown files. It's is passed the entire frontmatter
// as a byte array, and a writer send your optionally modified front matter
// too
type WriteFrontMatterFunc func(frontMatter []byte, w io.Writer) error

// WriteContentFunc are for callbacks which allow you to customize
// the contentOf of markdown files. It's is passed the entire contentOf
// as a byte array, and a writer send your optionally modified contentOf
// too
type WriteContentFunc func(content []byte, w io.Writer) error

func IsRecognizableMarkdown(path string) bool {
	fmt.Println("Validating", path)
	b, err := FSUtil.ReadFile(path) // just pass the file name
	if err != nil {
		return false
	}
	matches := MarkdownRe.FindSubmatch(b)
	return matches != nil
}

// Checks if a markdown file exists, if it doesn't create an empty one
func touch(path string) error {
	_, err := FSUtil.Stat(path)
	if os.IsNotExist(err) {
		f, err := FSUtil.Create(path)
		if err == nil {
			_, err = f.WriteString("---\n---\n")
			if err == nil {
				err = f.Close()
			}
		}
		return err
	}
	return nil
}

// ManipulateMarkdown Manipulate a markdown file with 2 optional callbacks
// matterFunc - callback to manipulate front matter
// contentFunc - callback to manipulate contentOf
// if the file doesn't exist, it will create
// if frontmatter is not present, then a front matter section is added
func ManipulateMarkdown(path string, matterFunc WriteFrontMatterFunc, contentFunc WriteContentFunc) error {

	err := touch(path)
	if err != nil {
		return err
	}

	b, err := FSUtil.ReadFile(path) // just pass the file name
	if err != nil {
		return err
	}
	matches := MarkdownRe.FindSubmatch(b)

	if matches == nil {
		matches = [][]byte{
			// we don't care about matches[0]
			[]byte{},
			[]byte("---\n"),
			[]byte{},
			[]byte("---\n"),
			b,
		}
	}
	f, err := FS.Create(path)
	if err != nil {
		return err
	}
	f.Write(matches[1])
	if matterFunc != nil {
		err := matterFunc(matches[2], f)
		if err != nil {
			return err
		}
	} else {
		f.Write(matches[2])
	}
	f.Write(matches[3])
	if contentFunc != nil {
		err := contentFunc(matches[4], f)
		if err != nil {
			return err
		}
	} else {
		f.Write(matches[4])
	}
	return f.Close()

}
