package markdown

import (
	"fmt"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/conversion/colors"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/conversion/html"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/spf13/viper"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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
	{Key: "replaceBaseUrl", Func: replaceBaseUrl},
	{Key: "replaceBaseUrlWithSpaces", Func: replaceBaseUrlWithSpaces},
	{Key: "fixImages", Func: fixImages},
	{Key: "removeTargetBlank", Func: removeTargetBlank},
	{Key: "fixImagesWithAttributes", Func: fixImagesWithAttributes},
	{Key: "removeRawTags", Func: removeRawTags},
	{Key: "replaceCallOuts", Func: replaceCallOuts},
	{Key: "replaceTooltips", Func: replaceTooltips},
	{Key: "replaceIfVariables", Func: replaceIfVariables},
	{Key: "replaceComments", Func: replaceComments},
}

// Run each operation on a path, making sure to check with viper to see if we must
func Operate(path string) error {
	for _, operation := range markdownFileOperations {
		_, err := os.Stat(path)
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
	info, err := os.Stat(filepath.Join(dir, img))
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// Make all image paths absolute by adding the {{%path%}} shortcode
func fixImages(path string) error {
	return ManipulateMarkdown(path, nil, func(content []byte, w io.Writer) error {
		var replacements []replacement
		for _, matches := range ImageWithoutAttributes.FindAllStringSubmatch(string(content), -1) {
			if !paths.IsAbsURL(matches[2]) {
				img := fmt.Sprintf("![%s]({{%%path%%}}/%s)\n", matches[1], matches[3])
				replacements = append(replacements, replacement{Find: matches[0], Replace: img})
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

// Changes images with attributes to use a shortcode
func fixImagesWithAttributes(path string) error {
	return ManipulateMarkdown(path, nil, func(content []byte, w io.Writer) error {
		var replacements []replacement
		for _, matches := range ImgWithAttributesRe.FindAllStringSubmatch(string(content), -1) {
			var alt, src, attributes = matches[1], matches[3], matches[5]
			if paths.IsAbsURL(matches[2]) {
				src = matches[2] + src
			}

			var replacementShortcode = fmt.Sprintf(`{{< img src="%s" alt="%s"`, src, alt)
			for _, attrMatches := range AttributesRe.FindAllStringSubmatch(attributes, -1) {
				var key, value = attrMatches[1], attrMatches[2]
				replacementShortcode = replacementShortcode + fmt.Sprintf(` %s="%s"`, key, value)
			}
			replacementShortcode = replacementShortcode + " >}}"
			replacements = append(replacements, replacement{Find: matches[0], Replace: replacementShortcode})
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

// Replaces references to site.baseurl with a shortcode
func replaceBaseUrl(path string) error {
	return simpleReplaceContentInMarkdown(path, []string{"{{site.baseurl}}"}, "{{% baseurl %}}")
}

func replaceBaseUrlWithSpaces(path string) error {
	return simpleReplaceContentInMarkdown(path, []string{"{{ site.baseurl }}"}, "{{% baseurl %}}")
}

// Removes target="_blank" (this is automatically for external links)
func removeTargetBlank(path string) error {
	return simpleReplaceContentInMarkdown(path, []string{`{: target="_blank"}`, `{:target="_blank"}`}, ``)
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
			innerContent := matches[4]
			openSymbol, closeSymbol := wrapSymbols(innerContent)
			replacements = append(replacements, replacement{Find: matches[0], Replace: fmt.Sprintf("{{%s callout level=\"%s\" title=\"%s\"%s}}\n  %s\n{{%s /callout %s}}", openSymbol, matches[1], title, closeSymbol, innerContent, openSymbol, closeSymbol)})
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
			replacements = append(replacements, replacement{Find: matches[0], Replace: "{{< tooltip \"" + matches[1] + "\" >}}"})
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
