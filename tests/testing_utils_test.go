package tests

import (
	"testing"
	"time"
)

func runWithTimeout(t *testing.T, timeout time.Duration, name string, run func(*testing.T)) {
	t.Run(name, func(t *testing.T) {
		timeoutEvent := time.After(timeout)
		done := make(chan bool)
		go func() {
			run(t)
			done <- true
		}()
		select {
		case <-timeoutEvent:
			t.Fatalf("%s timed out after %s", name, timeout)
		case <-done:
			return
		}
	})
}
