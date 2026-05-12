package proxy

import (
	"net/http"
	"net/http/httptest"
	"strings"
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

// TestRewriteMiddleware_ArticleParamVsFragmentIdentifier verifies that URLs with
// ?article=<id> query parameters and #<id> fragment identifiers produce the same
// server-side result. Fragment identifiers (#) are client-side only and never sent
// to the server, so both URLs should result in the same path after middleware processing.
func TestRewriteMiddleware_ArticleParamVsFragmentIdentifier(t *testing.T) {
	handler := RewriteMiddleware(&mockHandler{})

	// Test 1: URL with ?article= query parameter
	// This gets sent to the server and stripped by middleware
	req1 := httptest.NewRequest("GET", "/best-practices/?article=design-learning-objectives", nil)
	w1 := httptest.NewRecorder()
	handler.ServeHTTP(w1, req1)

	// Test 2: URL with # fragment identifier
	// Important: Fragment identifiers are NEVER sent to the server by the browser.
	// The server only sees the path before the #, so we test with just the path.
	req2 := httptest.NewRequest("GET", "/best-practices/", nil)
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, req2)

	// Both should produce the same server response: /best-practices/
	expected := "/best-practices/?"
	
	if w1.Body.String() != expected {
		t.Errorf("URL with ?article= query param: expected %q, got %q", expected, w1.Body.String())
	}
	
	if w2.Body.String() != expected {
		t.Errorf("URL with # fragment (client-side only): expected %q, got %q", expected, w2.Body.String())
	}
	
	// Verify both produce identical results
	if w1.Body.String() != w2.Body.String() {
		t.Errorf("URLs should produce same result. ?article= produced %q, # fragment produced %q",
			w1.Body.String(), w2.Body.String())
	}
}

// TestRewriteMiddleware_FormatMd_ComprehensiveTests provides comprehensive coverage
// for the ?format=md query parameter handling across various URL patterns.
func TestRewriteMiddleware_FormatMd_ComprehensiveTests(t *testing.T) {
	tests := []struct {
		name                string
		inputURL            string
		expectedPath        string
		expectedContentType string
	}{
		{
			name:                "root path with format=md",
			inputURL:            "/?format=md",
			expectedPath:        "//index.md?",
			expectedContentType: "text/markdown; charset=utf-8",
		},
		{
			name:                "single level path with format=md",
			inputURL:            "/overview?format=md",
			expectedPath:        "/overview/index.md?",
			expectedContentType: "text/markdown; charset=utf-8",
		},
		{
			name:                "nested path with format=md",
			inputURL:            "/reference/configuration?format=md",
			expectedPath:        "/reference/configuration/index.md?",
			expectedContentType: "text/markdown; charset=utf-8",
		},
		{
			name:                "deeply nested path with format=md",
			inputURL:            "/reference/front-matter/user-roles?format=md",
			expectedPath:        "/reference/front-matter/user-roles/index.md?",
			expectedContentType: "text/markdown; charset=utf-8",
		},
		{
			name:                "path with trailing slash and format=md",
			inputURL:            "/getting-started/?format=md",
			expectedPath:        "/getting-started/index.md?",
			expectedContentType: "text/markdown; charset=utf-8",
		},
		{
			name:                "format=md with additional query parameters",
			inputURL:            "/docs/page?format=md&version=1.0",
			expectedPath:        "/docs/page/index.md?version=1.0",
			expectedContentType: "text/markdown; charset=utf-8",
		},
		{
			name:                "format=md with multiple query parameters",
			inputURL:            "/api/endpoint?format=md&limit=10&offset=20",
			expectedPath:        "/api/endpoint/index.md?limit=10&offset=20",
			expectedContentType: "text/markdown; charset=utf-8",
		},
		{
			name:                "format=md combined with article parameter",
			inputURL:            "/best-practices/?format=md&article=design-learning",
			expectedPath:        "/best-practices/index.md?",
			expectedContentType: "text/markdown; charset=utf-8",
		},
		{
			name:                "format=md at root with trailing slash",
			inputURL:            "/?format=md",
			expectedPath:        "//index.md?",
			expectedContentType: "text/markdown; charset=utf-8",
		},
		{
			name:                "hyphenated path with format=md",
			inputURL:            "/key-concepts/?format=md",
			expectedPath:        "/key-concepts/index.md?",
			expectedContentType: "text/markdown; charset=utf-8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := RewriteMiddleware(&mockHandler{})
			req := httptest.NewRequest("GET", tt.inputURL, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			// Verify path rewrite
			if w.Body.String() != tt.expectedPath {
				t.Errorf("Path: expected %q, got %q", tt.expectedPath, w.Body.String())
			}

			// Verify Content-Type header
			contentType := w.Header().Get("Content-Type")
			if contentType != tt.expectedContentType {
				t.Errorf("Content-Type: expected %q, got %q", tt.expectedContentType, contentType)
			}
		})
	}
}

