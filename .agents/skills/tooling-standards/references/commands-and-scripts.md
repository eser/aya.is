# Commands and Scripts

Which tool runs what, which checks run where, where a project command is
defined and how a script is written. `make help` (root and `apps/services`)
prints every make target with one line each.

## Contents

- Tool per Place
- Checks and Git Hooks
- Entry Points
- Writing a Script
- Make Variables
- Build Output Stays Out of the Tree

---

## Tool per Place

Scope: Running, checking and installing anything in the repository

Rule: Each part of the repository has one toolchain; use it and do not swap or
re-wire a tool without the owner's approval:

| Place | Toolchain | Run |
| --- | --- | --- |
| Webclient (`apps/webclient`) | Node.js 24, pnpm (`packageManager` in the root `package.json`, activated by corepack), Vite, Vitest, TypeScript 7 native `tsc` | in `apps/webclient`: `pnpm dev` (port 3000), `pnpm build`, `pnpm lint` (`oxlint --type-aware`), `pnpm lint:fix`, `pnpm typecheck` (`tsc --noEmit`), `pnpm test`, `pnpm test:watch`, `pnpm test:update`, `pnpm test:ci` |
| Whole repository formatting | oxfmt (`.oxfmtrc.json`) | at the root: `pnpm fmt` writes, `pnpm fmt:check` verifies |
| Scripts (`scripts/*.ts`) | Node.js 24 | `pnpm lint:scripts` (`oxlint scripts`); run with `node scripts/<name>.ts` |
| Go services (`apps/services`) | Go from `go.mod` (`go 1.27.0`); tools pinned in `go.mod` and run as `go tool <name>` (golangci-lint, betteralign, govulncheck, sqlc, air) | in `apps/services`: `make dev`, `make run`, `make build`, `make test`, `make lint`, `make check`, `make fix`, `make ok`, `make generate`, `make migrate-up`, `make migrate-down` |
| Containers | Docker Compose (`compose.yml`, `compose.production.yml`) | at the root: `make up`, `make down`, `make restart`, `make logs`, `make cli`, `make build` |
| Git hooks and the full check | the `eser` CLI, installed with Deno from JSR | `eser codebase install` (hooks), `eser ok` (all checks), `eser fix` (checks with fixes); Deno is used for nothing else |
| Apple app (`apps/apple`) | Swift, Fastlane, its own `Makefile` | see `apps/apple/README.md` |

The root `pnpm dev`, `build`, `start`, `lint`, `typecheck` and `test` forward to
the webclient with `pnpm --filter @aya/webclient`. `pnpm install` runs once at
the root (the workspace is the root and `apps/webclient`); CI and the Docker
build use `pnpm install --frozen-lockfile`. Never run `npm`, `npx` or `yarn`,
and never `deno fmt`, `deno lint`, `deno test` or `deno task`: there is no
`deno.json`, and oxlint, oxfmt and Vitest are the checkers. The Nix flake
(`flake.nix`, loaded by direnv through `.envrc`) gives a shell with Go, Node,
pnpm, Deno and the Go tools; it is optional.

`golangci-lint` runs with `fix: true` in `apps/services/.golangci.yaml`, so
`make lint` and `make ok` rewrite files; use
`go tool golangci-lint run --fix=false ./...` when the tree must stay as it is
(`apps/services/AGENTS.md`, Linting and CI). oxfmt also writes by default;
verify with `pnpm fmt:check`.

**Why:** the lockfile belongs to pnpm, the formatter to oxfmt and the hooks to
`eser`; a second tool produces a different tree or a different verdict than
the one CI runs.

---

## Checks and Git Hooks

Scope: What must pass before a commit, a push and a merge

