
# Presidium Proxy Middleware

The Presidium proxy middleware provides Caddy-style URL rewriting capabilities for the development server, enabling advanced features without requiring external dependencies.

## Architecture

When you run `presidium server`, the following happens:

1. **Hugo Server Start**: Hugo's built-in server starts on an auto-selected internal port (typically 1313 or next available)
2. **Proxy Server Start**: A reverse proxy server starts on your specified port (default: 3131)
3. **Request Flow**: Incoming requests → Middleware Chain → Hugo Server
4. **Graceful Shutdown**: Both servers shut down cleanly on interrupt (Ctrl+C)

## Middleware Chain

Requests pass through the following middleware in order:

```
Client Request
    ↓
[NormalizePathMiddleware] - Clean up paths (remove //, ./, etc.)
    ↓
[RewriteMiddleware] - Handle query parameter rewrites
    ↓
[ReverseProxy] - Forward to Hugo server (handles WebSocket upgrades automatically)
    ↓
Hugo Response → Client
```

Note: WebSocket support is built into the reverse proxy and does not require a separate middleware.

## Rewrite Rules

### 1. Article Navigation (`?article=<id>`)

The `?article=<id>` parameter is passed through to Hugo unchanged — no redirect, no path rewrite. When Hugo returns the HTML response, the proxy injects a small `<script>` just before `</body>` that reads the parameter from `window.location.search` and scrolls the matching element (`document.getElementById(id)`, falling back to `document.getElementsByName(id)`) into view on load.

**Example:**

```
Request:  GET /docs/page?article=my-section
Forward:  GET /docs/page?article=my-section (unchanged)
Browser:  /docs/page?article=my-section (URL preserved), page scrolls to #my-section
```

The injection only happens for `text/html` responses where the request actually carries `?article=`, so other pages are unaffected.

**Use Case:** Deep-link to a specific element on a page without converting the link into a fragment URL via a server-side redirect. The original `?article=` URL stays in the address bar (useful when sharing or processing the URL programmatically).

### 2. Section Navigation (`?section=<section>`)

Locates the section slug inside the request path and masks the path to everything up to (and including) that slug. Any path segments after the slug are promoted to `?article=<remainder>` so the response-injected anchor-scroll script (see section 1) targets that element.

This is equivalent to: *"serve the section page, but scroll to the article that the original URL pointed at."* The browser's address bar is preserved.

**Example:**

```
Request:  GET /docs/foo/bar?section=foo
Forward:  GET /docs/foo/?article=bar (path masked, remainder promoted)
Browser:  /docs/foo/bar?section=foo (unchanged), page scrolls to #bar
```

**Edge cases:**

| Input                                  | Behavior                                                                  |
|----------------------------------------|---------------------------------------------------------------------------|
| Slug not in path                       | `?section=` is stripped and the request passes through unchanged.         |
| Slug appears multiple times in path    | Last (deepest) occurrence wins.                                           |
| Slug is the final path segment         | Path is masked to `/.../slug/`; no `?article=` is promoted (no scroll).   |
| Slug fails `^[a-z0-9_-]+$` validation  | `400 Bad Request` (path-traversal guard).                                 |

**Use Case:** Deep-link to a specific article on a section index page without the URL collapsing into the section's canonical path.

### 3. Markdown Export (`?format=md`)

Rewrites to the pre-built markdown output file.

**Example:**

```
Request:  GET /docs/page?format=md
Rewrite:  GET /docs/page/index.md
Headers:  Content-Type: text/markdown; charset=utf-8
```

**Use Case:** Export documentation as markdown for external tools or editing.

### 4. Embed Format (`?format=embed`)

Rewrites to the embed-optimized HTML output.

**Example:**

```
Request:  GET /docs/page?format=embed
Rewrite:  GET /docs/page/embed.html
Headers:  Content-Type: text/html; charset=utf-8
```

**Use Case:** Embed documentation pages in iframes or external sites.

## Usage

### Basic Server

Start with default settings (proxy on port 3131):

```bash
presidium server
```

### Custom Port

Run proxy on a different port:

```bash
presidium server --port 8080
```

### Disable Proxy

Run Hugo server directly without middleware:

```bash
presidium server --no-proxy
```

### Pass Hugo Flags

Additional flags are passed through to Hugo. The `--port` and `--no-proxy` flags are
handled by Presidium; all other flags go directly to Hugo:

```bash
# Pass Hugo flags (Presidium parses only --port and --no-proxy)
presidium server --buildDrafts --buildFuture

# Or be explicit with -- separator (both work the same)
presidium server -- --buildDrafts --buildFuture --baseURL http://example.com
```

## WebSocket Support

