package proxy

import (
	"net/http"
	"path"
	"regexp"
	"strings"
)

// RewriteMiddleware implements URL rewriting logic
// that handles special query parameters for article navigation, sections,
// and format conversions.
func RewriteMiddleware(next http.Handler) http.Handler {
	// Validate section parameter - only allow alphanumeric, hyphens, and underscores
	// This prevents path traversal attacks with values like "../", "./", or URL-encoded separators
	sectionPattern := regexp.MustCompile(`^[a-z0-9_-]+$`)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		// Handle ?article=<id>
		// Convert to fragment identifier (#) via redirect so browser scrolls to the element.
		// e.g., /docs/page?article=my-section → redirect to /docs/page#my-section
		if articleID := query.Get("article"); articleID != "" {
			// Build the redirect URL with fragment identifier
			query.Del("article")
			redirectURL := r.URL.Path
			if len(query) > 0 {
				redirectURL += "?" + query.Encode()
			}
			redirectURL += "#" + articleID

			// Issue 302 redirect to the URL with fragment
			http.Redirect(w, r, redirectURL, http.StatusFound)
			return
		}

		// Handle ?section=<section>
		// Extract the last path segment as the article id, rewrite to /<section>/.
		// e.g. /docs/my-module/some-article/?section=my-module → /my-module/
		if section := query.Get("section"); section != "" {
			// Validate section to prevent path traversal
			if !sectionPattern.MatchString(section) {
				http.Error(w, "Invalid section parameter", http.StatusBadRequest)
				return
			}
			r.URL.Path = "/" + section + "/"
			// Apply path.Clean to ensure no relative references remain
			r.URL.Path = path.Clean(r.URL.Path)
			if !strings.HasPrefix(r.URL.Path, "/") {
				r.URL.Path = "/" + r.URL.Path
			}
			if !strings.HasSuffix(r.URL.Path, "/") {
				r.URL.Path += "/"
			}
			query.Del("section")
			r.URL.RawQuery = query.Encode()
		}

		// Handle ?format=md
		// Rewrite to the pre-built index.md Hugo output.
		// Note: Content-Type headers for format conversions are set in the
		// reverse proxy's ModifyResponse based on file extension to avoid duplicate headers.
		if query.Get("format") == "md" {
			// Remove trailing slashes from path
			cleanPath := strings.TrimRight(r.URL.Path, "/")
			if cleanPath == "" {
				cleanPath = "/"
			}
			r.URL.Path = cleanPath + "/index.md"
			query.Del("format")
			r.URL.RawQuery = query.Encode()
		}

		// Handle ?format=embed
		// Rewrite to the pre-built embed.html Hugo output.
		// Note: Content-Type headers for format conversions are set in the
		// reverse proxy's ModifyResponse based on file extension to avoid duplicate headers.
		if query.Get("format") == "embed" {
			// Remove trailing slashes from path
			cleanPath := strings.TrimRight(r.URL.Path, "/")
			if cleanPath == "" {
				cleanPath = "/"
			}
			r.URL.Path = cleanPath + "/embed.html"
			query.Del("format")
			r.URL.RawQuery = query.Encode()
		}

		// Pass to next handler (Caddy reverse proxy)
		next.ServeHTTP(w, r)
	})
}

// NormalizePathMiddleware ensures paths are properly formatted before rewriting
func NormalizePathMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Preserve trailing slash before cleaning — Hugo uses directory-based URLs
		// that require trailing slashes. Stripping them causes Hugo to issue a
		// relative redirect (e.g. "getting-started/") which browsers resolve
		// relative to the current path, producing doubled segments.
		hadTrailingSlash := strings.HasSuffix(r.URL.Path, "/")
		r.URL.Path = path.Clean(r.URL.Path)
		if r.URL.Path == "." {
			r.URL.Path = "/"
		}
		if hadTrailingSlash && r.URL.Path != "/" {
			r.URL.Path += "/"
		}
		next.ServeHTTP(w, r)
	})
}

// ChainMiddleware chains multiple middleware functions together
func ChainMiddleware(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	// Apply middlewares in reverse order so they execute in the order specified
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}
