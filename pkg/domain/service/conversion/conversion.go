package conversion

import (
	"errors"
	"fmt"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/hugo"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/SPANDigital/presidium-hugo/pkg/configtranslation"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/conversion/fileactions"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/conversion/resources"
	"github.com/spf13/viper"

	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/conversion/colors"
	"github.com/SPANDigital/presidium-hugo/pkg/filesystem"
	"github.com/logrusorgru/aurora"

	"github.com/otiai10/copy"
)

var (
	ErrSourceIsNotDirectory        = errors.New("source is not directory")
	ErrConversionAlreadyRunning    = errors.New("conversion is already running")
	ErrSourceDirectoryDoesNotExist = errors.New("source directory does not exists")
	ErrUnableToCreateDestDirectory = errors.New("unable to create destination directory")
)

type Converter struct {
	ConvertJekyllConfig        bool
	EnableColorOutput          bool
	EraseMarkdownWithNoContent bool
	FixImages                  bool
	FixHtmlImages              bool
	LogoPrefix                 string
	RemoveRawTags              bool
	RemoveTargetBlank          bool
	ReplaceBaseUrl             bool
	ReplaceBaseUrlWithSpaces   bool
	ReplaceCallOuts            bool
	ReplaceIfVariables         bool
	ReplaceComments            bool
	ReplaceRootWith            string
	ReplaceToolTips            bool
	SlugBasedOnFileName        bool
	SilenceUserMessages        bool
	WeightBasedOnFileName      bool
	CommonMarkdownAttributes   bool
	CopyMediaToStatic          bool
	FixTables                  bool
	GenerateHugoModule         bool
	SiteModuleName             string

	// --- private state follows from here on: --
	stagingDir        string
	stagingContentDir string
	running           bool
	fs                filesystem.FileSystem

	// --- source directory structure (Jekyll)
	sourceDir               string
	sourceRepoContentDir    string
	sourceRepoStaticDir     string
	sourceRepoConfigYmlFile string

	// -- destination directory structure
	destinationRepoDir       string
	destinationContentDir    string
	destinationStaticDir     string
	destinationMediaDir      string
	destinationConfigYmlFile string
	hugoConfig               configtranslation.HugoConfig
}

func (c *Converter) IsRunning() bool {
	return c.running
}

func (c *Converter) GetSourceDir() string {
	return c.sourceDir
}

func (c *Converter) GetDestDir() string {
	return c.destinationRepoDir
}

func (c *Converter) initSourceDir(sourceDir string) error {

	if info, err := os.Stat(sourceDir); err != nil {
		return ErrSourceDirectoryDoesNotExist
	} else if !info.IsDir() {
		return ErrSourceIsNotDirectory
	}

	sourceDir, err := c.fs.AbsolutePath(sourceDir)
	if err != nil {
		return err
	}

	c.sourceDir = sourceDir
	c.sourceRepoContentDir = filepath.Join(c.sourceDir, "content")
	c.sourceRepoStaticDir = filepath.Join(c.sourceDir, "media")
	c.sourceRepoConfigYmlFile = filepath.Join(c.sourceDir, "_config.yml")

	return nil
}

func (c *Converter) initDestinationDir(destinationDir string) error {

	if destinationDir == "" {
		destinationDir = "."
	}

	destinationDir, err := c.fs.AbsolutePath(destinationDir)
	if err != nil {
		return err
	}

	if _, err := os.Stat(destinationDir); err == nil {
		if err = os.MkdirAll(destinationDir, 0666); err != nil {
			return ErrUnableToCreateDestDirectory
		}
	}

	c.destinationRepoDir = destinationDir
	c.destinationContentDir = filepath.Join(destinationDir, "content")
	c.destinationStaticDir = filepath.Join(destinationDir, "static")
	c.destinationMediaDir = filepath.Join(destinationDir, "media")
	c.destinationConfigYmlFile = filepath.Join(destinationDir, "config.yml")

	return nil
}

type terminalMessage struct {
	content      string
	label        string
	labelStyle   func(arg interface{}) aurora.Value
	contentStyle func(arg interface{}) aurora.Value
}

func message(label string, content string) *terminalMessage {
	return &terminalMessage{
		content:      content,
		label:        label,
		labelStyle:   colors.Labels.Underline,
		contentStyle: colors.Labels.Info,
	}
}