// TestRewriteMiddleware_FormatMd_PreservesOtherParams ensures that when format=md
// is processed, other query parameters are preserved in the correct order.
func TestRewriteMiddleware_FormatMd_PreservesOtherParams(t *testing.T) {
	handler := RewriteMiddleware(&mockHandler{})

	req := httptest.NewRequest("GET", "/docs/api?version=2.0&format=md&lang=en", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// format=md should be removed, others preserved
	result := w.Body.String()
	expectedPath := "/docs/api/index.md"
	
	if !strings.HasPrefix(result, expectedPath) {
		t.Errorf("Expected path to start with %q, got %q", expectedPath, result)
	}

	// Verify other params are still present (order may vary)
	if !strings.Contains(result, "version=2.0") {
		t.Error("Expected version=2.0 to be preserved")
	}
	if !strings.Contains(result, "lang=en") {
		t.Error("Expected lang=en to be preserved")
	}
	if strings.Contains(result, "format=md") {
		t.Error("format=md should be removed from query string")
	}
}

// TestRewriteMiddleware_FormatMd_WithoutFormatParam verifies that URLs without
// format=md are not affected by the format middleware.
func TestRewriteMiddleware_FormatMd_WithoutFormatParam(t *testing.T) {
	handler := RewriteMiddleware(&mockHandler{})

	tests := []struct {
		name     string
		inputURL string
		expected string
	}{
		{
			name:     "no query parameters",
			inputURL: "/docs/page",
			expected: "/docs/page?",
		},
		{
			name:     "other query parameters only",
			inputURL: "/docs/page?version=1.0&lang=en",
			expected: "/docs/page?version=1.0&lang=en",
		},
		{
			name:     "format with different value",
			inputURL: "/docs/page?format=html",
			expected: "/docs/page?format=html",
		},
		{
			name:     "format as part of other param value",
			inputURL: "/docs/page?data=format%3Dmd",
			expected: "/docs/page?data=format%3Dmd",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.inputURL, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Body.String() != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, w.Body.String())
			}

			// Should not set markdown Content-Type without format=md
			contentType := w.Header().Get("Content-Type")
			if contentType == "text/markdown; charset=utf-8" {
				t.Error("Should not set markdown Content-Type without format=md")
			}
		})
	}
}

