# Go Performance

Go tools and techniques for removing a measured bottleneck. The method (profile
first, rule out waiting on dependencies, fix the algorithm, change one thing at
a time, comment what an optimization buys) is coding-practices: Performance Claims Need Measurement. Everything below is a candidate to try after a profile
points at the code, not a default style.

## Contents

- Find Where the Time Goes
- Benchmark, Change, Compare
- Algorithms Before Constants
- Allocations
- Cache Work That Repeats
- CPU
- I/O and Network
- Runtime Settings
- Production Signals
- Common Mistakes

---

## Find Where the Time Goes

| Symptom                                 | Tool and signal                                                    |
| --------------------------------------- | ------------------------------------------------------------------ |
| High CPU                                | CPU profile: the function with the largest flat or cumulative time |
| Many allocations, high GC share         | Heap profile with `-sample_index=alloc_objects`                    |
| Memory keeps growing                    | Heap profile with `inuse_space`; compare two snapshots             |
| Latency high, CPU low                   | Goroutine profile: goroutines blocked in network or `database/sql` |
| Goroutine count rising                  | Goroutine profile: a stack that keeps accumulating                 |
| Throughput stops scaling with cores     | Mutex and block profiles                                           |
| Stages idle while waiting on each other | `go tool trace`                                                    |
| A value escapes to the heap             | `go build -gcflags=-m` (`-m -m` prints the reason)                 |
| A small helper is not inlined           | `go build -gcflags=-m`, look for `cannot inline`                   |

When goroutines mostly wait, the time is off-CPU and the CPU profile misses it:
the goroutine profile, the trace and the `logfx` spans and metrics
(errors-and-logging.md: Operations Telemetry) show which dependency is slow.

Getting profiles in `apps/services` (`*.out` is git-ignored; `tmp/` is the
scratch directory):

```bash
go test -run='^$' -bench=BenchmarkColored -cpuprofile=tmp/cpu.out -memprofile=tmp/mem.out ./pkg/ajan/logfx/
go tool pprof -http=:0 tmp/cpu.out
```

The service can expose `/debug/pprof/` through
`pkg/ajan/httpfx/modules/profiling`, only when the `profiling` setting is true
(default `false`), and it should then be guarded with the `PPROF_TOKEN`
environment variable. The endpoints reveal memory contents, requests in flight
and stacks: keep them off in production unless you are profiling, and otherwise
profile through benchmarks and tests (security-practices: Secrets and Personal Data Never Reach Output).

---

## Benchmark, Change, Compare

Scope: Every change made for performance

Rule:

1. Name the metric you are improving: latency, throughput, memory or CPU.
2. Write a benchmark that exercises one function, using `for b.Loop()` (Go
   1.24+), and report allocations.
3. Record a baseline, run the same command after one change, compare:

```bash
go test -run='^$' -bench=BenchmarkColored -benchmem -count=10 ./pkg/ajan/logfx/ > tmp/bench-1.txt
# apply one change
go test -run='^$' -bench=BenchmarkColored -benchmem -count=10 ./pkg/ajan/logfx/ > tmp/bench-2.txt
benchstat tmp/bench-1.txt tmp/bench-2.txt
```

`benchstat` (`golang.org/x/perf/cmd/benchstat`) is not in the `tool` block of
`go.mod`; adding it is a dependency change that needs the owner's approval
(tooling-standards: Adding a Dependency). Run benchmarks one after another on
an otherwise quiet machine: two benchmark runs sharing CPUs distort each
other, even when the variants were written in parallel. Keep the result files
under `tmp/` (git-ignored) and put the `benchstat` table in the task's report.
The module has one benchmark today (`BenchmarkColored` in
`pkg/ajan/logfx/colors_test.go`); add the benchmark before the change.

```go
func BenchmarkParseStory(b *testing.B) {
    body := sampleStoryBody(b)
    b.ReportAllocs()
    for b.Loop() {
        parseStory(body)
    }
}
```

---

## Algorithms Before Constants

| Pattern                                  | Cost            | Replace with                          |
| ---------------------------------------- | --------------- | ------------------------------------- |
| `slices.Contains` inside a loop          | O(n·m)          | Build a `map[T]struct{}` once: O(n+m) |
| Nested loops to match two collections    | O(n·m)          | Index one side by key in a map        |
| Repeated min, max or lookups on one list | O(n) per query  | Sort once, `slices.BinarySearch`      |
| `append` in a loop with a known size     | Repeated growth | `make([]T, 0, n)`                     |
| `s += part` in a loop                    | O(n²) copying   | `strings.Builder`                     |

