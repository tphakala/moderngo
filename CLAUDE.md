# moderngo — Project Memory

## What This Project Is

A collection of **go-ruleguard DSL** pattern files that detect outdated Go idioms and suggest modern replacements. Not a Go module — just `.go` files with `//go:build ruleguard` build constraints that integrate with **golangci-lint v2** via the gocritic/ruleguard checker.

## Project Structure

```
moderngo/
├── builtins.go      # MinMaxBuiltin, ClearBuiltin, RangeOverInteger, AppendWithoutValues, NewWithExpression
├── crypto.go        # DeprecatedCipherModes, WeakRSAKeySize, DeprecatedElliptic, DeprecatedRSAMultiPrime, DeprecatedPKCS1v15
├── errors.go        # ErrorsAsType (Go 1.26)
├── net.go           # JoinHostPort, FilepathIsLocal, DeprecatedReverseProxyDirector, ErrorBeforeUse
├── random.go        # RandV2Migration
├── reflect.go       # ReflectTypeAssert, ReflectPtrTo, ReflectTypeOf, DeprecatedReflectHeaders, ReflectFieldsIterator, ReflectMethodsIterator, ReflectInsOutsIterator
├── runtime.go       # SetFinalizerDeprecated, GorootDeprecated
├── slices.go        # SortInts, BytesClone, SlicesClone, BackwardIteration, MapKeysCollection, MapValuesCollection, SliceRepeat
├── strings.go       # StringsLinesIteration, StringsSplitIteration, StringsFieldsIteration, StringsFieldsFuncIteration
├── sync.go          # WaitGroupGo
├── testing.go       # BenchmarkLoop, TestingContext, TestingArtifactDir
├── time.go          # TimeDateTimeConstants, TimerChannelLen, DeferredTimeSince, DeferredTimeNow
├── test.sh          # Test runner (bash)
├── README.md        # Documentation
└── testdata/
    ├── .golangci.yml          # golangci-lint v2 config for tests
    ├── go.mod                 # Module: moderngo-testdata (go 1.26.0)
    ├── go.sum
    ├── builtins_check.go      # Tests for builtins rules
    ├── crypto_check.go        # Tests for crypto rules
    ├── errors_check.go        # Tests for errors rules
    ├── net_check.go           # Tests for net rules + false positive tests
    ├── random_check.go        # Tests for random rules
    ├── reflect_check.go       # Tests for reflect rules
    ├── runtime_check.go       # Tests for runtime rules
    ├── slices_check.go        # Tests for slices rules
    ├── strings_check.go       # Tests for strings rules (+ bytes variants)
    ├── sync_check.go          # Tests for sync rules
    ├── testing_check_test.go  # Tests for testing rules (must be _test.go)
    └── time_check.go          # Tests for time rules
```

## How Rules Work

Each `.go` file at root uses the ruleguard DSL:
```go
//go:build ruleguard
package gorules
import "github.com/quasilyte/go-ruleguard/dsl"

func RuleName(m dsl.Matcher) {
    m.Match(`pattern`).
        Where(/* optional type/file guards */).
        Report("message").
        Suggest("fix")  // optional auto-fix
}
```

Key DSL features:
- `$x` matches any expression, `$*x` matches zero or more
- `.Where(m["x"].Type.Is("pkg.Type"))` for type-aware matching
- `.Where(m.File().Name.Matches("_test\\.go$"))` for file-name filtering
- `.Where(!m["n"].Text.Matches("pattern"))` for text exclusions

## How Testing Works

`test.sh` runs golangci-lint on `testdata/` and verifies diagnostics against `// want: "fragment"` annotations.

```bash
./test.sh       # run all tests (124 expectations)
./test.sh -v    # verbose output
```

The test runner:
1. Runs `golangci-lint run --max-issues-per-linter=0 --max-same-issues=0 ./...` in testdata/
2. For each `// want: "fragment"` annotation, checks that golangci-lint output contains that fragment on that line
3. For each diagnostic on a line WITHOUT a want annotation, reports a false positive
4. Reports pass/fail counts

**Critical:** The `--max-issues-per-linter=0 --max-same-issues=0` flags are required. Without them, golangci-lint silently caps at 50 issues and tests fail.

## golangci-lint v2 Config

The testdata config (`testdata/.golangci.yml`):
```yaml
version: "2"
linters:
  default: none
  enable:
    - gocritic
  settings:
    gocritic:
      enabled-checks:
        - ruleguard
      settings:
        ruleguard:
          rules: "${config-path}/../*.go"
          failOn: "all"
```

Key v2 differences from v1:
- `linters-settings:` (v1) -> `linters: settings:` (v2)
- `${configDir}` (v1) -> `${config-path}` (v2)
- `disable: - default` (v1) -> `default: none` (v2)

## Known False Positives (Confirmed by Tests)

1. **FilepathIsLocal** — `strings.Contains(x, "..")` fires on ANY string, not just paths. Version ranges ("1.0..2.0"), ellipsis ("loading...") all trigger. The pattern is too broad. Not fixable via ruleguard DSL (no way to distinguish path strings from other strings).

