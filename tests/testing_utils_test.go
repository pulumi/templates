package tests

import (
	"fmt"
	"testing"
	"time"
)

func runWithTimeout(
	t *testing.T,
	timeout time.Duration,
	name string,
	prepare func(*testing.T),
	run func(*testing.T),
) {
	t.Run(name, func(t *testing.T) {
		prepare(t)
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

func parallel(t *testing.T) {
	t.Parallel()
}