Test the growth: benchmark with 100, 1,000 and 10,000 items; ten times the input
taking a hundred times as long means the algorithm is quadratic. In a hot path
prefer one direct loop over a chain of iterator adapters, each of which adds a
closure call per element.

---

## Allocations

Preallocation, cloning and map sizing are covered in
[data-structures.md](data-structures.md). In hot paths, also:

- **Reuse a slice's backing array** with `buf = buf[:0]` or
  `buf = append(buf[:0], items...)`, when nothing else still holds the old
  contents.
- **Take the value from `range`** (`for k, v := range m`) instead of indexing
  the map again inside the loop.
- **Return sentinel errors** for expected failures on a hot path; `fmt.Errorf`
  allocates on every call. Keep wrapping with context everywhere else
  (errors-and-logging.md: Error Handling).
- **Avoid boxing into `any`.** Passing a concrete value through `any` or a
  `...any` parameter usually allocates. Use a typed parameter or a type
  parameter.
- **Convert between `string` and `[]byte` once.** Each conversion copies. The
  `bytes` package mirrors most of `strings`, so work on `[]byte` directly when
  that is what you have.

### Pooling

`sync.Pool` keeps short-lived objects between uses. Use it for objects created
on every request or record (encoding buffers, scratch structs), and only after a
heap profile shows that allocation near the top:

```go
var bufferPool = sync.Pool{New: func() any { return new(bytes.Buffer) }}

const maxPooledBuffer = 64 << 10

func encodeFrame(frame *Frame) ([]byte, error) {
    buf := bufferPool.Get().(*bytes.Buffer)
    buf.Reset()
    defer func() {
        if buf.Cap() <= maxPooledBuffer { // a rare huge buffer is left to the GC
            bufferPool.Put(buf)
        }
    }()

    if err := frame.writeTo(buf); err != nil {
        return nil, fmt.Errorf("%w: %w", ErrFailedToEncodeFrame, err)
    }
    return bytes.Clone(buf.Bytes()), nil // pooled memory never leaves this function
}
```

- Reset an object when taking it out, and drop references it holds, so the pool
  does not keep large graphs alive.
- Never return pooled memory to a caller; return a copy.
- Cap what goes back into the pool. One oversized buffer, once pooled, keeps its
  memory for as long as the pool reuses it.
- Do not pool rarely created objects; the pool costs more than it saves.

### Memory Held by Small Views

A sub-slice or substring keeps its whole backing array alive; copy what you keep
from a large input (`bytes.Clone(frame[:16])`, `strings.Clone(line[:8])`).

---

## Cache Work That Repeats

Compile regular expressions and parse templates once (language.md: Regular
Expressions). For lazily computed values use `sync.OnceValue`. Any cache follows
coding-practices: Bound Every Resource (the repository cache has a TTL,
`storage.DefaultCacheTTL`).

When many callers miss the same cache key at once, they all do the same
expensive fetch. `golang.org/x/sync/singleflight` lets one of them fetch while
the others wait for its result. Use `DoChan` so each caller still honors its own
context, and give the shared fetch a context of its own: it must not fail
because the first caller gave up.

`golang.org/x/sync` is only an indirect dependency, so importing
`singleflight` directly goes through tooling-standards: Adding a Dependency
first. Without it, a map of in-flight fetches keyed by the cache key, guarded
by a mutex, does the same job. A cache holds data keyed by the revision or
version it was read at where staleness matters, so a newer revision misses
instead of serving stale state. The shared cache helper is `pkg/lib/caching`.

For eviction, a generic LRU such as `hashicorp/golang-lru/v2` (in `go.sum`
through another module) also needs that approval before code imports it.

---

## CPU

### Inlining

The compiler inlines small functions, which removes call overhead in tight
loops. A call to something that cannot be inlined, such as a logging call, can
push a small helper over the budget. Keep side effects out of hot helpers, and
check the decision:

```bash
go build -gcflags=-m ./pkg/api/business/stories/ 2>&1 | grep -E "can inline|cannot inline"
```

Receiver kind does not decide inlining: methods with pointer receivers inline as
well. Choose receivers by language.md: Value vs Pointer Parameters, not to chase
inlining.

### Logging in Hot Paths

A disabled log call still evaluates its arguments, and `...any` arguments are
boxed before the level check runs. In a loop, check the level once and build
attributes only when it is enabled; `LogAttrs` avoids the boxing (`logfx.Logger`
embeds `*slog.Logger`, so both are available):

