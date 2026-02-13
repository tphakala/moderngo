package testdata

import (
	"errors"
	"io/fs"
)

func checkErrorsAsType(err error) {
	// Should trigger: errors.As with address-of target
	var pathErr *fs.PathError
	if errors.As(err, &pathErr) { // want: "use errors.AsType"
		_ = pathErr
	}

	// Should NOT trigger: errors.Is (different function)
	if errors.Is(err, fs.ErrNotExist) {
		return
	}

	// Should NOT trigger: errors.As without address-of (rare but valid)
	var target error
	_ = errors.As(err, &target) // want: "use errors.AsType"
}
