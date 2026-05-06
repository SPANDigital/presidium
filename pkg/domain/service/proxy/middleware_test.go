package proxy

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockHandler is a simple handler that returns the request path and query
type mockHandler struct{}

func (h *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(r.URL.Path + "?" + r.URL.RawQuery))
}

func TestRewriteMiddleware_ArticleParam(t *testing.T) {
	handler := RewriteMiddleware(&mockHandler{})

	req := httptest.NewRequest("GET", "/docs/page?article=123", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// article param should be removed
	expected := "/docs/page?"
	if w.Body.String() != expected {
		t.Errorf("Expected %q, got %q", expected, w.Body.String())
	}
}

func TestRewriteMiddleware_SectionParam(t *testing.T) {
	handler := RewriteMiddleware(&mockHandler{})

	req := httptest.NewRequest("GET", "/docs/my-module/some-article?section=my-module", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Should rewrite to /<section>/
	expected := "/my-module/?"
	if w.Body.String() != expected {
		t.Errorf("Expected %q, got %q", expected, w.Body.String())
	}
}

func TestRewriteMiddleware_FormatMd(t *testing.T) {
	handler := RewriteMiddleware(&mockHandler{})

	req := httptest.NewRequest("GET", "/docs/page?format=md", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Should rewrite to /path/index.md
	expected := "/docs/page/index.md?"
	if w.Body.String() != expected {
		t.Errorf("Expected %q, got %q", expected, w.Body.String())
	}

	// Should set markdown content type
	contentType := w.Header().Get("Content-Type")
	if contentType != "text/markdown; charset=utf-8" {
		t.Errorf("Expected Content-Type text/markdown, got %q", contentType)
	}
}

func TestRewriteMiddleware_FormatEmbed(t *testing.T) {
	handler := RewriteMiddleware(&mockHandler{})

	req := httptest.NewRequest("GET", "/docs/page?format=embed", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Should rewrite to /path/embed.html
	expected := "/docs/page/embed.html?"
	if w.Body.String() != expected {
		t.Errorf("Expected %q, got %q", expected, w.Body.String())
	}

	// Should set HTML content type
	contentType := w.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("Expected Content-Type text/html, got %q", contentType)
	}
}

func TestRewriteMiddleware_FormatMdWithTrailingSlash(t *testing.T) {
	handler := RewriteMiddleware(&mockHandler{})

	req := httptest.NewRequest("GET", "/docs/page/?format=md", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Should remove trailing slash and add /index.md
	expected := "/docs/page/index.md?"
	if w.Body.String() != expected {
		t.Errorf("Expected %q, got %q", expected, w.Body.String())
	}
}

func TestRewriteMiddleware_NoParams(t *testing.T) {
	handler := RewriteMiddleware(&mockHandler{})

	req := httptest.NewRequest("GET", "/docs/page", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Should pass through unchanged
	expected := "/docs/page?"
	if w.Body.String() != expected {
		t.Errorf("Expected %q, got %q", expected, w.Body.String())
	}
}

func TestNormalizePathMiddleware(t *testing.T) {
	handler := NormalizePathMiddleware(&mockHandler{})

	tests := []struct {
		input    string
		expected string
	}{
		{"/docs//page", "/docs/page?"},
		{"/docs/./page", "/docs/page?"},
		{"/docs/../page", "/page?"},
		{"/", "/?"},
		// Trailing slashes must be preserved — Hugo directory URLs need them.
		// Without preservation, Hugo issues a relative redirect that browsers
		// resolve against the current path, producing doubled path segments.
		{"/getting-started/", "/getting-started/?"},
		{"/docs/section/", "/docs/section/?"},
	}

	for _, tt := range tests {
		req := httptest.NewRequest("GET", tt.input, nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Body.String() != tt.expected {
			t.Errorf("For input %q, expected %q, got %q", tt.input, tt.expected, w.Body.String())
		}
	}
}

func TestChainMiddleware(t *testing.T) {
	// Test that middleware is applied in correct order
	handler := ChainMiddleware(
		&mockHandler{},
		NormalizePathMiddleware,
		RewriteMiddleware,
	)

	req := httptest.NewRequest("GET", "/docs//page?format=md", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Path should be normalized, then format rewrite applied
	expected := "/docs/page/index.md?"
	if w.Body.String() != expected {
		t.Errorf("Expected %q, got %q", expected, w.Body.String())
	}
}
