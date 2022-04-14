package fileactions

import (
	"fmt"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/conversion/colors"
	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/conversion/markdown"
	"github.com/SPANDigital/presidium-hugo/pkg/filesystem"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
	"io/fs"
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
	filesystem.SetFileSystem(afero.NewMemMapFs())

	colors.Setup()
	When("removing Jekyll weight indicators", func() {
		stageDir := "/tester"
		mockFile(stageDir, "introduction/1.1-work-culture/01-dos.md")
		mockFile(stageDir, "introduction/1.1-work-culture/03-fun.md")
		mockFile(stageDir, "introduction/1.2-staff-facilities/_index.md")
		mockFile(stageDir, "introduction/1.2-staff-facilities/2.1-conferences.md")
		mockFile(stageDir, "introduction/1.2-staff-facilities/3-toilets-and-the-rest.md")

		It("should produce a clean tree of file paths with no with weight indicators", func() {
			regexNameStarsWithWeighIndicators := regexp.MustCompile(`^[\d+\-.]+`)
			err := RemoveWeightIndicatorsFromFilePaths(stageDir)
			Expect(err).ShouldNot(HaveOccurred())
			pathsWithWeightIndicators := make([]string, 0)
			err = filesystem.AFS.Walk(stageDir, func(path string, info fs.FileInfo, err error) error {
				if info == nil {
					return nil
				}

				local := strings.TrimPrefix(path, stageDir)
				nameIsWeighted := regexNameStarsWithWeighIndicators.FindStringSubmatch(info.Name()) != nil
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
			"content/onboard/enrolling-as-an-developer/enroll.md":           "Enroll",
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
		stageDir := "/wm"
		mockFile(stageDir, "test/_index.md")
		mockFile(stageDir, "test/00_a.md")
		mockFile(stageDir, "test/network/01_a.md")
		mockFile(stageDir, "test/network/0_a.jpg")
		mockFile(stageDir, "test/store/3a_a.md")
		mockFile(stageDir, "test/store/2a_c.md")

		It("should build a valid weight map", func() {
			expected := directoryMap{
				"/wm/test":         []string{"/wm/test/00_a.md"},
				"/wm/test/network": []string{"/wm/test/network/01_a.md"},
				"/wm/test/store":   []string{"/wm/test/store/2a_c.md", "/wm/test/store/3a_a.md"},
			}

			dm, err := buildWeightMap(stageDir)
			Expect(err).Should(BeNil())
			for path, dirs := range expected {
				Expect(dm[path]).Should(Equal(dirs))
			}
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

func mockFile(stagingContentDir, p string) {
	path := filepath.Join(stagingContentDir, p)
	if err := filesystem.AFS.MkdirAll(filepath.Dir(path), 0770); err != nil {
		Fail("failed to create test dir")
	}

	_, err := filesystem.AFS.Create(path)
	if err != nil {
		Fail("failed to create test dir")
	}
}
