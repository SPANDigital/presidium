package proxy

import (
	"bytes"
	"encoding/json"
	"net/http"
	"path"
	"regexp"
	"strings"
)

// articleAnchorScript builds an HTML <script> tag that scrolls the page to
// the element whose id matches target on load. The target is baked into the
// script literal (JSON-encoded so any value is safely escaped) so the script
// works regardless of what the browser's address bar shows — necessary for
// the ?section= case, where the article id is derived server-side and never
// appears in the visible URL.
func articleAnchorScript(target string) string {
	encoded, _ := json.Marshal(target)
	return `<script>(function(){var id=` + string(encoded) + `;function go(){var el=document.getElementById(id);if(!el){var n=document.getElementsByName(id);if(n&&n.length)el=n[0];}if(el)el.scrollIntoView();}if(document.readyState==='complete'||document.readyState==='interactive'){go();}else{document.addEventListener('DOMContentLoaded',go);}})();</script>`
}

// injectArticleAnchorScript returns body with a parameterized anchor-scroll
// script for target inserted just before the last </body> tag. If no </body>
// tag is present the script is appended at the end — browsers still execute
// it. If target is empty the body is returned unchanged.
func injectArticleAnchorScript(body []byte, target string) []byte {
	if target == "" {
		return body
	}
	script := []byte(articleAnchorScript(target))
	closingBody := []byte("</body>")
	idx := bytes.LastIndex(body, closingBody)
	if idx == -1 {
		return append(body, script...)
	}
	out := make([]byte, 0, len(body)+len(script))
	out = append(out, body[:idx]...)
	out = append(out, script...)
	out = append(out, body[idx:]...)
	return out
}

// RewriteMiddleware implements URL rewriting logic
// that handles special query parameters for section navigation
// and format conversions. The ?article=<id> parameter is passed
// through unchanged for downstream handling.
func RewriteMiddleware(next http.Handler) http.Handler {
	// Validate section parameter - only allow alphanumeric, hyphens, and underscores
	// This prevents path traversal attacks with values like "../", "./", or URL-encoded separators
	sectionPattern := regexp.MustCompile(`^[a-z0-9_-]+$`)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		// Handle ?section=<section>
		// Find the section slug in the request path and mask-rewrite the path
		// to everything up to (and including) that slug. Anything after the
		// slug becomes the article id — promoted to ?article=<remainder> so
		// the response-injected script scrolls to that element on load.
		// e.g. /docs/foo/bar?section=foo → mask to /docs/foo/?article=bar
		//
		// Edge cases:
		//   - slug not in path: strip ?section= and pass through unchanged
		//   - slug appears multiple times: last occurrence wins (deepest match)
		//   - slug is the final segment (no remainder): mask the path only
		if section := query.Get("section"); section != "" {
			// Validate section to prevent path traversal
			if !sectionPattern.MatchString(section) {
				http.Error(w, "Invalid section parameter", http.StatusBadRequest)
				return
			}
			trimmed := strings.Trim(r.URL.Path, "/")
			var segments []string
			if trimmed != "" {
				segments = strings.Split(trimmed, "/")
			}
			sectionIdx := -1
			for i, seg := range segments {
				if seg == section {
					sectionIdx = i
				}
			}
			query.Del("section")
			if sectionIdx >= 0 {
				r.URL.Path = "/" + strings.Join(segments[:sectionIdx+1], "/") + "/"
				if remainder := strings.Join(segments[sectionIdx+1:], "/"); remainder != "" {
					query.Set("article", remainder)
				}
			}
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
