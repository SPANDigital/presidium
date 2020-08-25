package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
		}
		return nil
	})
	for _, path := range idxRenameList {
		log.Printf("Renaming %v/index.md to _index.md...", path)
		os.Rename(path+"/index.md", path+"/_index.md")

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

// addIndex adds a directory index file to override the title of the folder, "unslugified"
func addIndex(path string) error {
	log.Printf("Adding an _index.md file to %v...", path)
	f, err := os.Create(path + "/_index.md")
	if err != nil {
		return err
	}
	defer f.Close()
	parts := strings.Split(path, "/")
	title := unSlugify(parts[len(parts)-1])
	f.WriteString("---\n")
	f.WriteString("title: " + title + "\n")
	f.WriteString("---\n")
	return nil
}

func main() {
	checkForDirRename("content")
	checkForDirIndex("content")
}
