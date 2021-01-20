package cmd

import (
	"fmt"
	"github.com/SPANDigital/presidium-hugo/internal/configtranslation"
	"path/filepath"
	"github.com/SPANDigital/presidium-hugo/internal/colors"
	"github.com/SPANDigital/presidium-hugo/internal/markdown"
	"github.com/SPANDigital/presidium-hugo/internal/filewalking"
	"github.com/SPANDigital/presidium-hugo/internal/paths"
	"github.com/SPANDigital/presidium-hugo/internal/settings"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/otiai10/copy"
	"github.com/spf13/cobra"
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert Jekyll to Hugo content",
	Long: `Convert Jekyll to Hugo content`,
	Run: func(cmd *cobra.Command, args []string) {

		sourceRepoDir := viper.GetString("sourceRepoDir")
		destinationRepoDir := viper.GetString("destDir")

		if destinationRepoDir == "" {
			destinationRepoDir, _ = os.Getwd()
		}

		if (sourceRepoDir != "") {

			stagingDir := settings.Flags.StagingDir

			sourceRepoContentDir := filepath.Join(sourceRepoDir, "content")
			sourceRepoStaticDir := filepath.Join(sourceRepoDir, "media")
			sourceRepoConfigYml := filepath.Join(sourceRepoDir, "_config.yml")


			stagingContentDir := filepath.Join(stagingDir, "content")

			destinationContentDir := filepath.Join(destinationRepoDir, "content")
			destinationStaticDir := filepath.Join(destinationRepoDir, "static")
			destinationMediaDir := filepath.Join(destinationStaticDir, "media")
			desinationConfigYml := filepath.Join(destinationRepoDir, "config.yml")

			fmt.Println()
			fmt.Println(colors.Labels.Underline("Source repo dir:"), colors.Labels.Info(sourceRepoDir))
			fmt.Println(colors.Labels.Underline("Destination repo dir:"), colors.Labels.Info(destinationRepoDir))
			fmt.Println(colors.Labels.Underline("Staging dir:"), colors.Labels.Info(stagingDir))
			fmt.Println()

			fmt.Println("Creating staging content directory:", colors.Labels.Info(stagingContentDir))
			err := os.MkdirAll(stagingContentDir, 0755)
			if err != nil {
				log.Fatal("Could not create staging directory", err)
			}
			fmt.Println("Emptying contents of staging directory:", colors.Labels.Info(stagingContentDir))
			err = RemoveContents(stagingDir)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Copying all contents from source content directory:", colors.Labels.Info(sourceRepoContentDir), " -> ", colors.Labels.Wanted(stagingContentDir))
			err = copy.Copy(sourceRepoContentDir, stagingContentDir, copy.Options{Skip: func(src string) (bool, error) {
				_, file := filepath.Split(src)
				if strings.HasPrefix(file, ".") {
					return true, nil
				}
				return false, nil
			}})
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println()
			fmt.Println(colors.Labels.Underline("With:"), colors.Labels.Wanted(stagingContentDir))
			fmt.Println()

			fmt.Println("Checking for directories to rename")
			err = filewalking.CheckForDirRename(stagingContentDir)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Gathering resources")
			err = paths.GatherResources(stagingContentDir)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Checking for directory index")
			err = filewalking.CheckForDirIndex(stagingContentDir)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Checking for missing index titles")
			err = filewalking.CheckIndexForTitles(stagingContentDir)
			if err != nil {
				log.Fatal(err)
			}

			os.MkdirAll(destinationContentDir, 0755)

			fmt.Println(colors.Labels.Underline("Emptying contents of:"), colors.Labels.Info(destinationContentDir))
			err = RemoveContents(destinationContentDir)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(colors.Labels.Underline("Copying content from"), colors.Labels.Info(stagingContentDir), " -> ", colors.Labels.Wanted(destinationContentDir))
			err = copy.Copy(stagingContentDir, destinationContentDir)
			if err != nil {
				log.Fatal(err)
			}

			if viper.GetBool("copyMediaToStatic") {

				err := os.MkdirAll(destinationStaticDir, 0755)
				if err != nil {
					log.Fatal("Could not create staging directory", err)
				}

				fmt.Println(colors.Labels.Underline("Emptying contents of:"), colors.Labels.Info(destinationStaticDir))
				err = RemoveContents(destinationStaticDir)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println(colors.Labels.Underline("Copying content from"), colors.Labels.Info(sourceRepoStaticDir), " -> ", colors.Labels.Wanted(destinationMediaDir))
				err = copy.Copy(sourceRepoStaticDir, destinationMediaDir)
				if err != nil {
					log.Fatal(err)
				}
			}

			if viper.GetBool("convertConfigYml") {
				jekyllConfig, err := configtranslation.ReadJekyllConfig(sourceRepoConfigYml)
				if err != nil {
					log.Fatal(err)
				}
				hugoConfig := configtranslation.ConvertConfig(jekyllConfig, map[string]string{})
				err = configtranslation.WriteHugoConfig(desinationConfigYml, hugoConfig)
				if err != nil {
					log.Fatal(err)
				}
			}

			fmt.Println(colors.Labels.Underline("Removing"), colors.Labels.Info(stagingDir))
			err = os.RemoveAll(stagingDir)
 			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

func currentWorkingDirectory() string {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return path
}

func init() {
	rootCmd.AddCommand(convertCmd)

	// Here you will define your flags and configuration settings.
	stagingDir, err := ioutil.TempDir(os.TempDir(), "staging")
	if err != nil {
		log.Fatal("Could not create staging directory", err)
	}
	settings.Flags.StagingDir = stagingDir

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// convertCmd.PersistentFlags().String("foo", "", "A help for foo")
	pflags := convertCmd.PersistentFlags()
	pflags.BoolVarP(&settings.Flags.EnableColor, "enableColor", "c", true, "Enable colorful output")
	pflags.StringVarP(&settings.Flags.SourceRepoDir, "sourceRepoDir", "s", "", "Source directory")
	pflags.StringVarP(&settings.Flags.DestinationRepoDir, "destDir", "d",  currentWorkingDirectory(),"Destination directory")
	pflags.BoolVarP(&settings.Flags.WeightBasedOnFilename, "weightBasedOnFilename", "w", true, "Base front matter weight on filename")
	pflags.BoolVarP(&settings.Flags.SlugBasedOnFileName, "slugBasedOnFileName", "g", true, "Base front matter slug on filename")
	pflags.BoolVarP(&settings.Flags.UrlBasedOnFilename, "urlBasedOnFilename", "u", true, "Base front matter url on filename")
	pflags.BoolVarP(&settings.Flags.CommonmarkAttributes, "commonmarkAttributes", "m", false, "Convert to commonmark attribute format")
	pflags.BoolVarP(&settings.Flags.ReplaceBaseUrl, "replaceBaseUrl", "b", true, "Replace {{site.baseurl}} with {{ site.BaseURL }}")
	pflags.BoolVarP(&settings.Flags.ReplaceBaseUrlWithSpaces, "replaceBaseUrlWithSpaces", "j", true, "Replace {{ site.baseurl }} with {{site.BaseURL}}")
	pflags.BoolVarP(&settings.Flags.RemoveTargetBlank, "removeTargetBlank", "t", true, `Remove target="blank" variants`)
	pflags.BoolVarP(&settings.Flags.FixImages, "fixImages", "i", true, "Fix images in same path")
	pflags.BoolVarP(&settings.Flags.FixImagesWithAttributes, "fixImagesWithAttributes", "a",  true,"Replace images with attributes with shortcodes")
	pflags.BoolVarP(&settings.Flags.EraseMarkdownWithNoContent, "eraseMarkdownWithNoContent", "e", true, "Erase markdown files with no content")
	pflags.BoolVarP(&settings.Flags.RemoveRawTags, "removeRawTags", "R", true, "Remove {% raw %} tags")
	pflags.StringVarP(&settings.Flags.ReplaceRoot, "replaceRoot", "p", "", "Replace this path with root")
	pflags.BoolVarP(&settings.Flags.ReplaceCallOuts, "replaceCallOuts", "o", true, "Replace callout HTML with callout shortcodes")
	pflags.BoolVarP(&settings.Flags.ReplaceTooltips, "replaceTooltips", "T", true, "Replace tooltip HTML with callout shortcodes")
	pflags.BoolVarP(&settings.Flags.ReplaceIfVariables, "replaceIfVariables", "V", true, "Replace {% if site.variable =} with shortcodes")
	pflags.BoolVarP(&settings.Flags.ReplaceComments, "replaceComments", "", true, "Replace {% comment %}...{% endcomment %} with HTML comments")
	pflags.BoolVarP(&settings.Flags.CopyMediaToStatic, "copyMediaToStatic", "C", true, "Copy Jekyll media to Hugo static folder")
	pflags.BoolVarP(&settings.Flags.ConvertConfigYml, "convertConfigYml", "y", true, "Convert jekyll _config.yml to hugo config.yml")
	viper.BindPFlags(pflags)

	colors.Setup()
	markdown.SetupExcludes()

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// convertCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}



func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}



