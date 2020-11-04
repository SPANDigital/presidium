package markdown

import (
	"fmt"
	"github.com/SPANDigital/presidium-hugo/internal/colors"
	"github.com/SPANDigital/presidium-hugo/internal/html"
	"github.com/spf13/viper"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// markdownFileOperationFunc represents an operation on a markdown file
type markdownFileOperationFunc func (path string) error

// operationInstruction maps a viper key to a markdownFileOperationFunc
type operationInstruction struct {
	Key string
	Func markdownFileOperationFunc
}

// Represents a discrete find and replace
type replacement struct {
	Find string
	Replace string
}

// All possible markdownfile operations in the order in which they will run
var markdownFileOperations = []operationInstruction{
	{Key: "eraseMarkdownWithNoContent", Func: eraseMarkdownWithNoContent},
	{Key: "commonmarkAttributes", Func: replaceCommonmarkAttributes},
	{Key: "fixImages", Func: fixImages},
	{Key: "replaceBaseUrl", Func: replaceBaseUrl},
	{Key: "removeTargetBlank", Func: removeTargetBlank},
	{Key: "fixImagesWithAttributes", Func: fixImagesWithAttributes},
	{Key: "removeRawTags", Func: removeRawTags},
	{Key: "replaceCallOuts", Func: replaceCallOuts},
	{Key: "replaceTooltips", Func: replaceTooltips},
	{Key: "replaceIfVariables", Func: replaceIfVariables},
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

// Fixes image urls of the form:
// {{ site.baseurl }}/pathtoimage.ext
func fixImages(path string) error {
	return ManipulateMarkdown(path, nil, func(content []byte, w io.Writer) error {
		var replacements []replacement
		for _, matches := range ImgWithBaseUrlRe.FindAllStringSubmatch(string(content), -1) {
			if imgIsInSameDir(path, matches[1]) {
				replacements = append(replacements, replacement{Find: matches[0], Replace: "(" + matches[1] + ")"})
			}
		}

		strContent := string(content)
		for _, replacement := range replacements {
			fmt.Println("Replacing",colors.Labels.Unwanted(replacement.Find), "with", colors.Labels.Wanted(replacement.Replace), "in", colors.Labels.Info(path))
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
			var alt, src, attributes = matches[1], matches[2], matches[4]
			var replacementShortcode = fmt.Sprintf(`{{< img src="%s" alt="%s"`, src, alt)
			for _, attrMatches := range AttributesRe.FindAllStringSubmatch(attributes, -1) {
				var key, value = attrMatches[1], attrMatches[2]
				replacementShortcode = replacementShortcode + fmt.Sprintf( ` %s="%s"`, key, value)
			}
			replacementShortcode = replacementShortcode + " >}}"
			replacements = append(replacements, replacement{Find: matches[0], Replace: replacementShortcode})
		}

		strContent := string(content)
		for _, replacement := range replacements {
			fmt.Println("Replacing",colors.Labels.Unwanted(replacement.Find), "with", colors.Labels.Wanted(replacement.Replace), "in", colors.Labels.Info(path))
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

// Removes target="_blank" (this is automatically for external links)
func removeTargetBlank(path string) error {
	return simpleReplaceContentInMarkdown(path, []string{`{: target="_blank"}`, `{:target="_blank"}`}, ``)
}

// Perform find and replace operations on markdown content
func simpleReplaceContentInMarkdown(path string, finds[] string, replace string) error {
	for _, find := range finds {
		err := ManipulateMarkdown(path, nil, func(content []byte, w io.Writer) error {
			strContent := string(content)
			if strings.Index(strContent, find) > -1 {
				if replace == "" {
					fmt.Println("Blanking", colors.Labels.Unwanted(find),"in",path)
				} else {
					fmt.Println("Replacing", colors.Labels.Unwanted(find),"with",colors.Labels.Wanted(replace),"in",colors.Labels.Info(path))
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
					fmt.Println("Blanking", colors.Labels.Unwanted(replacement.Find),"in",path)
				} else {
					fmt.Println("Replacing", colors.Labels.Unwanted(replacement.Find),"with",colors.Labels.Wanted(replacement.Replace),"in",colors.Labels.Info(path))
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
			fmt.Println("Found", colors.Labels.Unwanted(len(allMatches)),"callouts in", colors.Labels.Info(path))
		}
		for _, matches := range allMatches {
			fmt.Println("Replacing callout ", colors.Labels.Unwanted(matches[2]),"and level", colors.Labels.Unwanted(matches[1]), " with shortcode in ", colors.Labels.Info(path))
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
			fmt.Println("Found", colors.Labels.Unwanted(len(allMatches)),"callouts in", colors.Labels.Info(path))
		}
		for _, matches := range allMatches {
			fmt.Println("Creating tooltip for ", colors.Labels.Info(matches[1]))
			replacements = append(replacements, replacement{Find: matches[0], Replace: "{{< tooltip \"" + matches[1] + "\" >}}"})
		}
		for _, rep := range replacements {
			strContent = strings.ReplaceAll(strContent, rep.Find, rep.Replace)
		}
		_, err := io.WriteString(w, strContent)
		return err
	})
}

func replaceIfVariables(path string) error {
	return ManipulateMarkdown(path, nil, func(content []byte, w io.Writer) error {
		var replacements []replacement
		strContent := string(content)
		allMatches := IfVariablesRe.FindAllStringSubmatch(strContent, -1)
		if allMatches != nil {
			fmt.Println("Found", colors.Labels.Unwanted(len(allMatches)),"switch-params in", colors.Labels.Info(path))
		}
		for _, matches := range allMatches {
			param := matches[1]
			ifValue := matches[2]
			fmt.Println("Creating switch-param for ", colors.Labels.Info(param))
			ifInnerContent := matches[3]
			openSymbol, closeSymbol := wrapSymbols(ifInnerContent)
			innerContent := fmt.Sprintf(`{{%s when "%s" %s}}%s{{%s /when %s}}`, openSymbol, ifValue, closeSymbol, ifInnerContent, openSymbol, closeSymbol)
			if matches[4] != "" {
				elifValue := matches[6]
				elifInnerContent := matches[7]
				openSymbol, closeSymbol = wrapSymbols(elifInnerContent)
				innerContent += fmt.Sprintf(`{{%s when "%s" %s}}%s{{%s /when %s}}`, openSymbol, elifValue, closeSymbol, elifInnerContent, openSymbol, closeSymbol)
			}
			if matches[8] != "" {
				defaultInnerContent := matches[9]
				openSymbol, closeSymbol = wrapSymbols(defaultInnerContent)
				innerContent += fmt.Sprintf(`{{%s default %s}}%s{{%s /default %s}}`, openSymbol, closeSymbol, defaultInnerContent, openSymbol, closeSymbol)
			}
			replacements = append(replacements, replacement{Find: matches[0], Replace: fmt.Sprintf(`{{< with-param "%s" >}}%s{{< /with-param >}}`, param, innerContent)})

		}
		for _, rep := range replacements {
			strContent = strings.ReplaceAll(strContent, rep.Find, rep.Replace)
		}
		_, err := io.WriteString(w, strContent)
		return err
	})
}