# Services (Go backend)

Scoped guidance for `apps/services`. The root `AGENTS.md` and the `architecture-guidelines` / `go-practices` skills still
apply; this file only adds facts that are specific to this directory.

## Layout

- `pkg/api/business/` — pure business logic (no external dependencies), `pkg/api/adapters/` — implementations of its
  interfaces, `pkg/ajan/` — framework utilities, `cmd/` — entry points that wire everything together.
- File names are `snake_case` (enforced by the `validate-filenames` hook).

## SQL and generated code

- Queries live in `etc/data/default/queries/`, migrations in `etc/data/default/migrations/`.
- Migrations use goose and every file has both `-- +goose Up` and `-- +goose Down` sections; write the Down migration.
- `make generate` runs the code generators (sqlc and others), `make migrate-up` / `make migrate-down` apply migrations.
- Never hand-edit generated `*_gen.go` files (for example under `pkg/api/adapters/storage/`); change the SQL and regenerate.

## Local configuration

- `.env.local`, `.env.*.local` and `config.local.json` are git-ignored development overrides. Never commit secrets.

## Linting and CI

- Linters run through `go tool` and are pinned in `go.mod`: `golangci-lint` (v2.2.2), `betteralign`, `govulncheck`.
  `make lint` / `make ok` run them, CI runs the same hooks from `.manifest.yml`.
- **`.golangci.yaml` has `fix: true`, so linting rewrites files in place.** Use
  `go tool golangci-lint run --fix=false ./...` whenever the tree must stay untouched.
- `golangci-lint` v2.2.2 needs `honnef.co/go/tools` v0.8.2 or newer: older versions panic on Linux with Go 1.27
  (`buildir: package "poll"`), but not on macOS. To reproduce CI locally use
  `GOOS=linux GOARCH=amd64 $(go tool -n golangci-lint) run --fix=false ./...` (`GOOS=linux go tool ...` would build the
  tool itself for Linux).
- CI pins the Go toolchain patch version (`.github/workflows/integrity.yml`); `govulncheck` fails on stdlib
  vulnerabilities fixed in later patch releases, so bump it when it reports new `Standard library` findings.
- `modernize` is run from `golang.org/x/tools/go/analysis/passes/modernize/cmd/modernize` (it moved out of gopls).
