# Go Language Conventions

Day-to-day Go idioms this repository settles on: scoping, iteration, strings,
types, parameters, enums, resources, iterators, regexp and JSON.

## Contents

- Shared State Synchronization
- Scope Variables to the `if`
- Prefer `range` for Iteration
- Strings: Conversion, Formatting and Building
- Types: Narrow Conversions and Generics
- Value vs Pointer Parameters
- Small Interfaces and Focused Functions
- Enums
- Release Resources Right After Acquiring
- Iterators and Streaming
- Regular Expressions
- JSON Encoding
- Use the Current Standard Library

---

## Shared State Synchronization

Scope: Go code with goroutines (workers, caches, routers, circuit breakers)

Rule: State that more than one goroutine reads or writes is guarded by a
`sync.Mutex`/`sync.RWMutex`, owned by a single goroutine behind a channel, or
held in `sync/atomic` types. Document which one guards each field. Passing
`go test -race` is required (`make test` runs it), but it does not prove the
absence of races on paths the tests do not run. `CircuitBreaker` in
`pkg/ajan/httpclient/circuit_breaker.go` guards its state with a `sync.Mutex`.

Correct:

```go
type pendingApprovals struct {
    mu      sync.Mutex
    pending map[string]*approval // guarded by mu
}
```

Incorrect:

```go
type pendingApprovals struct {
    pending map[string]*approval // written by the relay and every prompt answer
}
```

---

## Scope Variables to the `if`

Scope: All Go code

Rule: When a value is only needed for one check, declare it in the `if`
statement itself (`if err := os.MkdirAll(dir, 0o750); err != nil {`), so it
cannot be reused by mistake further down. Declare it outside only when the code
after the check uses it.

---

## Prefer `range` for Iteration

Scope: All Go code

Rule: Iterate with `range` rather than an index loop, and use `for i := range n`
(Go 1.22+) for plain counting. Reach for a three-part `for` only when the index
arithmetic is the point, such as stepping by two or walking backwards.

---

## Strings: Conversion, Formatting and Building

Scope: All Go code

Rule:

- One value to or from a string: `strconv`, not `fmt.Sprintf("%d", n)`.
- A host and port: `net.JoinHostPort`, which brackets IPv6; never
  `fmt.Sprintf("%s:%d", ...)`.
- A user-supplied string in an error or log: `%q`, so empty values, spaces and
  control characters are visible (`slug ""` instead of a bare gap).
- A string built in a loop: `strings.Builder`.
- Bytes that are mutated or moved through I/O stay `[]byte`; convert to `string`
  once at the boundary, not per call.

Correct:

```go
port := strconv.Itoa(config.Port)
address := net.JoinHostPort(config.Host, port) // brackets IPv6: [::1]:8443

if !validSlug(slug) {
    return fmt.Errorf("%w: slug %q", ErrInvalidSlug, slug)
}

var b strings.Builder
for _, host := range hosts {
    b.WriteString(host)
    b.WriteByte('\n')
}
```

Incorrect:

```go
port := fmt.Sprintf("%d", config.Port) // strconv.Itoa is the direct tool

return fmt.Errorf("%w: slug %s", ErrInvalidSlug, slug) // "" is invisible

var out string
for _, host := range hosts {
    out += host + "\n" // copies the whole string every iteration
}
```

---

## Types: Narrow Conversions and Generics

Scope: All Go code

Rule: Keep values in the narrowest type that fits and convert explicitly.

- A type assertion uses the two-value form (`v, ok := x.(T)`) unless a failure
  is a programming error that should panic.
- A conversion to a smaller integer type checks the range first; `int32(n)`
  silently wraps on overflow.
- When a function works on several concrete types, write it with a type
  parameter and a constraint, not with `any`. The compiler then checks every
  call and the caller gets its own type back without an assertion. Use `any`
  only for values whose type is truly open, such as decoded JSON.
- Check `slices`, `maps` and `cmp` before writing a helper.

Correct:

```go
func Decode[T any](data []byte) (T, error)

story, err := Decode[Story](data) // story is a Story

if n > math.MaxInt32 || n < math.MinInt32 {
    return fmt.Errorf("timeout %d does not fit in int32", n)
}
timeout := int32(n)
```

Incorrect:

```go
func Decode(data []byte) (any, error)

story := decoded.(Story) // panics on a mismatch

timeout := int32(n) // wraps silently
```

---

## Value vs Pointer Parameters

Scope: Function parameters in Go (method receivers are a separate decision)

Rule: Pass by value by default, and take a pointer only for a reason the
signature should communicate. A pointer parameter tells the reader that the
function may change the value or that nil is a valid input; using one without
that reason sends the wrong signal.

Take a pointer when:

