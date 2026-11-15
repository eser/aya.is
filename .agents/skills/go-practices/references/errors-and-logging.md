# Go Errors, Logging and Telemetry

How Go code in `apps/services` defines, returns and matches errors, logs
through the injected `*logfx.Logger` (a `log/slog` logger) and records spans
and metrics. The general rules are in coding-practices (Error Wrapping, Handle
Each Error Once, Structured Logging); this file is their Go form.

## Contents

- Error Handling
- Coded Errors
- Errors on the Wire
- Logging Conventions
- Operations Telemetry
- Security and Resilience

---

## Error Handling

Scope: All Go code

Rule:

- Check every returned error; never discard one with `_` except on a cleanup
  path whose comment says why (`defer func() { _ = c.Request.Body.Close() }()`).
- An expected condition a caller branches on, and every failure a service
  reports, is a sentinel declared at package level in the business package's
  `errors.go` (`var ErrNoAccount = errors.New("account not found")`). Never
  `return errors.New("...")` or a message-only `fmt.Errorf("...")` inline:
  `err113` is on in `.golangci.yaml`, and a caller cannot match a string.
- Wrap the cause under a sentinel with `fmt.Errorf("%w: %w", ErrSentinel, err)`.
  Both are `%w`, so `errors.Is` finds the sentinel and the cause. A sentinel
  that needs data adds it after it: `fmt.Errorf("%w (type=%s)", ErrNotStruct, name)`.
  `wrapcheck` is on, so an error from another package is wrapped before it is
  returned; a bare `return err` loses the context.
- Sentinel messages are lowercase, without trailing punctuation, because they
  are joined into longer messages.
- Match with `errors.Is` for sentinels and `errors.AsType[T](err)` for types
  (`errors.As` is acceptable where the code already uses it). Never compare
  with `==`, use a bare type assertion on a wrapped error, or look at the
  message text (`strings.Contains(err.Error(), ...)`), in tests too.
- Combine independent failures (cleanup after an error, a batch) with
  `errors.Join`.
- Return an error for anything a caller can expect or handle. `panic` is only
  for a bug that should never happen, or a failed start-up that `main` cannot
  continue from.
- `return nil, nil` is for a lookup where "not found" is a normal answer; mark
  it `//nolint:nilnil` and document the case on the function.

Correct:

```go
var (
    ErrNoAccount          = errors.New("account not found")
    ErrFailedToGetAccount = errors.New("failed to get account")
)

func (s *Service) account(ctx context.Context, userID string) (*Account, error) {
    account, err := s.repo.GetAccount(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("%w: %w", ErrFailedToGetAccount, err)
    }

    if account == nil {
        return nil, ErrNoAccount
    }

    return account, nil
}

if errors.Is(err, privacy.ErrNoAccount) { /* answer 404 */ }
```

Incorrect:

```go
result, _ := operation()                         // error discarded
if err != nil { return err }                     // no context added
if err != nil { return errors.New("failed") }    // original error lost, inline sentinel
if strings.Contains(err.Error(), "not found") {} // matches text, breaks on rewording
```

---

## Coded Errors

Scope: Failures a caller branches on, and failures that carry data

Rule: One failure class is one sentinel (`ErrExportInProgress`). A failure that
carries data the caller needs is a small typed error with an `Error()` method,
matched with `errors.AsType[*T]` (or `errors.As`) and never by message:
`privacy.CooldownError{NextRequestAt}`, `privacy.OwnedProfilesError{Profiles}`
(`pkg/api/business/privacy/errors.go`). The stable, machine-readable identity of
a failure that reaches a client is the `code` string the HTTP adapter assigns
(Errors on the Wire), not the message, because messages and translations change.
Name the specific failure; avoid catch-all sentinels such as `ErrFailed` for
many conditions.

```go
// CooldownError reports that a new export may not be requested yet.
type CooldownError struct {
    NextRequestAt time.Time
}

func (e *CooldownError) Error() string {
    return "a new export can be requested after " + e.NextRequestAt.Format(time.RFC3339)
}

var cooldown *privacy.CooldownError
if errors.As(err, &cooldown) { /* answer 429 with cooldown.NextRequestAt */ }
```

A test asserts the sentinel or the type (`errors.Is(err, ErrX)`), not the text.

**Why:** the sentinel is the contract between the layers; text can be reworded
without breaking a caller.

---

## Errors on the Wire

Scope: HTTP handlers in `pkg/api/adapters/http/` and answers to other processes

Rule: The business layer returns sentinels; the HTTP adapter maps them to a
status and a stable `code` with `errors.Is` and `errors.As` in one `switch`
(`routes_privacy.go`). The body is `{"error": "<fixed message>", "code": "<snake_case_code>"}`
through `httpfx.WithJSON`, and a missing resource answers
`httpfx.WithErrorMessage("not found")` without detail. The `default` branch logs
the real error and answers `httpfx.WithSanitizedError(err)`, so the chain, with
its paths, hosts and identifiers, stays in the log. Clients branch on `code`,
never on `error`.

