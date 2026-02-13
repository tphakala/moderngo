package testdata

import (
	"log"
	"time"
)

// --- TimeDateTimeConstants ---

func checkTimeDateTimeConstants() {
	t := time.Now()

	// Should trigger: magic DateTime format
	_ = t.Format("2006-01-02 15:04:05") // want: "use t.Format(time.DateTime)"

	// Should trigger: magic DateOnly format
	_ = t.Format("2006-01-02") // want: "use t.Format(time.DateOnly)"

	// Should trigger: magic TimeOnly format
	_ = t.Format("15:04:05") // want: "use t.Format(time.TimeOnly)"

	// Should trigger: time.Parse with magic format
	_, _ = time.Parse("2006-01-02 15:04:05", "2024-01-01 00:00:00") // want: "use time.Parse(time.DateTime"

	// Should NOT trigger: custom format that doesn't match exactly
	_ = t.Format("2006-01-02T15:04:05")

	// Should NOT trigger: time.RFC3339
	_ = t.Format(time.RFC3339)
}

// --- DeferredTimeSince ---

func checkDeferredTimeSince() {
	start := time.Now()

	// Should trigger: time.Since in defer argument
	defer log.Println(time.Since(start)) // want: "time.Since(start) is evaluated at defer time"

	// Should NOT trigger: wrapped in closure (correct pattern)
	defer func() { log.Println(time.Since(start)) }()
}

// --- DeferredTimeNow ---

func checkDeferredTimeNow() {
	// Should trigger: time.Now() in defer argument
	defer log.Println(time.Now()) // want: "time.Now() is evaluated at defer time"

	// Should NOT trigger: wrapped in closure (correct pattern)
	defer func() { log.Println(time.Now()) }()
}
