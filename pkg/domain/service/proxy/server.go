package proxy

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
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

// NewServer creates a new proxy server instance
func NewServer(proxyPort int) *Server {
	hugoPort := findAvailablePort(1313) // Start with Hugo's default port
	return &Server{
		hugoPort:  hugoPort,
		proxyPort: proxyPort,
	}
}

// Start initializes Hugo server and reverse proxy
func (s *Server) Start(hugoArgs []string) error {
	// Prepare Hugo with themes
	hugoService := hugo.New()

	// Start Hugo server on internal port in a goroutine
	hugoDone := make(chan error, 1)
	s.hugoDone = hugoDone
	go func() {
		// Add port flag to Hugo args
		args := append(hugoArgs, "--port", strconv.Itoa(s.hugoPort))

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
	
	// Enhance the default transport for better performance
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
	}

	// Fix Hugo's relative redirects to be absolute with proxy port
	proxy.ModifyResponse = func(resp *http.Response) error {
		// Only handle redirects
		if resp.StatusCode < 300 || resp.StatusCode >= 400 {
			return nil
		}
		
		location := resp.Header.Get("Location")
		if location == "" {
			return nil
		}

		// Parse location URL
		locationURL, err := url.Parse(location)
		if err != nil || locationURL.IsAbs() {
			return nil // Already absolute or invalid
		}

		// Make it absolute using the proxy's scheme and host
		locationURL.Scheme = "http"
		locationURL.Host = fmt.Sprintf("localhost:%d", s.proxyPort)
		resp.Header.Set("Location", locationURL.String())
		
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

// waitForHugo polls the Hugo server until it's ready or times out
func (s *Server) waitForHugo() error {
	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout waiting for Hugo server to start")
		case <-ticker.C:
			conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", s.hugoPort), time.Second)
			if err == nil {
				conn.Close()
				return nil
			}
		}
	}
}

// findAvailablePort finds an available port starting from the given port
func findAvailablePort(startPort int) int {
	for port := startPort; port < startPort+100; port++ {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			listener.Close()
			return port
		}
	}
	return startPort // Fallback to original port
}

