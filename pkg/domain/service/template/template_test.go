package template

import (
	"testing"

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
