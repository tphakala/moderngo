package testdata

import (
	"bytes"
	"strings"
)

// --- StringsLinesIteration ---

func checkStringsLines(s string) {
	// Should trigger: split by \n for iteration
	for _, line := range strings.Split(s, "\n") { // want: "use for line := range strings.Lines"
		_ = line
	}

	// Should trigger: split by \r\n for iteration
	for _, line := range strings.Split(s, "\r\n") { // want: "use for line := range strings.Lines"
		_ = line
	}

	// Should trigger: bytes.Split by \n for iteration
	bs := []byte(s)
	for _, line := range bytes.Split(bs, []byte("\n")) { // want: "use for line := range bytes.Lines"
		_ = line
	}

	// Should trigger: bytes.Split by byte literal for iteration
	for _, line := range bytes.Split(bs, []byte{'\n'}) { // want: "use for line := range bytes.Lines"
		_ = line
	}

	// Should NOT trigger: split result used as slice (not just iteration)
	lines := strings.Split(s, "\n")
	_ = lines[0]
}

// --- StringsSplitIteration ---

func checkStringsSplitSeq(s string) {
	// Should trigger: split by comma for iteration
	for _, part := range strings.Split(s, ",") { // want: "use for part := range strings.SplitSeq"
		_ = part
	}

	// Should trigger: bytes.Split by comma for iteration
	bs := []byte(s)
	for _, part := range bytes.Split(bs, []byte(",")) { // want: "use for part := range bytes.SplitSeq"
		_ = part
	}

	// Should NOT trigger: newline split (caught by Lines rule instead)
	// (the rule explicitly excludes \n separators)
}

// --- StringsFieldsIteration ---

func checkStringsFieldsSeq(s string) {
	// Should trigger: Fields used for iteration
	for _, field := range strings.Fields(s) { // want: "use for field := range strings.FieldsSeq"
		_ = field
	}

	// Should trigger: bytes.Fields for iteration
	bs := []byte(s)
	for _, field := range bytes.Fields(bs) { // want: "use for field := range bytes.FieldsSeq"
		_ = field
	}
}

// --- StringsFieldsFuncIteration ---

func checkStringsFieldsFuncSeq(s string) {
	// Should trigger: FieldsFunc used for iteration
	for _, field := range strings.FieldsFunc(s, func(r rune) bool { return r == ',' }) { // want: "use for field := range strings.FieldsFuncSeq"
		_ = field
	}

	// Should trigger: bytes.FieldsFunc for iteration
	bs := []byte(s)
	for _, field := range bytes.FieldsFunc(bs, func(r rune) bool { return r == ',' }) { // want: "use for field := range bytes.FieldsFuncSeq"
		_ = field
	}
}
