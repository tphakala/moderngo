package testdata

import "runtime"

// --- SetFinalizerDeprecated ---

type resource struct{ fd int }

func checkSetFinalizer() {
	r := &resource{fd: 42}

	// Should trigger: runtime.SetFinalizer
	runtime.SetFinalizer(r, func(r *resource) { _ = r.fd }) // want: "consider using runtime.AddCleanup"

	// Should NOT trigger: runtime.KeepAlive (different function)
	runtime.KeepAlive(r)
}

// --- GorootDeprecated ---

func checkGoroot() {
	// Should trigger: runtime.GOROOT()
	_ = runtime.GOROOT() // want: "runtime.GOROOT() is deprecated"

	// Should NOT trigger: runtime.GOOS (not deprecated)
	_ = runtime.GOOS
}
