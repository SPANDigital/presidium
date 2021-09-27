package validate

import (
	"container/list"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	model "github.com/SPANDigital/presidium-hugo/pkg/domain/model/validate"
	"github.com/SPANDigital/presidium-hugo/pkg/log"
	"github.com/scylladb/go-set/strset"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type validation struct {
	path    string                      // The path being validated
	seen    *strset.Set                 // Keep track of paths we have seen
	isLocal bool                        // TODO: Flag to determine validation live web site, or local file path
	queue   *list.List                  // Keep track of links still to be processed per page
	tracked map[model.Status]*list.List // Keep track of collected links per status
}

func (validation validation) IsLocal() bool {
	return validation.isLocal
}

func New(path string) Validator {
	return validation{
		path:    path,
		isLocal: true,
		seen:    strset.New(),
		queue:   list.New(),
		tracked: make(map[model.Status]*list.List),
	}
}

func (validation validation) hasSeen(f string) bool {
	seen := validation.seen.Has(f)
	if !seen {
		validation.seen.Add(f)
	}
	return seen
}

func (validation validation) cleanUp() {
	validation.seen.Clear()
	for k, v := range validation.tracked {
		v.Init()
		delete(validation.tracked, k)
	}
}

func (validation validation) Validate() (model.Report, error) {

	validation.seen.Clear()

	err := filepath.Walk(validation.path, func(path string, info fs.FileInfo, err error) error {

		if err != nil {
			log.ErrorWithFields(err, log.Fields{"validation_path": path})
			return err
		}

		if !info.IsDir() {
			log.DebugWithFields("validation started", log.Fields{"validation_path": path})
			err = validation.process(path)
			if err != nil {
				log.ErrorWithFields(err, log.Fields{"validation_path": path})
				return err
			}
		}

		return nil
	})

	if err != nil {
		return model.Report{}, err
	} else {
		return validation.newReport(), err
	}

}

func (validation validation) newReport() model.Report {

	report := model.Report{
		Data:       make(map[model.Status][]model.Link),
		Valid:      0,
		Broken:     0,
		External:   0,
		Warning:    0,
		TotalLinks: 0,
	}

	for s, links := range validation.tracked {

		countedLinks := links.Len()
		report.TotalLinks += countedLinks
		collected := make([]model.Link, 0)

		var next *list.Element
		for e := links.Front(); e != nil; e = next {
			link := e.Value.(model.Link)
			collected = append(collected, link)
			next = e.Next()
		}

		report.Data[s] = collected

		switch s {
		case model.Valid:
			report.Valid = countedLinks
			break
		case model.Broken:
			report.Broken = countedLinks
			break
		case model.Warning:
			report.Warning = countedLinks
			break
		case model.External:
			report.External = countedLinks
			break
		}
	}

	return report
}

func (validation validation) process(path string) error {

	s := strings.TrimSpace(strings.ToLower(path))

	if validation.hasSeen(s) {
		return nil
	}

	if !(strings.HasSuffix(s, ".html")) {
		return nil
	}

	validation.queue.PushFront(model.Link{
		Uri:        path,
		Location:   path,
		IsExternal: false,
	})

	for {

		if validation.queue.Len() == 0 {
			break
		}

		todo := validation.queue.Front()
		validation.queue.Remove(todo)
		link := todo.Value.(model.Link)

		if validation.hasSeen(link.Uri) {
			continue
		}

		if link.Uri == "/" {
			continue
		}

		if link.IsExternal {
			validation.reportLink(link, model.External, "")
			continue
		}

		info, err := os.Stat(link.Uri)

		if err != nil {
			link.Message = err.Error()
			continue
		}

		if info.IsDir() {
			file := fmt.Sprintf("%s/index.html", link.Uri)
			info, err = os.Stat(file)
			if err == nil {
				continue
			}
			if info.IsDir() {
				link.Message = fmt.Sprintf("expected file here but found directory: %s", file)
				continue
			}
			link.Uri = file
		}

		file, err := os.OpenFile(link.Uri, os.O_RDONLY, 0666)

		if err != nil {
			validation.reportLink(link, model.Broken, fmt.Sprintf("Unable to open file %s: %s", link.Uri, err.Error()))
			continue
		}

		var doc *goquery.Document
		doc, err = goquery.NewDocumentFromReader(file)
		if err != nil {
			validation.reportLink(link, model.Broken, fmt.Sprintf("file %s is propably not a valid HTML file: %s", link.Uri, err.Error()))
		} else {
			validation.reportLink(link, model.Valid, "")
			// Find all links referenced by this page!
			doc.Find("a[href]").Each(func(i int, item *goquery.Selection) {
				href, ok := item.Attr("href")
				if !ok || len(href) == 0 || href == "/" {
					return
				}
				validationHref := strings.ToLower(href)
				validationHref = strings.TrimSpace(validationHref)
				if strings.HasPrefix(validationHref, "mailto:") ||
					strings.HasPrefix(validationHref, "tel:") {
					validation.reportLink(link, model.Warning, fmt.Sprintf("Unhandled url scheme: %s", href))
				} else if strings.Contains(validationHref, "#") {
					return
				}

				parsedLinkUrl, err := url.Parse(href)

				if err != nil {
					link.Message = fmt.Sprintf("%v", err.Error())
					return
				}

				link.IsExternal = len(parsedLinkUrl.Scheme) > 0

				finalUri := fmt.Sprintf("%s%s", validation.path, href)

				validation.reportLink(link, model.Valid, "")

				validation.queue.PushFront(model.Link{
					Uri:      finalUri,
					Location: link.Uri,
					Label:    strings.TrimSpace(item.Text()),
				})
			})
		}

		_ = file.Close()

	}

	return nil
}

func fileOnPath(path string, name string) (string, error) {
	file := fmt.Sprintf("%s/%s", path, name)
	info, err := os.Stat(file)
	if err != nil {
		return file, err
	}
	if info.IsDir() {
		return file, errors.New(fmt.Sprintf("expected file but foun directory: %s", file))
	}
	return file, nil
}

func (validation validation) reportLink(link model.Link, status model.Status, message string) {

	link.Status = status
	link.Message = message

	collection, found := validation.tracked[status]

	if !found {
		collection = list.New()
		validation.tracked[status] = collection
	}

	collection.PushBack(link)

}
