package fileactions

import (
	"fmt"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/conversion/colors"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/conversion/markdown"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

	When("deduceWeightAndSlug", func() {
		dirUrls = map[string]string{
			"content/_development":       "development-lead",
			"content/products/presidium": "presidium",
			"content/_best-practices":    "best-practices-online",
		}
		givenExpectations := map[string]markdown.FrontMatter{
			"content/_development/02_update_terraform.md": {
				Title:  "Terraform Deployment",
				URL:    "development-lead/terraform-deployment",
				Weight: "3", Slug: "terraform-deployment",
			},
			"content/products/presidium/differentiation.md": {
				Title:  "Differentiation",
				URL:    "presidium/differentiation",
				Weight: "", Slug: "differentiation",
			},
			"content/_best-practices/02a-improve-data-completeness.md": {
				Title:  "Improve data completeness",
				URL:    "best-practices-online/improve-data-completeness",
				Weight: "3", Slug: "improve-data-completeness",
			},
		}

		w := &contentWeightTracker{}
		for given, expecting := range givenExpectations {
			should := fmt.Sprintf("path of \"%s\" should be: \"%s\"", given, expecting.Slug)
			a, b := given, expecting
			It(should, func() {
				actual := deduceWeightAndSlug("", b, a, w)
				Expect(actual).Should(Equal(b))
			})
		}
	})

	When("title should be derived from the path", func() {
		givenExpectations := map[string]string{
			"content/onboard/developer-authorization/create-private-key.md": "Create Private Key",
			"content/onboard/enrolling-as-an-apple-developer/enroll.md":     "Enroll",
			"sample docs.md": "Sample Docs",
			"_overview":      "Overview",
		}

		for given, expecting := range givenExpectations {
			should := fmt.Sprintf("title of path \"%s\" should be: \"%s\"", given, expecting)
			a, b := given, expecting
			It(should, func() {
				actual := titleFromPath(a)
				Expect(actual).Should(Equal(b))
			})
		}
	})
})
