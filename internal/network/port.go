package network

import (
	"fmt"
	"net"
	"time"
)

// IsPortActive checks if an application is listening on 127.0.0.1 for the given port
func IsPortActive(port int) bool {
	address := fmt.Sprintf("127.0.0.1:%d", port)
	timeout := 1 * time.Second

	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false
	}
	
	if conn != nil {
		defer conn.Close()
		return true
	}
	
	return false
}
