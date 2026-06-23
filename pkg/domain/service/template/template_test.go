package template

import (
	"os"
	"path/filepath"
	"testing"

	model "github.com/SPANDigital/presidium-hugo/pkg/domain/model/generator"
)

func TestMain(m *testing.M) {
	if err := model.LoadTemplates(os.DirFS(testModuleRoot())); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func testModuleRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	for {
		if _, statErr := os.Stat(filepath.Join(dir, "go.mod")); statErr == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "."
		}
		dir = parent
	}
}

func TestNew_CanLocateTemplates(t *testing.T) {
	svc, err := New()
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	if len(model.SupportedTemplates) == 0 {
		t.Fatal("expected SupportedTemplates to be populated by LoadTemplates")
	}

	for _, tmpl := range model.SupportedTemplates {
		listing, err := svc.GetListing(tmpl.Code())
		if err != nil {
			t.Errorf("GetListing(%q) returned error: %v", tmpl.Code(), err)
			continue
		}
		if len(listing) == 0 {
			t.Errorf("GetListing(%q) returned no files", tmpl.Code())
		}
	}
}
