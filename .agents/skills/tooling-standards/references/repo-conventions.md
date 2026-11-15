# Repository Conventions

File and directory naming, file headers, and where a new package, module or
document goes. The module boundaries of the Go services and the webclient are
in AGENTS.md: Project Structure and the architecture-guidelines skill.

## Contents

- File Naming
- Source Files and Headers
- Adding a Package or a Command
- Where Other Files Go

---

## File Naming

Scope: Every file in the repository

Rule: File names follow the `validate-filenames` rules in `.manifest.yml`, which
the pre-commit hook enforces. The most specific directory rule wins:

| Place | Convention | Exempt from it |
| --- | --- | --- |
| `apps/webclient/` | kebab-case | `src/messages/` (locale files such as `pt-PT.json`), `src/routes/` (the router's names: `$locale`, `route.tsx`, `-home.module.css`), `src/routeTree.gen.ts`, `public/assets/`, `__snapshots__`, `Dockerfile`, `README.md` |
| `apps/services/` | snake_case | `AGENTS.md`, `Dockerfile`, `Makefile`, `README.md`, `etc/locales/`, `pkg/ajan/*/README.md` |
| everything else | kebab-case | `.agents/`, `.claude/`, `.github/`, `apps/apple/`, `docs/`, `scripts/`, and the root files `CHANGELOG.md`, `AGENTS.md`, `Dockerfile`, `LICENSE`, `Makefile`, `README.md`, `README.en.md`, `VERSION`, `MEMORY.md` |

Examples: `create-profile-form.tsx`, `locale-utils.ts`, `live-now.module.css`
and `use-mobile.ts` in the webclient; `bulletin.sql_gen.go` and
`event_audit.sql_gen.go` in the services; `push_current_release_tag.ts` in
`scripts/` (snake_case there by habit, since the directory is exempt). A
React component is `PascalCase` inside its kebab-case file
(`UserProfile` in `user-profile.tsx`); a Go test file ends in `_test.go`.

Because `.agents/` is exempt, a skill's files and folders keep the names the
skill format wants (`SKILL.md`, `references/`, `evals/`).

**Why:** a hook that rejects a wrong name costs seconds; a second casing in one
directory invites case-only renames that break on case-insensitive file
systems.

---

## Source Files and Headers

Scope: Source files

Rule: Every `.ts` and `.tsx` file starts with the one-line license header, then
a blank line or the imports:

```typescript
// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
```

Copy it into a new TypeScript file; the router plugin writes the same line
into `routeTree.gen.ts` (`routeTreeFileHeader` in `vite.config.ts`). Go files,
CSS modules, SQL and Markdown carry no header, and a Go file starts with its
package doc comment or the `package` clause. The `validate-licenses` step of
`.manifest.yml` runs on every commit.

---

## Adding a Package or a Command

Scope: New Go packages, executables, webclient modules and routes

Rule: Put new code where the matching code already is, and register it where
the neighbours are registered:

| What | Where |
| --- | --- |
| Go business module | `apps/services/pkg/api/business/<module>` (pure logic, no external dependencies) |
| Go adapter for an external system | `apps/services/pkg/api/adapters/<name>` (`github`, `s3client`, `storage`, ...) |
| Go framework utility | `apps/services/pkg/ajan/<name>` or `apps/services/pkg/lib/<name>` |
| Go executable | `apps/services/cmd/<name>/main.go` (`serve`, `migrate`, `manage`), only entry point and wiring |
| SQL query or migration | `apps/services/etc/data/default/queries/`, `.../migrations/` (then `make generate`) |
| Backend call from the webclient | `apps/webclient/src/modules/backend/<domain>/<verb-noun>.ts`, added to the `backend` object in `backend.ts`, with query options in `queries.ts` |
| Route | `apps/webclient/src/routes/` (file-based; `routeTree.gen.ts` regenerates) |
| Component | `apps/webclient/src/components/<feature>/`; generated shadcn components stay in `components/ui/` |
| Pure helper and its test | `apps/webclient/src/lib/foo-utils.ts` and `foo-utils.test.ts` (javascript-practices: Pure Function Extraction) |
| Hook, server-only code, translations | `src/hooks/`, `src/server/`, `src/messages/<locale>.json` for all 13 locales |

A new Go package or a new Go command builds with `make build` and passes
`make ok` in `apps/services`; a new webclient module passes `pnpm lint`,
`pnpm typecheck` and `pnpm test:ci`.

---

## Where Other Files Go

Scope: Documents, packaging, release and deployment files

Rule:

| What | Where |
| --- | --- |
| Architecture decisions | `docs/adr/` (`template.md` is the basis; `README.md` is the log) |
| Agent instructions and skills | `AGENTS.md` (nested ones scope a directory) and `.agents/skills/` (AGENTS.md: Agent tooling layout) |
| Workflows | `.github/workflows/` (`integrity.yml`, `deployment.yml`, `release_notes_sync.yml`, `apple.yml`) |
| Dockerfiles and compose files | `apps/webclient/Dockerfile`, `apps/services/Dockerfile`, `compose.yml`, `compose.production.yml` |
| Release version and notes | `VERSION` and `CHANGELOG.md` (`make tag`, `make release`) |
| Scripts | `scripts/` (javascript-practices: Node Scripts) |
| Dependency patches | `patches/` (tooling-standards: Overrides and Patches) |
| Local secrets and overrides | `.env.local`, `.env.*.local`, `config.local.json`: ignored, never committed |
