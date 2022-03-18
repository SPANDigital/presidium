package fileactions

import (
	"fmt"
	"github.com/SPANDigital/presidium-hugo/pkg/utils"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
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

var dirUrls map[string]string

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
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		fmt.Println("Walking", colors.Labels.Info(path))
		if isContentPath(path, stagingDir) {
			return nil
		}

		if filepath.Base(path) == "index.md" {
			newPath := filepath.Join(filepath.Dir(path), "_index.md")
			err = os.Rename(path, newPath)
			if err != nil {
				return err
			}
			return nil
		}

		if !info.IsDir() {
			return nil
		}

		indexPath := filepath.Join(path, "_index.md")
		fmt.Println(fmt.Sprintf("Checking %s for _index.md...\n", colors.Labels.Wanted(indexPath)))
		if utils.FileExists(indexPath) {
			return nil
		}

		return markdown.AddFrontMatter(indexPath, markdown.FrontMatter{})
	})
}

func AddFrontMatter(stagingDir, path string) error {
	weighTracker := newContentWeighTracker(path)
	dirUrls = map[string]string{}
	return filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, "_index.md") || isContentPath(path, stagingDir) {
			return nil
		}

		if info.IsDir() {
			path = filepath.Join(path, "_index.md")
		}

		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		md, err := markdownForPath(path)
		if err != nil {
			return nil
		}

		fm := deduceWeightAndSlug(stagingDir, md.FrontMatter, path, &weighTracker)
		err = markdown.AddFrontMatter(path, fm)
		if err != nil {
			return err
		}

		return markdown.Operate(path)
	})
}

func CheckForTitles(path string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		md, err := markdownForPath(path)
		if err != nil {
			return err
		}

		fmt.Println(path)
		if len(md.FrontMatter.Title) > 0 {
			return nil
		}

		if strings.HasSuffix(path, "_index.md") {
			base := filepath.Base(filepath.Dir(path))
			md.FrontMatter.Title = utils.UnSlugify(base)
			if err != nil {
				return err
			}
		} else {
			md.FrontMatter.Title = utils.TitleToSlug(path)
		}

		return markdown.AddFrontMatter(path, md.FrontMatter)
	})
}

func deduceWeightAndSlug(stagingDir string, fm markdown.FrontMatter, path string, weightTracker *contentWeightTracker) markdown.FrontMatter {
	tracked := weightTracker.tracking(path)
	fileName := markdown.SpaceRe.ReplaceAllLiteralString(filepath.Base(path), "-")
	matches := markdown.WeightAndSlugRe.FindStringSubmatch(fileName)
	weightStr := matches[2]
	var weightAdjustment int64
	adjusted := false

	if len(matches[4]) == 1 {
		ac := matches[4][0]
		if ac >= 'a' && ac <= 'z' {
			weightAdjustment = int64(int(ac) - int('a'))
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
		weight += weightAdjustment
	}

	if adjusted {
		weightTracker.update(tracked, weight)
	}

	if weight >= 0 {
		fm.Weight = fmt.Sprintf("%d", weight)
	}

	fm.Slug = utils.TitleToSlug(fm.Title)
	fm.URL = fm.Slug
	base := filepath.Dir(path)
	if parent, ok := dirUrls[base]; ok {
		fm.URL = filepath.Join(parent, fm.Slug)
	}
	fm.URL = replaceRoot(fm.URL)

	if strings.HasSuffix(path, "_index.md") {
		root := filepath.Dir(base)
		if !isContentPath(path, stagingDir) {
			if parent, ok := dirUrls[root]; ok {
				fm.URL = filepath.Join(parent, fm.URL)
			}
		}
		dirUrls[base] = replaceRoot(fm.URL)
	}
	return fm
}

func markdownForPath(path string) (*markdown.Markdown, error) {
	if !markdown.IsRecognizableMarkdown(path) {
		fmt.Println("Is not valid markdown", colors.Labels.Info(path))
		return &markdown.Markdown{}, nil
	}
	return markdown.Parse(path)
}

func isContentPath(path, stagingDir string) bool {
	contentDir := filepath.Join(stagingDir, "content")
	return path == contentDir
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

// titleFromPath generates an article title based on the filename
func titleFromPath(path string) string {
	base := filepath.Base(path)
	fileName := strings.TrimSuffix(base, filepath.Ext(base))
	return utils.UnSlugify(fileName)
}

func replaceRoot(url string) string {
	rootUrl := viper.GetString("replaceRoot")
	url = strings.TrimPrefix(url, strings.ToLower(rootUrl))
	if len(url) == 0 {
		return "/"
	}
	return url
}
