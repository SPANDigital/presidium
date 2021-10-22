package filesystem

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

func TestFilesystem(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Filesystem Suite")
}

func init() {
	FS = afero.NewMemMapFs()
	FSUtil = &afero.Afero{Fs: FS}
}

var (
	filesystem = New()
)

var _ = Describe("Filesystem", func() {
	Describe("copy functionality", func() {
		Context("when calling Copy()", func() {
			srcFileName := "testfile1.md"
			dstFileName := "testfile.md"
			testDir := "/home/testuser/testdata/copy/test"
			BeforeEach(func() {
				FS.MkdirAll(testDir, 0755)
				FSUtil.WriteFile(fmt.Sprintf("%s/%s", testDir, srcFileName), []byte("Hello World!"), 0644)
			})
			AfterEach(func() {
				// no need to clean up - memory mapped filesystem will just go away
			})
			It("Should correctly copy a file", func() {
				srcPath := filepath.Join(testDir, srcFileName)
				destPath := filepath.Join(testDir, "..", dstFileName)

				err := filesystem.Copy(srcPath, destPath, fs.ModePerm)
				Expect(err).NotTo(HaveOccurred())

				// Check the file exists and is the same file
				_, err = FS.Open(destPath)
				Expect(err).NotTo(HaveOccurred())
				// Cannot use Gomega's BeAnExistingFile() here
				_, err = FS.Stat(destPath)
				Expect(err).NotTo(HaveOccurred())

				srcFile, err := FSUtil.ReadFile(srcPath)
				Expect(err).NotTo(HaveOccurred())
				destFile, err := FSUtil.ReadFile(destPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(srcFile).To(Equal(destFile))

				// Clean up dir
				err = FS.Remove(destPath)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("when calling CopyDir()", func() {
			testDir := "/home/testuser/testdata/copydir/test"
			file1, file2, file3 := "file1.md", "file2.md", "file3.md"
			BeforeEach(func() {
				FS.MkdirAll(testDir, 0755)
				FSUtil.WriteFile(fmt.Sprintf("%s/%s", testDir, file1), []byte("Hello World!"), 0644)
				FSUtil.WriteFile(fmt.Sprintf("%s/%s", testDir, file2), []byte("Hello World!"), 0644)
				FSUtil.WriteFile(fmt.Sprintf("%s/%s", testDir, file3), []byte("Hello World!"), 0644)
			})
			AfterEach(func() {
				// no need to clean up - memory mapped filesystem will just go away
			})
			It("Should copy the contents of a directory", func() {
				srcPath := testDir
				destPath := filepath.Join(testDir, "..", "result")

				err := filesystem.CopyDir(srcPath, destPath)
				Expect(err).NotTo(HaveOccurred())

				srcFiles := make([]string, 0)
				_ = filepath.WalkDir(testDir, func(path string, d fs.DirEntry, err error) error {
					srcFiles = append(srcFiles, path)
					return nil
				})

				destFiles := make([]string, 0)
				_ = filepath.WalkDir(destPath, func(path string, d fs.DirEntry, err error) error {
					destFiles = append(destFiles, strings.ReplaceAll(path, "result", "test"))
					return nil
				})

				Expect(srcFiles).To(Equal(destFiles))
				_ = FS.RemoveAll(destPath)
			})
		})
	})
	Describe("Rename functionality", func() {
		Context("When calling Rename()", func() {
			testDir := "/home/testuser/testdata/rename/test"
			AfterEach(func() {
				// no need to clean up - memory mapped filesystem will just go away
			})
			It("Should rename a folder correctly", func() {
				old := filepath.Join(testDir, "old")
				newDir := filepath.Join(testDir, "new")

				err := FS.MkdirAll(old, os.ModePerm)
				Expect(err).NotTo(HaveOccurred())

				err = filesystem.Rename(old, newDir)
				Expect(err).NotTo(HaveOccurred())

				info, err := FS.Stat(newDir)
				Expect(err).NotTo(HaveOccurred())
				Expect(info.IsDir()).To(BeTrue())
			})
		})
	})
	Describe("Delete functionality", func() {
		Context("When calling EmptyDir", func() {
			testDir := "/home/testuser/testdata/delete/test"
			AfterEach(func() {
				// no need to clean up - memory mapped filesystem will just go away
			})
			It("Should delete everything in a folder, leaving an empty directory", func() {
				dirTree := []string{
					"/documents",
					"/documents/personal",
					"/archives/documents/1",
					"/archives/documents/2",
				}

				for _, dirPath := range dirTree {
					dir := filepath.Join(testDir, dirPath)
					err := FS.MkdirAll(dir, os.ModePerm)
					Expect(err).NotTo(HaveOccurred())
				}

				err := filesystem.EmptyDir(testDir)
				Expect(err).NotTo(HaveOccurred())
				parentDir, err := FS.Open(testDir)
				Expect(err).NotTo(HaveOccurred())
				defer parentDir.Close()
				dirNames, err := parentDir.Readdirnames(-1)
				Expect(err).NotTo(HaveOccurred())
				Expect(dirNames).To(BeEmpty())
			})
		})
	})
	Describe("Resolving absolute path", func() {
		Context("Wen calling AbsolutePath()", func() {
			It("should remove all relative path place holders and replace it with the actual path", func() {
				relativePath := "./jellyBabyOhBaby"
				absolutePath, err := filesystem.AbsolutePath(relativePath)
				Expect(err).NotTo(HaveOccurred())
				Expect(absolutePath).NotTo(ContainSubstring("."))
			})
		})
	})
	Describe("Working with current working dir", func() {
		var workingDir string
		BeforeEach(func() {
			var err error
			workingDir, err = os.Getwd()
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("GetWorkingDir() should return the correct working dir", func() {
			actualWorkingDir, err := filesystem.GetWorkingDir()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(actualWorkingDir).Should(Equal(workingDir))
		})
	})
})
