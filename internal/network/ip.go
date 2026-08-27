package network

import (
	"fmt"
	"net"
	"strings"
)

// GetLocalIP returns the best guess for the machine's local IPv4 address on the LAN.
func GetLocalIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	var fallback string
	for _, iface := range interfaces {
		// Skip down interfaces and loopback interfaces
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		// Skip common virtual/docker interfaces based on name
		name := strings.ToLower(iface.Name)
		if strings.Contains(name, "docker") || strings.Contains(name, "veth") || strings.Contains(name, "vmware") || strings.Contains(name, "wsl") || strings.Contains(name, "virtual") || strings.Contains(name, "hyper-v") {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}

			// Prefer private network ranges (10.x.x.x, 172.16.x.x-172.31.x.x, 192.168.x.x)
			if ip.IsPrivate() {
				return ip.String(), nil
			}

			// Record fallback if we find a public/non-private IPv4 just in case
			if fallback == "" {
				fallback = ip.String()
			}
		}
	}

	if fallback != "" {
		return fallback, nil
	}

	return "", fmt.Errorf("no suitable network interface found")
}
