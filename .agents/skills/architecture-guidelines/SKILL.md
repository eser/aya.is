---
name: architecture-guidelines
description: "System design for aya.is: apps/services as hexagonal Go (pkg/api/business, pkg/api/adapters, pkg/ajan, cmd/ composition), the apps/webclient layout and backend object, settings, request identity, module and file conventions, testing expectations and ADRs. Use when designing a package, module or service, adding a port, adapter or setting, placing code, or writing an ADR. Not for Go idioms (use go-practices)."
---

# Architecture Guidelines

How the Go services and the webclient are structured and composed, and how a
lasting decision is recorded. Repository layout is in AGENTS.md: Project
Structure; Go-specific layout facts are in `apps/services/AGENTS.md`.

## Always

- A command under `apps/services/cmd/<program>` holds only its entry point and
  wiring; logic lives in packages under `apps/services/pkg`.
- Business logic lives in `pkg/api/business/<module>` and imports no adapter;
  the module declares the ports it needs as Go interfaces beside that logic.
- Adapters in `pkg/api/adapters/<name>` implement those ports; `pkg/ajan` holds
  framework utilities and imports nothing from `pkg/api`.
- `adapters/appcontext` (`AppContext.Init`, called from `cmd/`) constructs every
  adapter and service and passes dependencies as constructor parameters; no
  `init()` registration, no adapter chosen by a string inside a package.
- A module uses another module through its exported types or a port, never its
  unexported state; an import cycle is a design bug.
- Settings are `conf:"..."` struct tags with a safe `default:"..."` loaded
  through `pkg/ajan/configfx`; secrets stay out of committed files.
- Request identity travels in `context.Context` through a typed key; business
  inputs and dependencies never do.
- The webclient reaches the API only through the `backend` object, and route
  data through the query option factories (AGENTS.md: Backend Object Pattern,
  AGENTS.md: React Query with SSR).
- Pure logic goes in its own file, away from environment-dependent code
  (design-principles: Pure Functions).
- Use ES Modules (never CommonJS); webclient directories and files are
  kebab-case, Go files and module directories are snake_case.
- Tests fake a port, never a concrete adapter (qa-practices: Fakes at the
  Port); table-driven tests cover multiple cases (qa-practices: Table-Driven
  Tests).
- A lasting architectural decision gets an ADR in `docs/adr/` from
  `docs/adr/template.md`, drafted for the owner to decide.

## References

| File                                                              | Read when                                                                                              |
| ----------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------ |
| [hexagonal-architecture.md](references/hexagonal-architecture.md) | Placing Go code, adding a port, adapter, module, setting or flag, wiring `AppContext`, passing context |
| [project-structure.md](references/project-structure.md)           | Adding webclient files or API calls, naming files and directories                                      |
| [adrs.md](references/adrs.md)                                     | Recording an architectural decision                                                                    |
