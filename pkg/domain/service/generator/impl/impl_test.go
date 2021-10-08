package impl

import (
	"fmt"
	"fs"
	"github.com/Masterminds/goutils"
	model "github.com/SPANDigital/presidium-hugo/pkg/domain/model/generator"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/generator"
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

var _ = Describe("Site Generation:", func() {

	var workDir string
	var target model.InitialSiteTarget
	var g generator.SiteGenerator

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
		target = model.InitialSiteTarget{
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
		if fs.DirExists(target.SiteTargetDirectory) {
			_ = fs.DeleteDir(target.SiteTargetDirectory)
		}
	})

	It("should fail when site directory exists already.", func() {
		makeDir(target.SiteTargetDirectory)
		err := g.Run(target)
		Expect(err).Should(HaveOccurred())
	})

	It("should overwrite the existing site if so configured.", func() {
		pathId, _ := goutils.RandomNumeric(6)
		up := func(s string) string { return strings.Replace(s, "*", pathId, 1) } // making a unique path here
		tree := makeTree("will be removed", []string{
			up("content-*/introduction/_index.md"),
			up("content-*/introduction/welcome-here.md"),
			up("content-*/polices/_index.md"),
			up("content-*/facilities/"),
		}, target.SiteTargetDirectory)
		target.WhenSiteExists = model.ReplaceTargetSiteIfExists
		err := g.Run(target)
		Expect(err).ShouldNot(HaveOccurred())
		found := findExisting(tree, target.SiteTargetDirectory)
		Expect(found).Should(BeEmpty())
	})

	Context("after running the generator, the target site", func() {

		BeforeEach(func() {
			err := g.Run(target)
			Expect(err).ShouldNot(HaveOccurred())
			_, local := filepath.Split(target.SiteTargetDirectory)
			fmt.Printf("site generated OK: \"%s\" [%s]\n", target.SiteName, local)
		})

		AfterEach(func() {
			err := fs.EmptyDir(target.SiteTargetDirectory)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should have a static \"assets\" folder", func() {
			staticFolder := filepath.Join(target.SiteTargetDirectory, "static")
			info, err := os.Stat(staticFolder)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(info.IsDir()).Should(BeTrue())
		})

		It("template should have been applied to the site", func() {

		})
	})

})

func listSiteGenerated(target model.InitialSiteTarget) []string {
	listing := make([]string, 0)
	err := filepath.WalkDir(target.SiteTargetDirectory, func(path string, d fs.DirEntry, err error) error {
		return nil
	})
	Expect(err).ShouldNot(HaveOccurred())
	return listing
}

func findExisting(tree []string, parent string) []string {
	found := make([]string, 0)
	for _, path := range tree {
		if _, err := os.Stat(filepath.Join(parent, path)); err == nil {
			found = append(found, path)
		}
	}
	return found
}

func makeDir(path string) {
	if fs.DirExists(path) {
		return
	} else {
		err := os.MkdirAll(path, 0755)
		Expect(err).ShouldNot(HaveOccurred())
	}
}

func makeTree(purpose string, tree []string, parent string) []string {
	println("Generating tree:", purpose)
	for i, local := range tree {
		wantDir := strings.HasSuffix(local, "/")
		if wantDir {
			makeDir(filepath.Join(parent, local))
		} else {
			makeFile(filepath.Join(parent, local))
		}
		fmt.Printf("%d: %s\n", i+1, local)
	}
	return tree
}

func makeFile(path string) string {
	dir, _ := filepath.Split(path)
	makeDir(dir)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0755)
	Expect(err).ShouldNot(HaveOccurred())
	_, _ = file.WriteString("dummy text!")
	_ = file.Close()
	return path
}
