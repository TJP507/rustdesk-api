package service

import (
	"sync"
	"testing"
)

// TestGetAppVersion
func TestGetAppVersion(t *testing.T) {
	s := &AppService{}
	v := s.GetAppVersion()
	// Print result
	t.Logf("App Version: %s", v)
}

func TestMultipleGetAppVersion(t *testing.T) {
	s := &AppService{}
	// Concurrency test
	// Use WaitGroup to wait for all goroutines to complete
	wg := sync.WaitGroup{}
	wg.Add(10) // Start 10 goroutines
	// Start 10 goroutines
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done() // Decrement counter when done
			v := s.GetAppVersion()
			// Print result
			t.Logf("App Version: %s", v)
		}()
	}
	// Wait for all goroutines to complete
	wg.Wait()
}
