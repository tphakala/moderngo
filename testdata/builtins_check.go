package testdata

import "math"

// --- NewWithExpression ---

func checkNewWithExpression() {
	// Should trigger: slice hack for pointer-to-value
	_ = &[]string{"hello"}[0] // want: "use new"
	_ = &[]int{42}[0]         // want: "use new"
	_ = &[]bool{true}[0]      // want: "use new"

	// Should NOT trigger: slice with multiple elements
	s := []string{"a", "b"}
	_ = s

	// Should NOT trigger: normal slice indexing
	items := []int{1, 2, 3}
	_ = &items[0]

	// Should NOT trigger: new without expression (existing usage)
	_ = new(int)
}

// --- ClearBuiltin (false-positive-prone) ---

func checkClearBuiltin() {
	// Should trigger: simple map clearing loop
	m := map[string]int{"a": 1}
	for k := range m { // want: "use clear"
		delete(m, k)
	}

	// Should NOT trigger: loop that does more than just delete
	m2 := map[string]int{"a": 1}
	for k := range m2 {
		if k != "keep" {
			delete(m2, k)
		}
	}

	// Should NOT trigger: delete with different map
	m3 := map[string]int{"a": 1}
	m4 := map[string]int{"b": 2}
	for k := range m3 {
		delete(m4, k)
	}
	_ = m4
}

// --- RangeOverInteger (false-positive-prone) ---

func checkRangeOverInteger() {
	n := 10

	// Should trigger: standard 0-to-n loop
	for i := 0; i < n; i++ { // want: "use for i := range n"
		_ = i
	}

	// Should NOT trigger: loop not starting at 0
	for i := 1; i < n; i++ {
		_ = i
	}

	// Should NOT trigger: decrementing loop
	for i := n; i > 0; i-- {
		_ = i
	}

	// Should NOT trigger: step by 2
	for i := 0; i < n; i += 2 {
		_ = i
	}
}

// --- MinMaxBuiltin ---

func checkMinMaxBuiltin(a, b int) {
	// Should trigger: int(math.Min(...))
	_ = int(math.Min(float64(a), float64(b))) // want: "use min"

	// Should trigger: int(math.Max(...))
	_ = int(math.Max(float64(a), float64(b))) // want: "use max"

	// Should NOT trigger: math.Min with actual floats
	x, y := 1.5, 2.5
	_ = math.Min(x, y)
}

// --- AppendWithoutValues ---

func checkAppendWithoutValues() {
	s := []int{1, 2, 3}

	// Should trigger: append with no values
	s = append(s) // want: "no-op append call"

	// Should NOT trigger: append with values
	s = append(s, 4)
	_ = s
}
