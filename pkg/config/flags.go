package config

type flags struct {
	EnableColor                bool
	SourceRepoDir              string
	BrandTheme                 string
	SyntaxStyle                string
	DestinationRepoDir         string
	StagingDir                 string
	WeightBasedOnFilename      bool
	SlugBasedOnFileName        bool
	UrlBasedOnFilename         bool
	CommonmarkAttributes       bool
	ReplaceBaseUrl             bool
	RemoveTargetBlank          bool
	FixImages                  bool
	FixHtmlImages              bool
	EraseMarkdownWithNoContent bool
	RemoveRawTags              bool
	ReplaceRoot                string
	ReplaceCallOuts            bool
	ReplaceTooltips            bool
	AddMissingTitles           bool
	ReplaceIfVariables         bool
	CopyMediaToStatic          bool
	ConvertConfigYml           bool
	ReplaceBaseUrlWithSpaces   bool
	ReplaceComments            bool
	LogoPrefix                 string
	FixTables                  bool
	GenerateHugoModule         bool
	HugoModuleName             string
}

var Flags flags
