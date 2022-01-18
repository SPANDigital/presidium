package markdown

import "regexp"

// This regex with capture groups can break a markdown file into it's
// front matter and contentOf sections
var MarkdownRe = regexp.MustCompile(`^(?s:(---\n)(.*)(---\n)(.*))$`)

// These regexes with capture groups are to assist in manipulating markdown contentOf
var ImageRe = regexp.MustCompile(`(?mi)!\[(.*)\]\(({{.+}})?(.*\/)*(.*\.(png|jpg|jpeg|gif|svg))\)(\{: (.*)\})?`)
var HtmlImageRe = regexp.MustCompile(`<img([^>]+)/>`)
var SourceRe = regexp.MustCompile(`(src)="({{.+}})?(.*\/)*(.*\.(png|jpg|jpeg|gif|svg))"`)
var AttributesRe = regexp.MustCompile(`(\w+)="([^\"]+)"`)
var CalloutRe = regexp.MustCompile(`<div class="presidium-([\w\-]+)">\s*(<span>(.*)<\/span>)?\s*(.*)\s*<\/div>`)
var TooltipRe = regexp.MustCompile(`(?m)\[([^(.]*)]\(#\s*'presidium-tooltip'\)`)
var FrontmatterRe = regexp.MustCompile(`([^:.]*)\s*:\s*(.*)\n?`)
var IfVariablesRe = regexp.MustCompile(`(?msU){% if ([^}]*?) %}(.+)({% elsif ([^}]*?) %}(.+))?({% else %}(.+))?{% endif %}`)
var IfConditionRe = regexp.MustCompile(`(?ms)(and|or)?\s??site.(\w+) ([!=]=) "(\w+)"`)
var IfConditionShortcodeRe = regexp.MustCompile(`(?ms)site.(\w+) == "(\w+)"`)
var ContainsShortcodeRe = regexp.MustCompile(`(?ms){{[%<].*?[%>]}}`)
var CommentRe = regexp.MustCompile(`(?ms){% comment %}(.*?){% endcomment %}`)
var TableBody = regexp.MustCompile(`(?:\|[^\n]+\|?\r?\n?){2,}`)
var TableHeader = regexp.MustCompile(`(\|[^\n]+\|?\r?\n)((?:\|:?[-]+:?)+\|?)`)
