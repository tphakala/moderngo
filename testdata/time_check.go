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

	// Should trigger: time.Parse with magic formats
	_, _ = time.Parse("2006-01-02 15:04:05", "2024-01-01 00:00:00") // want: "use time.Parse(time.DateTime"
	_, _ = time.Parse("2006-01-02", "2024-01-01")                   // want: "use time.Parse(time.DateOnly"
	_, _ = time.Parse("15:04:05", "12:00:00")                       // want: "use time.Parse(time.TimeOnly"

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

	// Should trigger: time.Since as second argument
	defer log.Printf("took %v", time.Since(start)) // want: "time.Since(start) is evaluated at defer time"

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

// --- TimerChannelLen ---

func checkTimerChannelLen() {
	timer := time.NewTimer(time.Second)

	// Should trigger: len() on timer.C
	_ = len(timer.C) // want: "len() on timer channel is always 0"

	// Should trigger: cap() on timer.C
	_ = cap(timer.C) // want: "cap() on timer channel is always 0"

	ticker := time.NewTicker(time.Second)

	// Should trigger: len() on ticker.C
	_ = len(ticker.C) // want: "len() on ticker channel is always 0"

	// Should trigger: cap() on ticker.C
	_ = cap(ticker.C) // want: "cap() on ticker channel is always 0"

	// Should NOT trigger: len/cap on regular channel
	ch := make(chan int, 1)
	_ = len(ch)
	_ = cap(ch)

	timer.Stop()
	ticker.Stop()
}
