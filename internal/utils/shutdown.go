package utils

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// WaitForShutdown sets up a channel to listen for specific shutdown signals and blocks until one of those signals is received.
func WaitForShutdown() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	sig := <-sc
	fmt.Printf("Received shutdown signal: %v\n", sig)
}