- the function mutates the caller's value;
- nil means something (an optional field in a partial update);
- the type must not be copied: it contains a `sync.Mutex`, `sync.WaitGroup` or
  similar (`go vet` reports these as copylocks), or its methods use pointer
  receivers;
- the struct is large (about 80 bytes or more), so copying it on every call
  costs more than the indirection.

Pass by value when:

- the type is a scalar or a small value type: `string`, `int`, `bool`,
  `float64`, `time.Time`, `time.Duration`. A string is already a pointer and a
  length, so `*string` only adds a nil case;
- the type is already a header over shared data: slice, map, channel, function
  or interface. A pointer to them is needed only to replace the header itself,
  and returning the new value is usually clearer;
- the function only reads a small struct.

Do not reach for a pointer "to save memory". Copying a small struct is cheap and
stays on the stack, while taking the address of a local value can make it escape
to the heap, which costs an allocation and GC work. Each dereference is also a
possible cache miss. For small read-only structs a value is usually faster; for
large structs or mutation a pointer wins. When the difference matters, measure
it (`go test -bench`, `go build -gcflags=-m` for escapes) as coding-practices: Performance Claims Need Measurement requires.

Correct:

```go
// Small value types by value
func storyLabel(id, title string, publishedAt time.Time) string

// Pointer because the function fills in the caller's config
func applyDefaults(config *Config)

// Pointer because nil means "leave the limit unchanged"
func updateStoryLimit(ctx context.Context, id string, limit *int) error

// Pointer because Router holds a sync.RWMutex and must not be copied
func routeNames(router *Router) []string

// Slice by value: it is already a header over shared data
func countActive(profiles []Profile) int
```

Incorrect:

```go
func storyLabel(id *string) string // no mutation, nil is not meaningful

func countActive(profiles *[]Profile) int // pointer to a slice header

func routeNames(router Router) []string // copies the mutex
```

---

## Small Interfaces and Focused Functions

Scope: All Go code

Rule: Write short functions with one responsibility, and prefer composition
over inheritance-style layering. Declare small, purpose-specific interfaces
(one to three methods) where they are used, and let public functions take
interfaces, not concrete types. Do not add a wrapper that only forwards a call:
call the method directly. Prefer the standard library over a third-party
package (the `depguard` list in `.golangci.yaml` makes each new dependency a
visible decision). A large interface is allowed only as one deliberate port
(`//nolint:interfacebloat // reason`, architecture.md: Modules, Ports and
Adapters).

Correct:

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

func ProcessData(r Reader) error {
    // works with any Reader implementation
}
```

Incorrect:

```go
type EverythingInterface interface {
    Read() error
    Write() error
    Delete() error
    Update() error
    // every consumer must implement and fake all of it
}
```

---

## Enums

Scope: `iota` constants

Rule: The zero value of a Go variable silently becomes the first enum member.
Start the enum at 1, or make 0 an explicit `Unknown`/`Unspecified` member, so an
unset value is detectable. Keep 0 as a real member only when it is the correct
default, and say so in a comment. Give the type a `String()` method (the repo generates it
with `//go:generate go tool stringer -type X -trimprefix X`, as
`pkg/ajan/results/types.go` does), and switch over it exhaustively
(coding-practices: Exhaustive Branching).

Correct:

```go
type KeyAlgorithm int

const (
    KeyAlgorithmUnknown KeyAlgorithm = iota // zero value: not set
    KeyAlgorithmEd25519
    KeyAlgorithmP256
)

type Outcome uint8

const (
    Completed Outcome = iota + 1 // starts at 1, so an unset outcome is detectable
    Rejected
    Failed
)
```

Incorrect:

```go
const (
    KeyAlgorithmEd25519 KeyAlgorithm = iota // an unset algorithm reads as Ed25519
    KeyAlgorithmP256
)
```

---

## Release Resources Right After Acquiring

Scope: Files, response bodies, rows, locks, tickers, listeners

Rule: Put `defer x.Close()` (or `Unlock`, `Stop`, `cancel`) on the line right
after a successful open, before any other code. A later edit that adds an early
return cannot skip it. For a writer, the error from `Close` or `Flush` can mean
lost data: return it instead of dropping it. Use `runtime.AddCleanup` rather
than `runtime.SetFinalizer` for the rare object whose cleanup the GC must
trigger. The general rule is coding-practices: Release Resources Where They Are Acquired. A file that a rename or another commit point depends on is
synced (`f.Sync()`) and its `Close` error checked before that point.

Correct:

```go
f, err := os.Open(path)
if err != nil {
    return fmt.Errorf("open %s: %w", path, err)
}
defer f.Close()

func writeReport(path string, data []byte) (err error) {
    f, err := os.Create(path)
    if err != nil {
        return fmt.Errorf("create %s: %w", path, err)
    }
    defer func() {
        if cerr := f.Close(); cerr != nil && err == nil {
            err = fmt.Errorf("close %s: %w", path, cerr)
        }
    }()
    _, err = f.Write(data)
    return err
}
```

