# Logging

Where and how code logs. The Go form is go-practices: Logging Conventions;
what never reaches a log is in security-practices: Secrets and Personal Data
Never Reach Output.

---

## Structured Logging

Scope: Services, adapters, workers and request handlers. Pure functions and
shared utilities do not log; they return values or errors, and the caller logs.

Rule: Log with structured fields, never with `fmt.Println`, `log.Printf` or
string-built messages. Go logs through the `*logfx.Logger` of
`pkg/ajan/logfx`, which embeds `log/slog`, with the `*Context` methods and
`slog.String` / `slog.Any` attributes. Server-side webclient code (SSR,
middleware) reports a failure with `console.error` and a bracketed scope
(`[custom-domain] Failed to lookup: ...`); browser code turns a failure into
words on the screen (errors.md: Errors the UI Shows) and does not log content.
A build or verification script may print its result to standard output.

- A failure entry carries the component or scope, the operation and the
  identifiers involved (`story_id`, `profile_id`, a trace id), and the error as
  an attribute, so one line pinpoints it. Identifiers, not display names or
  e-mail addresses.
- Pick the level by what an operator must do and by layer. The storage layer
  logs only `Warn`, `Debug` and `Trace`: `Debug` for data operations, `Warn`
  for a failure the process survived (an unreachable provider, a record set
  aside). The business layer logs `Info` for a completed operation an operator
  follows (a session created, a translation stored). `Error` is for a failure
  an operator has to act on, logged where it is handled. Detail that is off in
  production is `Debug`; `Trace` is for step-by-step flow.
- An error is logged once, where it is handled (errors.md: Handle Each Error
  Once). A retry loop logs on a change of state or through a rate limit, not on
  every attempt.
- Never log credentials, tokens, session identifiers, request or response
  bodies, personal text or full URLs with their query; log ids and codes
  instead.

Correct:

```go
p.logger.WarnContext(ctx, "Failed to fetch presentation page for PDF URL",
    slog.String("url", presentationURL),
    slog.Any("error", err))
```

```go
logger.InfoContext(ctx, "user created", slog.String("user_id", userID), slog.String("operation", "createUser"))
```

Incorrect:

```go
log.Printf("job failed: %v", err)                                  // no structure, no identifiers
logger.InfoContext(ctx, "request", slog.String("body", string(body))) // content in the log
logger.Error("error")                                                 // no context
```

**Why:** an operator diagnosing a failure has only the log; entries without
fields cannot be filtered or joined, and an entry with content turns the log
into a copy of what the platform is trusted to protect.