The proxy automatically handles WebSocket upgrade requests for Hugo's live-reload feature:

- Detects `Connection: Upgrade` and `Upgrade: websocket` headers
- Preserves upgrade headers through the proxy
- Maintains bidirectional communication for live reload

No special configuration needed - live reload works out of the box.

## Equivalent Caddy Configuration

This middleware replicates the following Caddyfile:

```caddy
:3131 {
 # ?article=<id> is passed through unchanged. The proxy injects an
 # anchor-scroll <script> into the HTML response based on resp.Request URL.
 # ?section=<slug> requires path-aware logic (locate the slug, mask the
 # path, promote the remainder to ?article=) that goes beyond a simple
 # Caddy rewrite — implemented in Go in RewriteMiddleware.

 @format_md {
  query format=md
  path_regexp md_path ^(.+?)/*$
 }
 rewrite @format_md {re.md_path.1}/index.md
 header @format_md Content-Type "text/markdown; charset=utf-8"

 @format_embed {
  query format=embed
  path_regexp embed_path ^(.+?)/*$
 }
 rewrite @format_embed {re.embed_path.1}/embed.html
 header @format_embed Content-Type "text/html; charset=utf-8"

 reverse_proxy localhost:{env.HUGO_PORT}
}
```

## Testing

The middleware includes comprehensive unit tests:

```bash
go test ./pkg/domain/service/proxy/... -v
```

Tests cover:

- Article parameter pass-through
- Section navigation rewriting
- Format conversion (md/embed)
- Path normalization
- Middleware chaining

## Implementation Details

### Path Normalization

Before rewriting, paths are cleaned to:

- Remove duplicate slashes (`//` → `/`)
- Resolve relative references (`./`, `../`)
- Ensure consistent formatting

### Trailing Slash Handling

Format conversions handle trailing slashes gracefully:

```
/docs/page/  → /docs/page/index.md
/docs/page   → /docs/page/index.md
```

### Port Selection

Hugo's internal port is auto-selected:

- Starts with port 1313 (Hugo's default)
- If occupied, tries 1314, 1315, etc.
- Continues until an available port is found

### Startup Sequence

1. Themes extracted and Hugo module replacements configured
2. Hugo server started asynchronously with selected port
3. Proxy polls Hugo's port until ready (timeout: 30s)
4. Proxy server starts listening on configured port
5. Both servers log their status

### Error Handling

- Hugo startup failures are captured and logged
- Proxy startup errors terminate gracefully
- Shutdown timeout: 10 seconds for cleanup
- Signal handling: SIGINT, SIGTERM

## Development

To modify the middleware:

1. **Edit middleware logic**: `pkg/domain/service/proxy/middleware.go`
2. **Add tests**: `pkg/domain/service/proxy/middleware_test.go`
3. **Run tests**: `go test ./pkg/domain/service/proxy/... -v`
4. **Build**: `make build`

### Adding New Rewrite Rules

To add a new query parameter handler:

```go
// In middleware.go, add to RewriteMiddleware:

if query.Get("yourparam") == "value" {
 // Modify r.URL.Path as needed
 r.URL.Path = "/your/new/path"
 query.Del("yourparam")
 r.URL.RawQuery = query.Encode()
 
 // Set headers if needed
 w.Header().Set("Your-Header", "value")
}
```

Then add tests in `middleware_test.go`:

```go
func TestRewriteMiddleware_YourParam(t *testing.T) {
 handler := RewriteMiddleware(&mockHandler{})
 req := httptest.NewRequest("GET", "/path?yourparam=value", nil)
 w := httptest.NewRecorder()
 
 handler.ServeHTTP(w, req)
 
 expected := "/your/new/path?"
 if w.Body.String() != expected {
  t.Errorf("Expected %q, got %q", expected, w.Body.String())
 }
}
```

## Performance

- **Middleware overhead**: <1ms per request
- **Hugo startup**: 1-3 seconds (one-time)
- **Proxy startup**: <100ms (one-time)
- **WebSocket latency**: No measurable impact on live-reload

## Troubleshooting

### "Timeout waiting for Hugo server to start"

- Check if port 1313+ is blocked by firewall
- Verify Hugo builds successfully: `presidium hugo version`
- Try with `--no-proxy` to test Hugo directly

### "Port already in use"

- Change proxy port: `presidium server --port 8080`
- Kill existing process: `lsof -ti:3131 | xargs kill`

### Live reload not working

- Check browser console for WebSocket errors
- Verify Hugo server started successfully
- Try disabling browser extensions that block WebSockets

### Rewrite not working as expected

- Check the request path format (trailing slashes matter)
- Verify query parameter spelling
- Enable debug mode: `presidium server --debug`
