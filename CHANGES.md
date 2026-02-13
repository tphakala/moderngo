# Changelog

## 2026-02-13

- **MapKeysCollection**: Added `.Where(m["m"].Type.Is("map[$k]$v"))` type guard to both patterns. Eliminates false positives on channel drains and iterator collection (11 false positives in birdnet-go).
- **MapValuesCollection**: Added `.Where(m["m"].Type.Is("map[$k]$v"))` type guard. Eliminates false positives on slice iteration with group-by patterns (1 false positive in birdnet-go).
- **ReflectFieldsIterator**: Added caveat to `reflect.Type` patterns about needing the loop index for `reflect.Value` field access. `reflect.Value` patterns are unaffected since `Value.Fields()` returns `(StructField, Value)` pairs.
- **SliceRepeat**: Added false-positive caveat to report message. When the appended expression depends on the loop variable (flatMap pattern), the expanded message now makes this obvious.
- **Test cases**: Added negative test cases for channel drain and slice iteration in slices_check.go.
- **Real-world validation**: Ran rules against vainu2 (11 findings: 10 correct, 1 SliceRepeat FP) and birdnet-go (41→30 findings after fixes, 12 false positives eliminated).
