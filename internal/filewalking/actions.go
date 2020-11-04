package filewalking

import (
	"fmt"
	"github.com/SPANDigital/presidium-hugo/internal/colors"
	"github.com/SPANDigital/presidium-hugo/internal/markdown"
	"github.com/SPANDigital/presidium-hugo/internal/settings"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var idxRenameList []string

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func CheckForDirRename(path string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() && strings.HasPrefix(file.Name(), "_") {
			old := path + "/" + file.Name()
			new := path + "/" + strings.TrimLeft(file.Name(), "_")
			fmt.Println("Renaming", colors.Labels.Unwanted(old), "to", colors.Labels.Wanted(new))
			err := os.Rename(old, new)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func CheckForDirIndex(path string) error {
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		fmt.Println("Walking", colors.Labels.Info(path))
		if info.IsDir() {
			fmt.Println(fmt.Sprintf("Checking %s for _index.md...\n", colors.Labels.Wanted(path)))
			if fileExists(path + "/_index.md") {
				return nil
			}

			// If we find an index.md, we'll rename it later, as if we rename it now, we'll affect the Walk function
			// and get a nil pointer exception when we try "walk" into the old directory
			if fileExists(path + "/index.md") {
				fmt.Println("Adding path",colors.Labels.Info(path),"for later index rename")
				idxRenameList = append(idxRenameList, path)
				return nil
			}
			addIndex(path)
		} else { // is a file
			if strings.HasSuffix(path, ".md") {
				if info.Name() != "_index.md" && info.Name() != "index.md" {
					err := injectSlugWeightAndURL(path)
					if err != nil {
						return err
					}
				}
				err := markdown.Operate(path)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	for _, path := range idxRenameList {
		fmt.Println("Renaming", colors.Labels.Unwanted(fmt.Sprintf("%v/index.md", path)), "to", colors.Labels.Wanted("_index.md"))
		os.Rename(path+"/index.md", path+"/_index.md")
		err := injectSlugWeightAndURLForIndex(path+"/_index.md")
		if err != nil {
			return err
		}
	}
	return nil
}

func CheckIndexForTitles(path string) error {
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path,"_index.md") {
			indexMarkdown, err := markdown.Parse(path)
			if err != nil {
				return err
			}
			if _, ok := indexMarkdown.FrontMatter["title"]; !ok {
				base := filepath.Base(filepath.Dir(path))
				title := unSlugify(base)
				markdown.AddFrontMatter(path, map[string]string{"title": title})
			}
		}
		return nil
	})
	return nil
}

// unNumerify turns "02-employment-contracts" into "employment-contracts" and "bill-add-customer" into "bill-add-customer"
func deduceWeightAndSlug(path string) (int64, string, string) {
	re := regexp.MustCompile(`(([\d\.]+)\-)?([^..]+)(\.[^\..]*)?`)
	base := filepath.Base(path)
	matches := re.FindStringSubmatch(base)
	weight, err := strconv.ParseInt(strings.ReplaceAll(matches[2], ".", ""), 10, 64)
	if err != nil {
		weight = 0
	}
	var slug = matches[3]
	var url string
	var contentDir = filepath.Join(settings.Flags.StagingDir, "content")
	if path == contentDir {
		url = ""
	} else {
		_, _, parentUrl := deduceWeightAndSlug(filepath.Dir(path))
		if parentUrl != "" {
			url = parentUrl + "/" + slug
		} else {
			url = slug
		}
		replaceRoot := viper.GetString("replaceRoot")
		if strings.HasPrefix(replaceRoot, replaceRoot) {
			url = strings.TrimPrefix(url, replaceRoot)
			if url == "" {
				url = "/"
			}
		}
	}

	return weight, slug, url
}

func injectSlugWeightAndURL(path string) error {
	if (markdown.IsRecognizableMarkdown(path)) {
		fmt.Println("Checking weight of ", colors.Labels.Info(path))
		weight, slug, url := deduceWeightAndSlug(path)

		m := make(map[string]string)
		m["slug"] = slug
		m["url"] = url
		if weight > 0 {
			m["weight"] = fmt.Sprintf("%d", weight)
		}

		return markdown.AddFrontMatter(path, m)
	}
	fmt.Println("Is not valid markdown", colors.Labels.Info(path))
	return nil
}

func injectSlugWeightAndURLForIndex(indexFile string) error {
	dir := filepath.Dir(indexFile)
	weight, slug, url := deduceWeightAndSlug(dir)
	m := make(map[string]string)
	m["slug"] = slug
	m["url"] = url
	if weight > 0 {
		m["weight"] =fmt.Sprintf("%d", weight)
	}
	return markdown.AddFrontMatter(indexFile, m)
}

// addIndex adds a directory index file to override the title of the folder, "unslugified"
func addIndex(path string) error {
	base := filepath.Base(path)
	title := unSlugify(base)
	weight, slug, url := deduceWeightAndSlug(path)

	fmt.Println("Adding an", colors.Labels.Unwanted("_index.md"),"file to ", colors.Labels.Wanted(path))

	if (url != "") {

		m := make(map[string]string)
		m["title"] = title
		m["slug"] = slug
		m["url"] = url
		if weight > 0 {
			m["weight"] = fmt.Sprintf("%d", weight)
		}

		return markdown.AddFrontMatter(filepath.Join(path, "_index.md"), m)
	}
	return nil
}

// unSlugify turns "something-like_this" into "Something Like This"
func unSlugify(name string) string {
	re := regexp.MustCompile(`(([\d\.]+)\s)?(.+)?`)
	name = strings.Replace(name, "-", " ", -1)
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	matches := re.FindStringSubmatch(name)
	if matches != nil {
		return matches[3]
	}
	return name
}