Incorrect:

```go
f, _ := os.Open(path)
data := parse(f) // an early return added here later leaks f
f.Close()
```

---

## Iterators and Streaming

Scope: Code that walks large or unbounded data

Rule: Produce and consume large sequences lazily. Return an `iter.Seq`/
`iter.Seq2` (Go 1.23+) instead of a fully built slice when a caller may stop
early or the data is large, and move bulk data with streams (`io.Copy`,
`bufio.Scanner`, `rows.Next`) instead of reading it all into memory. Memory then
stays constant however large the input grows. The general rule is
coding-practices: Stream Large Data.

Correct:

```go
func (r *journalReader) Events() iter.Seq2[Event, error] {
    return func(yield func(Event, error) bool) {
        scanner := bufio.NewScanner(r.file)
        scanner.Buffer(make([]byte, 0, initialLineBytes), maxEventBytes)
        for scanner.Scan() {
            var event Event
            err := json.Unmarshal(scanner.Bytes(), &event)
            if !yield(event, err) {
                return
            }
        }
        if err := scanner.Err(); err != nil {
            yield(Event{}, err)
        }
    }
}
```

Incorrect:

```go
data, err := os.ReadFile(journalPath) // loads the whole journal to read the last event
```

---

## Regular Expressions

Scope: All Go code using `regexp`

Rule: Compile a constant pattern once, in a package-level variable with
`regexp.MustCompile`; a panic there is a bug found at start-up, not an expected
error. Compile a pattern built from input with `regexp.Compile` and handle the
error. Never compile inside a loop or a hot function. The same holds for
`text/template` and `html/template`: parse once, at start-up or in a
`sync.OnceValues`, and execute many times.

Correct:

```go
// namePattern is the form of the name after the prefix of a key.
var namePattern = regexp.MustCompile(`^[a-z][a-z0-9_]{0,62}$`)
```

Incorrect:

```go
func valid(name string) bool {
    return regexp.MustCompile(`^[a-z][a-z0-9_]{0,62}$`).MatchString(name) // compiles every call
}
```

---

## JSON Encoding

Scope: All Go JSON operations

Rule:

- Use `encoding/json`. Every JSON file in the module imports it and none uses
  `encoding/json/v2`, which needs `GOEXPERIMENT=jsonv2`; neither the Makefile,
  the Dockerfile nor the workflows set it (modernization.md: JSON: Stay on
  `encoding/json` for Now).
- JSON keys are snake_case (`activity_kind`, `external_activity_uri`), written
  in the struct tags; sqlc generates the same style (`json_tags_case_style:
  "snake"` in `sqlc.yaml`). Keys come from the tags, never from hand-written key
  lists. `tagliatelle` is disabled, so review enforces this.
- A JSONB column (`properties`) decodes into a typed struct with json tags
  (`ActivityProperties` in `pkg/api/business/stories/types.go`) rather than
  `map[string]any`, so a misspelled key is a compile error.
- Use `omitempty` for strings, slices, maps and pointers, and `omitzero` (Go
  1.24+) for fields whose zero value is left out and `omitempty` does nothing:
  `time.Time` and struct values.
- Decode a request body with `json.NewDecoder` and bound untrusted bodies
  (`http.MaxBytesReader`); stream large payloads instead of reading them into
  memory first. Check the error of every `Unmarshal` and `Decode`.
- Wire bodies of the HTTP API answer with `httpfx.WithJSON`
  (errors-and-logging.md: Errors on the Wire).

Correct:

```go
type PolicyRetryView struct {
    CanRetry   bool      `json:"can_retry"`
    RetryAfter time.Time `json:"retry_after,omitzero"`
}

var props ActivityProperties
err := json.Unmarshal(raw, &props)
if err != nil {
    return fmt.Errorf("%w: %w", ErrInvalidProperties, err)
}
```

Incorrect:

```go
var state map[string]any
_ = json.Unmarshal(data, &state) // untyped, error dropped

type PolicyRetryView struct {
    RetryAfter time.Time `json:"retryAfter,omitempty"` // camelCase key; zero time still written
}
```

**Why:** the same record is written by the Go service and read by the web
client; one key style and typed structs make a disagreement a compile or decode
error instead of a silently lost field.

---

## Use the Current Standard Library

Scope: All Go code

Rule: Use the standard library feature that replaced an older idiom or a
third-party package, up to the Go version in the module's `go` directive
(`go 1.27.0` in `apps/services/go.mod`), and nothing newer. The replacement table and the upgrade
procedure are in [modernization.md](modernization.md). Run `go fix ./...` to
apply the automatic rewrites.
