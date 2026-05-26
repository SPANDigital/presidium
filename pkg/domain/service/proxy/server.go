package proxy

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/SPANDigital/presidium-hugo/pkg/domain/service/hugo"
	"github.com/SPANDigital/presidium-hugo/pkg/log"
)

// Server wraps Hugo server with a reverse proxy
// Uses httputil.ReverseProxy from stdlib - production-grade, used by Caddy internally
type Server struct {
	hugoPort   int
	proxyPort  int
	httpServer *http.Server
	hugoDone   <-chan error
}

// NewServer creates a new proxy server instance.
func NewServer(proxyPort int) *Server {
	return &Server{
		proxyPort: proxyPort,
	}
}

// Start initializes Hugo server and reverse proxy
func (s *Server) Start(hugoArgs []string) error {
	// If no proxy port was specified, pick a random pair N (Hugo) / N+1 (proxy).
	// Otherwise honor the explicit proxy port and pick Hugo independently.
	if s.proxyPort == 0 {
		hugoPort, proxyPort, err := pickRandomPortPair()
		if err != nil {
			return fmt.Errorf("failed to find a free port pair: %w", err)
		}
		s.hugoPort = hugoPort
		s.proxyPort = proxyPort
	} else {
		if err := checkPortAvailable(s.proxyPort); err != nil {
			return fmt.Errorf("proxy port %d is already in use — stop the existing process or omit --port to auto-select", s.proxyPort)
		}
		hugoPort, err := pickRandomHighPort()
		if err != nil {
			return fmt.Errorf("failed to find a free Hugo port: %w", err)
		}
		s.hugoPort = hugoPort
	}

	// Prepare Hugo with themes
	hugoService := hugo.New()

	// Start Hugo server on internal port in a goroutine
	hugoDone := make(chan error, 1)
	s.hugoDone = hugoDone
	go func() {
		// Check if user provided --baseURL flag
		hasBaseURL := false
		for _, arg := range hugoArgs {
			if arg == "--baseURL" || strings.HasPrefix(arg, "--baseURL=") {
				hasBaseURL = true
				break
			}
		}

		// Add port flag (Hugo always needs to know which port to use)
		args := append(hugoArgs, "--port", strconv.Itoa(s.hugoPort))

		// Only add baseURL if user hasn't provided one
		// When proxy is enabled, baseURL should match the proxy port for correct absolute URLs
		if !hasBaseURL {
			args = append(args, "--baseURL", fmt.Sprintf("http://localhost:%d/", s.proxyPort))
		}

		log.Info(fmt.Sprintf("Starting Hugo server on port %d...", s.hugoPort))
		err := hugoService.Execute(args...)
		hugoDone <- err
	}()

	// Wait for Hugo server to start
	if err := s.waitForHugo(); err != nil {
		return fmt.Errorf("failed to start Hugo server: %w", err)
	}

	log.Info(fmt.Sprintf("Hugo server ready on port %d", s.hugoPort))

	// Create reverse proxy to Hugo
	// httputil.ReverseProxy is production-grade: used by Caddy, handles HTTP/2, WebSockets, etc.
	hugoURL, err := url.Parse(fmt.Sprintf("http://localhost:%d", s.hugoPort))
	if err != nil {
		return fmt.Errorf("failed to parse Hugo URL: %w", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(hugoURL)

	// Customize Director to preserve the original path exactly
	// The default Director can add unwanted redirects
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		// Override the target to ensure we use the exact original path
		req.URL.Host = hugoURL.Host
		req.URL.Scheme = hugoURL.Scheme
		req.Host = hugoURL.Host
		// Don't modify the path - keep it exactly as received
	}

	// Enhance the default transport for better performance.
	// DisableCompression: ModifyResponse rewrites HTML/JS bodies by byte-replacing
	// Hugo's internal port with the proxy port. If Hugo returns gzip/br-encoded
	// content the byte replacement corrupts the response, so we always negotiate
	// an uncompressed body from the upstream.
	proxy.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true, // Enable HTTP/2
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DisableCompression:    true,
	}

	// Fix Hugo's URLs to use proxy port instead of Hugo's internal port
	proxy.ModifyResponse = func(resp *http.Response) error {
		// Set Content-Type for format conversions based on file extension
		// This overrides any Content-Type from Hugo and avoids duplicate headers
		if resp.Request != nil && resp.StatusCode == 200 {
			path := resp.Request.URL.Path
			if strings.HasSuffix(path, "/index.md") {
				resp.Header.Set("Content-Type", "text/markdown; charset=utf-8")
			} else if strings.HasSuffix(path, "/embed.html") {
				resp.Header.Set("Content-Type", "text/html; charset=utf-8")
			}
		}

		// Handle redirects
		if resp.StatusCode >= 300 && resp.StatusCode < 400 {
			location := resp.Header.Get("Location")
			if location != "" {
				// Parse location URL
				locationURL, err := url.Parse(location)
				if err == nil && !locationURL.IsAbs() {
					// Make it absolute using the proxy's scheme and host
					locationURL.Scheme = "http"
					locationURL.Host = fmt.Sprintf("localhost:%d", s.proxyPort)
					resp.Header.Set("Location", locationURL.String())
				}
			}
		}

		// Rewrite HTML content to replace Hugo's internal port with proxy port
		contentType := resp.Header.Get("Content-Type")
		if resp.StatusCode == 200 && (strings.Contains(contentType, "text/html") || strings.Contains(contentType, "application/javascript")) {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			resp.Body.Close()

			// Replace Hugo's internal port with proxy port in the HTML/JS
			hugoURL := fmt.Sprintf("http://localhost:%d", s.hugoPort)
			proxyURL := fmt.Sprintf("http://localhost:%d", s.proxyPort)
			modifiedBody := bytes.ReplaceAll(body, []byte(hugoURL), []byte(proxyURL))

			// Also fix WebSocket port in livereload script
			wsOldPort := fmt.Sprintf("port=%d", s.hugoPort)
			wsNewPort := fmt.Sprintf("port=%d", s.proxyPort)
			modifiedBody = bytes.ReplaceAll(modifiedBody, []byte(wsOldPort), []byte(wsNewPort))

			// Inject the article-anchor scroll script into HTML responses
			// where the request carried ?article=<id> (either passed in
			// directly or promoted from ?section= by RewriteMiddleware).
			// The id is baked into the script so it works even when the
			// browser's visible URL doesn't contain ?article= (the
			// ?section= case is masked, so the browser keeps showing the
			// original ?section= URL).
			if strings.Contains(contentType, "text/html") && resp.Request != nil {
				if article := resp.Request.URL.Query().Get("article"); article != "" {
					modifiedBody = injectArticleAnchorScript(modifiedBody, article)
				}
			}

			// Update Content-Length header
			resp.Body = io.NopCloser(bytes.NewReader(modifiedBody))
			resp.Header.Set("Content-Length", strconv.Itoa(len(modifiedBody)))
			resp.ContentLength = int64(len(modifiedBody))
		}

		return nil
	}

	// Wrap proxy with our custom middleware for query parameter rewrites
	handler := ChainMiddleware(
		proxy,
		NormalizePathMiddleware,
		RewriteMiddleware,
	)

	// Create HTTP server
	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.proxyPort),
		Handler: handler,
	}

	log.Info(fmt.Sprintf("Starting Presidium proxy server on port %d...", s.proxyPort))
	log.Info(fmt.Sprintf("Proxying requests to Hugo on port %d (HTTP/2, WebSockets enabled)", s.hugoPort))
	log.Info(fmt.Sprintf("Presidium is ready — open http://localhost:%d/ in your browser", s.proxyPort))

	// Start serving (this blocks)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("proxy server error: %w", err)
	}

	return nil
}

