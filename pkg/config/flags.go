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
	AddSlugAndUrl              bool
	CleanTarget                bool
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
	FixFigureCaptions          bool
}

var Flags flags
