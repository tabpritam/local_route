package proxy

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type Server struct {
	targetPort int
	listenPort int
	server     *http.Server
}

// NewServer creates a new reverse proxy server
func NewServer(targetPort int, listenPort int) *Server {
	return &Server{
		targetPort: targetPort,
		listenPort: listenPort,
	}
}

// Start runs the proxy in the background and sends errors to errChan. Returns the actual bound port.
func (s *Server) Start(errChan chan<- error) (int, error) {
	targetURL, _ := url.Parse(fmt.Sprintf("http://127.0.0.1:%d", s.targetPort))
	
	// Go's built-in ReverseProxy automatically supports WebSockets (Connection Upgrade)
	// and Server-Sent Events (SSE) streaming!
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		// Ensure Host header is preserved for frameworks like Next.js / Vite
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.Host = targetURL.Host
	}

	s.server = &http.Server{
		Handler: proxy,
	}

	// Try binding to the preferred listen port
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", s.listenPort))
	if err != nil {
		// Fallback to random available port if preferred port is already in use
		listener, err = net.Listen("tcp", "0.0.0.0:0")
		if err != nil {
			return 0, fmt.Errorf("failed to bind to any port: %v", err)
		}
	}

	actualPort := listener.Addr().(*net.TCPAddr).Port

	go func() {
		if err := s.server.Serve(listener); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("proxy server error: %v", err)
		}
	}()
	
	return actualPort, nil
}

// Stop cleanly shuts down the reverse proxy
func (s *Server) Stop() error {
	if s.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return s.server.Shutdown(ctx)
	}
	return nil
}
