package template

import (
	"testing"
	"testing/fstest"

	"github.com/spf13/afero"

	"github.com/SPANDigital/presidium-hugo/pkg/filesystem"
	model "github.com/SPANDigital/presidium-hugo/pkg/domain/model/generator"
)

func TestNew_CanLocateTemplates(t *testing.T) {
	svc, err := New()
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
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

func TestIsVerbatimPath(t *testing.T) {
	cases := []struct {
		name        string
		path        string
		templateDir string
		want        bool
	}{
		{"layouts file", "requirements/layouts/shortcodes/x.html", "requirements", true},
		{"static file", "requirements/static/img/logo.png", "requirements", true},
		{"data file", "requirements/data/discovery-stats.yaml", "requirements", true},
		{"scripts file", "requirements/scripts/snapshot-stats.py", "requirements", true},
		{"assets file", "requirements/assets/css/site.css", "requirements", true},
		{"discovery file", "requirements/discovery/foo.pdf", "requirements", true},
		{"resources file", "blog/resources/_gen/font.eot", "blog", true},
		{"public file", "blog/public/index.html", "blog", true},
		{"content file", "requirements/content/01-overview/_index.md", "requirements", false},
		{"archetypes file", "requirements/archetypes/requirements/scenario.md", "requirements", false},
		{"top-level config", "requirements/config.yaml", "requirements", false},
		{"go.mod.tmpl", "requirements/go.mod.tmpl", "requirements", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isVerbatimPath(tc.path, tc.templateDir); got != tc.want {
				t.Errorf("isVerbatimPath(%q, %q) = %v, want %v",
					tc.path, tc.templateDir, got, tc.want)
			}
		})
	}
}

// TestProcessDirTemplates_RealTemplatesGenerateCleanly walks every shipped
// template through the generator using the on-disk filesystem fallback. It
// catches any malformed Go template syntax in template files (e.g. an
// un-escaped Hugo shortcode under content/ or archetypes/) before they reach
// users.
func TestProcessDirTemplates_RealTemplatesGenerateCleanly(t *testing.T) {
	prevAFS := filesystem.AFS
	prevFS := filesystem.FS
	t.Cleanup(func() {
		filesystem.AFS = prevAFS
		filesystem.FS = prevFS
	})
	filesystem.SetFileSystem(afero.NewMemMapFs())

	svc, err := New()
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	params := model.TemplateParameters{
		Title:       "Smoke Test",
		ProjectName: "smoke-test",
	}

	for _, tmpl := range model.SupportedTemplates {
		t.Run(tmpl.Code(), func(t *testing.T) {
			out := "out/" + tmpl.Code()
			if err := svc.ProcessDirTemplates(tmpl.Code(), out, params); err != nil {
				t.Fatalf("ProcessDirTemplates(%q) returned error: %v", tmpl.Code(), err)
			}
		})
	}
}

// TestProcessDirTemplates_VerbatimLayoutsAreNotParsed proves that files under
// layouts/ — which contain Hugo template syntax that is not valid Go template
// syntax — are copied byte-for-byte to the output without invoking the
// text/template parser. Without this behavior, generation of any template
// shipping inline Hugo layouts would fail.
func TestProcessDirTemplates_VerbatimLayoutsAreNotParsed(t *testing.T) {
	hugoLayoutContent := `{{ .Page.Params.name }}
{{ if .Page.Params.attributes }}
{{ range $.Page.Params.attributes }}
  {{ if not (in $existing .) }}{{ end }}
{{ end }}
{{ end }}
`
	templatedContent := `title: "{{ .Title }}"
project: "{{ .ProjectName }}"
`

	fakeFS := fstest.MapFS{
		"templates/sample/layouts/shortcodes/entity-table.html": &fstest.MapFile{
			Data: []byte(hugoLayoutContent),
		},
		"templates/sample/content/_index.md": &fstest.MapFile{
			Data: []byte(templatedContent),
		},
	}

	prevTemplatesFS := templatesFS
	t.Cleanup(func() { templatesFS = prevTemplatesFS })
	SetFS(fakeFS)

	prevAFS := filesystem.AFS
	prevFS := filesystem.FS
	t.Cleanup(func() {
		filesystem.AFS = prevAFS
		filesystem.FS = prevFS
	})
	filesystem.SetFileSystem(afero.NewMemMapFs())

	svc, err := New()
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	params := model.TemplateParameters{
		Title:       "Test Site",
		ProjectName: "test-site",
	}

	if err := svc.ProcessDirTemplates("sample", "out", params); err != nil {
		t.Fatalf("ProcessDirTemplates returned error: %v", err)
	}

	gotLayout, err := filesystem.AFS.ReadFile("out/layouts/shortcodes/entity-table.html")
	if err != nil {
		t.Fatalf("reading verbatim output: %v", err)
	}
	if string(gotLayout) != hugoLayoutContent {
		t.Errorf("layouts file was not copied verbatim:\n  got:  %q\n  want: %q",
			string(gotLayout), hugoLayoutContent)
	}

	gotContent, err := filesystem.AFS.ReadFile("out/content/_index.md")
	if err != nil {
		t.Fatalf("reading templated output: %v", err)
	}
	want := `title: "Test Site"
project: "test-site"
`
	if string(gotContent) != want {
		t.Errorf("content file was not templated correctly:\n  got:  %q\n  want: %q",
			string(gotContent), want)
	}
}
