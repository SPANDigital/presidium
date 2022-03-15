package fileactions

import (
	"fmt"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/conversion/colors"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/conversion/markdown"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestFileActions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FileActions Suite")
}

var _ = Describe("Performing file actions", func() {
	colors.Setup()
	When("removing Jekyll weight indicators", func() {
		var stagingContentDir string
		var stagedContentFiles = []string{
			"introduction/",
			"introduction/1.1-work-culture/",
			"introduction/1.1-work-culture/_index.md",
			"introduction/1.1-work-culture/01-dos.md",
			"introduction/1.1-work-culture/02-donts.md",
			"introduction/1.1-work-culture/03-fun.md",
			"introduction/1.2-staff-facilities/",
			"introduction/1.2-staff-facilities/_index.md",
			"introduction/1.2-staff-facilities/1.1-meeting-room.md",
			"introduction/1.2-staff-facilities/2.1-conferences.md",
			"introduction/1.2-staff-facilities/2.2-kitchen.md",
			"introduction/1.2-staff-facilities/3-toilets-and-the-rest.md",
		}

		BeforeEach(func() {
			tempDir, err := ioutil.TempDir("", "stagedContentDir-*")
			Expect(err).ShouldNot(HaveOccurred())
			stagingContentDir = tempDir
			// Make staged content files and directories for to test on:
			for _, local := range stagedContentFiles {
				wantDir := strings.HasSuffix(local, "/")
				path := filepath.Join(stagingContentDir, local)
				if wantDir {
					err = os.MkdirAll(path, 0775)
					Expect(err).ShouldNot(HaveOccurred())
				} else /* WANT A FILE INSTEAD */ {
					f, e := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0755)
					Expect(e).ShouldNot(HaveOccurred())
					_, err := f.WriteString("dummy text")
					f.Close()
					Expect(err).ShouldNot(HaveOccurred())
				}
			}
		})

		AfterEach(func() { _ = os.RemoveAll(stagingContentDir) })

		It("should produce a clean tree of file paths with no with weight indicators", func() {
			regexNameStarsWithWeighIndicators := regexp.MustCompile(`^[\d+\-.]+`)
			err := RemoveWeightIndicatorsFromFilePaths(stagingContentDir)
			Expect(err).ShouldNot(HaveOccurred())
			pathsWithWeightIndicators := make([]string, 0)
			err = filepath.WalkDir(stagingContentDir, func(path string, d fs.DirEntry, err error) error {
				local := strings.TrimPrefix(path, stagingContentDir)
				nameIsWeighted := regexNameStarsWithWeighIndicators.FindStringSubmatch(d.Name()) != nil
				if nameIsWeighted {
					pathsWithWeightIndicators = append(pathsWithWeightIndicators, local)
				}
				return nil
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(pathsWithWeightIndicators).Should(BeEmpty())
		})
	})

	When("deriving article title", func() {
		// Just add more expectations here to have them tested
		givenExpectations := map[string]string{
			"0.1.1-sales-exercises": "Sales Exercises",
			"financial--activities": "Financial Activities",
			"1.1.1-the-happyð-path": "The Happyð Path",
		}
		for given, expecting := range givenExpectations {
			should := fmt.Sprintf("title of \"%s\" should be: \"%s\"", given, expecting)
			It(should, func() {
				actual := unSlugify(given)
				Expect(actual).Should(Equal(expecting))
			})
		}
	})

	When("deriving article slug", func() {
		givenExpectations := map[string]string{
			"v0 .18.6_8.":      "v0-18-6-8",
			"update_terraform": "update-terraform",
		}
		for given, expecting := range givenExpectations {
			should := fmt.Sprintf("slug of \"%s\" should be: \"%s\"", given, expecting)
			It(should, func() {
				actual := slugify(given)
				Expect(actual).Should(Equal(expecting))
			})
		}
	})

	When("deriving slug from title", func() {
		givenExpectations := map[string]string{
			"Troubleshooting":                    "troubleshooting",
			"Set up Development Environment":     "set-up-development-environment",
			"\"Set up Development Environment\"": "set-up-development-environment",
			"Introduction & Overview":            "introduction-and-overview",
		}
		for given, expecting := range givenExpectations {
			should := fmt.Sprintf("title of \"%s\" should be: \"%s\"", given, expecting)
			It(should, func() {
				actual := titleToSlug(given)
				Expect(actual).Should(Equal(expecting))
			})
		}
	})

	When("slug is already set", func() {
		It("should not override the slug", func() {
			actual := slugByPriority(markdown.FrontMatter{
				Slug: "demo",
			})
			Expect(actual).Should(BeNil())
		})
	})

	When("slugBasedOnFilename is false", func() {
		It("should be nil", func() {
			viper.Set("slugBasedOnFilename", true)
			actual := slugByPriority(markdown.FrontMatter{})
			Expect(actual).Should(BeNil())
		})
	})

	When("slugBasedOnFilename is false", func() {
		It("should be derived from the title", func() {
			viper.Set("slugBasedOnFilename", false)
			actual := slugByPriority(markdown.FrontMatter{
				Title: "Sample title",
			})
			Expect(*actual).Should(Equal("sample-title"))
		})
	})

	When("deduceWeightAndSlug", func() {
		givenExpectations := map[string]markdown.FrontMatter{
			"content/_development/02_update_terraform.md": {
				URL:    "development/update-terraform",
				Weight: "3", Slug: "update-terraform",
			},
			"content/products/presidium/differentiation.md": {
				URL:    "products/presidium/differentiation",
				Weight: "", Slug: "differentiation",
			},
			"content/_best-practices/02a-improve-data-completeness.md": {
				URL:    "best-practices/improve-data-completeness",
				Weight: "3", Slug: "improve-data-completeness",
			},
		}

		w := &contentWeightTracker{}
		for given, expecting := range givenExpectations {
			should := fmt.Sprintf("path of \"%s\" should be: \"%s\"", given, expecting.Slug)
			It(should, func() {
				actual := deduceWeightAndSlug("", given, w)
				Expect(actual).Should(Equal(expecting))
			})
		}
	})

	When("title should be derived from the path", func() {
		givenExpectations := map[string]string{
			"content/onboard/developer-authorization/create-private-key.md": "Create Private Key",
			"content/onboard/enrolling-as-an-apple-developer/enroll.md": "Enroll",
			"sample docs.md": "Sample Docs",
			"overview": "Ovewview",
		}

		for given, expecting := range givenExpectations {
			should := fmt.Sprintf("title of path \"%s\" should be: \"%s\"", given, expecting)
			It(should, func() {
				actual := titleFromPath(given)
				Expect(actual).Should(Equal(expecting))
			})
		}
	})
})
