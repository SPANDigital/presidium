package validation

import (
	"container/list"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/scylladb/go-set/strset"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	strings "strings"
)

type LinkListener = func(link Link)

type Validator interface {
	Validate() (Report, error)
	IsLocal() bool
	IsInProgress() bool
}

type Link struct {
	Uri        string
	Location   string
	Status     Status
	Message    string
	IsExternal bool
}

type Report struct {
	Broken   []Link
	Warnings []Link
	External []Link
	Valid    []Link
}

type Status string

const (
	Valid    = Status("valid")
	Broken   = Status("broken")
	Warning  = Status("warning")
	External = Status("external")
)

const (
	validationPending = iota
	validationInProgress
	validationCompleted
)

type validation struct {
	path    string      // The path being validated
	seen    *strset.Set // Keep track of paths we have seen
	isLocal bool        // Flag to determine validation live web site, or local file path
	report  Report
	state   int
}

func (validation validation) IsInProgress() bool {
	return validation.state == validationInProgress
}

func (validation validation) IsLocal() bool {
	return validation.isLocal
}

func New(path string) Validator {
	return validation{
		path:    path,
		isLocal: true,
		seen:    strset.New(),
	}
}

func (validation validation) hasSeen(f string) bool {
	seen := validation.seen.Has(f)
	if !seen {
		validation.seen.Add(f)
	}
	return seen
}

func (validation validation) Validate() (Report, error) {

	if validation.state == validationInProgress {
		return validation.report, nil
	}

	validation.state = validationInProgress

	validation.report = Report{
		Broken:   make([]Link, 0),
		Warnings: make([]Link, 0),
		External: make([]Link, 0),
		Valid:    make([]Link, 0),
	}

	validation.seen.Clear()

	err := filepath.Walk(validation.path, func(path string, info fs.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if !info.IsDir() && !validation.hasSeen(path) {
			err = validation.process(path)
			if err != nil {
				return err
			}
		}

		return nil
	})

	validation.state = validationCompleted

	return validation.report, err
}

func (validation validation) process(path string) error {

	s := strings.TrimSpace(strings.ToLower(path))

	if !(strings.HasSuffix(s, ".html")) {
		return nil
	}

	queue := list.New()

	queue.PushFront(Link{
		Uri:      path,
		Location: path,
	})

	for {

		if queue.Len() == 0 {
			break
		}

		todo := queue.Front()
		queue.Remove(todo)
		link := todo.Value.(Link)

		if link.IsExternal {
			validation.reportLink(link, External, "")
			continue
		}

		file, err := os.OpenFile(link.Uri, os.O_RDONLY, 0666)

		if err != nil {
			validation.reportLink(link, Broken, fmt.Sprintf("Unable to open file %s: %s", link.Uri, err.Error()))
			continue
		}

		var doc *goquery.Document
		doc, err = goquery.NewDocumentFromReader(file)
		if err != nil {
			validation.reportLink(link, Broken, fmt.Sprintf("file %s is propably not a valid HTML file: %s", link.Uri, err.Error()))
		} else {
			validation.reportLink(link, Valid, "")
			// Find all links referenced by this page!
			doc.Find("a[href]").Each(func(i int, item *goquery.Selection) {
				href, ok := item.Attr("href")
				if !ok {
					return
				}
				validationHref := strings.ToLower(href)
				validationHref = strings.TrimSpace(validationHref)
				if strings.HasPrefix(validationHref, "mailto:") ||
					strings.HasPrefix(validationHref, "tel:") {
					validation.reportLink(link, Warning, fmt.Sprintf("Unhandled url scheme: %s", href))
				} else if strings.HasPrefix(validationHref, "#") {
					return
				}
				_, err := url.Parse(href)

				var finalUri string

				if err == nil {
					finalUri = href
				} else if strings.HasPrefix(validationHref, "/") {
					finalUri = fmt.Sprintf("%s%s", validation.path, href)
				} else {
					finalUri = fmt.Sprintf("%s/%s", link.Uri, href)
				}

				queue.PushFront(Link{
					Uri:        finalUri,
					Location:   link.Uri,
					IsExternal: err == nil,
				})
			})
		}

		_ = file.Close()

	}

	return nil
}

func (validation validation) reportLink(link Link, status Status, message string) {

	link.Status = status

	if len(message) == 0 {
		link.Message = fmt.Sprintf("%v", status)
	} else {
		link.Message = fmt.Sprintf("%v: %s", status, message)
	}

	if link.Status == Broken {
		validation.report.Broken = append(validation.report.Broken, link)
	} else if link.Status == Valid {
		validation.report.Valid = append(validation.report.Valid, link)
	} else if link.Status == Warning {
		validation.report.Warnings = append(validation.report.Warnings, link)
	} else if link.Status == External {
		link.IsExternal = true
		validation.report.External = append(validation.report.External, link)
	}
}
