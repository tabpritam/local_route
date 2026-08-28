package network

import (
	"fmt"
	"net"
	"time"
)

// IsPortActive checks if an application is listening on the given port
func IsPortActive(port int) bool {
	timeout := 1 * time.Second

	// Try localhost first (resolves to IPv6 ::1 or IPv4 depending on the system/server)
	address := fmt.Sprintf("localhost:%d", port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err == nil && conn != nil {
		conn.Close()
		return true
	}

	// Fallback to explicit IPv4 loopback
	addressIPv4 := fmt.Sprintf("127.0.0.1:%d", port)
	connIPv4, errIPv4 := net.DialTimeout("tcp", addressIPv4, timeout)
	if errIPv4 == nil && connIPv4 != nil {
		connIPv4.Close()
		return true
	}
	
	return false
}
