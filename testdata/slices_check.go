package testdata

import "sort"

// --- SortInts ---

func checkSortInts() {
	nums := []int{3, 1, 2}
	sort.Ints(nums) // want: "use slices.Sort"

	strs := []string{"c", "a", "b"}
	sort.Strings(strs) // want: "use slices.Sort"

	floats := []float64{3.0, 1.0, 2.0}
	sort.Float64s(floats) // want: "use slices.Sort"

	_ = sort.IntsAreSorted(nums)       // want: "use slices.IsSorted"
	_ = sort.StringsAreSorted(strs)    // want: "use slices.IsSorted"
	_ = sort.Float64sAreSorted(floats) // want: "use slices.IsSorted"

	// Should NOT trigger: sort.Slice (custom comparison)
	sort.Slice(nums, func(i, j int) bool { return nums[i] < nums[j] })
}

// --- SlicesClone ---

func checkSlicesClone() {
	original := []int{1, 2, 3}

	// Should trigger: append([]T(nil), s...)
	_ = append([]int(nil), original...) // want: "use slices.Clone"

	// Should trigger: append([]T{}, s...)
	_ = append([]int{}, original...) // want: "use slices.Clone"

	// Should trigger: append(s[:0:0], s...)
	_ = append(original[:0:0], original...) // want: "use slices.Clone"

	// Should NOT trigger: append with additional elements
	_ = append([]int{0}, original...)
}

// --- BytesClone ---

func checkBytesClone() {
	data := []byte("hello")

	// Should trigger: append([]byte(nil), b...)
	_ = append([]byte(nil), data...) // want: "use slices.Clone"

	// Should trigger: append([]byte{}, b...)
	_ = append([]byte{}, data...) // want: "use slices.Clone"

	// Should NOT trigger: append with extra bytes
	_ = append([]byte{0xFF}, data...)
}

// --- BackwardIteration ---

func checkBackwardIteration() {
	s := []int{1, 2, 3}

	// Should trigger: standard reverse loop (>= 0)
	for i := len(s) - 1; i >= 0; i-- { // want: "use slices.Backward"
		_ = s[i]
	}

	// Should trigger: alternate reverse loop (> -1)
	for i := len(s) - 1; i > -1; i-- { // want: "use slices.Backward"
		_ = s[i]
	}

	// Should NOT trigger: loop with different condition (i > 0, skips index 0)
	for i := len(s) - 1; i > 0; i-- {
		_ = s[i]
	}
}

// --- MapKeysCollection (false-positive-prone) ---

func checkMapKeysCollection() {
	m := map[string]int{"a": 1, "b": 2}

	// Should trigger: simple key collection
	var keys []string
	for k := range m { // want: "use slices.Collect"
		keys = append(keys, k)
	}
	_ = keys

	// Should trigger: key collection with underscore value
	var keys2 []string
	for k, _ := range m { // want: "use slices.Collect"
		keys2 = append(keys2, k)
	}
	_ = keys2

	// Should NOT trigger: loop does more than append
	var filteredKeys []string
	for k := range m {
		if k != "a" {
			filteredKeys = append(filteredKeys, k)
		}
	}
	_ = filteredKeys

	// Should NOT trigger: collecting transformed keys
	var upperKeys []string
	for k := range m {
		upperKeys = append(upperKeys, k+"_suffix")
	}
	_ = upperKeys
}

// --- MapValuesCollection (false-positive-prone) ---

func checkMapValuesCollection() {
	m := map[string]int{"a": 1, "b": 2}

	// Should trigger: simple value collection
	var values []int
	for _, v := range m { // want: "use slices.Collect"
		values = append(values, v)
	}
	_ = values

	// Should NOT trigger: loop does more than append
	var filtered []int
	for _, v := range m {
		if v > 0 {
			filtered = append(filtered, v)
		}
	}
	_ = filtered
}

// --- SliceRepeat ---

func checkSliceRepeat() {
	s := []int{1, 2, 3}
	n := 5

	// Should trigger: manual repetition loop
	var result []int
	for i := 0; i < n; i++ { // want: "use for i := range n"
		result = append(result, s...)
	}
	_ = result
}
