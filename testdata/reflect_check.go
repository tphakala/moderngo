package testdata

import (
	"errors"
	"reflect"
	"unsafe"
)

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

	// Should trigger: SliceHeader with initialized fields
	_ = reflect.SliceHeader{Data: 0, Len: 0, Cap: 0} // want: "reflect.SliceHeader is deprecated"

	// Should trigger: StringHeader literal
	_ = reflect.StringHeader{} // want: "reflect.StringHeader is deprecated"

	// Should trigger: cast to SliceHeader
	s := []byte{1, 2, 3}
	_ = (*reflect.SliceHeader)(unsafe.Pointer(&s)) // want: "reflect.SliceHeader is deprecated"

	// Should trigger: cast to StringHeader
	str := "hello"
	_ = (*reflect.StringHeader)(unsafe.Pointer(&str)) // want: "reflect.StringHeader is deprecated"

	// Should NOT trigger: unrelated reflect usage
	_ = reflect.TypeOf(0).Kind()
}

// --- ReflectTypeAssert ---

func checkReflectTypeAssert() {
	v := reflect.ValueOf("hello")

	// Should trigger: v.Interface().(T) pattern
	_ = v.Interface().(string) // want: "reflect.TypeAssert"

	// Should trigger: comma-ok type assertion form
	_, _ = v.Interface().(string) // want: "reflect.TypeAssert"

	// Should NOT trigger: Interface() without type assertion
	_ = v.Interface()
}

// --- ReflectFieldsIterator ---

func checkReflectFieldsIterator() {
	t := reflect.TypeOf(struct{ X int }{})

	// Should trigger ReflectFieldsIterator: index-based loop over NumField
	for i := 0; i < t.NumField(); i++ { // want: "range t.Fields()"
		_ = t.Field(i)
	}

	// Should trigger ReflectFieldsIterator: range over NumField
	for i := range t.NumField() { // want: "range t.Fields()"
		_ = t.Field(i)
	}

	v := reflect.ValueOf(struct{ X int }{})

	// Should trigger ReflectFieldsIterator: index-based loop over NumField (value)
	for i := 0; i < v.NumField(); i++ { // want: "range v.Fields()"
		_ = v.Field(i)
	}

	// Should trigger ReflectFieldsIterator: range over Value.NumField
	for i := range v.NumField() { // want: "range v.Fields()"
		_ = v.Field(i)
	}

	// Should NOT trigger: NumField used without loop
	_ = t.NumField()
}

// --- ReflectMethodsIterator ---

func checkReflectMethodsIterator() {
	t := reflect.TypeOf((*error)(nil)).Elem() // want: "use reflect.TypeFor"

	// Should trigger ReflectMethodsIterator: index-based loop over NumMethod
	for i := 0; i < t.NumMethod(); i++ { // want: "range t.Methods()"
		_ = t.Method(i)
	}

	// Should trigger ReflectMethodsIterator: range over NumMethod
	for i := range t.NumMethod() { // want: "range t.Methods()"
		_ = t.Method(i)
	}

	// Use a concrete type value to get a reflect.Value with methods
	v := reflect.ValueOf(errors.New("test"))

	// Should trigger ReflectMethodsIterator: index loop on Value
	for i := 0; i < v.NumMethod(); i++ { // want: "range v.Methods()"
		_ = v.Method(i)
	}

	// Should trigger ReflectMethodsIterator: range over Value.NumMethod
	for i := range v.NumMethod() { // want: "range v.Methods()"
		_ = v.Method(i)
	}
}

// --- ReflectInsOutsIterator ---

func checkReflectInsOutsIterator() {
	t := reflect.TypeOf(func(int, string) bool { return false })

	// Should trigger ReflectInsOutsIterator: index-based loop over NumIn
	for i := 0; i < t.NumIn(); i++ { // want: "range t.Ins()"
		_ = t.In(i)
	}

	// Should trigger ReflectInsOutsIterator: index-based loop over NumOut
	for i := 0; i < t.NumOut(); i++ { // want: "range t.Outs()"
		_ = t.Out(i)
	}

	// Should trigger ReflectInsOutsIterator: range over NumIn
	for i := range t.NumIn() { // want: "range t.Ins()"
		_ = t.In(i)
	}

	// Should trigger ReflectInsOutsIterator: range over NumOut
	for i := range t.NumOut() { // want: "range t.Outs()"
		_ = t.Out(i)
	}
}
