package fileactions

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/conversion/colors"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/conversion/markdown"
	"github.com/spf13/viper"
)

var reNamedWeight = regexp.MustCompile(`^[\d\-.]+`)
var spaces = regexp.MustCompile(`\s+`)

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func RemoveUnderscoreDirPrefix(dirPath string) error {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() && strings.HasPrefix(file.Name(), "_") {
			oldPath := dirPath + "/" + file.Name()
			newPath := dirPath + "/" + strings.TrimLeft(file.Name(), "_")
			fmt.Println("Renaming", colors.Labels.Unwanted(oldPath), "to", colors.Labels.Wanted(newPath))
			err := os.Rename(oldPath, newPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func listMarkDownFiles(stagingDir string) ([]string, error) {

	markdownFiles := make([]string, 0)

	err := filepath.WalkDir(stagingDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() || strings.ToLower(filepath.Ext(path)) != ".md" {
			return nil
		}
		markdownFiles = append(markdownFiles, path)
		return nil

	})

	if err != nil {
		return markdownFiles, err
	}

	sort.Strings(markdownFiles)

	return markdownFiles, nil
}

func CheckForDirIndex(stagingDir, path string) error {

	var idxRenameList = make([]string, 0)
	var markdownFiles, err = listMarkDownFiles(path)

	if err != nil {
		panic(err)
	}

	println(len(markdownFiles))

	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		fmt.Println("Walking", colors.Labels.Info(path))
		if info.IsDir() {
			fmt.Println(fmt.Sprintf("Checking %s for _index.md...\n", colors.Labels.Wanted(path)))
			if fileExists(path + "/_index.md") {
				return nil
			}

			// If we find an index.md, we'll rename it later, as if we rename it now, we'll affect the Walk function
			// and get a nil pointer exception when we try "walk" into the old directory
			if fileExists(path + "/index.md") {
				fmt.Println("Adding path", colors.Labels.Info(path), "for later index rename")
				idxRenameList = append(idxRenameList, path)
				return nil
			}
			_ = addIndex(stagingDir, path, &markdownFiles)
		} else { // is a file
			if strings.HasSuffix(path, ".md") {
				if info.Name() != "_index.md" && info.Name() != "index.md" {
					err := injectSlugWeightAndURL(stagingDir, path, &markdownFiles)
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
	if err != nil {
		return err
	}
	for _, path := range idxRenameList {
		fmt.Println("Renaming", colors.Labels.Unwanted(fmt.Sprintf("%v/index.md", path)), "to", colors.Labels.Wanted("_index.md"))
		os.Rename(path+"/index.md", path+"/_index.md")
		err := injectSlugWeightAndURLForIndex(stagingDir, path+"/_index.md", &markdownFiles)
		if err != nil {
			return err
		}
	}
	return nil
}

func CheckIndexForTitles(path string) error {
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, "_index.md") {
			indexMarkdown, err := markdown.Parse(path)
			if err != nil {
				return err
			}
			if _, ok := indexMarkdown.FrontMatter["title"]; !ok {
				base := filepath.Base(filepath.Dir(path))
				title := unSlugify(base)
				markdown.AddFrontMatter(path, map[string]interface{}{"title": title})
			}
		}
		return nil
	})
	return nil
}

// unNumerify turns "02-employment-contracts" into "employment-contracts" and "bill-add-customer" into "bill-add-customer"
func deduceWeightAndSlug(stagingDir, path string, markdownFiles *[]string) (int64, string, string) {

	if strings.Contains(path, "2b-improve-data-completeness.md") {
		println("wha!")
	}

	//re := regexp.MustCompile(`(([\d\.]+)\-)?([^..]+)(\.[^\..]*)?`)

	re := regexp.MustCompile(`((([\d\.]+)([a-z]{0,1}))\-)?([^..]+)(\.[^\..]*)?`)

	fileName := spaces.ReplaceAllLiteralString(filepath.Base(path), "-")
	matches := re.FindStringSubmatch(fileName)
	weight, err := strconv.ParseInt(strings.ReplaceAll(matches[2], ".", ""), 10, 64)
	if err != nil {
		weight = -1
	} else {
		weight += 1
	}

	var slug = matches[3]
	var url string
	var contentDir = filepath.Join(stagingDir, "content")
	if path == contentDir {
		url = ""
	} else {
		_, _, parentUrl := deduceWeightAndSlug(stagingDir, filepath.Dir(path), markdownFiles)
		if parentUrl != "" {
			url = parentUrl + "/" + slug
		} else {
			url = slug
		}
		replaceRoot := viper.GetString("replaceRoot")
		url = strings.TrimPrefix(url, strings.ToLower(replaceRoot))
		if url == "" {
			url = "/"
		}
	}

	return weight, slug, url
}

func injectSlugWeightAndURL(stagingDir, path string, markdownFiles *[]string) error {
	if markdown.IsRecognizableMarkdown(path) {
		fmt.Println("Checking weight of ", colors.Labels.Info(path))
		weight, slug, url := deduceWeightAndSlug(stagingDir, path, markdownFiles)

		m := make(map[string]interface{})
		m["slug"] = slug
		m["url"] = url
		if weight >= 0 {
			m["weight"] = fmt.Sprintf("%d", weight)
		}

		return markdown.AddFrontMatter(path, m)
	}
	fmt.Println("Is not valid markdown", colors.Labels.Info(path))
	return nil
}

func injectSlugWeightAndURLForIndex(stagingDir, indexFile string, markdowFiles *[]string) error {
	dir := filepath.Dir(indexFile)
	weight, slug, url := deduceWeightAndSlug(stagingDir, dir, markdowFiles)
	m := make(map[string]interface{})
	m["slug"] = slug
	m["url"] = url
	if weight > 0 {
		m["weight"] = fmt.Sprintf("%d", weight)
	}
	return markdown.AddFrontMatter(indexFile, m)
}

// addIndex adds a directory index file to override the title of the folder, "unslugified"
func addIndex(stagingDir, path string, markdownFiles *[]string) error {
	base := filepath.Base(path)
	title := unSlugify(base)
	weight, slug, url := deduceWeightAndSlug(stagingDir, path, markdownFiles)

	fmt.Println("Adding an", colors.Labels.Unwanted("_index.md"), "file to ", colors.Labels.Wanted(path))

	if url != "" {
		m := make(map[string]interface{})
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

func RemoveWeightIndicatorsFromFilePaths(stagingContentDir string) error {
	return removeWeightIndicatorsFromFilePaths(stagingContentDir, "/")
}

func removeWeightIndicatorsFromFilePaths(contentDir string, dir string) error {

	parentDir, name := filepath.Split(dir)

	if strings.HasPrefix(name, ".") {
		return nil
	}

	nameIsWeighted := dir != "/" && reNamedWeight.FindStringSubmatch(name) != nil

	if nameIsWeighted {
		newName := reNamedWeight.ReplaceAllStringFunc(name, func(s string) string { return "" })
		newPath := filepath.Join(contentDir, parentDir, newName)
		oldPath := filepath.Join(contentDir, dir)
		if err := os.Rename(oldPath, newPath); err != nil {
			return err
		}
		fmt.Println("Removed weight from ", colors.Labels.Unwanted(name), " -> ", colors.Labels.Wanted(filepath.Join(dir, newName)))
		dir = filepath.Join(parentDir, newName)
	}

	path := filepath.Join(contentDir, dir)
	info, err := os.Stat(path)

	if err != nil {
		return err
	}

	if !info.IsDir() {
		return nil
	}

	entries, _ := os.ReadDir(path)

	if entries == nil || len(entries) == 0 {
		return nil
	}

	for _, entry := range entries {
		local := filepath.Join(dir, entry.Name())
		err := removeWeightIndicatorsFromFilePaths(contentDir, local)
		if err != nil {
			return err
		}
	}

	return nil
}

// unSlugify turns "something-like_this" into "Something Like This"
func unSlugify(name string) string {
	re := regexp.MustCompile(`(([\d.]+)\s)?(.+)?`)
	reDividers := regexp.MustCompile(`[\-_]+`)
	name = reDividers.ReplaceAllStringFunc(name, func(s string) string { return " " })
	name = strings.Title(name)
	matches := re.FindStringSubmatch(name)
	if matches != nil {
		return matches[3]
	}
	return name
}
