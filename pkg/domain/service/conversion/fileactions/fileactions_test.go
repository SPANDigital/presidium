package fileactions

import (
	"fmt"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/conversion/colors"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/conversion/markdown"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
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

	When("getSlugAndUrl", func() {
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

		for given, expecting := range givenExpectations {
			should := fmt.Sprintf("path of \"%s\" should be: \"%s\"", given, expecting.Slug)
			a, b := given, expecting
			It(should, func() {
				slug, url := getSlugAndUrl("", b.Title, a)
				Expect(slug).Should(Equal(b.Slug))
				Expect(url).Should(Equal(b.URL))
			})
		}
	})

	When("title should be derived from the path", func() {
		givenExpectations := map[string]string{
			"content/onboard/developer-authorization/create-private-key.md": "Create Private Key",
			"content/onboard/enrolling-as-an-developer/enroll.md":     "Enroll",
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

	When("calculating path weight", func() {
		dm := directoryMap{
			"_getting_started": []string{
				"_getting_started/01_before_install",
				"_getting_started/02_installation",
				"_getting_started/03_shield_networking",
			},
			"_getting_started/01_before_install": []string{
				"_getting_started/01_before_install/00_introduction.md",
				"_getting_started/01_before_install/01_configure_aws.md",
			},
			"_getting_started/03_shield_networking": []string{
				"_getting_started/02_installation/01_create_s3_bucket_for_terraform.md",
				"_getting_started/02_installation/02_bootstrap_tf_env.md",
			},
		}

		givenExpectations := map[string]string{
			"_getting_started/03_shield_networking":                 "3",
			"_getting_started/01_before_install/00_introduction.md": "1",
			"_getting_started/01_before_install/_index.md":          "1",
			"_getting_started/03_shield_networking/_index.md":       "3",
			"_getting_started/notfound/_index.md":                   "",
		}

		for given, expecting := range givenExpectations {
			should := fmt.Sprintf("weight of path \"%s\" should be: \"%s\"", given, expecting)
			a, b := given, expecting
			It(should, func() {
				actual := getPathWeight(dm, a)
				Expect(actual).Should(Equal(b))
			})
		}
	})

	When("building a weight map", func() {
		afFs = afero.NewMemMapFs()
		mockFile("test/_index.md")
		mockFile("test/00_a.md")
		mockFile("test/network/01_a.md")
		mockFile("test/network/0_a.jpg")
		mockFile("test/store/3a_a.md")
		mockFile("test/store/2a_c.md")

		It("should build a valid weight map", func() {
			dm, err := buildWeightMap(".")
			Expect(err).Should(BeNil())
			Expect(dm).Should(Equal(directoryMap{
				"test":[]string{"test/00_a.md"},
				"test/network":[]string{"test/network/01_a.md"},
				"test/store":[]string{"test/store/2a_c.md", "test/store/3a_a.md"},
			}))
		})
	})

	When("checking if a file is an index", func() {
		givenExpectations := map[string]bool{
			"_getting_started/03_shield_networking": false,
			"_getting_started/index.md":             false,
			"_getting_started/_index.md":            true,
			"index.md":                              false,
			"_index.md":                             true,
			"":                                      false,
			"index":                                 false,
		}

		for given, expecting := range givenExpectations {
			should := fmt.Sprintf("path \"%s\" should be: \"%v\"", given, expecting)
			a, b := given, expecting
			It(should, func() {
				actual := isIndex(a)
				Expect(actual).Should(Equal(b))
			})
		}
	})

	When("checking if a file is markdown", func() {
		givenExpectations := map[string]bool{
			"_getting_started/03_shield_networking.md": true,
			"_getting_started/index":                   false,
			"_getting_started/_index.md":               true,
			"":                                         false,
		}

		for given, expecting := range givenExpectations {
			should := fmt.Sprintf("path \"%s\" should be: \"%v\"", given, expecting)
			a, b := given, expecting
			It(should, func() {
				actual := isMdFile(a)
				Expect(actual).Should(Equal(b))
			})
		}
	})
})

func mockFile(p string) {
	if err := afFs.MkdirAll(filepath.Dir(p), 0770); err != nil {
		Fail("failed to create test dir")
	}

	_, err := afFs.Create(p)
	if err != nil {
		Fail("failed to create test dir")
	}
}