// TestRewriteMiddleware_FormatEmbed_ComprehensiveTests provides comprehensive coverage
// for the ?format=embed query parameter handling across various URL patterns.
func TestRewriteMiddleware_FormatEmbed_ComprehensiveTests(t *testing.T) {
	tests := []struct {
		name                string
		inputURL            string
		expectedPath        string
		expectedContentType string
	}{
		{
			name:                "root path with format=embed",
			inputURL:            "/?format=embed",
			expectedPath:        "//embed.html?",
			expectedContentType: "text/html; charset=utf-8",
		},
		{
			name:                "single level path with format=embed",
			inputURL:            "/overview?format=embed",
			expectedPath:        "/overview/embed.html?",
			expectedContentType: "text/html; charset=utf-8",
		},
		{
			name:                "nested path with format=embed",
			inputURL:            "/reference/configuration?format=embed",
			expectedPath:        "/reference/configuration/embed.html?",
			expectedContentType: "text/html; charset=utf-8",
		},
		{
			name:                "deeply nested path with format=embed",
			inputURL:            "/reference/front-matter/user-roles?format=embed",
			expectedPath:        "/reference/front-matter/user-roles/embed.html?",
			expectedContentType: "text/html; charset=utf-8",
		},
		{
			name:                "path with trailing slash and format=embed",
			inputURL:            "/getting-started/?format=embed",
			expectedPath:        "/getting-started/embed.html?",
			expectedContentType: "text/html; charset=utf-8",
		},
		{
			name:                "format=embed with additional query parameters",
			inputURL:            "/docs/page?format=embed&version=1.0",
			expectedPath:        "/docs/page/embed.html?version=1.0",
			expectedContentType: "text/html; charset=utf-8",
		},
		{
			name:                "format=embed with multiple query parameters",
			inputURL:            "/api/endpoint?format=embed&limit=10&offset=20",
			expectedPath:        "/api/endpoint/embed.html?limit=10&offset=20",
			expectedContentType: "text/html; charset=utf-8",
		},
		{
			name:                "format=embed combined with article parameter",
			inputURL:            "/best-practices/?format=embed&article=design-learning",
			expectedPath:        "/best-practices/embed.html?",
			expectedContentType: "text/html; charset=utf-8",
		},
		{
			name:                "format=embed at root with trailing slash",
			inputURL:            "/?format=embed",
			expectedPath:        "//embed.html?",
			expectedContentType: "text/html; charset=utf-8",
		},
		{
			name:                "hyphenated path with format=embed",
			inputURL:            "/key-concepts/?format=embed",
			expectedPath:        "/key-concepts/embed.html?",
			expectedContentType: "text/html; charset=utf-8",
		},
		{
			name:                "format=embed for iframe embedding",
			inputURL:            "/tools/importers/oapi3?format=embed",
			expectedPath:        "/tools/importers/oapi3/embed.html?",
			expectedContentType: "text/html; charset=utf-8",
		},
		{
			name:                "format=embed with section parameter",
			inputURL:            "/overview/?format=embed",
			expectedPath:        "/overview/embed.html?",
			expectedContentType: "text/html; charset=utf-8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := RewriteMiddleware(&mockHandler{})
			req := httptest.NewRequest("GET", tt.inputURL, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			// Verify path rewrite
			if w.Body.String() != tt.expectedPath {
				t.Errorf("Path: expected %q, got %q", tt.expectedPath, w.Body.String())
			}

			// Verify Content-Type header
			contentType := w.Header().Get("Content-Type")
			if contentType != tt.expectedContentType {
				t.Errorf("Content-Type: expected %q, got %q", tt.expectedContentType, contentType)
			}
		})
	}
}

