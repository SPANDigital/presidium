package configtranslation

import (
	"fmt"
	"github.com/SPANDigital/presidium-hugo/pkg/config"
	"io/ioutil"
	"regexp"

	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/conversion/colors"
	"github.com/SPANDigital/presidium-hugo/pkg/log"
	"gopkg.in/yaml.v2"
)

type JekyllShow struct {
	Status bool `yaml:"status""`
	Author bool `yaml:"author""`
	Roles  bool `yaml:"role"`
}

type JekyllExternal struct {
	AuthorsUrl string `yaml:"authors-url:`
}

type JekyllSectionItem struct {
	Title      string `yaml:"title""`
	Url        string `yaml:"url""`
	Collection string `yaml:"collection""`
}

type JekyllConfig struct {
	Name        string              `yaml:"name"`
	Description string              `yaml:"description"`
	Baseurl     string              `yaml:"baseurl"`
	Footer      string              `yaml:"footer"`
	Logo        string              `yaml:"logo"`
	Audience    string              `yaml:"audience"`
	Scope       string              `yaml:"scope"`
	AppleScope  string              `yaml:"apple_scope"`
	Show        interface{}         `yaml:"show"`
	External    JekyllExternal      `yaml:"external"`
	Sections    []JekyllSectionItem `yaml:"sections"`
	Roles       Roles               `yaml:"roles"`
}

func (j *JekyllConfig) reparsedShowOptionsAsSequenceDictionaries() bool {

	if seqOpts, ok := j.Show.([]interface{}); ok {
		parsed := JekyllShow{
			Status: true,
			Author: true,
			Roles:  true,
		}
		for _, opt := range seqOpts {
			if values, ok := opt.(map[interface{}]interface{}); ok {
				var name string
				var flagged bool
				for k, v := range values {
					if name, ok = k.(string); ok {
						if flagged, ok = v.(bool); ok {
							if name == "roles" {
								parsed.Roles = flagged
							} else if name == "author" {
								parsed.Author = flagged
							} else if name == "status" {
								parsed.Status = flagged
							} else {
								log.Debug(fmt.Sprintf("unsupported shop option: [%s:%v]", name, flagged))
							}
						}
					}
				}
			}
		}
		j.Show = parsed
		return true
	}

	return false
}
func (j *JekyllConfig) reparsedShowOptionsAsDictionary() bool {

	if dictOpts, ok := j.Show.(map[interface{}]interface{}); ok {
		parsed := JekyllShow{
			Status: true,
			Author: true,
			Roles:  true,
		}
		for k, v := range dictOpts {
			var option string
			var flagged bool
			if option, ok = k.(string); ok {
				if flagged, ok = v.(bool); ok {
					if option == "author" {
						parsed.Author = flagged
					} else if option == "roles" {
						parsed.Roles = flagged
					} else if option == "status" {
						parsed.Status = flagged
					} else {
						log.Debug(fmt.Sprintf("unsupported: [%s:%v]", option, flagged))
					}
				}
			}
		}
		j.Show = parsed
		return true
	}

	return false
}

func (j *JekyllConfig) reparseShowOptions() {

	if j.reparsedShowOptionsAsSequenceDictionaries() {
		return
	}

	if j.reparsedShowOptionsAsDictionary() {
		return
	}

	panic(fmt.Errorf("unsupported show options stylo : %v", j.Show))

}

type HugoRenderer struct {
	Unsafe bool `yaml:"Unsafe"`
}

type HugoGoldmark struct {
	Renderer HugoRenderer `yaml:"renderer"`
}

type HugoMarkupStyle struct {
	Style string `yaml:"style"`
}
type HugoMarkup struct {
	Goldmark  HugoGoldmark    `yaml:"goldmark"`
	Highlight HugoMarkupStyle `yaml:"highlight"`
}
type HugoMenuItem struct {
	Identifier string `yaml:"identifier"`
	Name       string `yaml:"name"`
	Url        string `yaml:"url"`
	Weight     int    `yaml:"weight"`
}

type HugoOutputFormat struct {
	BaseName  string `yaml:"baseName"`
	MediaType string `yaml:"mediaType"`
}

type HugoConfig struct {
	LanguageCode           string                      `yaml:"languageCode"`
	Title                  string                      `yaml:"title"`
	Copyright              string                      `yaml:"copyright"`
	AssetDir               string                      `yaml:"assetDir"`
	PluralizeListTitles    bool                        `yaml:"pluralizelisttitles"`
	EnableGitInfo          bool                        `yaml:"enableGitInfo"`
	Markup                 HugoMarkup                  `yaml:"markup"`
	Params                 map[string]interface{}      `yaml:"params"`
	SectionPagesMenu       string                      `yaml:"sectionPagesMenu"`
	Menu                   map[string][]HugoMenuItem   `yaml:"menu"`
	OutputFormats          map[string]HugoOutputFormat `yaml:"outputFormats"`
	Outputs                map[string][]string         `yaml:"outputs"`
	Module                 HugoModule                  `yaml:"module"`
	EnableInlineShortcodes bool                        `yaml:"enableInlineShortcodes"`
	Frontmatter            HugoFrontmatter             `yaml:"frontmatter"`
}

