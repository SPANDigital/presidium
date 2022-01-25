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

type contentWeightTracker struct {
	paths   []string
	indices map[string]int
	weights map[int]int64
	current int
}

func (t *contentWeightTracker) tracking(path string) int {
	if v, found := t.indices[path]; found {
		t.current = v
	} else {
		t.current = -1
	}
	return t.current
}

func (t *contentWeightTracker) update(tracked int, weight int64) {
	if tracked != -1 {
		t.weights[tracked] = weight
	}
}

func (t *contentWeightTracker) lookupSiblingWeight(tracked int) int64 {

	if tracked == -1 || tracked == 0 || tracked >= len(t.paths) || len(t.paths) == -1 {
		return 0
	}

	prev := tracked - 1
	prevParent, _ := filepath.Split(t.paths[prev])
	thisParent, _ := filepath.Split(t.paths[tracked])

	if thisParent != prevParent {
		return 0
	}

	w := t.weights[prev]
	return w
}

func newContentWeighTracker(root string) contentWeightTracker {

	cwt := contentWeightTracker{
		weights: make(map[int]int64),
		indices: map[string]int{},
		paths:   make([]string, 0),
		current: -1,
	}

	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {

		if d.IsDir() {
			return nil
		}

		if _, name := filepath.Split(path); strings.ToLower(filepath.Ext(name)) != ".md" {
			return nil
		}

		cwt.paths = append(cwt.paths, path)

		return nil
	})

	sort.Strings(cwt.paths)

	for i := 0; i < len(cwt.paths); i++ {
		cwt.indices[cwt.paths[i]] = i
	}

	return cwt
}

func CheckForDirIndex(stagingDir, path string) error {

	var idxRenameList = make([]string, 0)

	weighTracker := newContentWeighTracker(path)

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
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
			_ = addIndex(stagingDir, path, &weighTracker)
		} else { // is a file
			if strings.HasSuffix(path, ".md") {
				if info.Name() != "_index.md" && info.Name() != "index.md" {
					err := injectSlugWeightAndURL(stagingDir, path, &weighTracker)
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
		err := injectSlugWeightAndURLForIndex(stagingDir, path+"/_index.md", &weighTracker)
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

var weightAndSlugRegex = regexp.MustCompile(`((([\d.]+)([a-z]?))-)?([^..]+)(\.[^.]*)?`)

// unNumerify turns "02-employment-contracts" into "employment-contracts" and "bill-add-customer" into "bill-add-customer"
func deduceWeightAndSlug(stagingDir, path string, weightTracker *contentWeightTracker) (int64, string, string) {

	tracked := weightTracker.tracking(path)
	fileName := spaces.ReplaceAllLiteralString(filepath.Base(path), "-")
	matches := weightAndSlugRegex.FindStringSubmatch(fileName)
	weightStr := matches[2]
	weightAdjustment := 0
	adjusted := false

	if len(matches[4]) == 1 {
		ac := matches[4][0]
		if ac >= 'a' && ac <= 'z' {
			weightAdjustment := int64(int(ac) - int('a'))
			if weightAdjustment > 0 {
				weightAdjustment += weightAdjustment + weightTracker.lookupSiblingWeight(tracked)
			}
			adjusted = true
			weightStr = matches[3]
		}
	}

	weight, err := strconv.ParseInt(strings.ReplaceAll(weightStr, ".", ""), 10, 64)
	if err != nil {
		weight = -1
	} else {
		weight += 1
		weight += int64(weightAdjustment)
	}

	if adjusted {
		weightTracker.update(tracked, weight)
	}

	var slug = slugify(matches[5])
	var url string
	var contentDir = filepath.Join(stagingDir, "content")
	if path == contentDir {
		url = ""
	} else {
		_, _, parentUrl := deduceWeightAndSlug(stagingDir, filepath.Dir(path), weightTracker)
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

func injectSlugWeightAndURL(stagingDir, path string, weightTracker *contentWeightTracker) error {
	if markdown.IsRecognizableMarkdown(path) {
		fmt.Println("Checking weight of ", colors.Labels.Info(path))
		weight, slug, url := deduceWeightAndSlug(stagingDir, path, weightTracker)

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

func injectSlugWeightAndURLForIndex(stagingDir, indexFile string, tracker *contentWeightTracker) error {
	dir := filepath.Dir(indexFile)
	weight, slug, url := deduceWeightAndSlug(stagingDir, dir, tracker)
	m := make(map[string]interface{})
	m["slug"] = slug
	m["url"] = url
	if weight > 0 {
		m["weight"] = fmt.Sprintf("%d", weight)
	}
	return markdown.AddFrontMatter(indexFile, m)
}

// addIndex adds a directory index file to override the title of the folder, "unslugified"
func addIndex(stagingDir, path string, weightTracker *contentWeightTracker) error {

	base := filepath.Base(path)
	title := unSlugify(base)
	weight, slug, url := deduceWeightAndSlug(stagingDir, path, weightTracker)

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

// slugify replaces all non word chars with a "-"
// turns "v0 .18.6" into "v0-18-6"
func slugify(name string) string {
	var re = regexp.MustCompile(`(?m)\W+`)
	return re.ReplaceAllString(name, "-")
}