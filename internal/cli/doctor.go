package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// doctorCmd represents the doctor command
var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose system and network capabilities",
	Long:  `Diagnoses OS, architecture, network interfaces, port availability, and dependencies (like cloudflared).`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running diagnostic checks...")
		fmt.Println("OS: (pending implementation)")
		fmt.Println("Network Interfaces: (pending implementation)")
		fmt.Println("cloudflared: (pending implementation)")
		// Future milestone: actual diagnostic logic
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
