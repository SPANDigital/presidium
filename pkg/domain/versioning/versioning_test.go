package versioning

import (
	"github.com/SPANDigital/presidium-hugo/pkg/filesystem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

type VersioningTestSuite struct {
	suite.Suite
	versioning  Versioning
	projectPath string
	fs          filesystem.FileSystem
}

func TestRunVersioningTestSuite(t *testing.T) {
	suite.Run(t, new(VersioningTestSuite))
}

func (s *VersioningTestSuite) SetupSuite() {

	s.projectPath = filepath.Join(os.TempDir(), "versioning_test")

	fileSystem := filesystem.New()
	_ = fileSystem.MakeDirs(filepath.Join(s.projectPath, "content"))

	s.fs = fileSystem
	s.createContent("about.md", "# About")
	s.createContent("article.md", "# Best Practices")
	s.versioning = New(s.projectPath)

	_ = filepath.Walk(s.projectPath, func(path string, info fs.FileInfo, err error) error {
		println(path)
		return nil
	})
}

func (s *VersioningTestSuite) TearDownSuite() {
	_ = filesystem.New().DeleteDir(s.projectPath)
}

func (s *VersioningTestSuite) createContent(file string, content string) {

	fileName := filepath.Join(s.projectPath, "content", file)
	_ = ioutil.WriteFile(fileName, []byte(content), 0777)
}

func (s *VersioningTestSuite) TestVersioning_IsEnabled() {
	assert.True(s.T(), s.versioning.IsEnabled())
}

func (s *VersioningTestSuite) TestVersioning_NextVersion() {
	assert.False(s.T(), s.versioning.IsActivated())
	assert.Equal(s.T(), 0, s.versioning.GetLatestVersionNo())
	s.versioning.NextVersion()
	assert.True(s.T(), s.versioning.IsActivated())
	assert.Equal(s.T(), 1, s.versioning.GetLatestVersionNo())

	s.versioning.NextVersion()
	assert.Equal(s.T(), 2, s.versioning.GetLatestVersionNo())

}
