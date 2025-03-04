package markdown

import (
	"fmt"
	"github.com/SPANDigital/presidium-hugo/pkg/filesystem"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/conversion/colors"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/conversion/html"
	"github.com/gohugoio/hugo/helpers"
	"github.com/spf13/viper"
)

// markdownFileOperationFunc represents an operation on a markdown file
type markdownFileOperationFunc func(path string) error

// operationInstruction maps a viper key to a markdownFileOperationFunc
type operationInstruction struct {
	Key  string
	Func markdownFileOperationFunc
}

// Represents a discrete find and replace
type replacement struct {
	Find    string
	Replace string
}

// All possible markdownfile operations in the order in which they will run
var markdownFileOperations = []operationInstruction{
	{Key: "eraseMarkdownWithNoContent", Func: eraseMarkdownWithNoContent},
	{Key: "commonmarkAttributes", Func: replaceCommonmarkAttributes},
	{Key: "fixImages", Func: fixImages},
	{Key: "fixFigureCaptions", Func: fixFigureCaptions},
	{Key: "fixHtmlImages", Func: fixHtmlImages},
	{Key: "replaceBaseUrl", Func: replaceBaseUrl},
	{Key: "replaceBaseUrlWithSpaces", Func: replaceBaseUrlWithSpaces},
	{Key: "removeTargetBlank", Func: removeTargetBlank},
	{Key: "removeRawTags", Func: removeRawTags},
	{Key: "replaceCallOuts", Func: replaceCallOuts},
	{Key: "replaceTooltips", Func: replaceTooltips},
	{Key: "replaceIfVariables", Func: replaceIfVariables},
	{Key: "replaceComments", Func: replaceComments},
	{Key: "addTableHeaders", Func: addTableHeaders},
}

