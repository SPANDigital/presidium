package markdown

import (
	"fmt"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/conversion/colors"
	. "github.com/SPANDigital/presidium-hugo/pkg/filesystem"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

func init() {
	FS = afero.NewMemMapFs()
	FSUtil = &afero.Afero{Fs:FS}
}

var _ = BeforeSuite(func() {
	colors.Setup()
})

var (
	input = `
 As mentioned in the [Handbook Introduction]({{% baseurl %}}/#contribution), the Handbook is a [Knowledge Management](# 'presidium-tooltip') resource that is continually updated. This requires the active participation of both the consumers and creators of this information. This section outlines the [contribution process]({{% baseurl %}}/handbook-contribution/contribution-process/) which defines how feedback and content changes are controlled. It also covers the basics of the [content development]({{% baseurl %}}/handbook-contribution/content-development/) procedure to help employees get started contributing to the Handbook. In addition, several Handbook specific content [style guides]({{% baseurl %}}/handbook-contribution/style-guides/) are included to improve consistency throughout the document. 
`
	expected = ` ---
---

 As mentioned in the [Handbook Introduction]({{% baseurl %}}/#contribution), the Handbook is a {{< tooltip "Knowledge Management" >}} resource that is continually updated. This requires the active participation of both the consumers and creators of this information. This section outlines the [contribution process]({{% baseurl %}}/handbook-contribution/contribution-process/) which defines how feedback and content changes are controlled. It also covers the basics of the [content development]({{% baseurl %}}/handbook-contribution/content-development/) procedure to help employees get started contributing to the Handbook. In addition, several Handbook specific content [style guides]({{% baseurl %}}/handbook-contribution/style-guides/) are included to improve consistency throughout the document. 
 `
)
var _ = Describe("Markdown", func() {
	Describe("manipulate", func() {
		Context("ManipulateMarkdown", func() {
			path := "/home/testuser/testdata"
			filepath := fmt.Sprintf("%s/file.md", path)
			FS.MkdirAll(path, 0755)
			FSUtil.WriteFile(filepath, []byte(input), 0644)
			It("Should correctly translate tooltips", func() {
				replaceTooltips(filepath)
				content, err := FSUtil.ReadFile(filepath)
				if err != nil {
					Fail("Expected file to exist")
				}
				Expect(content, []byte(expected))
			})
		})
	})
})
