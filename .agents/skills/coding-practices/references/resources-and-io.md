# Resources and I/O

Timeouts, limits, cleanup, streaming, background work and performance claims,
in Go and TypeScript. Go specifics: go-practices context, data-structures and
performance references.

## Contents

- External Calls Have Timeouts
- Bound Every Resource
- Release Resources Where They Are Acquired
- Stream Large Data
- Owned Asynchronous Work
- Performance Claims Need Measurement

---

## External Calls Have Timeouts

Scope: Every network request, subprocess and other call that leaves the process

Rule: Every external call has a timeout, so one slow dependency cannot hang the
caller. Retry only operations that are safe to repeat (reads, idempotent
writes), with exponential backoff and a maximum number of attempts, and never
retry a validation or authorization failure. Between attempts, check whether the
caller has cancelled (Go `ctx.Err()`, JavaScript `signal.aborted`) and wait
with a timer that cancellation interrupts, so a retry loop never outlives its
caller. A server's `Retry-After` is a floor for the wait, never shortened. A
caller that makes many calls to the same dependency bounds its concurrency.
Outbound HTTP in Go goes through `pkg/ajan/httpclient`, which already brings
retries with backoff and a circuit breaker; configure it instead of writing a
second retry loop.

Correct:

```go
// providerRequestTimeout bounds one provider call, so a stalled
// provider cannot hold the caller.
const providerRequestTimeout = 10 * time.Second

ctx, cancel := context.WithTimeout(ctx, providerRequestTimeout)
defer cancel()

req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
```

```go
cmd := exec.CommandContext(ctx, binaryPath, arguments...) // killed when ctx ends
```

```ts
const response = await fetch(url, { signal: AbortSignal.timeout(FETCH_TIMEOUT_MS) });
```

Incorrect:

```go
response, err := http.Get(url) // the default client waits forever on a stalled server
```

**Why:** without a bound, the slowest dependency sets the latency of everything
that calls it, and a stuck call holds its connection and memory indefinitely.

---

## Bound Every Resource

Scope: Pools, queues, buffers, caches, concurrency, request and file sizes

Rule: Everything that can grow has a limit: connection and worker pools, queue
and channel depths, in-memory caches, concurrent goroutines or promises, and the
size of what is read from a request, a file or a subprocess. Pick the limit on
purpose, make it a setting with a safe default when callers may need another
value (`httpfx.Config` does: `MaxConnections`, `MaxRequestSizeMB`,
`RateLimitRequests`), and decide what happens at the limit: wait, reject with an
error, or drop the oldest entry. A map keyed per client or per user evicts idle
keys; a map that only grows is a memory exhaustion hole.

Correct:

```go
// maxQueuedJobs bounds the jobs waiting for a worker; producers block beyond it.
const maxQueuedJobs = 100

jobs := make(chan Job, maxQueuedJobs)
for range config.Workers { // a setting with a named default
    wg.Go(func() { work(ctx, jobs) })
}
body := http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
```

Incorrect:

```go
for _, story := range stories {
    go summarize(ctx, story) // one goroutine per story, however many there are
}
```

**Why:** an unbounded resource grows until the process runs out of memory or
file handles, usually under the load where it matters most.

---

## Release Resources Where They Are Acquired

Scope: Files, sockets, locks, timers, subprocesses, request slots

Rule: Release a resource in the same scope that acquired it, set up right after
the acquisition succeeds: `defer` in Go, `try`/`finally` in TypeScript, `trap
... EXIT` in a shell script. A failure to close a writer can mean lost data,
so report it.

Correct:

```go
resp, err := client.Do(req)
if err != nil {
    return fmt.Errorf("fetch feed: %w", err)
}
defer func() { _ = resp.Body.Close() }()
```

```ts
const release = await acquireSlot(signal);
try {
  return await send(path, method, body, signal);
} finally {
  release();
}
```

```sh
work=$(mktemp -d)
trap 'rm -rf "${work:?}"' EXIT
```

Incorrect:

```ts
const release = await acquireSlot(signal);
const value = await send(path, method, body, signal); // a throw here keeps the slot forever
release();
```

**Why:** cleanup at the end of a function is skipped by the first early return
or thrown error someone adds later.

---

## Stream Large Data

Scope: Reading, transforming or transferring data whose size is not small and
fixed

Rule: Process large or unbounded data as a stream: rows one at a time, files,
request and response bodies through readers, results through lazy iterators.
Load all at once only when the size is small and bounded
(`io.LimitReader`, `http.MaxBytesReader`).

Correct:

```go
rows, err := tx.Query(ctx, selectChunks, generation)
if err != nil {
    return fmt.Errorf("read chunks: %w", err)
}
defer rows.Close()
for rows.Next() {
    // handle one chunk
}
return rows.Err()
```

```go
body, err := io.ReadAll(io.LimitReader(resp.Body, maxPageBytes)) // bounded
```

Incorrect:

```go
body, err := io.ReadAll(resp.Body) // unbounded: the whole upstream body in memory
```

**Why:** memory stays constant however large the input grows, and the first
result arrives before the last byte is read.

---

## Owned Asynchronous Work

Scope: Goroutines in Go, promises in TypeScript

Rule: Every asynchronous operation that is started is either awaited or owned
by something that waits for it and handles its error. A long-lived worker is
started by its owner, stops with the owner's context and is waited for at
shutdown (the owner's `Close`; the background workers of the application are
registered in the `workerfx.Registry` of `pkg/ajan/workerfx`). Work that is deliberately detached gets its
own named timeout and reports its failure (an event, a log entry).

Correct:

```go
var wg sync.WaitGroup
for range config.Workers {
    wg.Go(func() { work(ctx) }) // work logs its own failures
}
wg.Wait()
```

```ts
await Promise.all(files.map((file) => upload(file, signal)));
```

Incorrect:

```go
go s.index(ctx, document) // nobody waits, the error is dropped
```

```ts
files.forEach((file) => upload(file)); // rejections are lost
```

**Why:** an unowned promise loses its rejection; an unowned goroutine leaks and
hides its failure.

---

## Performance Claims Need Measurement

Scope: Changes made for performance, and any code or review comment claiming one

Rule: A performance claim comes with a measurement: `go test -bench` before and
after on the same machine, or for the webclient the measured timing of the
interaction or request in the browser. Report the numbers and the command.
Without a measurement, write the plain version.

- Profile before choosing what to change. If most of the time is spent waiting
  on a database, an upstream service or the disk, fix that component; faster
  code in the process will not show.
- Fix the algorithm before the constants: replacing a linear search in a loop
  with a map lookup beats any tuning of the loop.
- Change one thing at a time and repeat each run several times; a single run
  before and after cannot tell a small gain from noise.
- Comment an optimization that makes the code less obvious with what it avoids
  and the measured gain, so nobody reverts it as needless complexity.
- Check new loops over I/O for per-item calls that can be batched (an `await`
  inside a loop, a query per row), and check that buffers and caches have a
  bound.

The Go procedure (`benchstat`, pprof) is in go-practices: Benchmark, Change, Compare.

**Why:** intuition about bottlenecks is usually wrong, and an unmeasured
optimization costs readability for a gain nobody can confirm.
