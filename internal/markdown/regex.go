package markdown

import "regexp"

// This regex with capture groups can break a markdown file into it's
// front matter and content sections
var MarkdownRe = regexp.MustCompile(`^(?s:(---\n)(.*)(---\n)(.*))$`)

// These regexes with capture groups are to assist in manipulating markdown content
var ImgWithBaseUrlRe = regexp.MustCompile(`\(\{\{\s?site.baseurl\s?\}\}/(.*\.(png|jpg|jpeg|gif|svg))\)`)
var ImgWithAttributesRe = regexp.MustCompile(`!\[(.*)\]\((.*\.(png|jpg|jpeg|gif|svg))\)\{: (.*)\}`)
var AttributesRe = regexp.MustCompile(`(\w+)="([^\"]+)"`)
var CalloutRe = regexp.MustCompile(`<div class="presidium-([\w\-]+)">\s*(<span>(.*)<\/span>)?\s*(.*)\s*<\/div>`)
var TooltipRe = regexp.MustCompile(`\[([^(.]*)]\(#\s*'presidium-tooltip'\)`)
var FrontmatterRe = regexp.MustCompile(`([^:.]*)\s*:\s*(.*)\n?`)
var IfVariablesRe = regexp.MustCompile(`(?msU){% if site.(\w+) == "(\w+)" %}(.+)({% elsif site.(\w+) == "(\w+)" %}(.+))?({% else %}(.+))?{% endif %}`)