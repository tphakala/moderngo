# Go 1.20+ Linter Rules

Custom ruleguard rules for modernizing Go code to use Go 1.20+ through Go 1.25+ features.

## Usage with golangci-lint v2

### Installation

```bash
# Recommended: binary installation
curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.8.0
```

> **Note:** `go install` is not recommended as it may produce unreliable builds. See the [official installation docs](https://golangci-lint.run/docs/welcome/install/local/) for details.

### Quick Start

1. Clone or copy the rules to your project:
   ```bash
   git clone https://github.com/yourorg/moderngo.git
   # Or copy the rules/ directory to your project
   ```

2. Add to your `.golangci.yml`:
   ```yaml
   version: "2"
   linters:
     enable:
       - gocritic
   linters-settings:
     gocritic:
       enabled-checks:
         - ruleguard
       settings:
         ruleguard:
           rules:
             - rules/*.go
   ```

3. Run the linter:
   ```bash
   golangci-lint run ./...
   ```

### Clearing Cache After Rule Changes

When you modify rules, clear the cache to ensure changes take effect:

```bash
golangci-lint cache clean && golangci-lint run ./...
```

### Running Specific Rules

To run only moderngo rules (via gocritic/ruleguard):

```bash
golangci-lint run --enable-only gocritic ./...
```

### Example Output

```
main.go:15:2: ruleguard: suggestion: use slices.Sort instead of sort.Ints (gocritic)
main.go:23:5: ruleguard: suggestion: use time.DateTime constant instead of "2006-01-02 15:04:05" (gocritic)
main.go:31:2: ruleguard: suggestion: use range over int (Go 1.22+) (gocritic)
```

---

## File Organization

Rules are organized by **package/topic** rather than Go version for easier maintenance.

| File | Topic | Rules |
|------|-------|-------|
| [strings.go](#stringsgo) | String iteration | Lines, SplitSeq, FieldsSeq |
| [time.go](#timego) | Time formatting & timers | DateTime constants, Timer len(), deferred time.Since |
| [slices.go](#slicesgo) | Slice operations | Sort, Clone, Backward, map keys/values, bytes.Clone |
| [sync.go](#syncgo) | Synchronization | WaitGroup.Go |
| [builtins.go](#builtinsgo) | Built-in functions | min/max, clear(), range-over-int, append no-op |
| [reflect.go](#reflectgo) | Reflection | TypeAssert, PointerTo, TypeFor, deprecated headers |
| [random.go](#randomgo) | Random numbers | math/rand/v2 migration, Seed/Read deprecation |
| [testing.go](#testinggo) | Testing utilities | b.Loop, t.Context |
| [net.go](#netgo) | Network & paths | JoinHostPort, filepath.IsLocal, error before use |
| [crypto.go](#cryptogo) | Cryptography | Cipher modes, RSA key size, elliptic deprecation |
| [runtime.go](#runtimego) | Runtime functions | SetFinalizer, GOROOT deprecation |

---

## strings.go

String iteration patterns using Go 1.23+ iterator APIs.

See: [strings package](https://pkg.go.dev/strings)

### strings.Lines Iteration

**Old pattern:**
```go
for _, line := range strings.Split(s, "\n") {
    process(line)
}
```

**New pattern (Go 1.23+):**
```go
for line := range strings.Lines(s) {
    process(line)
}
```

**Benefits:**
- No intermediate slice allocation
- Handles both `\n` and `\r\n` line endings automatically

### strings.SplitSeq Iteration

**Old pattern:**
```go
for _, part := range strings.Split(s, ",") {
    process(part)
}
```

**New pattern (Go 1.23+):**
```go
for part := range strings.SplitSeq(s, ",") {
    process(part)
}
```

### strings.FieldsSeq Iteration

**Old pattern:**
```go
for _, field := range strings.Fields(s) {
    process(field)
}
```

**New pattern (Go 1.23+):**
```go
for field := range strings.FieldsSeq(s) {
    process(field)
}
```

---

## time.go

Time formatting and timer patterns.

See: [time package](https://pkg.go.dev/time)

### time.DateTime/DateOnly/TimeOnly Constants (Go 1.20+)

**Old pattern:**
```go
t.Format("2006-01-02 15:04:05")
t.Format("2006-01-02")
t.Format("15:04:05")
```

**New pattern:**
```go
t.Format(time.DateTime)
t.Format(time.DateOnly)
t.Format(time.TimeOnly)
```

### Timer Channel len() Checks (Go 1.23+)

**Broken pattern (Go 1.23+):**
```go
timer := time.NewTimer(1 * time.Second)
if len(timer.C) > 0 {  // Always false - channels are now unbuffered
    <-timer.C
}
```

**Correct pattern:**
```go
timer := time.NewTimer(1 * time.Second)
select {
case <-timer.C:
    // timer fired
default:
    // timer not yet fired
}
```

**Why:** In Go 1.23+, timer and ticker channels have capacity 0 (unbuffered).

### Deferred time.Since Bug

**Broken pattern:**
```go
func foo() {
    start := time.Now()
    defer log.Println(time.Since(start))  // Evaluated NOW, not at exit!
    // ... work ...
}
```

**Correct pattern:**
```go
func foo() {
    start := time.Now()
    defer func() { log.Println(time.Since(start)) }()
    // ... work ...
}
```

---

## slices.go

Slice operations using Go 1.21-1.23+ APIs.

See: [slices package](https://pkg.go.dev/slices), [maps package](https://pkg.go.dev/maps)

### sort.Ints/Strings/Float64s → slices.Sort (Go 1.21+)

**Old patterns:**
```go
sort.Ints(nums)
sort.Strings(strs)
sort.Float64s(floats)
```

**New pattern:**
```go
slices.Sort(nums)
slices.Sort(strs)
slices.Sort(floats)
```

### Slice Clone Patterns (Go 1.21+)

**Old patterns:**
```go
clone := append([]T(nil), original...)
clone := append([]T{}, original...)
```

**New pattern:**
```go
clone := slices.Clone(original)
```

### bytes.Clone (Go 1.20+)

**Old patterns:**
```go
clone := append([]byte(nil), original...)
clone := append([]byte{}, original...)
```

**New pattern:**
```go
clone := bytes.Clone(original)
```

### Backward Iteration (Go 1.23+)

**Old pattern:**
```go
for i := len(s) - 1; i >= 0; i-- {
    process(s[i])
}
```

**New pattern:**
```go
for i, v := range slices.Backward(s) {
    process(v)
}
```

### Map Keys/Values Collection (Go 1.23+)

**Old pattern:**
```go
keys := make([]string, 0, len(m))
for k := range m {
    keys = append(keys, k)
}
```

**New pattern:**
```go
keys := slices.Collect(maps.Keys(m))
// Or for sorted keys:
keys := slices.Sorted(maps.Keys(m))
```

---

## sync.go

Synchronization patterns.

See: [sync.WaitGroup.Go](https://pkg.go.dev/sync#WaitGroup.Go)

### WaitGroup.Go Pattern (Go 1.25+)

**Old pattern:**
```go
wg.Add(1)
go func() {
    defer wg.Done()
    doSomething()
}()
```

**New pattern:**
```go
wg.Go(func() {
    doSomething()
})
```

**Note:** Only flags simple patterns without closure parameters. Patterns with closure parameters cannot be directly converted since `wg.Go()` only accepts `func()`.

---

## builtins.go

Built-in function patterns.

### min/max Built-in Functions (Go 1.21+)

**Old pattern:**
```go
result := int(math.Min(float64(a), float64(b)))
result := int(math.Max(float64(a), float64(b)))
```

**New pattern:**
```go
result := min(a, b)
result := max(a, b)
```

### clear() Built-in Function (Go 1.21+)

**Old pattern:**
```go
for k := range m {
    delete(m, k)
}
```

**New pattern:**
```go
clear(m)
```

### Range Over Integer (Go 1.22+)

**Old pattern:**
```go
for i := 0; i < n; i++ {
    process(i)
}
```

**New pattern:**
```go
for i := range n {
    process(i)
}
```

### Append Without Values

**Broken pattern:**
```go
slice = append(slice)  // No effect
```

---

## reflect.go

Reflection patterns.

See: [reflect.TypeAssert](https://pkg.go.dev/reflect#TypeAssert), [reflect.PointerTo](https://pkg.go.dev/reflect#PointerTo), [reflect.TypeFor](https://pkg.go.dev/reflect#TypeFor)

### reflect.TypeAssert Pattern (Go 1.25+)

**Old pattern (allocates):**
```go
val := v.Interface().(string)
```

**New pattern (no allocation):**
```go
val := reflect.TypeAssert[string](v)
```

### reflect.PtrTo → reflect.PointerTo (Go 1.22+)

**Deprecated pattern:**
```go
ptrType := reflect.PtrTo(t)
```

**New pattern:**
```go
ptrType := reflect.PointerTo(t)
```

### reflect.TypeFor Pattern (Go 1.22+)

**Old pattern:**
```go
t := reflect.TypeOf((*MyType)(nil)).Elem()
```

**New pattern:**
```go
t := reflect.TypeFor[MyType]()
```

### Deprecated reflect.SliceHeader/StringHeader (Go 1.21+)

See: [unsafe.Slice](https://pkg.go.dev/unsafe#Slice), [unsafe.String](https://pkg.go.dev/unsafe#String)

**Deprecated pattern:**
```go
sh := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
hdr := (*reflect.StringHeader)(unsafe.Pointer(&str))
```

**New pattern:**
```go
slice := unsafe.Slice(ptr, len)
str := unsafe.String(ptr, len)
```

---

## random.go

Random number generation patterns.

See: [math/rand/v2](https://pkg.go.dev/math/rand/v2)

### math/rand/v2 Migration (Go 1.22+)

**Method renames:**
- `rand.Intn(n)` → `rand.IntN(n)`
- `rand.Int31()` → `rand.Int32()`
- `rand.Int31n(n)` → `rand.Int32N(n)`
- `rand.Int63()` → `rand.Int64()`
- `rand.Int63n(n)` → `rand.Int64N(n)`

### rand.Seed/Read Deprecation (Go 1.20+)

**Deprecated patterns:**
```go
rand.Seed(time.Now().UnixNano())  // Auto-seeded since 1.20
rand.Read(buf)                     // Use crypto/rand for security
```

**New patterns:**
```go
// For reproducibility, use a local source
r := rand.New(rand.NewSource(42))

// For cryptographic randomness
crypto_rand.Read(buf)
```

---

## testing.go

Testing utilities.

See: [testing.B.Loop](https://pkg.go.dev/testing#B.Loop), [testing.T.Context](https://pkg.go.dev/testing#T.Context)

### Benchmark b.Loop Pattern (Go 1.24+)

**Old pattern:**
```go
func BenchmarkFoo(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // work
    }
}
```

**New pattern:**
```go
func BenchmarkFoo(b *testing.B) {
    for b.Loop() {
        // work
    }
}
```

**Benefits:**
- Setup/cleanup executes only once per `-count`
- Compiler cannot optimize away the loop body

### Testing t.Context Pattern (Go 1.24+)

**Old pattern:**
```go
func TestFoo(t *testing.T) {
    ctx := context.Background()
    result, err := doSomething(ctx)
}
```

**New pattern:**
```go
func TestFoo(t *testing.T) {
    ctx := t.Context()
    result, err := doSomething(ctx)
}
```

**Benefits:**
- Context is automatically canceled when test completes
- Resources are released promptly on test failure

---

## net.go

Network and path utilities.

See: [net.JoinHostPort](https://pkg.go.dev/net#JoinHostPort), [filepath.IsLocal](https://pkg.go.dev/path/filepath#IsLocal)

### net.JoinHostPort Pattern

**Old pattern:**
```go
addr := fmt.Sprintf("%s:%d", host, port)
```

**New pattern:**
```go
addr := net.JoinHostPort(host, strconv.Itoa(port))
```

**Why:** `net.JoinHostPort` properly handles IPv6 addresses by wrapping them in brackets.

### filepath.IsLocal (Go 1.20+)

**Old pattern:**
```go
if strings.Contains(userPath, "..") {
    return errors.New("invalid path")
}
```

**New pattern:**
```go
if !filepath.IsLocal(userPath) {
    return errors.New("invalid path")
}
```

**Benefits:**
- Comprehensive path validation
- Handles OS-specific path separators
- Prevents directory traversal attacks

### Error Before Use Pattern

**Broken pattern:**
```go
f, err := os.Open(path)
name := f.Name()  // PANICS if err != nil
if err != nil { ... }
```

**Correct pattern:**
```go
f, err := os.Open(path)
if err != nil { ... }
name := f.Name()
```

---

## crypto.go

Cryptography patterns.

See: [crypto/cipher](https://pkg.go.dev/crypto/cipher), [crypto/ecdh](https://pkg.go.dev/crypto/ecdh)

### Deprecated Cipher Modes (Go 1.24+)

**Deprecated:**
- `cipher.NewOFB` - OFB mode
- `cipher.NewCFBEncrypter` - CFB encryption
- `cipher.NewCFBDecrypter` - CFB decryption

**New pattern:**
```go
// Authenticated encryption (preferred)
aead, _ := cipher.NewGCM(block)
ciphertext := aead.Seal(nil, nonce, plaintext, additionalData)

// Stream cipher without authentication
stream := cipher.NewCTR(block, iv)
```

### Weak RSA Key Size

**Weak pattern:**
```go
key, _ := rsa.GenerateKey(rand.Reader, 1024)  // Weak
key, _ := rsa.GenerateKey(rand.Reader, 512)   // Rejected in Go 1.24+
```

**Recommended:**
```go
key, _ := rsa.GenerateKey(rand.Reader, 2048)  // Minimum recommended
key, _ := rsa.GenerateKey(rand.Reader, 4096)  // For long-term security
```

### Deprecated crypto/elliptic Functions (Go 1.21+)

**Deprecated pattern:**
```go
import "crypto/elliptic"
key, _ := elliptic.GenerateKey(curve, rand.Reader)
```

**New pattern:**
```go
import "crypto/ecdh"
key, _ := ecdh.P256().GenerateKey(rand.Reader)
```

### rsa.GenerateMultiPrimeKey Deprecated (Go 1.21+)

**Deprecated:**
```go
key, _ := rsa.GenerateMultiPrimeKey(rand.Reader, nprimes, bits)
```

**Use instead:**
```go
key, _ := rsa.GenerateKey(rand.Reader, bits)
```

---

## runtime.go

Runtime function patterns.

See: [runtime.AddCleanup](https://pkg.go.dev/runtime#AddCleanup)

### SetFinalizer → AddCleanup (Go 1.24+)

**Old pattern:**
```go
runtime.SetFinalizer(obj, func(o *Type) { cleanup(o) })
```

**New pattern:**
```go
runtime.AddCleanup(obj, func(arg ArgType) { cleanup(arg) }, arg)
```

**Benefits of AddCleanup:**
- Multiple cleanups per object
- Can attach to interior pointers
- No cycle leaks (SetFinalizer can leak cycles)
- Doesn't delay object freeing

### GOROOT Deprecated (Go 1.24+)

**Deprecated:**
```go
root := runtime.GOROOT()
```

**New pattern:**
```bash
go env GOROOT
```

**Why:** `runtime.GOROOT()` may not reflect the actual GOROOT when the binary is moved or when using toolchains.

---

## Configuration

The rules are configured in `.golangci.yaml`:

```yaml
gocritic:
  settings:
    ruleguard:
      rules: "rules/*.go"
```

## Testing Rules

To test if a rule catches a pattern:

```bash
# Create test file
cat > /tmp/test.go << 'EOF'
package main
// ... test code ...
EOF

# Run linter (must clear cache after rule changes)
golangci-lint cache clean && golangci-lint run /tmp/test.go
```

## Adding New Rules

1. Create or update a file in `rules/` with `//go:build ruleguard` constraint
2. Import `github.com/quasilyte/go-ruleguard/dsl`
3. Write rule functions that take `dsl.Matcher` parameter
4. Add Go doc reference links in comments (e.g., `// See: https://pkg.go.dev/...`)
5. Document the old and new patterns clearly
6. Run `golangci-lint cache clean && golangci-lint run` to test

See [go-ruleguard documentation](https://go-ruleguard.github.io/by-example/) for pattern syntax.
