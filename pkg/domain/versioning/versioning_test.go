package versioning

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVersioning(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Versioning")
}

var _ = Describe("project", func() {

	var project string

	BeforeSuite(func() {
		var localSite = []string{
			"about.md",
			"best-practice/article.md",
		}
		project = mustMakeWorkDir()
		mustMakeTree(project,
			localSite)
		listContent(project)
	})

	AfterSuite(func() { _ = os.RemoveAll(project) })

	Describe("after enabling versioning", func() {

		var v Versioning

		BeforeEach(func() {
			v = New(project)
			v.SetEnabled(true)
			Expect(v.IsEnabled()).Should(BeTrue())
		})

		When("calling NextVersion()", func() {
			It("should increase the version number by one", func() {
				versionBefore := v.GetLatestVersionNo()
				v.NextVersion()
				versionAfter := v.GetLatestVersionNo()
				Expect(versionAfter).Should(Equal(versionBefore + 1))
			})
		})
	})

})

func mustMakeWorkDir() string {
	workDir, err := ioutil.TempDir("", "versioning-test-dir-*")
	Expect(err).ShouldNot(HaveOccurred())
	return workDir
}

func listContent(path string) {
	println("Project")
	println("-------")
	listContentErr := filepath.WalkDir(path, func(path string, e fs.DirEntry, err error) error {
		println(path)
		return nil
	})
	Expect(listContentErr).ShouldNot(HaveOccurred())
}

func mustMakeTree(workDir string, template []string) []string {
	tree := make([]string, 0)
	for _, local := range template {
		path := filepath.Join(workDir, local)
		if strings.HasSuffix(local, "/") {
			mustHaveDir(path)
		} else {
			mustMakeFile(path)
		}
		tree = append(tree, path)
	}
	return tree
}

func mustHaveDir(path string) {
	dirInfo, dirInfoErr := os.Stat(path)
	if dirInfoErr != nil {
		if os.IsNotExist(dirInfoErr) {
			makeDirErr := os.MkdirAll(path, os.ModePerm)
			Expect(makeDirErr).ShouldNot(HaveOccurred())
		} else {
			Expect(dirInfoErr).ShouldNot(HaveOccurred())
		}
	} else {
		Expect(dirInfo.IsDir()).Should(BeTrue())
	}
}

func mustHaveParentDir(filePath string) {
	dir, _ := filepath.Split(filePath)
	mustHaveDir(dir)
}

func mustMakeFile(path string) {
	mustHaveParentDir(path)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	_, _ = file.WriteString("The resurrection of respecting creators is popular.")
	file.Close()
	Expect(err).ShouldNot(HaveOccurred())
}
