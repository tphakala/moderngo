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

	// Should trigger: pointer receiver pattern
	wg2 := &sync.WaitGroup{}
	wg2.Add(1) // want: "instead of manual Add/Done pattern"
	go func() {
		defer wg2.Done()
		_ = 43
	}()
	wg2.Wait()

	// Should trigger: wg passed as parameter to goroutine
	var wg3 sync.WaitGroup
	wg3.Add(1) // want: "instead of manual Add/Done pattern"
	go func(w *sync.WaitGroup) {
		defer w.Done()
		_ = 44
	}(&wg3)
	wg3.Wait()

	// Should NOT trigger: just wg.Add with no immediate goroutine
	wg.Add(1)
	wg.Done()
}
