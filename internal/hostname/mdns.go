package hostname

import (
	"fmt"
	"net"

	"github.com/hashicorp/mdns"
)

type Broadcaster struct {
	server *mdns.Server
}

// Start registers a .local domain (e.g. appname.local) on the network pointing to our IP.
func Start(name string, port int, lanIP string) (*Broadcaster, error) {
	// Format must end in a dot
	hostName := fmt.Sprintf("%s.local.", name)
	
	info := []string{"RouteLocal Proxy"}
	ips := []net.IP{net.ParseIP(lanIP)}

	// We use the application name for both the Instance and the HostName.
	// This allows it to be discovered as a service AND resolved via ping/browser.
	service, err := mdns.NewMDNSService(name, "_http._tcp", "local.", hostName, port, ips, info)
	if err != nil {
		return nil, fmt.Errorf("failed to create mDNS service: %v", err)
	}

	server, err := mdns.NewServer(&mdns.Config{Zone: service})
	if err != nil {
		return nil, fmt.Errorf("failed to start mDNS server: %v", err)
	}

	return &Broadcaster{server: server}, nil
}

// Stop cleanly shuts down the mDNS broadcaster
func (b *Broadcaster) Stop() error {
	if b.server != nil {
		return b.server.Shutdown()
	}
	return nil
}
