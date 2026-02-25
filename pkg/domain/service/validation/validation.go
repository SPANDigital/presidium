package validation

import (
	"container/list"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/SPANDigital/presidium-hugo/pkg/filesystem"
	"github.com/SPANDigital/presidium-hugo/pkg/log"
	"github.com/scylladb/go-set/strset"
)

type LinkListener = func(link Link)

type Validator interface {
	Validate() (Report, error)
	IsLocal() bool
}

type Link struct {
	Uri        string
	Location   string
	Status     Status
	Message    string
	IsExternal bool
	Label      string
}

type Report struct {
	Data       map[Status][]Link
	Valid      int // How many valid links we have
	Broken     int // How many broken links we have
	External   int // How many external links we have
	Warning    int // How many warning links we have
	TotalLinks int // The total number of links processed
}

type Status string

const (
	Valid    = Status("valid")
	Broken   = Status("broken")
	Warning  = Status("warning")
	External = Status("external")
)

type validation struct {
	path    string                // The path being validated
	seen    *strset.Set           // Keep track of paths we have seen
	isLocal bool                  // TODO: Flag to determine validation live web site, or local file path
	queue   *list.List            // Keep track of links still to be processed per page
	tracked map[Status]*list.List // Keep track of collected links per status
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
		tracked: make(map[Status]*list.List),
	}
}

func (v validation) hasSeen(f string) bool {
	seen := v.seen.Has(f)
	if !seen {
		v.seen.Add(f)
	}
	return seen
}

func (v validation) Validate() (Report, error) {
	v.seen.Clear()

	err := filesystem.AFS.Walk(v.path, func(path string, info fs.FileInfo, err error) error {

		if err != nil {
			log.ErrorWithFields(err, log.Fields{"validation_path": path})
			return err
		}

		if !info.IsDir() {
			log.DebugWithFields("validation started", log.Fields{"validation_path": path})
			err = v.process(path)
			if err != nil {
				log.ErrorWithFields(err, log.Fields{"validation_path": path})
				return err
			}
		}

		return nil
	})

	if err != nil {
		return Report{}, err
	}
	return v.newReport(), nil
}

func (v validation) newReport() Report {

	report := Report{
		Data:       make(map[Status][]Link),
		Valid:      0,
		Broken:     0,
		External:   0,
		Warning:    0,
		TotalLinks: 0,
	}

	for s, links := range v.tracked {

		countedLinks := links.Len()
		report.TotalLinks += countedLinks
		collected := make([]Link, 0)

		var next *list.Element
		for e := links.Front(); e != nil; e = next {
			link := e.Value.(Link)
			collected = append(collected, link)
			next = e.Next()
		}

		report.Data[s] = collected

		switch s {
		case Valid:
			report.Valid = countedLinks
		case Broken:
			report.Broken = countedLinks
		case Warning:
			report.Warning = countedLinks
		case External:
			report.External = countedLinks
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

	v.queue.PushFront(Link{
		Uri:        path,
		Location:   path,
		IsExternal: false,
	})

	for v.queue.Len() > 0 {

		todo := v.queue.Front()
		v.queue.Remove(todo)
		link := todo.Value.(Link)

		if v.hasSeen(link.Uri) {
			continue
		}

		if link.Uri == "/" {
			continue
		}

		if link.IsExternal {
			v.reportLink(link, External, "")
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
			v.reportLink(link, Broken, fmt.Sprintf("Unable to open file %s: %s", link.Uri, err.Error()))
			continue
		}

		doc, err := goquery.NewDocumentFromReader(file)
		if err != nil {
			v.reportLink(link, Broken, fmt.Sprintf("file %s is propably not a valid HTML file: %s", link.Uri, err.Error()))
		} else {
			v.reportLink(link, Valid, "")
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
					v.reportLink(link, Warning, fmt.Sprintf("Unhandled url scheme: %s", href))
				} else if strings.Contains(validationHref, "#") {
					return
				}

				parsedLinkUrl, parseErr := url.Parse(href)

				if parseErr != nil {
					link.Message = parseErr.Error()
					return
				}

				link.IsExternal = len(parsedLinkUrl.Scheme) > 0

				finalUri := fmt.Sprintf("%s%s", v.path, href)

				v.reportLink(link, Valid, "")

				v.queue.PushFront(Link{
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

func (v validation) reportLink(link Link, status Status, message string) {

	link.Status = status
	link.Message = message

	collection, found := v.tracked[status]

	if !found {
		collection = list.New()
		v.tracked[status] = collection
	}

	collection.PushBack(link)

}
