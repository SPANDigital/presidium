package markdown

import (
	"github.com/SPANDigital/presidium-hugo/pkg/filesystem"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)


var _ = Describe("AddFrontMatter", func() {
	filesystem.FS = afero.NewMemMapFs()
	filesystem.FSUtil = &afero.Afero{Fs:  filesystem.FS}

	BeforeEach(func() {
		filesystem.FS.Remove("test.md")
	})

	When("adding front matter", func() {
		It("should contain front matter", func() {
			err := AddFrontMatter("test.md", FrontMatter{
				Title: "test",
				Author: "steve",
				Status: "draft",
				Slug: "test",
				URL:"/test",
				Github: "@test",
				Roles: "Developer",
				Weight: "1",
			})
			Expect(err).Should(BeNil())
			expected := "---\ntitle: test\nslug: test\nurl: /test\nweight: \"1\"\nauthor: steve\ngithub: '@test'\nstatus: draft\nroles: Developer\n---\n"
			content, err := afero.ReadFile(filesystem.FS, "test.md")
			Expect(err).Should(BeNil())
			Expect(string(content)).Should(Equal(expected))
		})
	})
})