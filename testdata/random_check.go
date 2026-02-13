package testdata

import "math/rand"

// --- RandV2Migration ---

func checkRandV2() {
	// Should trigger: rand.Intn
	_ = rand.Intn(10) // want: "rand.IntN"

	// Should trigger: rand.Int31
	_ = rand.Int31() // want: "rand.Int32()"

	// Should trigger: rand.Int63
	_ = rand.Int63() // want: "rand.Int64()"

	// Should trigger: rand.Seed (deprecated)
	rand.Seed(42) // want: "rand.Seed is deprecated"

	// Should NOT trigger: rand.Int (not renamed in v2)
	_ = rand.Int()

	// Should NOT trigger: rand.Float64 (not renamed)
	_ = rand.Float64()
}
