package testdata

import "sync"

// --- WaitGroupGo ---

func checkWaitGroupGo() {
	var wg sync.WaitGroup

	// Should trigger: Add(1) + go func() { defer Done() }
	wg.Add(1) // want: "instead of manual Add/Done pattern"
	go func() {
		defer wg.Done()
		_ = 42
	}()

	wg.Wait()

	// Should NOT trigger: just wg.Add with no immediate goroutine
	wg.Add(1)
	wg.Done()
}