// Shutdown gracefully shuts down the proxy server and waits for Hugo to exit.
// Hugo receives the same OS signal that triggered shutdown and cleans up on its
// own; we just wait to ensure its deferred tmpDir removal completes before the
// process exits.
func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown proxy server: %w", err)
		}
	}
	if s.hugoDone != nil {
		select {
		case err := <-s.hugoDone:
			if err != nil {
				log.Warn(fmt.Sprintf("Hugo server stopped with error: %v", err))
			}
		case <-ctx.Done():
			log.Warn("Hugo server did not stop within shutdown timeout")
		}
	}
	return nil
}

// waitForHugo polls until Hugo is ready, then verifies the listener on s.hugoPort is actually a Hugo dev server by fetching /livereload.js — a Hugo-specific endpoint that non-Hugo servers won't have.
func (s *Server) waitForHugo() error {
	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	httpClient := &http.Client{Timeout: 2 * time.Second}
	livereloadURL := fmt.Sprintf("http://localhost:%d/livereload.js", s.hugoPort)

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout waiting for Hugo server to start")
		case err := <-s.hugoDone:
			if err != nil {
				return fmt.Errorf("hugo server failed to start: %w", err)
			}
			return fmt.Errorf("hugo server exited unexpectedly before becoming ready")
		case <-ticker.C:
			conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", s.hugoPort), time.Second)
			if err != nil {
				continue
			}
			conn.Close()

			resp, err := httpClient.Get(livereloadURL)
			if err != nil {
				continue
			}
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			if resp.StatusCode != 200 || !bytes.Contains(body, []byte("LiveReload")) {
				return fmt.Errorf("port %d is serving HTTP but does not appear to be a Hugo dev server (no /livereload.js — likely a port collision)", s.hugoPort)
			}
			return nil
		}
	}
}

// pickRandomHighPort returns a free port in a wide high range; concurrent-instance collisions are caught by identity verification in waitForHugo.
func pickRandomHighPort() (int, error) {
	const (
		rangeStart = 20000
		rangeEnd   = 30000
		attempts   = 50
	)
	for i := 0; i < attempts; i++ {
		port := rangeStart + rand.Intn(rangeEnd-rangeStart)
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			listener.Close()
			return port, nil
		}
	}
	return 0, fmt.Errorf("could not find a free port in range %d-%d after %d attempts", rangeStart, rangeEnd, attempts)
}

// pickRandomPortPair returns adjacent free ports (N, N+1) for the Hugo backend and the proxy.
func pickRandomPortPair() (hugoPort, proxyPort int, err error) {
	const (
		rangeStart = 20000
		rangeEnd   = 29999 // leave room for N+1
		attempts   = 50
	)
	for i := 0; i < attempts; i++ {
		n := rangeStart + rand.Intn(rangeEnd-rangeStart)
		if checkPortAvailable(n) != nil || checkPortAvailable(n+1) != nil {
			continue
		}
		return n, n + 1, nil
	}
	return 0, 0, fmt.Errorf("could not find an adjacent free port pair in range %d-%d after %d attempts", rangeStart, rangeEnd, attempts)
}

// checkPortAvailable returns an error if the given port is already in use
func checkPortAvailable(port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	listener.Close()
	return nil
}
