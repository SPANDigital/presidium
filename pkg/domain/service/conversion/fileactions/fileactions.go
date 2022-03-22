package fileactions

import (
	"fmt"
	"github.com/SPANDigital/presidium-hugo/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/conversion/colors"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/conversion/markdown"
	"github.com/spf13/viper"
)

type directoryMap map[string][]string
var dirUrls map[string]string
var afFs = afero.NewOsFs()

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
	pm, err := buildWeightMap(path)
	if err != nil {
		return errors.Wrap(err, "path map")
	}

	dirUrls = map[string]string{}
	return filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if isIndex(path) || isContentPath(path, stagingDir) {
			return nil
		}

		if info.IsDir() {
			path = filepath.Join(path, "_index.md")
		}

		if !isMdFile(path) {
			return nil
		}

		md, err := markdownForPath(path)
		if err != nil {
			return nil
		}

		fm := md.FrontMatter
		fm.Weight = getPathWeight(pm, path)
		fm.Slug, fm.URL = getSlugAndUrl(stagingDir, md.FrontMatter.Title, path)
		err = markdown.AddFrontMatter(path, fm)
		if err != nil {
			return err
		}

		return markdown.Operate(path)
	})
}

func CheckForTitles(path string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || !isMdFile(path) {
			return nil
		}

		md, err := markdownForPath(path)
		if err != nil {
			return err
		}

		if len(md.FrontMatter.Title) > 0 {
			return nil
		}

		if isIndex(path) {
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

func RemoveWeightIndicatorsFromFilePaths(stagingContentDir string) error {
	return removeWeightIndicatorsFromFilePaths(stagingContentDir, "/")
}

func getSlugAndUrl(stagingDir string, title string, path string) (slug string, url string) {
	slug = utils.TitleToSlug(title)
	url = slug
	base := filepath.Dir(path)
	if parent, ok := dirUrls[base]; ok {
		url = filepath.Join(parent, slug)
	}
	url = replaceRoot(url)

	if isIndex(path) {
		root := filepath.Dir(base)
		if !isContentPath(path, stagingDir) {
			if parent, ok := dirUrls[root]; ok {
				url = filepath.Join(parent, url)
			}
		}
		dirUrls[base] = replaceRoot(url)
	}
	return slug, url
}

func buildWeightMap(path string) (directoryMap, error) {
	dirMap := directoryMap{}
	err := afero.Walk(afFs, path, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && !isMdFile(path) {
			return nil
		}

		fileName := filepath.Base(path)
		if !markdown.WeightRe.MatchString(fileName) {
			return nil
		}

		dir := filepath.Dir(path)
		if _, ok := dirMap[dir]; ok {
			dirMap[dir] = append(dirMap[dir], path)
		} else {
			dirMap[dir] = []string{path}
		}
		return nil
	})
	return dirMap, err
}

func getPathWeight(dm directoryMap, path string) string {
	if isIndex(path) {
		path = filepath.Dir(path)
	}

	dir := filepath.Dir(path)
	if paths, ok := dm[dir]; ok {
		for i, p := range paths {
			if p == path {
				return strconv.Itoa(i + 1)
			}
		}
	}
	return ""
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

func removeWeightIndicatorsFromFilePaths(contentDir string, dir string) error {
	parentDir, name := filepath.Split(dir)
	if strings.HasPrefix(name, ".") {
		return nil
	}

	nameIsWeighted := dir != "/" && markdown.WeightRe.FindStringSubmatch(name) != nil
	if nameIsWeighted {
		newName := markdown.WeightRe.ReplaceAllString(name, "")
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
		err = removeWeightIndicatorsFromFilePaths(contentDir, local)
		if err != nil {
			return err
		}
	}

	return nil
}

// isIndex checks if the file is a hugo md index
func isIndex(path string) bool {
	return strings.HasSuffix(path, "_index.md")
}

// isMdFile checks if a file has the markdown ext
func isMdFile(path string) bool {
	return strings.HasSuffix(path, ".md")
}

// titleFromPath generates an article title based on the filename
func titleFromPath(path string) string {
	base := filepath.Base(path)
	fileName := strings.TrimSuffix(base, filepath.Ext(base))
	return utils.UnSlugify(fileName)
}

// replaceRoot removes the root prefix from an url
func replaceRoot(url string) string {
	rootUrl := viper.GetString("replaceRoot")
	url = strings.TrimPrefix(url, strings.ToLower(rootUrl))
	if len(url) == 0 {
		return "/"
	}
	return url
}
