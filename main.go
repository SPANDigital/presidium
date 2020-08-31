package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var idxRenameList []string

func checkForDirRename(path string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() && strings.HasPrefix(file.Name(), "_") {
			old := path + "/" + file.Name()
			new := path + "/" + strings.TrimLeft(file.Name(), "_")
			log.Printf("Renaming %v to %v...\n", old, new)
			err := os.Rename(old, new)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func addFrontMatter(path string, params map[string]string) error {
	log.Println("Adding front matter",params,"to", path)
	b, err := ioutil.ReadFile(path) // just pass the file name
	if err != nil {
		return err
	}
	//re := regexp.MustCompile(`(\-\-\-)(.+)(\-\-\-)(.*)`)
	re := regexp.MustCompile(`^(?s:(---\n)(.*)(---\n)(.*))$`)
	matches := re.FindSubmatch(b)

	if matches != nil {
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		f.Write(matches[1])
		f.Write(matches[2])
		for key, value := range params {
			f.WriteString(fmt.Sprintf("%s: %s\n", key, value))
		}
		f.Write(matches[3])
		f.Write(matches[4])

		return f.Close()
	}
	return nil
}

func injectWeight(path string) error {
	parts := strings.Split(path, "/")
	lastPart := parts[len(parts)-1]
	weight, _ := deduceWeight(lastPart)
	if weight > 0 {
		return addFrontMatter(path, map[string]string{
			"weight": fmt.Sprintf("%d", weight),
		})
	}
	return nil
}

func injectSlugAndWeightForIndex(indexFile string) error {
	dir := filepath.Dir(indexFile)
	parts := strings.Split(dir, "/")
	lastPart := parts[len(parts)-1]
	weight, slug := deduceWeight(lastPart)
	m := make(map[string]string)
	m["slug"] = slug
	if weight > 0 {
		m["weight"] =fmt.Sprintf("%d", weight)
	}
	return addFrontMatter(indexFile, m)
}

func checkForDirIndex(path string) error {
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			log.Printf("Checking %v for _index.md...\n", path)
			if fileExists(path + "/_index.md") {
				return nil
			}

			// If we find an index.md, we'll rename it later, as if we rename it now, we'll affect the Walk function
			// and get a nil pointer exception when we try "walk" into the old directory
			if fileExists(path + "/index.md") {
				log.Printf("Adding path %v for later index rename...\n", path)
				idxRenameList = append(idxRenameList, path)
				return nil
			}
			addIndex(path)
		} else {
			if info.Name() != "_index.md" && info.Name() != "index.md" {
				log.Println("Checking weight of ", path)
				return injectWeight(path)
			}
		}
		return nil
	})
	for _, path := range idxRenameList {
		log.Printf("Renaming %v/index.md to _index.md...", path)
		os.Rename(path+"/index.md", path+"/_index.md")
		err := injectSlugAndWeightForIndex(path+"/_index.md")
		if err != nil {
			return err
		}
	}
	return nil
}

// unSlugify turns "something-like_this" into "Something Like This"
func unSlugify(name string) string {
	name = strings.Replace(name, "-", " ", -1)
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	return name
}

// unNumerify turns "02-employment-contracts" into "employment-contracts" and "bill-add-customer" into "bill-add-customer"
func deduceWeight(filename string) (int64, string) {
	re := regexp.MustCompile(`((?P<numeric>\d+)\-)?(?P<rest>.*)(?P<extension>\..*)?`)
	matches := re.FindStringSubmatch(filename)
	weight, err := strconv.ParseInt(matches[2], 10, 64)
	if err != nil {
		weight = 0
	}
	return weight, matches[3]
}

// addIndex adds a directory index file to override the title of the folder, "unslugified"
func addIndex(path string) error {
	log.Printf("Adding an _index.md file to %v...", path)
	f, err := os.Create(path + "/_index.md")
	if err != nil {
		return err
	}
	defer f.Close()
	parts := strings.Split(path, "/")
	lastPart := parts[len(parts)-1]
	title := unSlugify(lastPart)
	weight, slug := deduceWeight(lastPart)
	f.WriteString("---\n")
	f.WriteString("title: " + title + "\n")
	f.WriteString("slug: " + slug + "\n")
	if weight > 0 {
		f.WriteString(fmt.Sprintf("weight: %d\n", weight))
	}
	f.WriteString("---\n")
	return nil
}

func main() {
	checkForDirRename("content")
	checkForDirIndex("content")
}