2. **TestingContext** — `context.Background()` in `_test.go` fires even in helper functions that don't have `*testing.T` in scope. Rule only checks filename, not scope. Not fixable via ruleguard DSL (no way to check enclosing function signature).

3. **SliceRepeat** — `for i := range xs { result = append(result, f(xs[i])...) }` fires even when the appended expression depends on the loop variable (flatMap pattern, not repetition). Not fixable via ruleguard DSL (no way to check if `$s` references `$i`). Mitigated: report message now includes "false positive if $s depends on the loop variable" with `$s` expanded, making it obvious to both humans and LLMs when the match is spurious.

4. **ReflectFieldsIterator (Type patterns only)** — `for i := range t.NumField()` where `t` is `reflect.Type` fires even when the loop body uses the index `i` for `reflect.Value.Field(i)` access. `reflect.Type.Fields()` yields only `StructField` (no index), so the developer can't access the corresponding Value field. Not fixable via ruleguard DSL (can't inspect whether loop body uses index for Value access). Mitigated: report message now includes caveat suggesting to range over the Value instead. Note: `reflect.Value` patterns are NOT affected — `Value.Fields()` returns `iter.Seq2[StructField, Value]` which provides both.

## Known Rule-Ordering Issues

1. **RangeOverInteger vs SliceRepeat** — C-style `for i := 0; i < n; i++ { result = append(result, s...) }` fires RangeOverInteger before SliceRepeat. Users get a two-step path: first convert to `for i := range n`, then SliceRepeat catches it. The `for range n` form is directly caught by SliceRepeat.

### Previously Fixed

- **MapKeysCollection on channels/iterators** — Fixed by adding `.Where(m["m"].Type.Is("map[$k]$v"))` type guard. Was firing on channel drains and iterator collection (11 false positives in birdnet-go).
- **MapValuesCollection on slices** — Fixed by adding `.Where(m["m"].Type.Is("map[$k]$v"))` type guard. Was firing on slice iteration with group-by patterns.
- **RangeOverInteger vs reflect iterators** — Fixed by adding `.Where(!m["n"].Text.Matches(`\.(NumField|NumMethod|NumIn|NumOut)\(\)$`))` exclusion to RangeOverInteger. Reflect-specific rules now fire directly.
- **SlicesClone vs BytesClone** — Fixed by reordering BytesClone before SlicesClone in slices.go.
- **DeprecatedReverseProxyDirector false positive** — Fixed by adding `.Where(m["proxy"].Type.Is("*httputil.ReverseProxy"))` type guard to the assignment pattern.

## Rules Added for Go 1.26

These were added in the most recent session:
- `ErrorsAsType` (errors.go) — `errors.As(err, &target)` -> `errors.AsType[T](err)`
- `NewWithExpression` (builtins.go) — `&[]T{v}[0]` -> `new(v)`
- `ReflectFieldsIterator` (reflect.go) — index loops -> `range t.Fields()`
- `ReflectMethodsIterator` (reflect.go) — index loops -> `range t.Methods()`
- `ReflectInsOutsIterator` (reflect.go) — index loops -> `range t.Ins()/Outs()`
- `DeprecatedPKCS1v15` (crypto.go) — PKCS#1 v1.5 -> OAEP
- `DeprecatedReverseProxyDirector` (net.go) — Director -> Rewrite
- `TestingArtifactDir` (testing.go) — `os.MkdirTemp` in tests -> `t.ArtifactDir()`

## Test Coverage Status

125 test expectations across 12 fixture files. All pass.

### Rules NOT tested (confirmed non-functional)

- **ErrorBeforeUse** (net.go) — 3 patterns matching `f, err := os.Open(); f.Method(); if err != nil`. Confirmed: ruleguard's multi-statement matching does not fire for these patterns. Test fixture exists but without `// want:` annotations.

### Notes on specific rules

- **AppendWithoutValues** — Caught by gocritic's built-in `badCall` checker (message: "no-op append call"), not our ruleguard rule. Our rule's message ("append with single argument has no effect") never appears because `badCall` fires first.
- **BytesClone** — All three patterns (`[]byte(nil)`, `[]byte{}`, `b[:0:0]`) now correctly fire after reordering BytesClone before SlicesClone. The `b[:0:0]` pattern uses a `.Where(m["b"].Type.Is("[]byte"))` guard.

## Pending Work / Ideas

- FilepathIsLocal and TestingContext false positives are NOT fixable via ruleguard DSL (see Known False Positives section)
- ErrorBeforeUse rule confirmed non-functional — consider removing it or waiting for ruleguard multi-statement support
- Consider adding more Go 1.26 rules as the ecosystem matures
- BytesClone `append(b[:0:0], b...)` pattern now works correctly after reordering (tested)
- AppendWithoutValues rule is dead code (gocritic `badCall` fires first) — consider removing

# Maintain a ChangeLog

To track our changes, maintain changelog in CHANGES.md.
