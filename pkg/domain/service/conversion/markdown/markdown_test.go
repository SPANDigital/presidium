package markdown

import (
	"fmt"
	"github.com/Masterminds/goutils"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/conversion/colors"
	"github.com/SPANDigital/presidium-hugo/pkg/filesystem"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestMarkdown(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Markdown Suite")
}

var _ = Describe("Processing markdown content", func() {
	var workDir string
	filesystem.FS = afero.NewMemMapFs()
	filesystem.FSUtil = &afero.Afero{Fs: filesystem.FS}

	BeforeSuite(func() {
		colors.Setup()
		var workDirErr error
		workDir, workDirErr = ioutil.TempDir("", "markdown-processing-test")
		Expect(workDirErr).ShouldNot(HaveOccurred())
	})

	AfterSuite(func() { _ = os.RemoveAll(workDir) })

	When("Replacing callouts", func() {
		var markdownText = `<div class="presidium-warning">
    		<span>Do Not Copy Working Documents from Other Projects</span>
    		<p>Working Documents are continually being improved, so it is critical that each new project uses the latest versions.</p>

      		Working Document Templates are stored in the following Google Drive directories:
  			<ul style="list-style: none; margin-left: -20px;">
    			<li><a href="https://drive.google.com/drive/folders/0B2ynAoQvesi5Zk9ZRFhnYWRrLU0" target="_blank">&#x2605; SPAN Services Templates</a></li>
    			<li><a href="https://drive.google.com/drive/folders/0B32NHmira-_eWnNGQWFKUHh5Mlk" target="_blank">&#x2605; SPAN Templates</a></li>
  			</ul>
			</div>
		`

		It("Should replace it as expected", func() {
			markdownFile := mustHaveMarkdownInputFile(workDir, markdownText)
			err := replaceCallOuts(markdownFile)
			Expect(err).ShouldNot(HaveOccurred())
			actual := contentOf(markdownFile)
			Expect(actual).Should(ContainSubstring(`{{< callout level="warning" title="Do Not Copy Working Documents from Other Projects">}}`))
			Expect(actual).Should(Not(MatchRegexp(EmptyLineRe.String())))
		})
	})

	When("Replacing tooltips", func() {

		var markdownText = "As mentioned in the [Handbook Introduction]({{% baseurl %}}/#contribution), the Handbook " +
			"is a [Knowledge Management](#top-context-menu 'presidium-tooltip') resource that is continually updated. This requires " +
			"the active participation of both the consumers and creators of this information. This section outlines " +
			"the [contribution process]({{% baseurl %}}/handbook-contribution/contribution-process/) " +
			"which defines how feedback and content changes are controlled. It also covers the basics " +
			"of the [content development]({{% baseurl %}}/handbook-contribution/content-development/) " +
			"procedure to help employees get started contributing to the Handbook. " +
			"In addition, several Handbook specific content " +
			"[style guides]({{% baseurl %}}/handbook-contribution/style-guides/) " +
			"are included to improve consistency throughout the document."

		It("Should replace it as expected", func() {
			markdownFile := mustHaveMarkdownInputFile(workDir, markdownText)
			err := replaceTooltips(markdownFile)
			Expect(err).ShouldNot(HaveOccurred())
			actual := contentOf(markdownFile)
			Expect(actual).Should(ContainSubstring("{{< tooltip \"Knowledge Management\" \"top-context-menu\" >}}"))
		})
	})
})

func mustHaveDir(path string) {
	pathInfo, pathErr := os.Stat(path)
	if pathErr == nil {
		Expect(pathInfo.IsDir()).Should(BeTrue())
	} else if os.IsNotExist(pathErr) {
		pathErr = os.MkdirAll(path, os.ModePerm)
	}
	Expect(pathErr).ShouldNot(HaveOccurred())
}

func mustHaveMarkdownInputFile(dir string, content string) string {
	fileId, fileIdErr := goutils.RandomAlphaNumeric(4)
	Expect(fileIdErr).ShouldNot(HaveOccurred())
	mustHaveDir(dir)
	name := fmt.Sprintf("contentOf-%s.md", fileId)
	path := filepath.Join(dir, name)
	file, err := filesystem.FS.OpenFile(path, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	Expect(err).ShouldNot(HaveOccurred())
	_, err = file.WriteString(content)
	file.Close()
	Expect(err).ShouldNot(HaveOccurred())
	return path
}

func contentOf(path string) string {
	bytes, err := filesystem.FSUtil.ReadFile(path)
	Expect(err).ShouldNot(HaveOccurred())
	return string(bytes)
}
