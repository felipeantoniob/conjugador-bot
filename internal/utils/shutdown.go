package utils

import (
	"os"
	"os/signal"
	"syscall"
)

// WaitForShutdown listens for specific shutdown signals and returns the received signal.
func WaitForShutdown(sigCh chan os.Signal) os.Signal {
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	sig := <-sigCh
	return sig
}
