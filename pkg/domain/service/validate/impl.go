package validate

import (
	"container/list"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	model "github.com/SPANDigital/presidium-hugo/pkg/domain/model/validate"
	"github.com/SPANDigital/presidium-hugo/pkg/filesystem"
	"github.com/SPANDigital/presidium-hugo/pkg/log"
	"github.com/scylladb/go-set/strset"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type validation struct {
	path    string                      // The path being validated
	seen    *strset.Set                 // Keep track of paths we have seen
	isLocal bool                        // TODO: Flag to determine validation live web site, or local file path
	queue   *list.List                  // Keep track of links still to be processed per page
	tracked map[model.Status]*list.List // Keep track of collected links per status
}

func (v validation) IsLocal() bool {
	return v.isLocal
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

func (v validation) hasSeen(f string) bool {
	seen := v.seen.Has(f)
	if !seen {
		v.seen.Add(f)
	}
	return seen
}

func (v validation) cleanUp() {
	v.seen.Clear()
	for k, val := range v.tracked {
		val.Init()
		delete(v.tracked, k)
	}
}

func (v validation) Validate() (model.Report, error) {
	v.seen.Clear()

	err := filesystem.AFS.Walk(v.path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.ErrorWithFields(err, log.Fields{"validation_path": path})
			return err
		}

		if !info.IsDir() {
			log.DebugWithFields("v started", log.Fields{"validation_path": path})
			err = v.process(path)
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
		return v.newReport(), err
	}
}

func (v validation) newReport() model.Report {
	report := model.Report{
		Data:       make(map[model.Status][]model.Link),
		Valid:      0,
		Broken:     0,
		External:   0,
		Warning:    0,
		TotalLinks: 0,
	}

	for s, links := range v.tracked {
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

func (v validation) process(path string) error {
	s := strings.TrimSpace(strings.ToLower(path))
	if v.hasSeen(s) {
		return nil
	}

	if !(strings.HasSuffix(s, ".html")) {
		return nil
	}

	v.queue.PushFront(model.Link{
		Uri:        path,
		Location:   path,
		IsExternal: false,
	})

	for {
		if v.queue.Len() == 0 {
			break
		}

		todo := v.queue.Front()
		v.queue.Remove(todo)
		link := todo.Value.(model.Link)

		if v.hasSeen(link.Uri) {
			continue
		}

		if link.Uri == "/" {
			continue
		}

		if link.IsExternal {
			v.reportLink(link, model.External, "")
			continue
		}

		if strings.Contains(link.Uri, "#") {
			if strings.Contains(link.Uri, "github") {
				fmt.Println("dadasd")
			}

			if err := v.validateRemoteAnchor(link); err != nil {
				continue
			}
			continue
		}

		info, err := filesystem.AFS.Stat(link.Uri)
		if err != nil {
			link.Message = err.Error()
			continue
		}

		if info.IsDir() {
			file := fmt.Sprintf("%s/index.html", link.Uri)
			info, err = filesystem.AFS.Stat(file)
			if err == nil {
				continue
			}
			if info.IsDir() {
				link.Message = fmt.Sprintf("expected file here but found directory: %s", file)
				continue
			}
			link.Uri = file
		}

		file, err := filesystem.AFS.OpenFile(link.Uri, os.O_RDONLY, 0666)
		if err != nil {
			v.reportLink(link, model.Broken, fmt.Sprintf("Unable to open file %s: %s", link.Uri, err.Error()))
			continue
		}

		var doc *goquery.Document
		doc, err = goquery.NewDocumentFromReader(file)
		if err != nil {
			v.reportLink(link, model.Broken, fmt.Sprintf("file %s is propably not a valid HTML file: %s", link.Uri, err.Error()))
		} else {
			v.reportLink(link, model.Valid, "")
			// Find all links referenced by this page!
			doc.Find(".article.child a[href]").Each(func(i int, item *goquery.Selection) {

				anchor := item.Closest(".article.child").Find("span.anchor[data-id]")
				dataId, _ := anchor.Attr("data-id")

				href, ok := item.Attr("href")
				if !ok || len(href) == 0 || href == "/" {
					return
				}
				validationHref := strings.ToLower(href)
				validationHref = strings.TrimSpace(validationHref)
				if strings.HasPrefix(validationHref, "mailto:") ||
					strings.HasPrefix(validationHref, "tel:") {
					v.reportLink(link, model.Warning, fmt.Sprintf("Unhandled url scheme: %s", href))
				}

				if strings.HasPrefix(href, "#") {
					v.validateAnchor(doc, link, href)
					return
				}

				parsedLinkUrl, err := url.Parse(href)
				if err != nil {
					link.Message = fmt.Sprintf("%v", err.Error())
					return
				}

				isExternal := len(parsedLinkUrl.Scheme) > 0
				if !isExternal {
					href = filepath.Clean(fmt.Sprintf("%s/%s", v.path, href))
				}

				v.queue.PushFront(model.Link{
					Uri:        href,
					Location:   link.Uri,
					DataId:     dataId,
					IsExternal: isExternal,
					Label:      strings.TrimSpace(item.Text()),
				})
			})
		}

		_ = file.Close()
	}

	return nil
}

func (v validation) validateAnchor(doc *goquery.Document, link model.Link, anchor string) {
	link.Uri = strings.Replace(link.Uri, "index.html", anchor, 1)
	anchor = strings.Replace(anchor, ".", "\\.", -1)
	if len(doc.Find(anchor).Nodes) == 0 {
		v.reportLink(link, model.Broken, "broken anchor reference")
		return
	}
	v.reportLink(link, model.Valid, "")
}

func (v validation) validateRemoteAnchor(link model.Link) error {
	anchorRe := regexp.MustCompile(`#[\w-.]+$`)
	anchor := anchorRe.FindString(link.Uri)
	path := strings.Replace(link.Uri, anchor, "index.html", 1)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		v.reportLink(link, model.Broken, "path does not exist: ")
		return err
	}

	file, err := filesystem.AFS.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		v.reportLink(link, model.Broken, "failed to open file")
		return err
	}

	defer file.Close()

	doc, err := goquery.NewDocumentFromReader(file)
	v.validateAnchor(doc, link, anchor)

	return nil
}

func fileOnPath(path string, name string) (string, error) {
	file := fmt.Sprintf("%s/%s", path, name)
	info, err := filesystem.AFS.Stat(file)
	if err != nil {
		return file, err
	}
	if info.IsDir() {
		return file, errors.New(fmt.Sprintf("expected file but found directory: %s", file))
	}
	return file, nil
}

func (v validation) reportLink(link model.Link, status model.Status, message string) {

	link.Status = status
	link.Message = message

	collection, found := v.tracked[status]

	if !found {
		collection = list.New()
		v.tracked[status] = collection
	}

	collection.PushBack(link)

}
