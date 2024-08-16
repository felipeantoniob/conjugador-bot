package utils

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// WaitForShutdown sets up a channel to listen for specific shutdown signals and blocks until one of those signals is received.
func WaitForShutdown() {
	// Create a channel to receive OS signals
	sc := make(chan os.Signal, 1)

	// Notify the channel for specific shutdown signals
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	// Block until one of the shutdown signals is received
	sig := <-sc

	// Print the received signal (for debugging or logging purposes)
	fmt.Printf("Received shutdown signal: %v\n", sig)
}

func NullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
