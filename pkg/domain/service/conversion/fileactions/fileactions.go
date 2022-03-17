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

type contentWeightTracker struct {
	paths   []string
	indices map[string]int
	weights map[int]int64
	current int
}

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
					fm := deduceWeightAndSlug(stagingDir, path, &weighTracker)
					err = injectFrontMatter(path, fm)
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

	for _, idxPath := range idxRenameList {
		oldPath := fmt.Sprintf("%s/%s", idxPath, "index.md")
		newPath := fmt.Sprintf("%s/%s", idxPath, "_index.md")
		fmt.Println("Renaming", colors.Labels.Unwanted(fmt.Sprintf("%v/index.md", idxPath)), "to", colors.Labels.Wanted("_index.md"))

		err = os.Rename(oldPath, newPath)
		if err != nil {
			return err
		}

		fm := deduceWeightAndSlug(stagingDir, idxPath, &weighTracker)
		err = injectFrontMatter(newPath, fm)
		if err != nil {
			return err
		}
	}
	return nil
}

func CheckForTitles(path string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		md, err := markdown.Parse(path)
		if err != nil {
			return err
		}

		if len(md.FrontMatter.Title) > 0 {
			return nil
		}

		if strings.HasSuffix(path, "_index.md") {
			base := filepath.Base(filepath.Dir(path))
			md.FrontMatter.Title = unSlugify(base)
			if err != nil {
				return err
			}
		} else {
			md.FrontMatter.Title = titleFromPath(path)
		}

		return markdown.AddFrontMatter(path, md.FrontMatter)
	})
}

func replaceRoot(url string) string {
	rootUrl := viper.GetString("replaceRoot")
	url = strings.TrimPrefix(url, strings.ToLower(rootUrl))
	if len(url) == 0 {
		return "/"
	}
	return url
}

func injectFrontMatter(path string, fm markdown.FrontMatter) error {
	if !markdown.IsRecognizableMarkdown(path) {
		fmt.Println("Is not valid markdown", colors.Labels.Info(path))
		return nil
	}

	fmt.Println("Checking weight of ", colors.Labels.Info(path))
	indexMarkdown, err := markdown.Parse(path)
	if err != nil {
		return err
	}

	if newSlug := slugByPriority(indexMarkdown.FrontMatter); newSlug != nil {
		if dir, _ := filepath.Split(fm.URL); len(dir) > 0 {
			fm.URL = dir + *newSlug
		} else {
			fm.URL = *newSlug
		}
		fm.Slug = *newSlug
	}
	return markdown.AddFrontMatter(path, fm)
}

func deduceWeightAndSlug(stagingDir, path string, weightTracker *contentWeightTracker) markdown.FrontMatter {
	tracked := weightTracker.tracking(path)
	fileName := markdown.SpaceRe.ReplaceAllLiteralString(filepath.Base(path), "-")
	matches := markdown.WeightAndSlugRe.FindStringSubmatch(fileName)
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

	slug := slugify(matches[5])
	fm := markdown.FrontMatter{
		Slug: slug,
	}

	if weight >= 0 {
		fm.Weight = fmt.Sprintf("%d", weight)
	}

	contentDir := filepath.Join(stagingDir, "content")
	if path == contentDir {
		return fm
	}

	fm.URL = slug
	parentFm := deduceWeightAndSlug(stagingDir, filepath.Dir(path), weightTracker)
	if len(parentFm.URL) > 0 {
		fm.URL = parentFm.URL + "/" + slug
	}

	fm.URL = replaceRoot(fm.URL)
	return fm
}

// Gets the title from the Front Matter and uses it to create a slug
func slugByPriority(fm markdown.FrontMatter) *string {
	if len(fm.Slug) > 0 {
		return nil
	} else if fromFile := viper.GetBool("slugBasedOnFilename"); fromFile {
		return nil
	} else if title := fm.Title; len(title) > 0 {
		titleSlug := titleToSlug(title)
		return &titleSlug
	}
	return nil
}

// addIndex adds a directory index file to override the title of the folder, "unslugified"
func addIndex(stagingDir, path string, weightTracker *contentWeightTracker) error {
	base := filepath.Base(path)
	fm := deduceWeightAndSlug(stagingDir, path, weightTracker)
	fm.Title = unSlugify(base)

	fmt.Println("Adding an", colors.Labels.Unwanted("_index.md"), "file to ", colors.Labels.Wanted(path))
	if len(fm.URL) == 0 {
		return nil
	}

	return markdown.AddFrontMatter(filepath.Join(path, "_index.md"), fm)
}

func RemoveWeightIndicatorsFromFilePaths(stagingContentDir string) error {
	return removeWeightIndicatorsFromFilePaths(stagingContentDir, "/")
}

func removeWeightIndicatorsFromFilePaths(contentDir string, dir string) error {
	parentDir, name := filepath.Split(dir)
	if strings.HasPrefix(name, ".") {
		return nil
	}

	nameIsWeighted := dir != "/" && markdown.WeightRe.FindStringSubmatch(name) != nil
	if nameIsWeighted {
		newName := markdown.WeightRe.ReplaceAllStringFunc(name, func(s string) string { return "" })
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

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

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
	name = reDividers.ReplaceAllString(name, " ")
	name = strings.Title(name)
	matches := re.FindStringSubmatch(name)
	if matches != nil {
		return strings.TrimSpace(matches[3])
	}
	return strings.TrimSpace(name)
}

// slugify replaces all non word chars with a "-"
// turns "v0 .18.6_8." into "v0-18-6-8"
func slugify(name string) string {
	var nonWordRe = regexp.MustCompile(`(?m)(\W|_)+`)
	slug := nonWordRe.ReplaceAllString(name, "-")
	return strings.Trim(slug, "-")
}

// Take a capitalized title and turn it into a slug
func titleToSlug(title string) string {
	title = strings.ToLower(title)
	title = strings.Replace(title, "&", "and", -1)
	title = slugify(title)
	return title
}

func titleFromPath(path string) string {
	base := filepath.Base(path)
	fileName := strings.TrimSuffix(base, filepath.Ext(base))
	return unSlugify(fileName)
}
