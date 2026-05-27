package themes

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"

	"github.com/SPANDigital/presidium-hugo/pkg/filesystem"
)

// createTestZip creates a zip file in memory with the specified files
func createTestZip(files map[string][]byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	for name, content := range files {
		f, err := w.Create(name)
		if err != nil {
			return nil, err
		}
		if _, err := f.Write(content); err != nil {
			return nil, err
		}
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func TestNew_WithEmbeddedZip(t *testing.T) {
	// Create a mock zip with theme files
	zipData, err := createTestZip(map[string][]byte{
		"presidium-styling-base/config.yml": []byte("name: test"),
	})
	if err != nil {
		t.Fatal(err)
	}

	SetZip(zipData)
	defer SetZip(nil)

	svc, err := New()
	if err != nil {
		t.Fatalf("New() returned unexpected error: %v", err)
	}
	if svc.zipData == nil {
		t.Fatal("New() returned service with nil zip data")
	}
}

func TestNew_WithNilZip(t *testing.T) {
	// Ensure themesZip is nil
	originalZip := GetThemesZip()
	SetZip(nil)
	defer func() {
		if originalZip != nil {
			SetZip(originalZip)
		}
	}()

	_, err := New()
	if err == nil {
		t.Fatal("New() should return error when no zip data available")
	}
	if !strings.Contains(err.Error(), "no themes zip data available") {
		t.Errorf("expected 'no themes zip data available' error, got: %v", err)
	}
}

func TestExtract_WritesFilesAndProducesReplacements(t *testing.T) {
	// Create a mock zip with theme files
	zipData, err := createTestZip(map[string][]byte{
		"presidium-styling-base/config.yml": []byte("name: styling"),
		"presidium-styling-base/go.mod":     []byte("module github.com/spandigital/presidium-styling-base"),
		"presidium-styling-base/go.sum":     []byte("some-dep v1.0.0 h1:abc="),
		"presidium-layouts-base/config.yml": []byte("name: layouts-base"),
		"presidium-layouts-base/go.mod":     []byte("module github.com/spandigital/presidium-layouts-base"),
		"presidium-layouts-blog/config.yml": []byte("name: layouts-blog"),
		"presidium-layouts-blog/go.mod":     []byte("module github.com/spandigital/presidium-layouts-blog"),
	})
	if err != nil {
		t.Fatal(err)
	}

	svc := Service{zipData: zipData}

	tmpDir, replacements, err := svc.Extract()
	if err != nil {
		t.Fatalf("Extract() returned unexpected error: %v", err)
	}
	defer func() { _ = filesystem.AFS.RemoveAll(tmpDir) }()

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