```go
if logger.Enabled(ctx, slog.LevelDebug) {
    logger.LogAttrs(ctx, slog.LevelDebug, "item processed",
        slog.Int("index", i), slog.Duration("took", time.Since(start)))
}
```

### Reflection

`reflect` and `reflect.DeepEqual` inspect types at run time and allocate. In
production code compare with `slices.Equal`, `maps.Equal`, `bytes.Equal` or a
typed `Equal` method; keep `DeepEqual` for tests. Dispatch on a dynamic type
with one type switch instead of a series of type assertions.

### Memory Layout and Locality

- Struct field order is checked by `betteralign` (`go tool betteralign ./...`
  in `make check`, a hook in `.manifest.yml`); `make fix` applies it
  (`go tool betteralign -apply ./...`). Do not hand-reorder fields against it, and
  do not chase padding by hand in types that are not held in large slices or maps.
- Keep hot data contiguous: one slice for a matrix, nodes in a slice linked by
  index, a struct of slices when a hot loop reads one field of many structs.

### Last-Resort Techniques

Each of these makes code harder to read. Use one only with a benchmark on a
profiled hot path, and leave a comment with the numbers:

- **False sharing.** Counters updated by different goroutines on the same 64
  byte cache line slow each other down. Pad them apart (`_ [56]byte` after an
  `atomic.Int64`, or `cpu.CacheLinePad` from `golang.org/x/sys/cpu`). The sign
  is a concurrent benchmark that gets slower as goroutines are added.
- **Several accumulators.** A tight arithmetic loop with one accumulator waits
  on each previous addition; splitting it into four independent sums lets the
  CPU run them in parallel.
- **SIMD and assembly.** Rely on the compiler's vectorization and `math/bits`.
  `simd/archsimd` is a Go 1.26 experiment behind `GOEXPERIMENT=simd`; do not use
  it here. Hand-written assembly needs an ADR (`docs/adr`). Pick a CPU-specific path in a
  package-level variable initializer (`golang.org/x/sys/cpu`), never in
  `init()`, and keep the portable version tested.
- **`unsafe`.** Only for platform glue and layout (data-structures.md).

### The Scheduler

Goroutines are preempted asynchronously since Go 1.14. Do not add
`//go:noinline` or artificial calls to create preemption points.

---

## I/O and Network

### HTTP Clients and Servers

The zero `http.Client` has no timeout, and `http.Transport` keeps only 2 idle
connections per host by default, so concurrent calls to one service keep
opening new connections. Build each client once, in its owner's constructor or
the composition root, with a configured `Transport`, a `Timeout` or a
per-request context deadline, and `CheckRedirect` returning
`http.ErrUseLastResponse` where a redirect must not be followed. The shared
client is `appcontext.HTTPClient` (`pkg/ajan/httpclient`): its
`httpclient.Config` sets `max_idle_conns` (200), `max_idle_conns_per_host`
(100) and `idle_conn_timeout` (60s), plus retry and a circuit breaker. Prefer it
for outbound calls; an adapter that builds its own `http.Client` sets at least
`Timeout`. The server is built from `httpfx.Config`
(`pkg/ajan/httpfx/config.go`): `ReadHeaderTimeout`, `ReadTimeout`,
`WriteTimeout`, `IdleTimeout` and `MaxHeaderBytes` all have defaults and are
set on the `http.Server` (`httpfx/http_service.go`); a hand-built `http.Server`
sets the same fields (coding-practices: External Calls Have Timeouts).

A connection returns to the pool only after its body is read to the end and
closed. When the body does not matter, drain a bounded amount:

```go
defer resp.Body.Close()
_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 64<<10)) // best effort: reuse the connection
```

### Streaming and Buffering

Read large or unbounded input as a stream (coding-practices: Stream Large Data):
`io.Copy` between reader and writer, `bufio.Scanner` for lines, and
`json.NewDecoder` with `More()` for a large JSON array. `io.ReadAll` is fine for
input already bounded by `io.LimitReader` or `http.MaxBytesReader`.

Many small reads or writes cost one system call each. Wrap files and sockets in
`bufio.Reader` or `bufio.Writer`, and check the error of the final `Flush`.

### Batching

One round trip per item multiplies latency by the item count. Send many items
per call:

- PostgreSQL through pgx: `pgx.Batch` with `SendBatch` for many statements
  (`PgxConnection.SendBatch` in `pkg/ajan/connfx/adapter_pgx.go`), `CopyFrom` for
  bulk inserts. Both take `ctx` and return errors to check. A repository method that
  runs one sqlc query per row is the first thing to look for in a slow endpoint:
  write one query with `= ANY(sqlc.arg(ids)::TEXT[])` instead
  (`etc/data/default/queries/profile_points.sql`).
