package versioning

import (
	"bufio"
	"fmt"
	"github.com/SPANDigital/presidium-hugo/pkg/filesystem"
	"os"
	"path/filepath"
	"strconv"
)

/*
Versioning Manages local versioning of presidium site.

NOTE: Only keeps up to 5 recent versions. Latest versions are appended. Oldest version are dropped off to make place
for a next version.

Workflow is the following:

1. User enables versioning

2. User activates next version

3. User works on site

4. User can any time work on update the latest version.

5. If the user wants to go back he can either reclaim the last version before the update, or restore a previous 1st,
2nd, 3rd, or 4th version

*/
type Versioning interface {
	IsEnabled() bool              // check if the versioning has been enabled or not
	SetEnabled(enabled bool) bool // enables version if not set.
	NextVersion()                 // activates the next version. Drop the oldest if exceeds more than 5 versions.
	GrabLatest()                  // Updates the last version backup
	IsActivated() bool            // Do we have a version?
	GetLatestVersionNo() int      // returns the latest version no
}

type versioning struct {
	fileSystem       filesystem.FileSystem // Local file system for copy actions.
	projectRoot      string                // The root the project
	siteContent      string                // path to the site convent
	versionsRootPath string                // path where all versions live. this is "versions/"
	enabled          bool                  // If this feature has been enabled ot not.
	versionNo        int                   // version number from 1 to 5
	versionLocalPath string                // location of this version, for example "versions/3"
	activated        bool                  // flag to version set has been activated
	statusFile       string                // keep track of the status
}

func (v *versioning) GetLatestVersionNo() int {
	return v.versionNo
}

func (v *versioning) IsActivated() bool {
	return v.activated
}

const maxVersionsToKeep = 5
const versionInfoFile = ".version"
const versionLocation = "versions"
const siteContentDirName = "content"

func (v *versioning) activateLatest() {

	v.load()
	v.activated = false
	v.versionNo = 0
	updated := false

	for i := maxVersionsToKeep; i > 0; i-- {
		pathToVersion := filepath.Join(v.versionsRootPath, fmt.Sprintf("%d", i))
		if _, err := os.Stat(pathToVersion); !os.IsNotExist(err) {
			v.activated = true
			v.versionLocalPath = pathToVersion
			v.versionNo = i
			updated = true
			break
		}
	}

	if updated {
		v.persist()
	}
}

func New(projectRoot string) Versioning {

	v := &versioning{
		fileSystem:       filesystem.New(),
		projectRoot:      projectRoot,
		versionsRootPath: filepath.Join(projectRoot, versionLocation),
		enabled:          false,
		versionNo:        0,
		statusFile:       filepath.Join(projectRoot, versionLocation, versionInfoFile),
		siteContent:      filepath.Join(projectRoot, siteContentDirName),
	}

	if err := v.fileSystem.MakeDirs(v.versionsRootPath); !os.IsNotExist(err) {
		v.enabled = true
		v.activateLatest()
	}

	return v
}

func (v *versioning) IsEnabled() bool {
	return v.enabled
}

func (v versioning) SetEnabled(enabled bool) bool {
	v.enabled = enabled
	return v.enabled
}

func (v *versioning) persist() {

	if _, err := os.Stat(v.versionsRootPath); os.IsNotExist(err) {
		if err = os.MkdirAll(v.versionsRootPath, os.ModePerm); err != nil {
			panic(err)
		}
	}

	file, err := os.Create(v.statusFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, _ = file.WriteString(fmt.Sprintf("%s\n", strconv.FormatBool(v.enabled)))
	_, _ = file.WriteString(fmt.Sprintf("%d\n", v.versionNo))
	_ = file.Sync()
}

func (v *versioning) load() {

	file, err := os.OpenFile(v.statusFile, os.O_RDONLY, 0666)

	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		panic(err)
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for i := 0; i < 2 && scanner.Scan(); i++ {
		switch i {
		case 0:
			if b, err := strconv.ParseBool(scanner.Text()); err == nil {
				v.enabled = b
			}
			break
		case 1:
			if i, err := strconv.ParseInt(scanner.Text(), 10, 16); err == nil {
				v.versionNo = int(i)
			}
			break
		}
	}
}

func (v *versioning) NextVersion() {

	if !v.activated {
		v.activateLatest()
	}

	v.activated = true

	oldVersion := v.versionNo
	nextVersion := oldVersion + 1

	if nextVersion > 5 {
		v.dropVersion(1)
		nextVersion = 5
	}

	v.setCurrentVersionNo(nextVersion)
	v.grabContent()
	v.persist()
}

func (v *versioning) GrabLatest() {
	if !v.enabled || v.activated {
		return
	}
	v.grabContent()
}

func (v *versioning) dropVersion(versionNo int) {
	for i := versionNo; i < maxVersionsToKeep+1; i++ {
		versionContentPath := filepath.Join(v.versionsRootPath, fmt.Sprintf("%d", i))
		if i == versionNo {
			_ = v.fileSystem.EmptyDir(versionContentPath)
		} else {
			newName := filepath.Join(v.versionsRootPath, fmt.Sprintf("%d", i-1))
			_ = v.fileSystem.Rename(versionContentPath, newName)
		}
	}
}

func (v *versioning) grabContent() {
	_ = v.fileSystem.MakeDirs(v.versionLocalPath)
	_ = v.fileSystem.CopyDir(v.siteContent, v.versionLocalPath)
}

func (v *versioning) setCurrentVersionNo(version int) {
	if v.versionNo != version {
		v.versionNo = version
		v.versionLocalPath = filepath.Join(v.versionsRootPath, fmt.Sprintf("%d", version))
		_ = v.fileSystem.MakeDirs(v.versionLocalPath)
	}
}