type HugoImport struct {
	Path     string  `yaml:"path"`
	Disabled bool    `yaml:"disabled"`
	Mounts   []Mount `yaml:"mounts,omitempty"`
}

type Mount struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}

type HugoModule struct {
	Imports []HugoImport `yaml:"imports"`
}

type HugoFrontmatter struct {
	Lastmod []string `yaml:"lastmod"`
}

type Roles struct {
	Label   string   `yaml:"label"`
	All     string   `yaml:"all"`
	Options []string `yaml:options`
}

func ReadJekyllConfig(path string) (*JekyllConfig, error) {

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &JekyllConfig{}
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}
	config.reparseShowOptions()
	return config, nil
}

func WriteHugoConfig(path string, config *HugoConfig) error {
	b, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, b, 0755)
}

// LogoImgRe This regex for grabbing the logo file name from the config file
var LogoImgRe = regexp.MustCompile(`(?m).*\/(.*\.(?:png|jpg|jpeg|gif|svg))$`)

func convertLogoPath(logoPrefix string, logoPath string) string {
	strings := LogoImgRe.FindStringSubmatch(logoPath)
	if strings == nil || len(strings) != 2 {
		fmt.Println("Searching", colors.Labels.Warning(logoPath), "for", colors.Labels.Wanted("Logo"), colors.Labels.Warning("unsuccessful"))
		return ""
	}
	logo := strings[1]
	fmt.Println("Found logo", colors.Labels.Normal(logo), "in", colors.Labels.Info(logoPath))
	// Don't want to hardcode the file prefix, but don't know where to put it
	return fmt.Sprintf("%s%s", logoPrefix, logo)
}

func ConvertConfig(jekyllConfig *JekyllConfig, logoPrefix string, additionalParams map[string]interface{}) *HugoConfig {
	var mainMenu []HugoMenuItem
	for idx, item := range jekyllConfig.Sections {
		mainMenu = append(mainMenu, HugoMenuItem{
			Identifier: item.Collection,
			Name:       item.Title,
			Url:        item.Url,
			Weight:     idx + 1,
		})
	}

	hugoConfig := &HugoConfig{
		LanguageCode:  "en-us",
		Title:         jekyllConfig.Name,
		EnableGitInfo: false,
		Markup: HugoMarkup{
			Goldmark: HugoGoldmark{
				Renderer: HugoRenderer{
					Unsafe: true,
				},
			},
			Highlight: HugoMarkupStyle{
				Style: config.Flags.SyntaxStyle,
			},
		},

		Menu: map[string][]HugoMenuItem{
			"Main": mainMenu,
		},
		OutputFormats: map[string]HugoOutputFormat{
			"MenuIndex": {BaseName: "menu", MediaType: "application/json"},
			"SearchMap": {BaseName: "searchmap", MediaType: "application/json"},
		},
		Outputs: map[string][]string{
			"home": {
				"HTML",
				"RSS",
				"MenuIndex",
				"SearchMap",
			},
		},
		Module: HugoModule{
			Imports: []HugoImport{
				{
					Path:     "github.com/spandigital/presidium-theme-website",
					Disabled: false,
				},
				//{ This theme doesn't work atm
				//	Path:     "github.com/spandigital/presidium-theme-pdf",
				//	Disabled: true,
				//},
			},
		},
	}

	if len(config.Flags.BrandTheme) > 0 {
		mod := getImportModule(config.Flags.BrandTheme)
		hugoConfig.Module.Imports = append([]HugoImport{mod}, hugoConfig.Module.Imports...)
	}

	if additionalParams != nil {
		hugoConfig.Params = additionalParams
	}
	hugoConfig.Params["audience"] = jekyllConfig.Audience
	hugoConfig.Params["scope"] = jekyllConfig.Scope
	hugoConfig.Params["appleScope"] = jekyllConfig.AppleScope
	hugoConfig.Params["description"] = jekyllConfig.Description
	hugoConfig.Params["logo"] = convertLogoPath(logoPrefix, jekyllConfig.Logo)
	hugoConfig.Copyright = jekyllConfig.Footer
	hugoConfig.Params["show"] = jekyllConfig.Show
	hugoConfig.Params["roles"] = jekyllConfig.Roles
	hugoConfig.AssetDir = "static"

	hugoConfig.EnableInlineShortcodes = true

	hugoConfig.Frontmatter.Lastmod = []string{"lastmod", ":fileModTime", ":default"}

	return hugoConfig
}

func getImportModule(path string) HugoImport {
	return HugoImport{
		Path:     path,
		Disabled: false,
		Mounts: []Mount{
			{
				Source: "assets",
				Target: "assets",
			},
			{
				Source: "static",
				Target: "static",
			},
		},
	}
}