- HTTP: use the API's batch endpoint when it has one.
- Queue and worker calls: pass a whole list per message, not one item.

A consumer that batches a stream flushes on size and on time, and stops with its
context:

```go
func (w *AuditWriter) Run(ctx context.Context, events <-chan AuditEvent) error {
    batch := make([]AuditEvent, 0, w.batchSize)
    ticker := time.NewTicker(w.flushInterval)
    defer ticker.Stop()

    flush := func() error {
        if len(batch) == 0 {
            return nil
        }
        err := w.store.InsertEvents(ctx, batch)
        batch = batch[:0]
        return err
    }

    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case e, ok := <-events:
            if !ok {
                return flush()
            }
            batch = append(batch, e)
            if len(batch) == w.batchSize {
                if err := flush(); err != nil {
                    return fmt.Errorf("flush audit batch: %w", err)
                }
            }
        case <-ticker.C:
            if err := flush(); err != nil {
                return fmt.Errorf("flush audit batch: %w", err)
            }
        }
    }
}
```

### Concurrent Pipeline Stages

Run stages concurrently only when each is limited by a different resource (CPU,
disk, network), a trace shows one idle while another is busy, and the sequential
version with batching was measured and is slower. Bound the channels between
stages (coding-practices: Bound Every Resource) and run the stages under one
owner that waits for them (coding-practices: Owned Asynchronous Work).

---

## Runtime Settings

- `GOMAXPROCS` follows cgroup CPU limits since Go 1.25; do not set it by hand.
- When a Go container gets a memory limit (a deployment setting, not the Dockerfile), set `GOMEMLIMIT` to about 80-90% of it, so the GC works harder before the
  container is killed. `apps/services/Dockerfile` sets neither `GOMEMLIMIT` nor `GOGC` today.
- Leave `GOGC` at its default unless a profile shows GC time is the problem.
- Profile-guided optimization: a `default.pgo` CPU profile from production in
  the main package is used by `go build` automatically. Commit one only with a
  benchmark showing the gain, and refresh it when the code changes a lot.

---

## Production Signals

`runtime/metrics` exposes the runtime's own counters without a dependency. Read
them in a benchmark or a diagnostic run. Application metrics go through
`logger.NewMetricsBuilder` and are exported over OTLP (`EnableOTLP`) with the
rest of `logfx`. Watch trends, not single values:

| Metric                               | Watch for                                      |
| ------------------------------------ | ---------------------------------------------- |
| `/gc/heap/allocs:bytes`              | Allocation rate rising after a deploy          |
| `/gc/cycles/total:gc-cycles`         | GC frequency rising with steady load           |
| `/cpu/classes/gc/total:cpu-seconds`  | Share of CPU spent in the GC                   |
| `/memory/classes/heap/objects:bytes` | Live heap climbing under constant load: a leak |
| `/sched/goroutines:goroutines`       | Goroutines growing independently of traffic    |
| `/sched/latencies:seconds`           | Goroutines waiting long for a CPU              |
| `/sync/mutex/wait/total:seconds`     | Lock contention                                |

Request latency comes from the `httpfx` metrics and tracing middlewares. Continuous
profiling services (Pyroscope, Parca and similar) add infrastructure and a
dependency; they need an ADR first.

---

## Common Mistakes

| Mistake                                           | Fix                                                 |
| ------------------------------------------------- | --------------------------------------------------- |
| One benchmark run before and after                | `-count=10` and `benchstat`                         |
| A bare `http.Client` or `http.Server`             | Configure the transport, timeouts and pool once     |
| Not draining response bodies                      | Bounded `io.Copy(io.Discard, ...)`, then close      |
| `io.ReadAll` on unbounded input                   | Stream, or bound with `io.LimitReader`              |
| Log calls with computed arguments in a hot loop   | `Enabled` check and `LogAttrs`                      |
| `reflect.DeepEqual` in production code            | Typed equality                                      |
| `panic` and `recover` for control flow            | Return errors                                       |
| Returning pooled buffers to callers               | Return a copy                                       |
| Keeping a small sub-slice of a large buffer       | `bytes.Clone` or `strings.Clone`                    |
| `//go:noinline` to help the scheduler             | Nothing; preemption is asynchronous since Go 1.14   |
| `unsafe`, padding or assembly without a benchmark | Measure on the hot path, or leave it out            |