func infoMessage(content string) *terminalMessage {
	m := message("INFO: ", content)
	m.withContentStyle(colors.Labels.Normal)
	m.withLabelStyle(colors.Labels.Info)
	return m
}

func (m *terminalMessage) withLabelStyle(style colors.StyleLabel) *terminalMessage {
	m.labelStyle = style
	return m
}

func (m *terminalMessage) withContentStyle(style colors.StyleLabel) *terminalMessage {
	m.contentStyle = style
	return m
}

func (m *terminalMessage) getLabel() aurora.Value {
	return m.labelStyle(m.label)
}

func (m *terminalMessage) getContent() aurora.Value {
	return m.contentStyle(m.content)
}

func (c *Converter) messageUser(messages ...*terminalMessage) {

	if c.SilenceUserMessages || len(messages) == 0 {
		return
	}

	for _, m := range messages {
		fmt.Println(m.getLabel(), m.getContent())
	}
}

func (c *Converter) Execute(sourceDir string, destDir string) error {

	if c.running {
		return ErrConversionAlreadyRunning
	}
	if err := c.initStaging(); err != nil {
		return err
	}
	if err := c.initSourceDir(sourceDir); err != nil {
		return err
	}
	if err := c.initDestinationDir(destDir); err != nil {
		return err
	}

	c.messageUser(
		message("Source repo dir: ", c.sourceDir),
		message("Destination repo: ", c.destinationRepoDir),
		message("Staging directory: ", c.stagingDir),
		infoMessage(fmt.Sprintf("starting to convert %s", c.sourceDir)))

	// NB: ⚠️ order may be important here:
	c.prepareStaging()
	c.gatherResources()
	c.performFileActions()
	c.processStaticMedia()
	c.convertConfig()
	c.generateHugoModule()
	c.finalize()

	return nil
}

func (c *Converter) initStaging() error {

	workDir, err := ioutil.TempDir(os.TempDir(), "staging")
	if err != nil {
		log.Fatalf("could not create staging directrory [%s]: %s", workDir, err.Error())
		return err
	}

	c.stagingDir = workDir
	c.stagingContentDir = filepath.Join(c.stagingDir, "content")
	if err := c.fs.MakeDirs(c.stagingDir); err != nil {
		log.Fatalf("could not create staging content dir [%s]: %s", c.stagingContentDir, err.Error())
		return err
	}

	return nil
}

// New returns an uninitialized Converter configured with default values.
func New() *Converter {
	return &Converter{
		ConvertJekyllConfig:        true,
		EnableColorOutput:          true,
		EraseMarkdownWithNoContent: true,
		FixImages:                  true,
		FixHtmlImages:              true,
		LogoPrefix:                 "/images/",
		RemoveRawTags:              true,
		RemoveTargetBlank:          true,
		ReplaceBaseUrl:             true,
		ReplaceBaseUrlWithSpaces:   true,
		ReplaceCallOuts:            true,
		ReplaceComments:            true,
		ReplaceIfVariables:         true,
		ReplaceRootWith:            "",
		ReplaceToolTips:            true,
		SlugBasedOnFileName:        false,
		SilenceUserMessages:        false,
		WeightBasedOnFileName:      true,
		CommonMarkdownAttributes:   false,
		CopyMediaToStatic:          true,
		FixTables:                  true,
		fs:                         filesystem.New(),
		GenerateHugoModule:         true,
	}
}

var Defaults = New()

