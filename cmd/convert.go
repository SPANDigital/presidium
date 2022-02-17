package cmd

import (
	"github.com/SPANDigital/presidium-hugo/pkg/config"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/conversion"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/conversion/colors"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/conversion/markdown"
	"github.com/SPANDigital/presidium-hugo/pkg/filesystem"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert Jekyll to Hugo content",
	Run: func(cmd *cobra.Command, args []string) {

		c := conversion.New()

		source := viper.GetString("sourceRepoDir")
		destination := viper.GetString("destDir")
		moduleName := viper.GetString("hugoModuleName")

		c.SiteModuleName = moduleName

		err := c.Execute(source, destination)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {

	rootCmd.AddCommand(convertCmd)

	fs := filesystem.New()

	cwd, err := fs.GetWorkingDir()
	if err != nil {
		log.Fatal("could not get working dir")
	} else {
		cwd, err = fs.AbsolutePath(".")
		if err != nil {
			log.Fatal("something bad happened to get working dir")
		} else {
			cwd = "."
		}
	}

	stagingDir, err := ioutil.TempDir(os.TempDir(), "staging")
	if err != nil {
		log.Fatal("Could not create staging directory", err)
	}
	config.Flags.StagingDir = stagingDir

	pflags := convertCmd.PersistentFlags()
	pflags.BoolVarP(&config.Flags.EnableColor, "enableColor", "c", conversion.Defaults.EnableColorOutput, "Enable colorful output")
	pflags.StringVarP(&config.Flags.SourceRepoDir, "sourceRepoDir", "s", "", "Source directory")
	pflags.StringVarP(&config.Flags.DestinationRepoDir, "destDir", "d", cwd, "Destination directory")
	pflags.BoolVarP(&config.Flags.WeightBasedOnFilename, "weightBasedOnFilename", "w", conversion.Defaults.WeightBasedOnFileName, "Base front matter weight on filename")
	pflags.BoolVarP(&config.Flags.SlugBasedOnFileName, "slugBasedOnFileName", "g", conversion.Defaults.SlugBasedOnFileName, "Base front matter slug on filename")
	pflags.BoolVarP(&config.Flags.UrlBasedOnFilename, "urlBasedOnFilename", "u", conversion.Defaults.ReplaceBaseUrl, "Base front matter url on filename")
	pflags.BoolVarP(&config.Flags.CommonmarkAttributes, "commonmarkAttributes", "m", conversion.Defaults.CommonMarkdownAttributes, "Convert to commonmark attribute format")
	pflags.BoolVarP(&config.Flags.ReplaceBaseUrl, "replaceBaseUrl", "b", conversion.Defaults.ReplaceBaseUrl, "Replace {{site.baseurl}} with {{ site.BaseURL }}")
	pflags.BoolVarP(&config.Flags.ReplaceBaseUrlWithSpaces, "replaceBaseUrlWithSpaces", "j", conversion.Defaults.ReplaceBaseUrlWithSpaces, "Replace {{ site.baseurl }} with {{site.BaseURL}}")
	pflags.BoolVarP(&config.Flags.RemoveTargetBlank, "removeTargetBlank", "t", conversion.Defaults.RemoveTargetBlank, `Remove target="blank" variants`)
	pflags.BoolVarP(&config.Flags.FixImages, "fixImages", "i", conversion.Defaults.FixImages, "Fix images in same path")
	pflags.BoolVar(&config.Flags.FixHtmlImages, "fixHtmlImages", conversion.Defaults.FixHtmlImages, "Fix the source of html images")
	pflags.BoolVarP(&config.Flags.EraseMarkdownWithNoContent, "eraseMarkdownWithNoContent", "e", conversion.Defaults.EraseMarkdownWithNoContent, "Erase markdown files with no content")
	pflags.BoolVarP(&config.Flags.RemoveRawTags, "removeRawTags", "R", conversion.Defaults.RemoveRawTags, "Remove {% raw %} tags")
	pflags.StringVarP(&config.Flags.ReplaceRoot, "replaceRoot", "p", conversion.Defaults.ReplaceRootWith, "Replace this path with root")
	pflags.StringVarP(&config.Flags.LogoPrefix, "logoPrefix", "l", conversion.Defaults.LogoPrefix, "Use this as the prefix to locate logo")
	pflags.BoolVarP(&config.Flags.ReplaceCallOuts, "replaceCallOuts", "o", conversion.Defaults.ReplaceCallOuts, "Replace callout HTML with callout shortcodes")
	pflags.BoolVarP(&config.Flags.ReplaceTooltips, "replaceTooltips", "T", conversion.Defaults.ReplaceToolTips, "Replace tooltip HTML with callout shortcodes")
	pflags.BoolVarP(&config.Flags.ReplaceIfVariables, "replaceIfVariables", "V", conversion.Defaults.ReplaceIfVariables, "Replace {% if site.variable =} with with-param shortcodes")
	pflags.BoolVarP(&config.Flags.ReplaceComments, "replaceComments", "", conversion.Defaults.ReplaceComments, "Replace {% comment %}...{% endcomment %} with HTML comments")
	pflags.BoolVarP(&config.Flags.CopyMediaToStatic, "copyMediaToStatic", "C", conversion.Defaults.CopyMediaToStatic, "Copy Jekyll media to Hugo static folder")
	pflags.BoolVarP(&config.Flags.ConvertConfigYml, "convertConfigYml", "y", conversion.Defaults.ConvertJekyllConfig, "Convert jekyll _config.yml to hugo config.yml")
	pflags.BoolVarP(&config.Flags.FixTables, "addTableHeaders", "F", conversion.Defaults.FixTables, "Add empty table headers to tables without headers")
	pflags.BoolVarP(&config.Flags.GenerateHugoModule, "generateHugoModule", "M", conversion.Defaults.GenerateHugoModule, "Generate Hugo (Go) module")
	pflags.StringVarP(&config.Flags.HugoModuleName, "hugoModuleName", "N", "", "Use a specific hugo module name (instead of one derived from the site title)")

	err = viper.BindPFlags(pflags)
	if err != nil {
		log.Fatal(err)
	}

	colors.Setup()
	markdown.SetupExcludes()

}