Rule: `.manifest.yml` is the single definition of the checks. Its `default`
workflow runs on `precommit` and `prepush`: the file validators
(end of file, trailing whitespace, line endings, large files, merge conflicts,
JSON, TOML, YAML, symlinks, secrets, file names, licenses), then `oxfmt`,
`oxlint`, `oxlint-scripts`, the webclient typecheck and tests (`pnpm test:ci`),
then the Go steps (`go vet`, `go build`, golangci-lint, `govulncheck`,
`betteralign`, `modernize`, `go test -failfast -race -count 1 ./...`). Its
`commit-validation` workflow runs on `commitmsg` and checks the one-line
conventional message (types `ci`, `chore`, `docs`, `feat`, `fix`, `perf`,
`refactor`, `revert`, `test`; workflow-practices has the commit rules). Steps
that can repair themselves carry a `fixCommand`, which `eser fix` runs.

The Integrity Pipeline (`.github/workflows/integrity.yml`) installs Node 24,
pnpm and Go, installs the `eser` CLI with Deno, and runs `eser ok`: the check
CI runs is the check that runs locally. `make ok` at the root runs the same
checks by hand (services `make ok`, then in the webclient `pnpm lint`,
`pnpm typecheck`, `pnpm test:ci`, then `pnpm fmt:check` and `pnpm lint:scripts`).
A new check is a step in `.manifest.yml`, not only a line in a document. Never
skip the hooks (`--no-verify`) to get a commit through.

**Why:** one definition cannot drift between the hook, CI and a developer's
habit.

---

## Entry Points

Scope: Project commands that people and CI run

Rule: A command belongs to the toolchain of the code it runs, and has one
documented entry:

- webclient and script checks are `pnpm` scripts in `apps/webclient/package.json`
  or the root `package.json`;
- Go and database commands are targets of `apps/services/Makefile`;
- container and release commands are targets of the root `Makefile`;
- cross-cutting checks are steps or `scripts:` entries of `.manifest.yml`
  (`eser ok`, `eser init`, `eser go-test`).

A new make target gets a `.PHONY` line and a `## description` after its name,
because `make help` prints exactly those. The target stays a short call; real
logic is a script under `scripts/` (`tag` runs `node scripts/push_current_release_tag.ts`).
Do not add a target or a `package.json` script that only renames an existing
command.

Correct:

```make
.PHONY: tag
tag: ## Creates and pushes a release tag from VERSION file.
	node scripts/push_current_release_tag.ts
```

Incorrect:

```make
lint-web:
	cd apps/webclient && pnpm lint # no .PHONY, no help text, duplicates `pnpm lint`
```

---

## Writing a Script

Scope: New or changed files under `scripts/`

Rule: A script is TypeScript that Node.js runs directly
(javascript-practices: Node Scripts has the code rules). It is started from a
make target, finds the repository root from its own file, checks its inputs
first and fails with a clear message and a non-zero status. File names are
snake_case (`push_current_release_tag.ts`, `sync_release_notes_from_changelog.ts`);
`scripts/` is exempt from the kebab-case rule (tooling-standards: File
Naming). A tool that needs a real parser, tests or more than a page of logic is
a tested module, not a long script.

---

## Make Variables

Scope: Arguments of a make target

Rule: The root `Makefile` forwards the words after the target to docker
compose through `ARGS` (`make up services` runs
`docker compose up --detach services`), and a catch-all `%:` rule keeps make
from treating them as targets. Reuse `$(ARGS)` for a new container target
instead of adding a variable per option. Every target documents what it needs
in its `##` text.

**Why:** one forwarding mechanism keeps `make up`, `make logs` and `make cli`
consistent, and the arguments are visible in the command line instead of
hidden in the environment.

---

## Build Output Stays Out of the Tree

Scope: Bundles, binaries and generated files

Rule: The webclient build writes `.output/` and the router plugin writes
`.tanstack/`; both are ignored (`apps/webclient/.gitignore`). `air` builds the
Go dev binary into `apps/services/tmp/`, ignored except `.gitkeep`. Never
commit or hand-edit build output. Generated files that are committed are
regenerated, never edited: `apps/webclient/src/routeTree.gen.ts` (by the
TanStack plugin) and `*_gen.go` under `apps/services/pkg/api/adapters/storage/`
(by `make generate`, which runs sqlc and the other generators).
