package generator

import (
	"fmt"
	"github.com/Masterminds/goutils"
	model "github.com/SPANDigital/presidium-hugo/pkg/domain/model/generator"
	"github.com/SPANDigital/presidium-hugo/pkg/filesystem"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGeneratorImpl(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Default Site SiteGenerator Suite")
}

var (
	fs = filesystem.New()
)

var _ = Describe("Site generation behaviour:", func() {

	var workDir string
	var t model.InitialSiteTarget
	var g SiteGenerator

	BeforeSuite(func() {
		if tempDir, err := ioutil.TempDir("", "presidium-site-generator-test-*"); err != nil {
			panic(err) // make sure it fails!
		} else {
			workDir = tempDir
		}
	})

	AfterSuite(func() { _ = os.RemoveAll(workDir) })

	BeforeEach(func() {
		g = New()
		t = model.InitialSiteTarget{
			SiteTargetDirectory: filepath.Join(workDir, "testSite"),
			SiteName:            "Test Site",
			SiteTitle:           "A Test site",
			BrandingModelUrl:    "",
			Theme:               model.PresidiumTheme,
			Template:            model.SpanTemplate,
			WhenSiteExists:      model.AbortWhenTargetSiteExists,
		}
	})

	AfterEach(func() {
		if fs.DirExists(t.SiteTargetDirectory) {
			_ = fs.DeleteDir(t.SiteTargetDirectory)
		}
	})

	It("should fail when site directory exists already.", func() {
		mustMakeDir(t.SiteTargetDirectory)
		siteGenerationErr := g.Run(t)
		Expect(siteGenerationErr).Should(HaveOccurred())
	})

	It("should overwrite the existing site if so configured.", func() {
		pathId, _ := goutils.RandomNumeric(6)
		up := func(s string) string { return strings.Replace(s, "*", pathId, 1) } // making a unique path here
		removablePats := mustMakeTree("will be removed", []string{
			up("content-*/introduction/_index.md"),
			up("content-*/introduction/welcome-here.md"),
			up("content-*/polices/_index.md"),
			up("content-*/facilities/"),
		}, t.SiteTargetDirectory)
		t.WhenSiteExists = model.ReplaceTargetSiteIfExists
		siteGenerationErr := g.Run(t)
		Expect(siteGenerationErr).ShouldNot(HaveOccurred())
		Expect(remainingOf(removablePats, t.SiteTargetDirectory)).Should(BeEmpty())
	})

	Context("after running the generator, the target site", func() {

		BeforeEach(func() {
			siteGenerationErr := g.Run(t)
			Expect(siteGenerationErr).ShouldNot(HaveOccurred())
			_, local := filepath.Split(t.SiteTargetDirectory)
			fmt.Printf("site generated OK: \"%s\" [%s]\n", t.SiteName, local)
		})

		AfterEach(func() {
			errAfterSiteRemoved := fs.EmptyDir(t.SiteTargetDirectory)
			Expect(errAfterSiteRemoved).ShouldNot(HaveOccurred())
		})

		It("should have a static assets folder", func() {
			mustHaveDir(t.AssetsDir())
		})

		Context("generated content", func() {
			var generatedContent []string
			BeforeEach(func() { generatedContent = listSiteContent(t) })
			It("should exists", func() {
				Expect(generatedContent).ShouldNot(BeEmpty())
			})
		})

	})
})

func listSiteContent(t model.InitialSiteTarget) []string {
	content := make([]string, 0)
	contentDir := mustHaveDir(t.ContentDir())
	contentListingErr := filepath.WalkDir(contentDir, func(path string, d os.DirEntry, err error) error {
		if !d.IsDir() {
			local := strings.TrimPrefix(contentDir, path)
			content = append(content, local)
		}
		return nil
	})
	Expect(contentListingErr).ShouldNot(HaveOccurred())
	return content
}

func mustHaveDir(path string) string {
	pathInfo, pathErr := os.Stat(path)
	Expect(pathErr).ShouldNot(HaveOccurred())
	Expect(pathInfo.IsDir()).Should(BeTrue())
	return path
}

func remainingOf(tree []string, parent string) []string {
	found := make([]string, 0)
	for _, path := range tree {
		if _, err := os.Stat(filepath.Join(parent, path)); err == nil {
			found = append(found, path)
		}
	}
	return found
}

func mustMakeDir(path string) {
	if fs.DirExists(path) {
		return
	} else {
		err := os.MkdirAll(path, 0755)
		Expect(err).ShouldNot(HaveOccurred())
	}
}

func mustMakeTree(purpose string, tree []string, parent string) []string {
	println("Generating tree:", purpose)
	for i, local := range tree {
		wantDir := strings.HasSuffix(local, "/")
		if wantDir {
			mustMakeDir(filepath.Join(parent, local))
		} else {
			mustMakeFile(filepath.Join(parent, local))
		}
		fmt.Printf("%d: %s\n", i+1, local)
	}
	return tree
}

func mustMakeFile(path string) string {
	dir, _ := filepath.Split(path)
	mustMakeDir(dir)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0755)
	Expect(err).ShouldNot(HaveOccurred())
	_, _ = file.WriteString("dummy text!")
	_ = file.Close()
	return path
}