// TestRewriteMiddleware_FormatEmbed_PreservesOtherParams ensures that when format=embed
// is processed, other query parameters are preserved in the correct order.
func TestRewriteMiddleware_FormatEmbed_PreservesOtherParams(t *testing.T) {
	handler := RewriteMiddleware(&mockHandler{})

	req := httptest.NewRequest("GET", "/docs/widget?theme=dark&format=embed&lang=en", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// format=embed should be removed, others preserved
	result := w.Body.String()
	expectedPath := "/docs/widget/embed.html"
	
	if !strings.HasPrefix(result, expectedPath) {
		t.Errorf("Expected path to start with %q, got %q", expectedPath, result)
	}

	// Verify other params are still present (order may vary)
	if !strings.Contains(result, "theme=dark") {
		t.Error("Expected theme=dark to be preserved")
	}
	if !strings.Contains(result, "lang=en") {
		t.Error("Expected lang=en to be preserved")
	}
	if strings.Contains(result, "format=embed") {
		t.Error("format=embed should be removed from query string")
	}
}

// TestRewriteMiddleware_FormatEmbed_WithoutFormatParam verifies that URLs without
// format=embed are not affected by the embed middleware.
func TestRewriteMiddleware_FormatEmbed_WithoutFormatParam(t *testing.T) {
	handler := RewriteMiddleware(&mockHandler{})

	tests := []struct {
		name     string
		inputURL string
		expected string
	}{
		{
			name:     "no query parameters",
			inputURL: "/docs/page",
			expected: "/docs/page?",
		},
		{
			name:     "other query parameters only",
			inputURL: "/docs/page?version=1.0&lang=en",
			expected: "/docs/page?version=1.0&lang=en",
		},
		{
			name:     "format with different value",
			inputURL: "/docs/page?format=json",
			expected: "/docs/page?format=json",
		},
		{
			name:     "format as part of other param value",
			inputURL: "/docs/page?data=format%3Dembed",
			expected: "/docs/page?data=format%3Dembed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.inputURL, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Body.String() != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, w.Body.String())
			}

			// Should not set HTML Content-Type without format=embed
			contentType := w.Header().Get("Content-Type")
			if contentType == "text/html; charset=utf-8" {
				t.Error("Should not set HTML Content-Type without format=embed")
			}
		})
	}
}

// TestRewriteMiddleware_FormatEmbed_VsFormatMd verifies that format=embed and
// format=md produce different outputs and cannot be combined.
func TestRewriteMiddleware_FormatEmbed_VsFormatMd(t *testing.T) {
	handler := RewriteMiddleware(&mockHandler{})

	// Test format=embed
	req1 := httptest.NewRequest("GET", "/docs/page?format=embed", nil)
	w1 := httptest.NewRecorder()
	handler.ServeHTTP(w1, req1)

	// Test format=md
	req2 := httptest.NewRequest("GET", "/docs/page?format=md", nil)
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, req2)

	// Should produce different paths
	embedPath := w1.Body.String()
	mdPath := w2.Body.String()

	if embedPath == mdPath {
		t.Errorf("format=embed and format=md should produce different results, both got %q", embedPath)
	}

	expectedEmbed := "/docs/page/embed.html?"
	expectedMd := "/docs/page/index.md?"

	if embedPath != expectedEmbed {
		t.Errorf("format=embed: expected %q, got %q", expectedEmbed, embedPath)
	}

	if mdPath != expectedMd {
		t.Errorf("format=md: expected %q, got %q", expectedMd, mdPath)
	}

	// Verify different Content-Types
	embedContentType := w1.Header().Get("Content-Type")
	mdContentType := w2.Header().Get("Content-Type")

	if embedContentType != "text/html; charset=utf-8" {
		t.Errorf("format=embed Content-Type: expected text/html, got %q", embedContentType)
	}

	if mdContentType != "text/markdown; charset=utf-8" {
		t.Errorf("format=md Content-Type: expected text/markdown, got %q", mdContentType)
	}
}

// TestRewriteMiddleware_FormatEmbed_MultipleTrailingSlashes verifies handling
// of URLs with multiple trailing slashes.
func TestRewriteMiddleware_FormatEmbed_MultipleTrailingSlashes(t *testing.T) {
	handler := RewriteMiddleware(&mockHandler{})

	tests := []struct {
		name     string
		inputURL string
		expected string
	}{
		{
			name:     "double trailing slash",
			inputURL: "/docs/page//?format=embed",
			expected: "/docs/page/embed.html?",
		},
		{
			name:     "triple trailing slash",
			inputURL: "/docs/page///?format=embed",
			expected: "/docs/page/embed.html?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.inputURL, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Body.String() != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, w.Body.String())
			}
		})
	}
}
