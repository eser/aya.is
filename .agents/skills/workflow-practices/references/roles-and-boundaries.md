# Roles and Boundaries

What an agent may do on its own in this repository, what needs the owner's
approval, and who owns which area.

## Contents

- Agent Roles
- File Ownership
- Never Do
- Ask First
- Language
- Changes Across Layers

---

## Agent Roles

Scope: Every agent task

Rule: An agent works in one role at a time.

- **Implementer:** changes code inside the area the task names, follows the
  conventions of the language skill and AGENTS.md, writes tests, and runs the
  quality gates (task-workflow.md: Quality Gates) before calling the work done.
- **Reviewer:** reads and reports; changes nothing.
  - Reports semantic issues. Formatting and lint-level issues are left to
    `gofmt`, `golangci-lint`, `oxlint` and `oxfmt`.
  - Marks each finding blocking or informational. Blocking: injection, missing
    authorization, unsynchronized shared state, non-exhaustive branching,
    secrets in output or logs, a migration without its `-- +goose Down`
    section, a new table or personal column missing from the privacy export or
    redaction lists (AGENTS.md: Privacy (export, hide, delete)), a
    hand-edited generated file. Everything else is informational.
  - Cites file:line for every finding. Before merge each finding is fixed,
    deferred with a named follow-up, or accepted with a reason; no blocking
    finding stays open.
- **Architect:** plans changes across layers with the protocol below, reads
  `apps/services/AGENTS.md` and architecture-guidelines, and records
  significant decisions as the next ADR in `docs/adr/` from `template.md`.

---

## File Ownership

Scope: Deciding whether a change stays inside its area

Rule: A task stays inside the area it was given. The areas marked high impact
need the Ask First step before any change.

| Area                                                  | Concern                                                                                                              |
| ----------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| `apps/services/pkg/api/business/`                     | Pure business logic with no external dependency; it reaches the outside only through its interfaces                  |
| `apps/services/pkg/api/adapters/`                     | Implementations of those interfaces (storage, http, workers, providers); `*_gen.go` files are generated              |
| `apps/services/pkg/ajan/`                             | Framework utilities shared by the services                                                                           |
| `apps/services/cmd/`                                  | Entry points (`serve`, `migrate`, `manage`) that wire everything together                                            |
| `apps/services/etc/data/default/migrations/`          | **High impact**: goose migrations, every file with an Up and a Down section; queries sit in `queries/`, seeds in `seed/` |
| `apps/webclient/src/`                                 | The TanStack Start frontend: `routes/`, `components/`, `modules/`, `lib/`, `messages/` (13 locale files)             |
| `apps/webclient/src/components/ui/`                   | Generated shadcn code; do not modify its inline Tailwind (AGENTS.md: Shadcn UI Components)                           |
| `apps/apple/`                                         | The Apple apps; their CI is `.github/workflows/apple.yml`                                                            |
| `patches/`                                            | pnpm patches of pinned TanStack packages; they change with the pinned versions only                                  |
| `.github/workflows/`, `.manifest.yml`                 | **High impact**: CI, deployment and the hooks that gate every commit                                                 |
| `compose.yml`, `compose.production.yml`, `Dockerfile` | **High impact** when production or deployment settings change                                                        |
| `VERSION`, `CHANGELOG.md`, `scripts/`                 | The release flow; owner only (Never Do: Release commands)                                                            |
| `docs/adr/`                                           | Architecture decision records, from `template.md`                                                                    |
| `.agents/skills/`                                     | Agent rules; update them when a convention changes (eser-rules-manager)                                              |

---

## Never Do

Scope: All agent actions

Rule: These are forbidden unless the owner asks for that exact action in the
current task. A general "go ahead" does not cover them.

- **Git history and working tree:** no `commit`, `push`, `tag`, `reset`,
  `rebase`, `checkout`, `restore`, `stash` or force push. The owner commits and
  pushes; the agent proposes a one-line commit message with no AI attribution
  (task-workflow.md: Git Commit Policy).
- **Release commands:** `make tag` creates and pushes the release tag from
  `VERSION`, `make release` syncs `CHANGELOG.md` to GitHub Releases, and a
  `v*` tag starts the deployment pipeline (`.github/workflows/deployment.yml`).
  Run them only when the owner asks for a release.
- **Versions:** never bump `VERSION` or write a `CHANGELOG.md` release entry
  unless the owner asks for the release.
- **Generated files and lock files:** `*_gen.go` files change only through
  `make generate`, `src/routeTree.gen.ts` only through the route generator,
  `go.sum` only through `go mod tidy` and `pnpm-lock.yaml` only through
  `pnpm install`.
- **Checks:** never skip hooks (`--no-verify`), and never disable an
  `oxlint` or `golangci-lint` rule, skip a failing test or loosen a check to
  make it pass.
- **Secrets:** never commit secrets; `.env.local`, `.env.*.local` and
  `config.local.json` are git-ignored development overrides.
- **Scope:** do not modify files outside the assigned area.

Stop and reconsider when the plan includes "just this once I'll skip the
checks", "this change is too small for a test", "the Down migration can come
later" or "it is only one line of the lock file".

**Why:** these actions are either irreversible (a pushed `v*` tag starts a
deployment) or hide a problem the checks exist to catch.

---

## Ask First

Scope: Actions that need the owner's confirmation

Rule: Describe the change and wait for approval before:

- adding a dependency to `apps/services/go.mod` or a `package.json`, with its
  purpose, license and the alternative (tooling-standards: Adding a
  Dependency);
- replacing or adding native tooling (Node.js, pnpm, Go, Make, Vitest, oxlint,
  oxfmt): each part of the repository has its own tooling and keeps it;
- changing `.github/workflows/`, `.manifest.yml` or the deployment settings
  (release-and-ci);
- adding a dependency between business modules that did not depend on each
  other;
- removing or changing the meaning of an API route, field or flag that the
  webclient, the Apple apps or another person calls;
- anything a later release cannot undo: a migration that drops or rewrites
  data, and persisted or public formats.

---

## Language

Scope: Everything written into the repository

Rule: Source code, comments, commit messages, ADRs, skills and task plans are
written in English, whatever language the conversation uses. A Turkish request
produces English code and commit messages; never leave a non-English fragment
in code or a comment because the owner wrote the request that way. User-facing
product text is the exception: it goes through i18n with a proper translation
in all 13 locales, and never English text in a non-English locale file
(AGENTS.md: Internationalization & Locale Fallback).

**Why:** the code base has one working language, and the translation keys are
English text.

---

## Changes Across Layers

Scope: A change that spans more than one layer, such as SQL, the Go business
and adapter packages, and the webclient

Rule:

1. **Justify:** say why the parts must change together. One area per change is
   the default.
2. **Plan:** list each affected package or directory and its part of the
   change.
3. **Order:** change the SQL (a migration with its Down section, then the
   query) and run `make generate` first; then the business interfaces before
   the adapters that implement them; then the HTTP route; then the webclient
   (`src/modules/backend`, query options in `queries.ts`, the route, the
   component); wire the entry point in `cmd/` last.
4. **Locales:** every new user-facing string comes with all 13 locale entries
   in the same change.
5. **Test:** run the tests of every changed package and file once all are
   changed, then `make ok` (task-workflow.md: Quality Gates).
