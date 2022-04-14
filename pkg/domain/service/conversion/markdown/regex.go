package markdown

import "regexp"

// MarkdownRe matches front matter and content see https://regex101.com/r/pOZOTG/1
var MarkdownRe = regexp.MustCompile(`(?m)^(?s:(^---\n)(.*?)(^---(?:\n|$))(.*))$`)

// FrontMatterRe matches front matter, see https://regex101.com/r/bXta8n/1
var FrontMatterRe = regexp.MustCompile(`([^:.]*)\s*:\s*(.*)\n?`)

// ImageRe matches markdown images, see https://regex101.com/r/Xen6Cp/2
var ImageRe = regexp.MustCompile(`(?mi)!\[(.*)\]\(({{.+}})?([^)\s]*\/)*([^)]*\.(png|jpg|jpeg|gif|svg))\)(\{:\s?(.*)\})?`)

// HtmlImageRe matches html images, see https://regex101.com/r/Vepscq/1
var HtmlImageRe = regexp.MustCompile(`<img([^>]+)/>`)

// FigureRe matches a figure caption that is preceded by a blank line https://regex101.com/r/fWZxR0/1
var FigureRe = regexp.MustCompile(`(\n)(\n)(\*Figure([^\*]+)\*)`)

// SourceRe matches the src attribute of images, see https://regex101.com/r/UfNttw/1
var SourceRe = regexp.MustCompile(`(src)="({{.+}})?([^"\s]*\/)*([^"]*\.(png|jpg|jpeg|gif|svg))"`)

// AttributesRe matches html attributes, see https://regex101.com/r/FgkHDJ/1
var AttributesRe = regexp.MustCompile(`(\w+)="([^\"]+)"`)

// CalloutRe matches callouts https://regex101.com/r/8CnxPG/1
var CalloutRe = regexp.MustCompile(`(?ms)<div class="presidium-([\w\-]+)">\s*(<span>(.*?)</span>)?(.*?)</div>`)

// TooltipRe matches tooltips, see https://regex101.com/r/McE770/1
var TooltipRe = regexp.MustCompile(`(?m)\[([^(.]*)]\(({{.+}})?([^)\s]*?)\s*'presidium-tooltip'\)`)

// TableBodyRe matches tables, see https://regex101.com/r/K9xbsc/1
var TableBodyRe = regexp.MustCompile(`(?:\|[^\n]+\|?\r?\n?\s*){2,}`)

// TableHeaderRe matches the header section of a table, see https://regex101.com/r/zeaXvT/1
var TableHeaderRe = regexp.MustCompile(`(\|[^\n]+\|?\r?\n)?(\s*(?:\|:?\s*[-]+:?\s*)+\|?)`)

// IfVariablesRe matches jekyll conditional logic https://regex101.com/r/n1vbLY/1
var IfVariablesRe = regexp.MustCompile(`(?msU){% if ([^}]*?) %}(.+)({% elsif ([^}]*?) %}(.+))?({% else %}(.+))?{% endif %}`)

// IfConditionRe matches the entire if condition, see https://regex101.com/r/BHIU0V/1
var IfConditionRe = regexp.MustCompile(`(?ms)(and|or)?\s??site.(\w+) ([!=]=) "(\w+)"`)

// IfConditionShortcodeRe matches individual conditions, see https://regex101.com/r/BvAtyC/1
var IfConditionShortcodeRe = regexp.MustCompile(`(?ms)site.(\w+) == "(\w+)"`)

// ContainsShortcodeRe matches shortcodes, see https://regex101.com/r/aKgUVk/1
var ContainsShortcodeRe = regexp.MustCompile(`(?ms){{[%<].*?[%>]}}`)

// CommentRe matches comments, see https://regex101.com/r/31XAPX/1
var CommentRe = regexp.MustCompile(`(?ms){% comment %}(.*?){% endcomment %}`)

// EmptyLineRe matches empty lines, see https://regex101.com/r/rNxkKo/3
var EmptyLineRe = regexp.MustCompile(`(?m)^(?:[\t ]*(?:\r?\n|\r))+`)

// WeightAndSlugRe matches the weight and slug in a filename, see https://regex101.com/r/DZmKgp/2
var WeightAndSlugRe = regexp.MustCompile(`((([\d.]+)([a-z]?))[-_])?(.+?)(\.[^.\s]+)?$`)

// WeightRe matches the weight in a filename, see https://regex101.com/r/90itMQ/1
var WeightRe = regexp.MustCompile(`^([\d.]+)([a-z]?)[-_]`)

// SpaceRe matches one or more spaces, see https://regex101.com/r/rRXM3s/1
var SpaceRe = regexp.MustCompile(`\s+`)
