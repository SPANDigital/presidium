package markdown

import (
	"github.com/SPANDigital/presidium-hugo/pkg/filesystem"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
	"io/fs"
)

var testFrontMatter = `
---
title: Scope
author: john.r.harris@spandigital.com
github: virtualtraveler
status: published
slug: scope
url: introduction/scope
weight: 3
roles: Developer
---
`
var testMarkdown = `
# Header
* Test 1
* Test 2
* Test 2`

var _ = Describe("Parse", func() {
	filesystem.SetFileSystem(afero.NewMemMapFs())

	BeforeEach(func() {
		filesystem.AFS.Remove("test.md")
	})

	When("parsing markdown file", func() {
		It("should parse the front matter", func() {
			mockFile("test.md", []byte(testFrontMatter))
			md, err := Parse("test.md")
			Expect(err).Should(BeNil())
			Expect(md.FrontMatter).Should(Equal(FrontMatter{
				Title:  "Scope",
				Author: "john.r.harris@spandigital.com",
				Github: "virtualtraveler",
				Status: "published",
				Slug:   "scope",
				URL:    "introduction/scope",
				Weight: "3",
				Roles:  "Developer",
			}))
		})

		It("should parse the content", func() {
			mockFile("test.md", []byte(testFrontMatter+testMarkdown))
			md, err := Parse("test.md")
			Expect(err).Should(BeNil())
			Expect(md.Content).Should(Equal(testMarkdown))
		})
	})
})

func mockFile(name string, data []byte) {
	err := filesystem.AFS.WriteFile(name, data, fs.ModePerm)
	if err != nil {
		Fail("failed to create test file")
	}
}
