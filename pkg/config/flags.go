package config

type flags struct {
	EnableColor                bool
	SourceRepoDir              string
	DestinationRepoDir         string
	StagingDir                 string
	WeightBasedOnFilename      bool
	SlugBasedOnFileName        bool
	UrlBasedOnFilename         bool
	CommonmarkAttributes       bool
	ReplaceBaseUrl             bool
	RemoveTargetBlank          bool
	FixImages                  bool
	FixImagesWithAttributes    bool
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
	LogoPrefix				   string
	FixTables                  bool
}

var Flags flags