// Run each operation on a path, making sure to check with viper to see if we must
func Operate(path string) error {
	for _, operation := range markdownFileOperations {
		_, err := filesystem.AFS.Stat(path)
		if !os.IsNotExist(err) {
			if viper.GetBool(operation.Key) {
				err := operation.Func(path)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func eraseMarkdownWithNoContent(path string) error {
	if !IsRecognizableMarkdown(path) {
		fmt.Println("Erasing", colors.Labels.Info(path))
		return os.Remove(path)
	}
	return nil
}

// Tells us if an image is in same directory as markdown path
func imgIsInSameDir(path string, img string) bool {
	dir := filepath.Dir(path)
	info, err := filesystem.AFS.Stat(filepath.Join(dir, img))
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func fixHtmlImages(path string) error {
	return ManipulateMarkdown(path, nil, func(content []byte, w io.Writer) error {
		var replacements []replacement
		images := HtmlImageRe.FindAllStringSubmatch(string(content), -1)
		for _, image := range images {
			srcAttr := SourceRe.FindStringSubmatch(image[1])
			src := parseSource(path, srcAttr[3], srcAttr[4], false)
			findSource := fmt.Sprintf("src=\"%s\"", src)
			replacement := replacement{Find: srcAttr[0], Replace: findSource}
			replacements = append(replacements, replacement)
		}

		strContent := string(content)
		for _, replacement := range replacements {
			fmt.Println("Replacing", colors.Labels.Unwanted(replacement.Find), "with", colors.Labels.Wanted(replacement.Replace), "in", colors.Labels.Info(path))
			strContent = strings.ReplaceAll(strContent, replacement.Find, replacement.Replace)
		}
		_, err := io.WriteString(w, strContent)
		return err
	})
}

func fixImages(path string) error {
	return ManipulateMarkdown(path, nil, func(content []byte, w io.Writer) error {
		images := ImageRe.FindAllStringSubmatch(string(content), -1)
		var replacements []replacement
		for _, image := range images {
			hasTags := len(image[6]) > 0
			src := parseSource(path, image[3], image[4], hasTags)
			if hasTags {
				replacements = append(replacements, parseImageWithTags(src, image))
			} else {
				replacements = append(replacements, parseImageWithoutTags(src, image))
			}
		}

		strContent := string(content)
		for _, replacement := range replacements {
			fmt.Println("Replacing", colors.Labels.Unwanted(replacement.Find), "with", colors.Labels.Wanted(replacement.Replace), "in", colors.Labels.Info(path))
			strContent = strings.ReplaceAll(strContent, replacement.Find, replacement.Replace)
		}
		_, err := io.WriteString(w, strContent)
		return err
	})
}

// fixFigureCaptions removes the blank line between figure captions and images
func fixFigureCaptions(path string) error {
	return ManipulateMarkdown(path, nil, func(content []byte, w io.Writer) error {
		figures := FigureRe.FindAllStringSubmatch(string(content), -1)
		var replacements []replacement
		for _, figure := range figures {
			nlRe := regexp.MustCompile("\n+")
			caption := nlRe.ReplaceAllString(figure[0], "\n")
			replacements = append(replacements, replacement{
				Find:    figure[0],
				Replace: caption,
			})
		}

		strContent := string(content)
		for _, replacement := range replacements {
			fmt.Println("Replacing", colors.Labels.Unwanted(replacement.Find), "with", colors.Labels.Wanted(replacement.Replace), "in", colors.Labels.Info(path))
			strContent = strings.ReplaceAll(strContent, replacement.Find, replacement.Replace)
		}
		_, err := io.WriteString(w, strContent)
		return err
	})
}

// shortcodes in strings are not supported atm
// https://github.com/gohugoio/hugo/issues/6703
func parseSource(path string, dir string, filename string, rawSource bool) string {
	src := dir + filename
	ps := helpers.PathSpec{}
	if isAbs, _ := ps.IsAbsURL(src); isAbs {
		return src
	}
	if imgIsInSameDir(path, filename) {
		if rawSource {
			return filename
		}
		return fmt.Sprintf("{{%%path%%}}/%s", filename)
	}

	idx := strings.LastIndex(src, "/images/")
	if idx > -1 {
		src = src[idx:]
	}
	if rawSource {
		return src
	}
	return filepath.Clean(fmt.Sprintf("{{%% baseurl %%}}/%s", src))
}

func parseImageWithoutTags(src string, image []string) replacement {
	img := fmt.Sprintf("![%s](%s)", image[1], src)
	caption := image[10]
	if len(caption) > 0 {
		img += caption
	}
	return replacement{Find: image[0], Replace: img}
}

func parseImageWithTags(src string, image []string) replacement {
	var alt, attributes, caption = image[1], image[7], image[11]
	var imgShortcode = fmt.Sprintf(`{{< img src="%s"`, src)
	var styleAttributes = []string{"width", "height"}
	var styles string
	for _, attrMatches := range AttributesRe.FindAllStringSubmatch(attributes, -1) {
		var key, value = attrMatches[1], attrMatches[2]
		if contains(styleAttributes, key) {
			styles += fmt.Sprintf(`%s:%s;`, key, value)
		} else {
			imgShortcode = imgShortcode + fmt.Sprintf(` %s="%s"`, key, value)
		}
	}

	if alt != caption {
		imgShortcode += fmt.Sprintf(` alt="%s"`, alt)
	}

	if len(caption) > 0 {
		imgShortcode = imgShortcode + fmt.Sprintf(` caption="%s"`, caption)
	}

	if len(styles) > 0 {
		imgShortcode = imgShortcode + fmt.Sprintf(` style="%s"`, styles)
	}

	imgShortcode = imgShortcode + " >}}"
	return replacement{Find: image[0], Replace: imgShortcode}
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

// Adds empty table headers for all tables with partial or no headers
func addTableHeaders(path string) error {
	return ManipulateMarkdown(path, nil, func(content []byte, w io.Writer) error {
		var replacements []replacement

		tables := TableBodyRe.FindAllString(string(content), -1)
		for _, table := range tables {
			header := TableHeaderRe.FindStringSubmatch(table)
			if header == nil || len(header[1]) == 0 {
				headerTable := appendHeader(table, header != nil)
				replacements = append(replacements, replacement{Find: table, Replace: headerTable})
			}
		}

		strContent := string(content)
		for _, replacement := range replacements {
			strContent = strings.ReplaceAll(strContent, replacement.Find, replacement.Replace)
		}

		_, err := io.WriteString(w, strContent)
		return err
	})
}

func appendHeader(table string, partialHeader bool) string {
	pipeSelector := regexp.MustCompile(`\|`)
	var mostCols []string
	for _, line := range strings.Split(table, "\n") {
		lineCols := pipeSelector.FindAllString(line, -1)
		if len(lineCols) > len(mostCols) {
			mostCols = lineCols
		}
	}

	header := strings.Join(mostCols, " ")
	divider := strings.Join(mostCols, "-")
	if partialHeader {
		return fmt.Sprintf("%s\n%s", header, table)
	}
	return fmt.Sprintf("%s\n%s\n%s", header, divider, table)
}

// Replaces references to site.baseurl with a shortcode
func replaceBaseUrl(path string) error {
	return simpleReplaceContentInMarkdown(path, []string{"{{site.baseurl}}"}, "{{% baseurl %}}")
}

func replaceBaseUrlWithSpaces(path string) error {
	return simpleReplaceContentInMarkdown(path, []string{"{{ site.baseurl }}"}, "{{% baseurl %}}")
}

// Removes target="_blank" (this is automatically for external links)
func removeTargetBlank(path string) error {
	return simpleReplaceContentInMarkdown(path, []string{
		`{: target="_blank"}`,
		`{:target="_blank"}`,
		`{:target="\_blank"}`,
		`{: target="\_blank"}`}, ``)
}

// Perform find and replace operations on markdown contentOf
func simpleReplaceContentInMarkdown(path string, finds []string, replace string) error {
	for _, find := range finds {
		err := ManipulateMarkdown(path, nil, func(content []byte, w io.Writer) error {
			strContent := string(content)
			if strings.Index(strContent, find) > -1 {
				if replace == "" {
					fmt.Println("Blanking", colors.Labels.Unwanted(find), "in", path)
				} else {
					fmt.Println("Replacing", colors.Labels.Unwanted(find), "with", colors.Labels.Wanted(replace), "in", colors.Labels.Info(path))
				}
				_, err := io.WriteString(w, strings.ReplaceAll(string(content), find, replace))
				return err
			} else {
				_, err := w.Write(content)
				return err
			}
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func replaceContentInMarkdown(path string, replacements []replacement) error {
	for _, replacement := range replacements {
		err := ManipulateMarkdown(path, nil, func(content []byte, w io.Writer) error {
			strContent := string(content)
			if strings.Index(strContent, replacement.Find) > -1 {
				if replacement.Replace == "" {
					fmt.Println("Blanking", colors.Labels.Unwanted(replacement.Find), "in", path)
				} else {
					fmt.Println("Replacing", colors.Labels.Unwanted(replacement.Find), "with", colors.Labels.Wanted(replacement.Replace), "in", colors.Labels.Info(path))
				}
				_, err := io.WriteString(w, strings.ReplaceAll(string(content), replacement.Find, replacement.Replace))
				return err
			} else {
				_, err := w.Write(content)
				return err
			}
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func removeRawTags(path string) error {
	return replaceContentInMarkdown(path, []replacement{{Find: "{% raw %}", Replace: ""}, {Find: "{% endraw %}", Replace: ""}})
}

// Convert to use commonmark attributes
func replaceCommonmarkAttributes(path string) error {
	return simpleReplaceContentInMarkdown(path, []string{"{: "}, "{")
}

func wrapSymbols(content string) (string, string) {
	if html.ContainsHTML(content) {
		return "<", ">"
	} else {
		return "%", "%"
	}
}

func replaceCallOuts(path string) error {
	return ManipulateMarkdown(path, nil, func(content []byte, w io.Writer) error {
		var replacements []replacement
		strContent := string(content)
		allMatches := CalloutRe.FindAllStringSubmatch(string(content), -1)
		if allMatches != nil {
			fmt.Println("Found", colors.Labels.Unwanted(len(allMatches)), "callouts in", colors.Labels.Info(path))
		}
		for _, matches := range allMatches {
			fmt.Println("Replacing callout ", colors.Labels.Unwanted(matches[2]), "and level", colors.Labels.Unwanted(matches[1]), " with shortcode in ", colors.Labels.Info(path))
			title := matches[3]
			innerContent := EmptyLineRe.ReplaceAllString(matches[4], "")
			innerContent = strings.TrimSpace(innerContent)
			openSymbol, closeSymbol := wrapSymbols(innerContent)
			replacements = append(replacements, replacement{Find: matches[0], Replace: fmt.Sprintf("{{%s callout level=\"%s\" title=\"%s\"%s}}\n%s\n{{%s /callout %s}}", openSymbol, matches[1], title, closeSymbol, innerContent, openSymbol, closeSymbol)})
		}
		for _, rep := range replacements {
			strContent = strings.ReplaceAll(strContent, rep.Find, rep.Replace)
		}
		_, err := io.WriteString(w, strContent)
		return err
	})
}

func replaceTooltips(path string) error {
	return ManipulateMarkdown(path, nil, func(content []byte, w io.Writer) error {
		var replacements []replacement
		strContent := string(content)
		allMatches := TooltipRe.FindAllStringSubmatch(string(content), -1)
		if allMatches != nil {
			fmt.Println("Found", colors.Labels.Unwanted(len(allMatches)), "tooltips in", colors.Labels.Info(path))
		}

		for _, matches := range allMatches {
			fmt.Println("Creating tooltip for ", colors.Labels.Info(matches[1]))
			fmt.Println("matches[0]: ", matches[0])
			ref := strings.Replace(matches[3], "#", "", 1)
			tooltip := fmt.Sprintf("{{< tooltip \"%s\" \"%s\" >}}", matches[1], ref)
			replacements = append(replacements, replacement{Find: matches[0], Replace: tooltip})
		}
		for _, rep := range replacements {
			strContent = strings.ReplaceAll(strContent, rep.Find, rep.Replace)
		}
		_, err := io.WriteString(w, strContent)
		return err
	})
}

func ensureCamelCase(input string) string {
	var snake = regexp.MustCompile("_([A-Za-z])")
	return snake.ReplaceAllStringFunc(input, func(s string) string {
		return strings.ToUpper(strings.Replace(s, "_", "", -1))
	})
}

func stripTooltips(strContent string) string {
	var TooltipRe = regexp.MustCompile(`(?ms){{< tooltip "(.*?)" >}}`)
	return TooltipRe.ReplaceAllString(strContent, "$1")
}

func parseIfStatements(strContent string) string {
	replacements := []replacement{}
	allMatches := IfVariablesRe.FindAllStringSubmatch(strContent, -1)
	if allMatches != nil {
		fmt.Println("Found", len(allMatches), "switch-params")
	}
	for _, matches := range allMatches {
		ifContent := matches[1]
		fmt.Println("Creating conditional statements for for ", ifContent)
		ifInnerContent := matches[2]
		innerContent := ""
		containsShortCode := ContainsShortcodeRe.MatchString(matches[0])

		// This appears to be more compatible than `wrapSymbols` when dealing with mixed contentOf
		openSymbol, closeSymbol := "%", "%"

		if containsShortCode {
			innerContent = fmt.Sprintf("{{%s when %s %s}}%s{{%s /when %s}}", openSymbol, parseIfConditionShortcode(ifContent), closeSymbol, ifInnerContent, openSymbol, closeSymbol)
			if matches[3] != "" {
				elifContent := matches[4]
				elifInnerContent := matches[5]
				innerContent += fmt.Sprintf("{{%s when %s %s}}%s{{%s /when %s}}", openSymbol, parseIfConditionShortcode(elifContent), closeSymbol, elifInnerContent, openSymbol, closeSymbol)
			}
			if matches[6] != "" {
				defaultInnerContent := matches[7]
				innerContent += fmt.Sprintf("{{%s default %s}}%s{{%s /default %s}}", openSymbol, closeSymbol, defaultInnerContent, openSymbol, closeSymbol)
			}
		} else {
			innerContent = fmt.Sprintf(`{{ if %s }}%s`, parseIfCondition(ifContent), ifInnerContent)
			if matches[3] != "" {
				elifContent := matches[4]
				elifInnerContent := matches[5]
				innerContent += fmt.Sprintf(`{{ else if %s }}%s`, parseIfCondition(elifContent), elifInnerContent)
			}
			if matches[6] != "" {
				defaultInnerContent := matches[7]
				innerContent += fmt.Sprintf(`{{ else }}%s`, defaultInnerContent)
			}
			innerContent += `{{ end }}`
			openSymbol, closeSymbol := wrapSymbols(innerContent)
			innerContent = fmt.Sprintf("{{%s if.inline %s}}%s{{%s /if.inline %s}}", openSymbol, closeSymbol, innerContent, openSymbol, closeSymbol)
		}
		replacements = append(replacements, replacement{Find: matches[0], Replace: innerContent})
	}
	for _, rep := range replacements {
		strContent = strings.ReplaceAll(strContent, rep.Find, rep.Replace)
	}

	return strContent
}

func formatConditional(matches []string) string {
	operator := "eq"
	if matches[3] == "!=" {
		operator = "ne"
	}
	return fmt.Sprintf("(%s $.Page.Site.Params.%s \"%s\")", operator, ensureCamelCase(matches[2]), matches[4])
}

func parseIfCondition(content string) string {
	allMatches := IfConditionRe.FindAllStringSubmatch(content, -1)
	output := []string{}
	for _, matches := range allMatches {
		if matches[1] != "" {
			output = append(output, fmt.Sprintf("| %s", matches[1]))
		}
		output = append(output, formatConditional(matches))
	}
	return strings.Join(output, " ")
}

func parseIfConditionShortcode(content string) string {
	allMatches := IfConditionShortcodeRe.FindAllStringSubmatch(content, -1)
	output := []string{}
	for _, matches := range allMatches {
		output = append(output, fmt.Sprintf(`"%s" "%s"`, ensureCamelCase(matches[1]), matches[2]))
	}
	return strings.Join(output, " ")
}

func replaceIfVariables(path string) error {
	return ManipulateMarkdown(path, nil, func(content []byte, w io.Writer) error {
		strContent := string(content)
		newContent := parseIfStatements(strContent)
		_, err := io.WriteString(w, newContent)
		return err
	})
}

func parseComments(strContent string) string {
	return CommentRe.ReplaceAllString(strContent, "{{< comment.inline >}}\n{{/*$1*/}}\n{{< /comment.inline >}}")
}

func replaceComments(path string) error {
	return ManipulateMarkdown(path, nil, func(content []byte, w io.Writer) error {
		strContent := string(content)
		strContent = parseComments(strContent)
		_, err := io.WriteString(w, strContent)
		return err
	})
}
