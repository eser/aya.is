---
name: go-practices
description: "Go conventions of apps/services: cmd/ and pkg/api layout, sentinel errors wrapped with %w, logfx logging, context, idioms, data structures, performance, sqlc/goose and golangci-lint. Use when writing or reviewing .go files, adding a package, port, query or migration, handling errors or logs, profiling, or fixing lint findings. Not for test mechanics (use qa-practices)."
---

# Go Practices

How Go is written in this repository: one module,
`github.com/eser/aya.is/services` (`apps/services`, `go 1.27.0`), with the
executables under `cmd/` and the code under `pkg/`. Test mechanics live in
qa-practices; scoped facts of the directory in `apps/services/AGENTS.md`.

## Always

- Binaries are `cmd/<executable>/` (`serve`, `migrate`, `manage`); business
  logic is `pkg/api/business/<module>/`, integrations are
  `pkg/api/adapters/<adapter>/`, framework code is `pkg/ajan/<name>fx/`. Never
  create `internal/`, `api/`, `configs/` or `test/` (architecture-guidelines: Hexagonal Architecture).
- Business packages import no adapter, driver or vendor SDK; they declare their
  ports as interfaces, and `appcontext` (the composition root) builds the
  adapters and hands them in.
- Go files are snake_case (`validate-filenames` hook). Never hand-edit a
  generated file (`*_gen.go`, `*_string.go`); change the SQL and run
  `make generate`.
- Failures are package-level sentinels (`ErrXxx`) wrapped with
  `fmt.Errorf("%w: %w", ErrSentinel, err)`. No `errors.New` or message-only
  `fmt.Errorf` inline, no `==`, no matching on message text: `errors.Is` and
  `errors.AsType`.
- Log through the injected `*logfx.Logger` with `*Context` methods and `slog`
  attributes, once, where the error is handled; never `fmt.Println`, secrets or
  personal data.
- `ctx context.Context` comes first; pass the caller's context down and cancel
  every context you derive; `context.Background()` only in `main` and the
  composition root.
- A constructor takes its dependencies (port interfaces, the logger, a
  `*Config`), returns the concrete type, opens nothing and reads no globals;
  clocks and id generators are fields with a production default.
- Numbers with a meaning are named constants with a comment saying why
  (coding-practices: Magic Values); struct literals set every field
  (`exhaustruct`).
- Use the standard library up to the `go` directive; JSON is `encoding/json`
  (`encoding/json/v2` needs `GOEXPERIMENT`, which the repository does not set).
- Run `make ok` in `apps/services` before handing over. `.golangci.yaml` has
  `fix: true`: use `go tool golangci-lint run --fix=false ./...` when the tree
  must stay untouched.

## References

| File                                                      | Read when                                                              |
| --------------------------------------------------------- | ---------------------------------------------------------------------- |
| [architecture.md](references/architecture.md)             | Adding a package or port, wiring, constructors, config, file names     |
| [errors-and-logging.md](references/errors-and-logging.md) | Defining, returning or matching errors, HTTP error bodies, logs, spans |
| [context.md](references/context.md)                       | Passing, deriving or detaching a context; handlers; outbound calls     |
| [language.md](references/language.md)                     | Strings, types, generics, parameters, enums, iterators, JSON, regexp   |
| [data-structures.md](references/data-structures.md)       | Choosing slices, maps, containers, buffers or generic types            |
| [performance.md](references/performance.md)               | Profiling, benchmarks, allocations, caching, I/O and batching          |
| [modernization.md](references/modernization.md)           | Replacing an old idiom, raising the `go` directive, Go tooling         |
| [tooling.md](references/tooling.md)                       | Make targets, lint rules, Go test shape, sqlc, goose, generated files  |

General rules (errors, logging, timeouts, config, interfaces) live in
coding-practices and design-principles; tests in qa-practices; the layer rules
in architecture-guidelines.
