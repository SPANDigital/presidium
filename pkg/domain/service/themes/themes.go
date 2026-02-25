package themes

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/SPANDigital/presidium-hugo/pkg/filesystem"
)

// themesFS holds the embedded filesystem set from main via SetFS.
var themesFS fs.FS

// SetFS sets the embedded filesystem used to read themes.
func SetFS(fsys fs.FS) {
	themesFS = fsys
}

// GetThemesFS returns the current themes filesystem (useful for testing).
func GetThemesFS() fs.FS {
	return themesFS
}

// moduleMap maps Hugo module paths to their corresponding theme directory names.
var moduleMap = map[string]string{
	"github.com/spandigital/presidium-styling-base": "presidium-styling-base",
	"github.com/spandigital/presidium-layouts-base": "presidium-layouts-base",
	"github.com/spandigital/presidium-layouts-blog": "presidium-layouts-blog",
}

type Service struct {
	themes fs.FS
}

func New() (Service, error) {
	fsys := themesFS
	if fsys == nil {
		// Fallback for tests: read themes from the module root filesystem.
		root, err := findModuleRoot()
		if err != nil {
			return Service{}, fmt.Errorf("locating module root: %w", err)
		}
		fsys = os.DirFS(root)
	}
	// Sub-FS into "themes" so callers can use theme names directly.
	sub, err := fs.Sub(fsys, "themes")
	if err != nil {
		return Service{}, fmt.Errorf("themes directory not found: %w", err)
	}
	return Service{
		themes: sub,
	}, nil
}

// findModuleRoot walks up from the working directory to find the module root.
func findModuleRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("unable to determine working directory: %w", err)
	}
	startDir := dir
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("no go.mod found when searching upwards from %s", startDir)
		}
		dir = parent
	}
}

// Extract extracts all embedded themes to a temporary directory and returns:
// - tmpDir: the temporary directory path (caller must clean up)
// - replacements: comma-separated string for HUGO_MODULE_REPLACEMENTS env var
// - err: any error encountered during extraction
func (s Service) Extract() (tmpDir string, replacements string, err error) {
	// Create a temporary directory for themes
	tmpDir, err = filesystem.AFS.TempDir("", "presidium-themes-*")
	if err != nil {
		return "", "", fmt.Errorf("creating temp directory: %w", err)
	}

	// Extract each theme (sorted for deterministic output)
	modulePaths := make([]string, 0, len(moduleMap))
	for modulePath := range moduleMap {
		modulePaths = append(modulePaths, modulePath)
	}
	sort.Strings(modulePaths)

	var replacementPairs []string
	for _, modulePath := range modulePaths {
		themeName := moduleMap[modulePath]
		// Create the theme directory in temp
		themeDir := filepath.Join(tmpDir, themeName)
		if err := filesystem.AFS.MkdirAll(themeDir, 0755); err != nil {
			_ = filesystem.AFS.RemoveAll(tmpDir)
			return "", "", fmt.Errorf("creating theme directory %s: %w", themeName, err)
		}

		// Copy theme files from embedded FS to temp directory
		err := fs.WalkDir(s.themes, themeName, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			// Calculate destination path
			destPath := filepath.Join(tmpDir, p)

			if d.IsDir() {
				return filesystem.AFS.MkdirAll(destPath, 0755)
			}

			// Read file from embedded FS
			content, err := fs.ReadFile(s.themes, p)
			if err != nil {
				return fmt.Errorf("reading %s: %w", p, err)
			}

			// Handle go.mod.tmpl -> go.mod and go.sum.tmpl -> go.sum renaming
			if strings.HasSuffix(destPath, "go.mod.tmpl") || strings.HasSuffix(destPath, "go.sum.tmpl") {
				destPath = strings.TrimSuffix(destPath, ".tmpl")
			}

			// Ensure parent directory exists
			if err := filesystem.AFS.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
				return fmt.Errorf("creating directory for %s: %w", destPath, err)
			}

			// Write file to temp directory
			f, err := filesystem.AFS.Create(destPath)
			if err != nil {
				return fmt.Errorf("creating file %s: %w", destPath, err)
			}
			defer f.Close()

			if _, err := f.Write(content); err != nil {
				return fmt.Errorf("writing %s: %w", destPath, err)
			}

			return nil
		})

		if err != nil {
			_ = filesystem.AFS.RemoveAll(tmpDir)
			return "", "", fmt.Errorf("extracting theme %s: %w", themeName, err)
		}

		// Add to replacements: "modulePath -> localPath"
		replacementPairs = append(replacementPairs, fmt.Sprintf("%s -> %s", modulePath, themeDir))
	}

	// Join all replacements with commas for HUGO_MODULE_REPLACEMENTS
	replacements = strings.Join(replacementPairs, ",")

	return tmpDir, replacements, nil
}
