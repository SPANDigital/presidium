package proxy

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// mockHandler is a simple handler that returns the request path and query
type mockHandler struct{}

func (h *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(r.URL.Path + "?" + r.URL.RawQuery))
}

func TestRewriteMiddleware_ArticleParam(t *testing.T) {
	handler := RewriteMiddleware(&mockHandler{})

	req := httptest.NewRequest("GET", "/docs/page?article=my-section", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Should pass through unchanged — no redirect, article param preserved
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	expected := "/docs/page?article=my-section"
	if w.Body.String() != expected {
		t.Errorf("Expected %q, got %q", expected, w.Body.String())
	}
}

func TestRewriteMiddleware_ArticleParam_WithOtherQueryParams(t *testing.T) {
	handler := RewriteMiddleware(&mockHandler{})

	req := httptest.NewRequest("GET", "/docs/page?article=my-section&version=2.0", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Should pass through unchanged, preserving all query params and their order
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	expected := "/docs/page?article=my-section&version=2.0"
	if w.Body.String() != expected {
		t.Errorf("Expected %q, got %q", expected, w.Body.String())
	}
}

func TestRewriteMiddleware_SectionParam(t *testing.T) {
	handler := RewriteMiddleware(&mockHandler{})

	req := httptest.NewRequest("GET", "/docs/my-module/some-article?section=my-module", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Path is masked to /<...up-to-section>/ and the remainder is promoted
	// to ?article=<remainder> so the response-injected script scrolls to it.
	expected := "/docs/my-module/?article=some-article"
	if w.Body.String() != expected {
		t.Errorf("Expected %q, got %q", expected, w.Body.String())
	}
}

func TestRewriteMiddleware_SectionParam_SlugNotInPath(t *testing.T) {
	handler := RewriteMiddleware(&mockHandler{})

	req := httptest.NewRequest("GET", "/docs/other/page?section=foo", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Slug not in path → strip ?section= and pass through unchanged.
	expected := "/docs/other/page?"
	if w.Body.String() != expected {
		t.Errorf("Expected %q, got %q", expected, w.Body.String())
	}
}

func TestRewriteMiddleware_SectionParam_LastOccurrenceWins(t *testing.T) {
	handler := RewriteMiddleware(&mockHandler{})

	req := httptest.NewRequest("GET", "/docs/foo/foo/bar?section=foo", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Slug appears twice → last occurrence wins (deepest match).
	expected := "/docs/foo/foo/?article=bar"
	if w.Body.String() != expected {
		t.Errorf("Expected %q, got %q", expected, w.Body.String())
	}
}

func TestRewriteMiddleware_SectionParam_NoRemainder(t *testing.T) {
	handler := RewriteMiddleware(&mockHandler{})

	req := httptest.NewRequest("GET", "/docs/foo?section=foo", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Slug is the final segment → mask the path only, no ?article= injected.
	expected := "/docs/foo/?"
	if w.Body.String() != expected {
		t.Errorf("Expected %q, got %q", expected, w.Body.String())
	}
}

func TestRewriteMiddleware_SectionParam_PathTraversalRejected(t *testing.T) {
	handler := RewriteMiddleware(&mockHandler{})

	req := httptest.NewRequest("GET", "/docs/foo?section=../etc/passwd", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for invalid section, got %d", http.StatusBadRequest, w.Code)
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

	// Note: Content-Type headers are set by the proxy's ModifyResponse, not the middleware
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

	// Note: Content-Type headers are set by the proxy's ModifyResponse, not the middleware
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

func TestInjectArticleAnchorScript_BeforeClosingBody(t *testing.T) {
	body := []byte("<html><body><h1>Hi</h1></body></html>")
	script := articleAnchorScript("my-section")
	out := injectArticleAnchorScript(body, "my-section")

	if !bytes.Contains(out, []byte(script)) {
		t.Fatalf("expected output to contain the anchor script")
	}

	scriptIdx := bytes.Index(out, []byte(script))
	closeIdx := bytes.Index(out, []byte("</body>"))
	if scriptIdx == -1 || closeIdx == -1 || scriptIdx > closeIdx {
		t.Errorf("expected script to appear before </body>, got scriptIdx=%d closeIdx=%d", scriptIdx, closeIdx)
	}

	if !bytes.Contains(out, []byte("<h1>Hi</h1>")) {
		t.Errorf("expected original content to be preserved")
	}
}

func TestInjectArticleAnchorScript_InsertsBeforeLastClosingBody(t *testing.T) {
	// Stray "</body>" inside a string in the page should not be the insertion
	// point — we use LastIndex so the real closing tag wins.
	body := []byte(`<html><body><pre>echo "</body>"</pre></body></html>`)
	out := injectArticleAnchorScript(body, "anchor")

	// The script must appear after the literal "</body>" in the <pre>
	// (i.e. the LastIndex match was the real closing tag).
	preIdx := bytes.Index(out, []byte(`echo "</body>"`))
	scriptIdx := bytes.Index(out, []byte(articleAnchorScript("anchor")))
	if scriptIdx <= preIdx {
		t.Errorf("expected script to be injected at the last </body>, not the one inside <pre>")
	}
}

func TestInjectArticleAnchorScript_NoClosingBodyAppends(t *testing.T) {
	body := []byte("<html><h1>Hi</h1>")
	script := articleAnchorScript("anchor")
	out := injectArticleAnchorScript(body, "anchor")

	if !bytes.HasSuffix(out, []byte(script)) {
		t.Errorf("expected script to be appended when </body> is missing")
	}
	if !bytes.HasPrefix(out, body) {
		t.Errorf("expected original content to remain at the start")
	}
}

func TestInjectArticleAnchorScript_EmptyTargetSkipsInjection(t *testing.T) {
	body := []byte("<html><body></body></html>")
	out := injectArticleAnchorScript(body, "")

	if !bytes.Equal(out, body) {
		t.Errorf("expected body to be returned unchanged when target is empty")
	}
}

func TestArticleAnchorScript_BakesTargetIntoScript(t *testing.T) {
	script := articleAnchorScript("my-section")

	// The id must be baked into the script as a JSON-encoded string literal.
	if !strings.Contains(script, `var id="my-section"`) {
		t.Errorf("expected script to bake the target id as a JSON string literal, got: %s", script)
	}
}

func TestArticleAnchorScript_EscapesSpecialCharacters(t *testing.T) {
	// JSON encoding must escape characters that would break out of the
	// string literal (quotes, backslashes, </script> via < escape isn't
	// strictly required for our case but the JSON encoder is conservative).
	script := articleAnchorScript(`evil"; alert(1); var x="`)

	if strings.Contains(script, `evil";`) {
		t.Errorf("expected dangerous characters to be escaped, got: %s", script)
	}
	// JSON encodes a quote as \" — the escaped form should appear instead.
	if !strings.Contains(script, `\"`) {
		t.Errorf("expected escaped quote in script, got: %s", script)
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
