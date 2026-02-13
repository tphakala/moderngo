package testdata

import "reflect"

// --- ReflectPtrTo ---

func checkReflectPtrTo() {
	t := reflect.TypeOf(0)

	// Should trigger: deprecated PtrTo
	_ = reflect.PtrTo(t) // want: "reflect.PtrTo is deprecated"

	// Should NOT trigger: PointerTo (the replacement)
	_ = reflect.PointerTo(t)
}

// --- ReflectTypeOf ---

func checkReflectTypeOf() {
	// Should trigger: TypeOf((*T)(nil)).Elem() pattern
	_ = reflect.TypeOf((*int)(nil)).Elem() // want: "use reflect.TypeFor"

	// Should NOT trigger: TypeOf with a real value
	_ = reflect.TypeOf(42)
}

// --- DeprecatedReflectHeaders ---

func checkReflectHeaders() {
	// Should trigger: SliceHeader literal
	_ = reflect.SliceHeader{} // want: "reflect.SliceHeader is deprecated"

	// Should trigger: StringHeader literal
	_ = reflect.StringHeader{} // want: "reflect.StringHeader is deprecated"

	// Should NOT trigger: unrelated reflect usage
	_ = reflect.TypeOf(0).Kind()
}
