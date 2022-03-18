package utils

import (
	"regexp"
	"strings"
)

// UnSlugify turns "something-like_this" into "Something Like This"
func UnSlugify(name string) string {
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

// Slugify replaces all non word chars with a "-"
// turns "v0 .18.6_8." into "v0-18-6-8"
func Slugify(name string) string {
	var nonWordRe = regexp.MustCompile(`(?m)(\W|_)+`)
	slug := nonWordRe.ReplaceAllString(name, "-")
	return strings.Trim(slug, "-")
}

// TitleToSlug Take a capitalized title and turn it into a slug
func TitleToSlug(title string) string {
	title = strings.ToLower(title)
	title = strings.Replace(title, "&", "and", -1)
	title = Slugify(title)
	return title
}
