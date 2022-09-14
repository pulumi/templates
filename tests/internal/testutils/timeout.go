package testutils

import (
	"testing"
	"time"
)

func RunWithTimeout(
	t *testing.T,
	timeout time.Duration,
	name string,
	prepare func(*testing.T),
	run func(*testing.T),
) {
	t.Run(name, func(t *testing.T) {
		if prepare != nil {
			prepare(t)
		}
		timeoutEvent := time.After(timeout)
		done := make(chan bool)
		go func() {
			defer func() {
				done <- true
			}()
			run(t)
		}()
		select {
		case <-timeoutEvent:
			t.Fatalf("%s timed out after %s", name, timeout)
		case <-done:
			return
		}
	})
}