func (c *Converter) prepareStaging() {

	if err := c.fs.EmptyDir(c.stagingDir); err != nil {
		log.Fatalf("unable to clean out staging [%s]: %s", c.stagingDir, err.Error())
	}

	err := c.fs.CopyWithOptions(c.sourceRepoContentDir, c.stagingContentDir, copy.Options{
		Skip: func(src string) (bool, error) {
			_, file := filepath.Split(src)
			if strings.HasPrefix(file, ".") {
				return true, nil
			}
			return false, nil
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	c.messageUser(message(
		"With: ",
		c.stagingContentDir,
	).withLabelStyle(colors.Labels.Underline).withContentStyle(colors.Labels.Wanted))

}

func (c *Converter) gatherResources() {
	if err := resources.GatherResources(c.stagingContentDir); err != nil {
		log.Fatal(err)
	}
}

func (c *Converter) performFileActions() {
	c.messageUser(infoMessage("check directories for rename"))
	if err := fileactions.RemoveUnderscoreDirPrefix(c.stagingContentDir); err != nil {
		log.Fatal(err)
	}

	c.messageUser(infoMessage("check directories for indexes"))
	if err := fileactions.CheckForDirIndex(c.stagingDir, c.stagingContentDir); err != nil {
		log.Fatal(err)
	}

	c.messageUser(infoMessage("check directories for titles"))
	if err := fileactions.CheckForTitles(c.stagingContentDir); err != nil {
		log.Fatal(err)
	}

	c.messageUser(infoMessage("add front matter"))
	if err := fileactions.AddFrontMatter(c.stagingDir, c.stagingContentDir); err != nil {
		log.Fatal(err)
	}

	c.messageUser(infoMessage("unslugify remaining directories and articles"))
	if err := fileactions.RemoveWeightIndicatorsFromFilePaths(c.stagingContentDir); err != nil {
		log.Fatal(err)
	}

	c.messageUser(infoMessage("preparing to copy content over"))
	if err := c.fs.MakeDirs(c.destinationContentDir); err != nil {
		log.Fatal(err)
	}

	c.messageUser(infoMessage(fmt.Sprintf("copying: %s -> %s", c.stagingDir, c.destinationContentDir)))
	if err := copy.Copy(c.stagingContentDir, c.destinationContentDir); err != nil {
		log.Fatal(err)
	}
}

func (c *Converter) processStaticMedia() {

	if !viper.GetBool("copyMediaToStatic") {
		return
	}

	if !c.fs.DirExists(c.sourceRepoStaticDir) {
		return
	}

	_ = c.fs.MakeDirs(c.destinationStaticDir)
	_ = c.fs.CopyDir(c.sourceRepoStaticDir, c.destinationStaticDir)

}

func (c *Converter) convertConfig() {

	if !viper.GetBool("convertConfigYml") {
		return
	}

	jekylConfig, _ := configtranslation.ReadJekyllConfig(c.sourceRepoConfigYmlFile)
	hugoConfig := configtranslation.ConvertConfig(jekylConfig, viper.GetString("logoPrefix"), map[string]interface{}{})
	err := configtranslation.WriteHugoConfig(c.destinationConfigYmlFile, hugoConfig)
	if err != nil {
		log.Fatal(err)
	}

	c.hugoConfig = *hugoConfig
}

func (c *Converter) finalize() {

	_ = c.fs.DeleteDir(c.stagingDir)

	copyOver("package.json", c.destinationRepoDir)
	copyOver(".gitignore", c.destinationRepoDir)
	copyOver("package-lock.json", c.destinationRepoDir)

	c.messageUser(infoMessage("Completed").withContentStyle(colors.Labels.Wanted))
}

func (c *Converter) generateHugoModule() {

	if !c.GenerateHugoModule {
		return
	}

	c.messageUser(infoMessage("Adding Hugo GO module to site").withContentStyle(colors.Labels.Wanted))
	hugo.New().Execute("--source", c.stagingDir, "mod", "init", c.moduleName())
	srcModFile := filepath.Join(c.stagingDir, "go.mod")
	dstModFile := filepath.Join(c.destinationRepoDir, "go.mod")
	_ = c.fs.Copy(srcModFile, dstModFile, fs.ModePerm)
	c.messageUser(infoMessage("Copied over hugo mod file"))

}

func (c *Converter) moduleName() string {
	siteModuleName := c.SiteModuleName
	if len(siteModuleName) == 0 {
		siteModuleName = filepath.Base(c.sourceDir)
	}
	return siteModuleName
}

// copyOver ensures that these files exists on at the destinationDir.
//
// This function will create such a file it does not exist.
func copyOver(file string, destinationDir string) {
	from, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer from.Close()

	to, err := os.OpenFile(filepath.Join(destinationDir, file), os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatal(err)
	}
	defer to.Close()
	_, err = io.Copy(to, from)
	if err != nil {
		log.Fatal(err)
	}
}
