package validation

import (
	"container/list"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"strings"
)

type Link struct {
	Uri     string
	Valid   bool
	Message string
	Level   int
	Text    string
	// --- private fields here (no need to expose it to the caller): ---
	page string // Used to turn hrefs into URLS which can be references from the page.
}

type Validation struct {
	// Private fields  - no need to expose it to the caller!
	seed    string          // The first address to look at
	depth   int             // The depth to scan for links
	respond func(link Link) // Configures in New() to respond back to the caller.
}

func New(
	seed string,
	depth int,
	listener func(link Link)) (Validation, error) {

	if len(seed) == 0 {
		return Validation{}, errors.New("seed url cannot be empty")
	}

	_, err := url.Parse(seed)
	if err != nil {
		return Validation{}, err
	}

	if depth < 0 {
		return Validation{}, errors.New("depth must be 1 or more")
	}

	validation := Validation{
		seed:    seed,
		depth:   depth,
		respond: listener,
	}

	return validation, nil
}

func (validation *Validation) Start() {

	if validation.depth == 0 {
		return
	}

	level := 0
	queue := list.New()

	queue.PushFront(Link{
		Uri:     validation.seed,
		Valid:   false,
		Message: "pending",
		Level:   level + 1, // Start at level 1
		Text:    "",
		page:    validation.seed,
	})

	var link Link

	for {
		if queue.Len() == 0 {
			break
		}

		front := queue.Front()
		queue.Remove(front)
		link = front.Value.(Link)
		resp, err := http.Get(link.Uri)
		if err != nil {
			link.Valid = false
			link.Message = fmt.Sprintf("err: %v", err.Error())
			validation.respond(link)
		} else {
			link.Message = "ok"
			link.Valid = true
			validation.respond(link)
			drillOneLevelDown := level < validation.depth
			if drillOneLevelDown {
				level += 1
				doc, err := goquery.NewDocumentFromReader(resp.Body)
				if err != nil {
					continue
				}
				doc.Find("a[href]").Each(func(i int, item *goquery.Selection) {
					href, ok := item.Attr("href")
					if !ok {
						return
					}
					validationHref := strings.ToLower(href)
					validationHref = strings.TrimSpace(validationHref)
					if strings.HasPrefix(validationHref, "mailto:") ||
						strings.HasPrefix(validationHref, "tel:") ||
						strings.HasPrefix(validationHref, "/#") || validationHref == "/" {
						return
					}
					nextUrl := nextPageUrl(link, href)
					queue.PushBack(Link{
						page:    link.page,
						Uri:     nextUrl,
						Valid:   false,
						Message: "",
						Level:   level,
						Text:    strings.TrimSpace(item.Text()),
					})
				})
			}
		}
	}
}

func nextPageUrl(link Link, uri string) string {
	if strings.HasPrefix(uri, "/") {
		return fmt.Sprintf("%s%s", link.page, uri)
	} else {
		return uri
	}
}
