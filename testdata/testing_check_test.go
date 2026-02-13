package testdata

import (
	"context"
	"os"
	"testing"
)

// --- TestingArtifactDir ---

func TestArtifactDir(t *testing.T) {
	// Should trigger: os.MkdirTemp in test file
	_, _ = os.MkdirTemp("", "test-output-*") // want: "consider t.ArtifactDir"

	// Should trigger: os.MkdirTemp with different args
	_, _ = os.MkdirTemp(os.TempDir(), "prefix-*") // want: "consider t.ArtifactDir"

	_ = t
}

// --- TestingContext (false-positive-prone) ---

func TestContextBackground(t *testing.T) {
	// Should trigger: context.Background() assigned in test
	ctx := context.Background() // want: "use t.Context() instead of context.Background()"

	// Should trigger: reassignment form (= not :=)
	ctx = context.Background() // want: "use t.Context() instead of context.Background()"
	_ = ctx

	// Should trigger: context.TODO() assigned in test
	ctx2 := context.TODO() // want: "use t.Context() instead of context.TODO()"
	_ = ctx2
	_ = t
}

func TestContextPassedDirectly(t *testing.T) {
	// Should trigger: context.Background() passed directly
	doSomethingWithCtx(context.Background(), "test") // want: "use t.Context() instead of context.Background()"

	// Should trigger: context.TODO() passed directly
	doSomethingWithCtx(context.TODO(), "test") // want: "use t.Context() instead of context.TODO()"
	_ = t
}

// FALSE POSITIVE: context.Background() in non-Test helper functions
// triggers because the rule only checks filename (_test.go), not scope.
func testHelperNoT() {
	ctx := context.Background() // want: "use t.Context() instead of context.Background()"
	_ = ctx
}

// Should NOT trigger: context.WithCancel (not Background/TODO)
func TestContextWithCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	_ = ctx
}

func doSomethingWithCtx(ctx context.Context, s string) { _ = ctx; _ = s }

// --- BenchmarkLoop ---

func BenchmarkOldLoop(b *testing.B) {
	// Should trigger: old b.N pattern
	for i := 0; i < b.N; i++ { // want: "use for b.Loop"
		_ = i
	}
}

func BenchmarkRangeN(b *testing.B) {
	// Should trigger: range over b.N pattern (Go 1.22+ style)
	for i := range b.N { // want: "use for b.Loop"
		_ = i
	}
}

func BenchmarkForRangeN(b *testing.B) {
	// Should trigger: for range b.N with no variable
	for range b.N { // want: "use for b.Loop"
		_ = 42
	}
}
