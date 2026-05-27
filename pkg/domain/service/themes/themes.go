package themes

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/SPANDigital/presidium-hugo/pkg/filesystem"
)

// themesZip holds the embedded zip file set from main via SetZip.
var themesZip []byte

// SetZip sets the embedded zip file used to extract themes.
func SetZip(zipData []byte) {
	themesZip = zipData
}

// GetThemesZip returns the current themes zip data (useful for testing).
func GetThemesZip() []byte {
	return themesZip
}

// moduleMap maps Hugo module paths to their corresponding theme directory names.
var moduleMap = map[string]string{
	"github.com/spandigital/presidium-styling-base": "presidium-styling-base",
	"github.com/spandigital/presidium-layouts-base": "presidium-layouts-base",
	"github.com/spandigital/presidium-layouts-blog": "presidium-layouts-blog",
}

type Service struct {
	zipData []byte
}

func New() (Service, error) {
	zipData := themesZip
	if zipData == nil {
		return Service{}, fmt.Errorf("no themes zip data available")
	}
	return Service{
		zipData: zipData,
	}, nil
}

// Extract extracts all embedded themes from the zip to a temporary directory and returns:
// - tmpDir: the temporary directory path (caller must clean up)
// - replacements: comma-separated string for HUGO_MODULE_REPLACEMENTS env var
// - err: any error encountered during extraction
func (s Service) Extract() (tmpDir string, replacements string, err error) {
	// Create a temporary directory for themes
	tmpDir, err = filesystem.AFS.TempDir("", "presidium-themes-*")
	if err != nil {
		return "", "", fmt.Errorf("creating temp directory: %w", err)
	}

	// Open the zip archive from memory
	reader, err := zip.NewReader(bytes.NewReader(s.zipData), int64(len(s.zipData)))
	if err != nil {
		_ = filesystem.AFS.RemoveAll(tmpDir)
		return "", "", fmt.Errorf("opening zip archive: %w", err)
	}

	// Extract all files from zip
	// Guard against zip-slip: file.Name may contain "../" or absolute paths.
	tmpDirClean := filepath.Clean(tmpDir) + string(os.PathSeparator)
	for _, file := range reader.File {
		// Calculate destination path
		destPath := filepath.Join(tmpDir, file.Name)
		if !strings.HasPrefix(destPath+string(os.PathSeparator), tmpDirClean) && destPath != filepath.Clean(tmpDir) {
			_ = filesystem.AFS.RemoveAll(tmpDir)
			return "", "", fmt.Errorf("invalid zip entry path (zip-slip): %s", file.Name)
		}

		if file.FileInfo().IsDir() {
			// Create directory
			if err := filesystem.AFS.MkdirAll(destPath, 0755); err != nil {
				_ = filesystem.AFS.RemoveAll(tmpDir)
				return "", "", fmt.Errorf("creating directory %s: %w", destPath, err)
			}
			continue
		}

		// Ensure parent directory exists
		if err := filesystem.AFS.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			_ = filesystem.AFS.RemoveAll(tmpDir)
			return "", "", fmt.Errorf("creating directory for %s: %w", destPath, err)
		}

		// Open file from zip
		rc, err := file.Open()
		if err != nil {
			_ = filesystem.AFS.RemoveAll(tmpDir)
			return "", "", fmt.Errorf("opening file %s from zip: %w", file.Name, err)
		}

		// Create destination file
		outFile, err := filesystem.AFS.Create(destPath)
		if err != nil {
			rc.Close()
			_ = filesystem.AFS.RemoveAll(tmpDir)
			return "", "", fmt.Errorf("creating file %s: %w", destPath, err)
		}

		// Copy contents
		_, err = io.Copy(outFile, rc)
		rc.Close()
		outFile.Close()

		if err != nil {
			_ = filesystem.AFS.RemoveAll(tmpDir)
			return "", "", fmt.Errorf("extracting file %s: %w", file.Name, err)
		}
	}

	// Build replacements string for Hugo modules (sorted for deterministic output)
	modulePaths := make([]string, 0, len(moduleMap))
	for modulePath := range moduleMap {
		modulePaths = append(modulePaths, modulePath)
	}
	sort.Strings(modulePaths)

	var replacementPairs []string
	for _, modulePath := range modulePaths {
		themeName := moduleMap[modulePath]
		themeDir := filepath.Join(tmpDir, themeName)

		// Verify theme directory exists
		if exists, _ := filesystem.AFS.Exists(themeDir); !exists {
			_ = filesystem.AFS.RemoveAll(tmpDir)
			return "", "", fmt.Errorf("theme %s not found in extracted zip", themeName)
		}

		// Add to replacements: "modulePath -> localPath"
		replacementPairs = append(replacementPairs, fmt.Sprintf("%s -> %s", modulePath, themeDir))
	}

	// Join all replacements with commas for HUGO_MODULE_REPLACEMENTS
	replacements = strings.Join(replacementPairs, ",")

	return tmpDir, replacements, nil
}
