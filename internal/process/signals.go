package process

import (
	"os"
	"os/signal"
	"syscall"
)

// WaitForInterrupt blocks execution until an OS interrupt (Ctrl+C) is received
func WaitForInterrupt() {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan // Block here until signal received
}
