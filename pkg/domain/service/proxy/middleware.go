package proxy

import (
	"net/http"
	"path"
	"strings"
)

// RewriteMiddleware implements URL rewriting logic
// that handles special query parameters for article navigation, sections,
// and format conversions.
func RewriteMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		// Handle ?article=<id>
		// Rewrite to the same path, stripping the query param.
		// Browser URL stays unchanged via proxy; Hugo sees only the clean path.
		if query.Has("article") {
			// Remove the article query parameter
			query.Del("article")
			r.URL.RawQuery = query.Encode()
		}

		// Handle ?section=<section>
		// Extract the last path segment as the article id, rewrite to /<section>/.
		// e.g. /docs/my-module/some-article/?section=my-module → /my-module/
		if section := query.Get("section"); section != "" {
			r.URL.Path = "/" + section + "/"
			query.Del("section")
			r.URL.RawQuery = query.Encode()
		}

		// Handle ?format=md
		// Rewrite to the pre-built index.md Hugo output.
		if query.Get("format") == "md" {
			// Remove trailing slashes from path
			cleanPath := strings.TrimRight(r.URL.Path, "/")
			if cleanPath == "" {
				cleanPath = "/"
			}
			r.URL.Path = cleanPath + "/index.md"
			query.Del("format")
			r.URL.RawQuery = query.Encode()
			w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
		}

		// Handle ?format=embed
		// Rewrite to the pre-built embed.html Hugo output.
		if query.Get("format") == "embed" {
			// Remove trailing slashes from path
			cleanPath := strings.TrimRight(r.URL.Path, "/")
			if cleanPath == "" {
				cleanPath = "/"
			}
			r.URL.Path = cleanPath + "/embed.html"
			query.Del("format")
			r.URL.RawQuery = query.Encode()
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
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
