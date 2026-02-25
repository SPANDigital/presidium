package themes

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/SPANDigital/presidium-hugo/pkg/filesystem"
)

func TestNew_WithEmbeddedFS(t *testing.T) {
	// Create a mock embedded FS with a themes directory
	mockFS := fstest.MapFS{
		"themes/presidium-styling-base/config.yml": &fstest.MapFile{Data: []byte("name: test")},
	}
	SetFS(mockFS)
	defer SetFS(nil)

	svc, err := New()
	if err != nil {
		t.Fatalf("New() returned unexpected error: %v", err)
	}
	if svc.themes == nil {
		t.Fatal("New() returned service with nil themes FS")
	}
}

func TestNew_WithNilFS_NoModuleRoot(t *testing.T) {
	// Ensure themesFS is nil
	originalFS := GetThemesFS()
	SetFS(nil)
	defer func() {
		if originalFS != nil {
			SetFS(originalFS)
		}
	}()

	// Change to a temp dir with no go.mod above it
	origDir, _ := os.Getwd()
	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	_, err := New()
	if err == nil {
		t.Fatal("New() should return error when no go.mod found and no embedded FS")
	}
	if !strings.Contains(err.Error(), "no go.mod found") {
		t.Errorf("expected 'no go.mod found' error, got: %v", err)
	}
}

func TestExtract_WritesFilesAndProducesReplacements(t *testing.T) {
	// Create a mock FS with theme files including go.mod.tmpl and go.sum.tmpl
	mockFS := fstest.MapFS{
		"presidium-styling-base/config.yml":  &fstest.MapFile{Data: []byte("name: styling")},
		"presidium-styling-base/go.mod.tmpl": &fstest.MapFile{Data: []byte("module github.com/spandigital/presidium-styling-base")},
		"presidium-styling-base/go.sum.tmpl": &fstest.MapFile{Data: []byte("some-dep v1.0.0 h1:abc=")},
		"presidium-layouts-base/config.yml":  &fstest.MapFile{Data: []byte("name: layouts-base")},
		"presidium-layouts-base/go.mod.tmpl": &fstest.MapFile{Data: []byte("module github.com/spandigital/presidium-layouts-base")},
		"presidium-layouts-blog/config.yml":  &fstest.MapFile{Data: []byte("name: layouts-blog")},
		"presidium-layouts-blog/go.mod.tmpl": &fstest.MapFile{Data: []byte("module github.com/spandigital/presidium-layouts-blog")},
	}

	svc := Service{themes: mockFS}

	tmpDir, replacements, err := svc.Extract()
	if err != nil {
		t.Fatalf("Extract() returned unexpected error: %v", err)
	}
	defer func() { _ = filesystem.AFS.RemoveAll(tmpDir) }()

	// Verify go.mod.tmpl was renamed to go.mod on disk
	goModPath := filepath.Join(tmpDir, "presidium-styling-base", "go.mod")
	if exists, _ := filesystem.AFS.Exists(goModPath); !exists {
		t.Error("Expected go.mod.tmpl to be extracted as go.mod, but file does not exist")
	}

	// Verify go.sum.tmpl was renamed to go.sum on disk
	goSumPath := filepath.Join(tmpDir, "presidium-styling-base", "go.sum")
	if exists, _ := filesystem.AFS.Exists(goSumPath); !exists {
		t.Error("Expected go.sum.tmpl to be extracted as go.sum, but file does not exist")
	}

	// Verify config.yml was extracted
	configPath := filepath.Join(tmpDir, "presidium-styling-base", "config.yml")
	if exists, _ := filesystem.AFS.Exists(configPath); !exists {
		t.Error("Expected config.yml to be extracted, but file does not exist")
	}

	// Verify replacements string contains all three modules
	for _, mod := range []string{
		"github.com/spandigital/presidium-styling-base",
		"github.com/spandigital/presidium-layouts-base",
		"github.com/spandigital/presidium-layouts-blog",
	} {
		if !strings.Contains(replacements, mod) {
			t.Errorf("Expected replacements to contain %q, got: %s", mod, replacements)
		}
	}

	// Verify replacements are in sorted order (deterministic)
	parts := strings.Split(replacements, ",")
	if len(parts) != 3 {
		t.Fatalf("Expected 3 replacement pairs, got %d: %s", len(parts), replacements)
	}
	for i := 1; i < len(parts); i++ {
		if parts[i-1] > parts[i] {
			t.Errorf("Replacement pairs are not sorted: %q > %q", parts[i-1], parts[i])
		}
	}
}

func TestExtract_EmptyTheme(t *testing.T) {
	// Test with a theme that has no files (just a directory entry)
	mockFS := fstest.MapFS{
		"presidium-styling-base/config.yml": &fstest.MapFile{Data: []byte("name: test")},
		"presidium-layouts-base/config.yml": &fstest.MapFile{Data: []byte("name: test")},
		"presidium-layouts-blog/config.yml": &fstest.MapFile{Data: []byte("name: test")},
	}

	svc := Service{themes: mockFS}

	tmpDir, replacements, err := svc.Extract()
	if err != nil {
		t.Fatalf("Extract() returned unexpected error: %v", err)
	}
	defer func() { _ = filesystem.AFS.RemoveAll(tmpDir) }()

	if replacements == "" {
		t.Error("Expected non-empty replacements string")
	}
}

func TestFindModuleRoot(t *testing.T) {
	// Save and restore working directory
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()

	// Create a temp directory structure with go.mod
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "a", "b", "c")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module test"), 0644); err != nil {
		t.Fatal(err)
	}

	// Change to the deeply nested directory
	if err := os.Chdir(subDir); err != nil {
		t.Fatal(err)
	}

	root, err := findModuleRoot()
	if err != nil {
		t.Fatalf("findModuleRoot() returned unexpected error: %v", err)
	}
	// Resolve symlinks for comparison (macOS /var -> /private/var)
	wantResolved, _ := filepath.EvalSymlinks(tmpDir)
	gotResolved, _ := filepath.EvalSymlinks(root)
	if gotResolved != wantResolved {
		t.Errorf("findModuleRoot() = %q, want %q", gotResolved, wantResolved)
	}
}

func TestFindModuleRoot_NotFound(t *testing.T) {
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()

	// Use a temp directory with no go.mod anywhere above
	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	_, err := findModuleRoot()
	if err == nil {
		t.Fatal("findModuleRoot() should return error when no go.mod found")
	}
	if !strings.Contains(err.Error(), "no go.mod found") {
		t.Errorf("expected 'no go.mod found' error, got: %v", err)
	}
}