```go
switch {
case errors.Is(err, privacy.ErrExportInProgress):
    return ctx.Results.Error(http.StatusConflict, httpfx.WithJSON(map[string]string{
        "error": "an export is already in progress", "code": "export_in_progress",
    }))
default:
    logger.Error("privacy request failed", "error", err)

    return ctx.Results.Error(http.StatusInternalServerError, httpfx.WithSanitizedError(err))
}
```

**Why:** an error chain can carry paths, SQL and request text; a fixed message
per code leaks nothing and still tells the client what to do
(security-practices: Error Sanitization).

---

## Logging Conventions

Scope: All logging statements

Rule: Log through the `*logfx.Logger` the constructor received, with the
`*Context` methods (`ErrorContext`, `InfoContext`, ...) and `slog` attributes;
`logfx` is built on `log/slog` and adds trace correlation and OTLP export. Never
use `fmt.Println`, `log.Printf` or a third-party logger. `*logfx.Logger` is
injected, not a global; `slog.Default()` is a fallback only where no logger
reaches the code.

- Repository layer: only `warn`, `debug` and `trace`.
- Service layer: successful operations at `info`; a failure that is handled
  and not returned at `error`/`warn`.
- Put the failure in the `error` attribute (`slog.Any("error", err)`) and the
  identifiers in their own attributes (`slog.String("export_id", id)`), so a
  log line is searchable without parsing the message.
- An error is logged once, where it is handled (coding-practices: Handle Each Error Once): either log it and carry on, or wrap and return it.
- Never log secrets, tokens, cookies, passwords or personal data
  (security-practices: Secrets and Personal Data Never Reach Output).
- In hot loops, see performance.md: Logging in Hot Paths.

Correct:

```go
// Service layer: a failure handled here, with identifiers as attributes
s.logger.ErrorContext(ctx, "failed to mark unscheduled export as failed",
    slog.String("export_id", export.ID),
    slog.Any("error", failErr))

// Repository layer: debug only
s.logger.DebugContext(ctx, "fetching user from database", slog.String("user_id", id))
```

Incorrect:

```go
s.logger.Info("executing query")                       // wrong level for the repository layer
s.logger.Error("password is: " + password)             // a secret in the log
fmt.Println("user created")                            // unstructured, no level, no correlation
```

---

## Operations Telemetry

Scope: Spans and metrics

Rule: Tracing and metrics go through `logfx`, which wraps OpenTelemetry and
exports through OTLP (`EnableOTLP`): `logger.StartSpan(ctx, name, attrs...)`
(`httpfx/middlewares/tracing_middleware.go`), `logger.NewMetricsBuilder(name)`
(`httpfx/http_service.go`), and `PropagatorExtract` / `PropagatorInject` to
carry a trace across a process boundary. Business and adapter packages do not
import `go.opentelemetry.io/otel` directly; only the `pkg/ajan` framework does.

- Start spans at service boundaries and pass the span's context on, so logs,
  spans and metrics of one request correlate.
- Attach identifiers that help diagnosis (record id, operation, error); never
  attach secrets, tokens or content.
- A new metric or span name is a fixed string, never built from user input.

Correct:

```go
ctx, span := logger.StartSpan(ctx, "HTTP Request", attrs...)
defer span.End()
```

**Why:** one wrapper keeps the exporter, the sampling and the attribute rules in
one place, and request data cannot reach a span by accident.

---

## Security and Resilience

Scope: Calls that leave the process and input that enters it

Rule: Validate and sanitize every external input at the adapter
(security-practices has the rules). Every outbound call has a timeout and a
bound on retries. The shared client is `appcontext.HTTPClient`
(`httpclient.NewClient`, `pkg/ajan/httpclient`): retry with exponential
backoff and a circuit breaker, configured in `httpclient.Config`
(`retry_strategy`, `circuit_breaker`, `transport`). A one-off client built by an
adapter (`coolify`, `resend`) sets at least `Timeout` and carries
`//nolint:exhaustruct // only Timeout needed`. Inbound request limits are the
`httpfx` rate-limit middleware (`pkg/ajan/httpfx/middlewares/rate_limit_middleware.go`).

Correct:

```go
// pkg/api/adapters/appcontext/appcontext.go
a.HTTPClient = httpclient.NewClient(httpclient.WithConfig(&a.Config.HTTPClient))

client := &http.Client{Timeout: 10 * time.Second} //nolint:exhaustruct // only Timeout needed
```

Incorrect:

```go
resp, err := http.Get(url) // no timeout, no retry, ignores the caller's context
```
