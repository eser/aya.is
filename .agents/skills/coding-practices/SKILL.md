---
name: coding-practices
description: "Code quality rules for aya.is in Go and TypeScript: errors with stable identities, logging, explicit checks, identity comparison, naming, magic values, bounds and timeouts, background work, input validation, settings, performance claims. Use when writing or reviewing code, handling errors, adding logs, choosing a limit or calling an external service. Not for syntax (use go-practices or javascript-practices)."
---

# Coding Practices

Rules that hold in every language in this repository. Language idioms live in
go-practices and javascript-practices; function and module design in
design-principles.

## Always

- Handle every failure that can happen; catch narrowly and let the rest
  propagate. No handling for states the types rule out. Never ignore an error.
- Every failure has a stable identity (a sentinel or typed error in Go, a
  `code` on the wire) and names the failing check, the record and the reason;
  callers and tests match the identity, never the message text.
- An error is handled or returned, never both: return with context, log once
  where it is handled.
- People see words chosen from the code in their language, never an internal
  message, identifier or chain (AGENTS.md: Internationalization & Locale
  Fallback).
- Logs carry structured fields and ids, never credentials, session ids,
  bodies or personal text; Go logs through `logfx` (`log/slog`).
- In TypeScript, implicit truthy checks only on booleans; compare `null`,
  `undefined`, `""`, `0` and `.length` explicitly; `??` for defaults.
- Compare entities by IDs, never by slugs/usernames/strings.
- Take a value after a prefix with `CutPrefix` or `startsWith` plus
  `slice(prefix.length)`, never `replace(prefix, "")` or `slice(7)`.
- Limits, timeouts and durations are named constants with a reason and the unit
  in the name, in tests too; reuse the existing validator of a value before
  choosing a bound.
- Every external call has a timeout, every pool, queue and cache a limit, and
  every started goroutine or promise an owner.
- Release resources in the scope that acquired them (`defer`, `finally`).
- No absolute paths or machine-specific values; a setting is one `conf` field of
  the config struct, never a second `os.Getenv` name.
- Comments explain why; every `//nolint` and every other shortcut says why.
- Performance changes come with before and after measurements.

## References

| File                                                    | Read when                                                                                |
| ------------------------------------------------------- | ---------------------------------------------------------------------------------------- |
| [errors.md](references/errors.md)                       | Raising, wrapping, coding, matching or showing errors                                    |
| [logging.md](references/logging.md)                     | Adding or reviewing log calls, levels and fields                                         |
| [readability.md](references/readability.md)             | Naming, comments, magic values, checks, identity, prefixes, early returns, line breaks   |
| [resources-and-io.md](references/resources-and-io.md)   | Network calls, subprocesses, retries, limits, cleanup, streaming, async work, benchmarks |
| [inputs-and-config.md](references/inputs-and-config.md) | Parsing input, resolving paths, reading settings                                         |
