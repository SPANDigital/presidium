package filesystem

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type FileSystemTestSuite struct {
	suite.Suite
	f       FileSystem
	testDir string
}

func (s *FileSystemTestSuite) SetupSuite() {
	s.f = New()
	s.testDir = "../../test/data/pkg/filesystem/testdir"
}

func TestRunFileSystemSuite(t *testing.T) {
	suite.Run(t, new(FileSystemTestSuite))
}

func (s *FileSystemTestSuite) TestFileSystem_MakeDirs() {
	err := s.f.MakeDirs(filepath.Join(s.testDir, "some/dirs/to/create"))
	assert.NoError(s.T(), err)
}

func (s *FileSystemTestSuite) TestFileSystem_DeleteDir() {

	dirTree := []string{
		"/documents",
		"/documents/personal",
		"/archives/documents/1",
		"/archives/documents/2",
	}

	for _, dirPath := range dirTree {
		dir := filepath.Join(s.testDir, dirPath)
		err := os.MkdirAll(dir, os.ModePerm)
		assert.NoError(s.T(), err)
	}

	err := s.f.DeleteDir(s.testDir)
	assert.NoError(s.T(), err)
}

func (s *FileSystemTestSuite) TestFileSystem_Rename() {

	var err error

	old := filepath.Join(s.testDir, "old")
	newDir := filepath.Join(s.testDir, "new")

	err = os.MkdirAll(old, os.ModePerm)
	assert.NoError(s.T(), err)

	err = s.f.Rename(old, newDir)
	assert.NoError(s.T(), err)

	info, err := os.Stat(newDir)
	assert.NoError(s.T(), err)
	assert.True(s.T(), info.IsDir())

}

func (s *FileSystemTestSuite) TestFileSystem_Copy() {
	srcPath := filepath.Join(s.testDir, "testfile1.md")
	destPath := filepath.Join(s.testDir, "..", "testfile.md")

	err := s.f.Copy(srcPath, destPath, fs.ModePerm)
	assert.NoError(s.T(), err, err)

	// Check the file exists and is the same file
	_, err = os.Open(destPath)
	assert.True(s.T(), !os.IsNotExist(err), "the file should exist")

	srcFile, err := ioutil.ReadFile(srcPath)
	assert.NoError(s.T(), err, err)
	destFile, err := ioutil.ReadFile(destPath)
	assert.NoError(s.T(), err, err)
	assert.True(s.T(), bytes.Equal(srcFile, destFile))

	// Clean up dir
	err = os.Remove(destPath)
	assert.NoError(s.T(), err, err)
}

func (s *FileSystemTestSuite) TestFileSystem_CopyDir() {
	srcPath := s.testDir
	destPath := filepath.Join(s.testDir, "..", "result")

	err := s.f.CopyDir(srcPath, destPath)
	assert.NoError(s.T(), err, err)

	srcFiles := make([]string, 0)
	_ = filepath.WalkDir(s.testDir, func(path string, d fs.DirEntry, err error) error {
		srcFiles = append(srcFiles, path)
		return nil
	})

	destFiles := make([]string, 0)
	_ = filepath.WalkDir(destPath, func(path string, d fs.DirEntry, err error) error {
		destFiles = append(destFiles, strings.ReplaceAll(path, "result", "testdir"))
		return nil
	})

	assert.EqualValues(s.T(), srcFiles, destFiles, "the folder structure should be the same")
	_ = os.RemoveAll(destPath)
}
