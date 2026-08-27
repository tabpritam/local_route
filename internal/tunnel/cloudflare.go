package tunnel

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"time"
)

type CloudflareTunnel struct {
	cmd    *exec.Cmd
	cancel context.CancelFunc
}

func ensureCloudflared() (string, error) {
	// 1. Check if it's already installed globally
	if path, err := exec.LookPath("cloudflared"); err == nil {
		return path, nil
	}

	// 2. Determine local cache path
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = os.TempDir()
	}
	binDir := filepath.Join(cacheDir, "routelocal", "bin")
	_ = os.MkdirAll(binDir, 0755)

	exeName := "cloudflared"
	if runtime.GOOS == "windows" {
		exeName = "cloudflared.exe"
	}
	binPath := filepath.Join(binDir, exeName)

	if _, err := os.Stat(binPath); err == nil {
		return binPath, nil // Already downloaded
	}

	// 3. Prompt and download
	fmt.Printf("\n⚠ cloudflared is not installed.\n⠋ Downloading cloudflared for %s to local cache...\n", runtime.GOOS)

	var downloadURL string
	switch runtime.GOOS {
	case "windows":
		downloadURL = "https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-windows-amd64.exe"
	case "darwin":
		downloadURL = "https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-darwin-amd64" // Assumes AMD64/Rosetta for MVP
	case "linux":
		downloadURL = "https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64"
	default:
		return "", fmt.Errorf("automatic download not supported for %s", runtime.GOOS)
	}

	resp, err := http.Get(downloadURL)
	if err != nil {
		return "", fmt.Errorf("failed to download cloudflared: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download cloudflared (HTTP %d)", resp.StatusCode)
	}

	out, err := os.OpenFile(binPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755) // Ensure executable permissions
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return "", err
	}

	fmt.Println("✓ Download complete!")
	return binPath, nil
}

// Start launches the cloudflared process and extracts the trycloudflare.com URL
func Start(port int) (*CloudflareTunnel, string, error) {
	binPath, err := ensureCloudflared()
	if err != nil {
		return nil, "", err
	}

	ctx, cancel := context.WithCancel(context.Background())
	
	// Start cloudflared targeting the developer's application
	cmd := exec.CommandContext(ctx, binPath, "tunnel", "--url", fmt.Sprintf("http://localhost:%d", port))
	
	// cloudflared logs its output to stderr
	stderr, err := cmd.StderrPipe()
	if err != nil {
		cancel()
		return nil, "", err
	}

	if err := cmd.Start(); err != nil {
		cancel()
		return nil, "", err
	}

	urlChan := make(chan string)
	errChan := make(chan error)

	// Scan stderr asynchronously for the public URL
	go func() {
		scanner := bufio.NewScanner(stderr)
		urlRegex := regexp.MustCompile(`https://[a-zA-Z0-9-]+\.trycloudflare\.com`)
		
		for scanner.Scan() {
			line := scanner.Text()
			if match := urlRegex.FindString(line); match != "" {
				urlChan <- match
				// We don't return here because we want cloudflared to keep running
				// but we stop listening for the URL.
				return
			}
		}
		errChan <- fmt.Errorf("failed to extract URL from cloudflared output")
	}()

	// Wait for the URL to be generated or timeout after 15 seconds
	select {
	case publicURL := <-urlChan:
		return &CloudflareTunnel{cmd: cmd, cancel: cancel}, publicURL, nil
	case err := <-errChan:
		cancel()
		return nil, "", err
	case <-time.After(15 * time.Second):
		cancel()
		return nil, "", fmt.Errorf("timeout waiting for Cloudflare Tunnel to start")
	}
}

// Stop terminates the cloudflared child process
func (c *CloudflareTunnel) Stop() {
	if c.cancel != nil {
		c.cancel()
	}
	if c.cmd != nil && c.cmd.Process != nil {
		c.cmd.Process.Kill()
		c.cmd.Wait()
	}
}
