package configtranslation

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type JekyllShow struct {
	Status bool `yaml:"status""`
	Author bool `yaml:"author""`
}

type JekyllExternal struct {
	AuthorsUrl string `yaml:"authors-url:`
}

type JekyllSectionItem struct {
	Title string `yaml:"title""`
	Url   string `yaml:"url""`
}

type JekyllConfig struct {
	Name     string              `yaml:"name"`
	Baseurl  string              `yaml:"baseurl"`
	Footer   string              `yaml:"footer"`
	Logo     string              `yaml:"logo"`
	Show     JekyllShow          `yaml:"show"`
	External JekyllExternal      `yaml:"external"`
	Sections []JekyllSectionItem `yaml::sections"`
}

type HugoRenderer struct {
	Unsafe bool `yaml:"Unsafe"`
}

type HugoGoldmark struct {
	Renderer HugoRenderer `yaml:"renderer"`
}

type HugoMarkup struct {
	Goldmark HugoGoldmark `yanl:"goldmark"`
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
	BaseURL             string                      `yaml:"baseURL"`
	LanguageCode        string                      `yaml:"languageCode"`
	Title               string                      `yaml:"title"`
	PluralizeListTitles bool                        `yaml:"pluralizelisttitles"`
	Markup              HugoMarkup                  `yaml:"markup"`
	Params              map[string]string           `yaml:"params"`
	SectionPagesMenu    string                      `yaml:"sectionPagesMenu"`
	Menu                map[string][]HugoMenuItem   `yaml:"menu"`
	OutputFormats       map[string]HugoOutputFormat `yaml:"outputFormats"`
	Outputs             map[string][]string         `yaml:"outputs"`
	Module              HugoModule                  `yaml:"module"`
}

type HugoImport struct {
	Path     string `yaml:"path"`
	Disabled bool   `yaml:"disabled"`
}

type HugoModule struct {
	Imports []HugoImport `yaml:"imports"`
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
	return config, nil
}

func WriteHugoConfig(path string, config *HugoConfig) error {
	b, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, b, 0755)
}

func ConvertConfig(config *JekyllConfig, additionalParams map[string]string) *HugoConfig {

	mainMenu := []HugoMenuItem{}
	for idx, item := range config.Sections {
		mainMenu = append(mainMenu, HugoMenuItem{
			Identifier: item.Title,
			Name:       item.Title,
			Url:        item.Url,
			Weight:     idx + 1,
		})
	}

	hugoConfig := &HugoConfig{
		BaseURL:      config.Baseurl,
		LanguageCode: "en-us",
		Title:        config.Name,
		Markup: HugoMarkup{
			Goldmark: HugoGoldmark{
				Renderer: HugoRenderer{
					Unsafe: true,
				},
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
				{
					Path:     "github.com/spandigital/presidium-theme-pdf",
					Disabled: true,
				},
			},
		},
	}
	if additionalParams != nil {
		hugoConfig.Params = additionalParams
	}
	return hugoConfig
}
