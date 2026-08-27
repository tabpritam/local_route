package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"routelocal/internal/network"
	"routelocal/internal/process"
	"routelocal/internal/proxy"
	"routelocal/internal/hostname"
	"routelocal/internal/tunnel"
)

var (
	port   int
	public bool
	name   string
	debug  bool
)

func printLogo() {
	fmt.Println(`
 ____             _       _                     _ 
|  _ \ ___  _   _| |_ ___| |    ___   ___  __ _| |
| |_) / _ \| | | | __/ _ \ |   / _ \ / __|/ _` + "`" + ` | |
|  _ < (_) | |_| | ||  __/ |__| (_) | (__| (_| | |
|_| \_\___/ \__,_|\__\___|_____\___/ \___|\__,_|_|
`)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "routelocal",
	Short: "RouteLocal makes locally running applications easy to share",
	Long: `RouteLocal is a cross-platform CLI developer tool that makes locally running applications easy to share.

It provides multiple ways to access your application:
- Local (localhost)
- Local Network / LAN (192.168.x.x)
- Local Hostname (myapp.local)
- Public Internet (Cloudflare Tunnel)`,
	Example: `  # Share a local application running on port 3000
  routelocal --port 3000
  
  # Share your app and get all 3 URLs at once (LAN IP, Hostname, and Public URL)
  routelocal --port 3000 --name myapp --public`,
	Run: func(cmd *cobra.Command, args []string) {
		if port == 0 {
			printLogo()
			cmd.Help()
			return
		}

		printLogo()
		fmt.Println("──────────────────────────────────────────────────")

		// 1. Detect Application
		if !network.IsPortActive(port) {
			fmt.Printf("\n✗ No application detected on localhost:%d\n\n", port)
			fmt.Println("Make sure your development server is running and try again.")
			os.Exit(1)
		}

		// 3. Start Proxy (Try target port first)
		proxyServer := proxy.NewServer(port, port)
		errChan := make(chan error, 1)
		
		listenPort, err := proxyServer.Start(errChan)
		if err != nil {
			fmt.Printf("\n[Error] Could not start proxy: %v\n", err)
			os.Exit(1)
		}

		// Print UI
		fmt.Printf("\n✓ Application detected\n  localhost:%d\n", port)
		fmt.Printf("\n✓ Local\n  http://localhost:%d\n", port)

		var mdnsBroadcaster *hostname.Broadcaster
		lanIP, err := network.GetLocalIP()
		if err == nil {
			fmt.Printf("\n✓ Network\n  http://%s:%d\n", lanIP, listenPort)

			// 4. Start mDNS Broadcaster
			if name != "" {
				mdnsBroadcaster, err = hostname.Start(name, listenPort, lanIP)
				if err == nil {
					fmt.Printf("\n✓ Hostname\n  http://%s.local:%d\n", name, listenPort)
				} else {
					fmt.Printf("\n⚠ Hostname failed to broadcast: %v\n", err)
				}
			}
		}

		var cfTunnel *tunnel.CloudflareTunnel
		if public {
			fmt.Printf("\n⠋ Starting Cloudflare Tunnel...\n")
			var publicURL string
			cfTunnel, publicURL, err = tunnel.Start(port)
			if err != nil {
				fmt.Printf("✗ Failed to start tunnel: %v\n", err)
			} else {
				fmt.Printf("✓ Public tunnel ready\n  %s\n", publicURL)
				fmt.Println("\n⚠ PUBLIC ACCESS ENABLED")
				fmt.Println("  Anyone with the public URL may be able to access your application.")
			}
		}

		if debug {
			fmt.Printf("\nDebug Info:\n")
			fmt.Printf("Public: %v\n", public)
			fmt.Printf("Name: %s\n", name)
			fmt.Printf("Debug: %v\n", debug)
		}

		fmt.Println("\n──────────────────────────────────────────────────")
		fmt.Println("\nPress Ctrl+C to stop.")
		
		// Listen for proxy background errors
		go func() {
			if err := <-errChan; err != nil {
				fmt.Printf("\n[Error] %v\n", err)
				os.Exit(1)
			}
		}()

		// Block and wait for OS interrupt (Ctrl+C)
		process.WaitForInterrupt()

		// Graceful Shutdown
		fmt.Print("\nStopping RouteLocal... ")
		if mdnsBroadcaster != nil {
			mdnsBroadcaster.Stop()
		}
		if cfTunnel != nil {
			cfTunnel.Stop()
		}
		proxyServer.Stop()
		fmt.Println("✓ Cleanup complete")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

const customUsageTemplate = `Usage:
  routelocal [flags]

Examples:{{.Example}}
`

func init() {
	rootCmd.SetUsageTemplate(customUsageTemplate)
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// Root command flags
	rootCmd.Flags().IntVarP(&port, "port", "p", 0, "Port of the local application (required)")
	rootCmd.Flags().BoolVar(&public, "public", false, "Expose the application to the public internet via Cloudflare Tunnel")
	rootCmd.Flags().StringVarP(&name, "name", "n", "", "Custom local hostname (e.g., myapp -> myapp.local)")
	rootCmd.Flags().BoolVar(&debug, "debug", false, "Enable verbose debugging output")
}
