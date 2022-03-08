package fileactions

import (
	"fmt"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/conversion/colors"
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

	When("deriving article title", func() {
		// Just add more expectations here to have them tested
		givenExpectations := map[string]string{
			"0.1.1-sales-exercises": "Sales Exercises",
			"financial--activities": "Financial Activities",
			"1.1.1-the-happyð-path": "The Happyð Path",
		}
		for given, expecting := range givenExpectations {
			should := fmt.Sprintf("title of \"%s\" should be: \"%s\"", given, expecting)
			name := given       // make copy of given & expectation as the
			wanted := expecting // underlying function changes the actual string.
			It(should, func() {
				actual := unSlugify(name)
				Expect(actual).Should(Equal(wanted))
			})
		}
	})
})
