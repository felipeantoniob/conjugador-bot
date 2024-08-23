package utils

import (
	"os"
	"syscall"
	"testing"
	"time"
)

func TestWaitForShutdown(t *testing.T) {
	sigCh := make(chan os.Signal, 1)

	go func() {
		time.Sleep(100 * time.Millisecond)
		sigCh <- syscall.SIGTERM
	}()

	sig := WaitForShutdown(sigCh)

	if sig != syscall.SIGTERM {
		t.Errorf("Expected SIGTERM, got %v", sig)
	}
}
