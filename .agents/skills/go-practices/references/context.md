# context.Context

`context.Context` carries cancellation, deadlines and request-scoped values
through a unit of work: one HTTP request, one worker run, one queue item. Every
operation that belongs to that unit receives the same context, so when the
caller gives up, everything it started stops. The general principle is
design-principles: Cancellation Flows With the Work; this file is its Go form.

## Contents

- Summary
- Cancellation
- Listening for Cancellation
- HTTP Handlers
- Outbound Calls and Databases
- Context Values
- Trace Correlation

---

## Summary

1. Pass the caller's context through the whole chain: handler, service,
   repository, outbound client. A link that starts a fresh context keeps working
   after the client has gone.
2. `ctx context.Context` is the first parameter. Do not store a request's
   context in a struct: a struct outlives the request that filled it. A
   long-lived worker may keep the context of its own lifetime
   (`processfx.Process.Ctx`).
3. Never pass `nil`; use `context.TODO()` where the caller has no context yet,
   which marks the gap to fix.
4. Call the `cancel` of every derived context on every path, normally with
   `defer cancel()` on the next line (`go vet` reports `lostcancel`).
5. `context.Background()` only at entry points: `main` (`cmd/serve/main.go`,
   `cmd/migrate/main.go`), the composition root, and tests through
   `t.Context()`.
6. Context values carry request-scoped metadata under a private key type, never
   business parameters.
7. Work that must outlive the request uses `context.WithoutCancel` plus its own
   named timeout.

---

## Cancellation

Work started for one owner shares its context and is waited for before the
owner returns. `processfx.Process.StartGoroutine` is that owner for the server:
it runs the function with the process context, logs an error unless the context
was cancelled, and `Process.Wait` blocks until shutdown. Inside a package,
`sync.WaitGroup.Go` starts workers that share the owner's context (coding-practices: Owned Asynchronous Work):

```go
// pkg/ajan/processfx/process.go
wg.Go(func() {
    err := fn(p.Ctx)
    if err != nil && p.BaseCtx.Err() == nil && !errors.Is(err, context.Canceled) {
        p.Logger.ErrorContext(p.BaseCtx, "Goroutine error", "name", name, "error", err)
    }
})
```

The number of workers is a setting with a named default, never one goroutine
per item. A bare `go func() { ... }()` in a handler needs a reason in a comment
and its own context (HTTP Handlers).

---

## Listening for Cancellation

A loop that waits on channels selects on `ctx.Done()` and returns `ctx.Err()`.
A loop that repeats work checks `ctx.Err()` between items
(`for ctx.Err() == nil { ... }`) and says how far it got
(`stopped after %d of %d items: %w`). A timeout set low in the chain is only a
ceiling: a child context never outlives its parent's deadline. For cleanup when
the context ends, use `context.AfterFunc` and call the returned `stop` when it
is no longer needed.

---

## HTTP Handlers

A handler reads the request context from `ctx.Request.Context()`
(`httpfx.Context.Request`); it ends when the client disconnects or the handler
returns. A middleware that adds a value replaces the request context with
`ctx.UpdateContext(...)`, derived from `ctx.Request.Context()`
(`httpfx/middlewares/auth_middleware.go`). Handlers do not create
`context.Background()`. When the work failed because the client left, there is
nobody to answer, so do not log it as a server failure.

Work that must finish after the response (a visit record, an audit entry, a
webhook that must answer quickly) uses `context.WithoutCancel`, which keeps the
request's values and drops the cancellation and the deadline. Give it its own
named timeout:

```go
// auditWriteTimeout bounds the write of one audit entry, so a slow store never
// holds the answer to the client.
const auditWriteTimeout = 2 * time.Second

func (s *Service) recordAfterResponse(ctx context.Context, entry Entry) {
    ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), auditWriteTimeout)
    defer cancel()

    err := s.repo.WriteAudit(ctx, entry)
    if err != nil {
        s.logger.WarnContext(ctx, "the audit entry was not recorded", slog.Any("error", err))
    }
}
```

`context.Background()` for the write would lose the request's values;
`ctx.Request.Context()` alone would cancel the write as soon as the handler
returns. Existing fire-and-forget code (`go telegram.Bot.HandleUpdate(context.Background(), &update)`
in `routes_telegram.go`) predates this rule; new code uses the pattern above.

---

## Outbound Calls and Databases

Every call that leaves the process takes the context, so a disconnect cancels
the whole downstream chain and the call is bounded (coding-practices: External Calls Have Timeouts).

- HTTP: `http.NewRequestWithContext(ctx, ...)`, never `http.NewRequest` or
  `http.Get`.
- Database: every sqlc query method and repository method takes `ctx` first,
  and the `database/sql` calls are the `*Context` variants (`QueryContext`,
  `ExecContext`); pgx methods take `ctx` first as well.
- Subprocesses: `exec.CommandContext(ctx, ...)`.

Correct:

```go
func (c *Client) Health(ctx context.Context) (Health, error) {
    request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/health", nil)
    if err != nil {
        return Health{}, fmt.Errorf("%w: %w", ErrFailedToBuildRequest, err)
    }

    response, err := c.httpClient.Do(request)
    if err != nil {
        return Health{}, fmt.Errorf("%w: %w", ErrUnavailable, err)
    }
    defer response.Body.Close()
    // ...
}
```

Incorrect:

```go
response, err := http.Get(c.baseURL + "/health") // ignores cancellation and deadline
rows, err := db.Query("SELECT ...", id)          // same for the query
```

---

## Context Values

Context values carry request-scoped metadata that crosses API boundaries: the
authenticated claims, the resolved client address, a heartbeat. Keys are a
private named type (`httpfx.ContextKey`, or an unexported empty struct such as
`heartbeatContextKey` in `pkg/ajan/processfx/heartbeat.go`), never a bare
string, and the package that
owns a key exposes typed accessors:

```go
type heartbeatContextKey struct{}

func ContextWithHeartbeat(ctx context.Context, hb HeartbeatSender) context.Context {
    return context.WithValue(ctx, heartbeatContextKey{}, hb)
}
```

Incorrect:

```go
ctx = context.WithValue(ctx, "identity", session) // any package can collide with "identity"
```

| Data                                                     | Context value?                          |
| -------------------------------------------------------- | --------------------------------------- |
| Request ids, trace correlation, the resolved client address | Yes                                  |
| The authenticated principal (JWT claims)                 | Yes                                     |
| Connection pools, HTTP clients, repositories             | No: inject through the composition root |
| Settings and flags                                       | No: pass explicitly or inject           |
| Business inputs (profile id, story content)              | No: function parameters                 |
| Logger                                                   | No: a constructor parameter             |

---

## Trace Correlation

`logfx.StartSpan` derives the span from the context it is given, and
`Logger.PropagatorExtract` / `PropagatorInject` move a trace parent in and out
of HTTP headers (errors-and-logging.md: Operations Telemetry). Pass the
context returned by `StartSpan` on, so logs written with the `*Context`
methods carry the same trace. Do not copy request headers into spans or
contexts wholesale.
