# READMEs and Project Docs

Where the repository's documents live, how they are shaped, and when they
change.

## Contents

- Where Documentation Lives
- Technical Documents
- Package README Structure
- Descriptive Text Preservation
- Docs Change With the Code
- Documentation Boundary
- Product Names in Documents

---

## Where Documentation Lives

Scope: Adding or moving a document

Rule: Put a document where its kind already lives, and extend an existing
document before adding a new one.

| Location                                      | Holds                                                                       |
| --------------------------------------------- | --------------------------------------------------------------------------- |
| `README.md`, `README.en.md`                   | Project overview in Turkish and its English counterpart, kept in step       |
| `CHANGELOG.md`                                | Release notes, one `## <version> - <YYYY-MM-DD>` entry per release         |
| `docs/adr/`                                   | Architectural decisions (architecture-guidelines: Architectural Decision Records) |
| `docs/*.md`                                   | Technical documents: architecture, contracts, behavior, deployment          |
| `.github/`                                    | Contributor and community files: `CONTRIBUTING.md`, `CONTRIBUTING.en.md`, `SECURITY.md`, `CODE_OF_CONDUCT.md`, `PULL_REQUEST_TEMPLATE.md`, `ARCHITECTURE.md` |
| `MEMORY.md`                                   | Notes on the story content model and the decisions behind it                |
| A README beside code                          | Instructions for that directory (`apps/webclient/README.md`, `apps/services/pkg/ajan/logfx/README.md`) |
| `AGENTS.md`, `apps/services/AGENTS.md`, `.agents/skills/` | Agent instructions (AGENTS.md: Agent tooling layout)            |

`CHANGELOG.md` is also the source of the GitHub release notes: `make release`
runs `scripts/sync_release_notes_from_changelog.ts` over its top entry, so the
heading format is part of the interface (release-and-ci covers the release
steps).

File names under `docs/` are kebab-case and name the subject
(`server-deployment.md`); ADR files carry their number
(`docs/adr/template.md`: `0005-example.md`).

---

## Technical Documents

Scope: Every document under `docs/`

Rule: A technical document has a sentence-case `#` title naming its subject and
a first paragraph that says what it covers and where neighboring material is.
It describes architecture, contracts, behavior, configuration, procedures and
reproducible tests with their results. It names real limits, keeps proposed
and implemented apart, and links the ADRs and documents it relies on with
relative links. Procedures give commands that run as pasted from the directory
the document names (the repository root unless it says otherwise).

Correct:

```markdown
# Coolify deployment

This document lists what a deployment consists of, what the repository
provides for it (`compose.production.yml` and the deployment workflow) and what
the operator still prepares. For a local stack on one machine, use `make up`.
```

Incorrect:

```markdown
# Deployment (in progress)

Status: Phase 2 done, Coolify support lands next sprint.
```

---

## Package README Structure

Scope: A README.md beside code

Rule: Add a README to a directory only when someone working there needs
instructions the code cannot give: what the package is for, how to build, run or
test it, where its output goes, or where pinned upstream data comes from. The
framework packages under `apps/services/pkg/ajan/` each describe themselves this
way (`pkg/ajan/logfx/README.md`); packages under `pkg/api/` are described by
`apps/services/AGENTS.md: Layout` and need none. A README contains:

1. A `#` title naming the directory's subject (`# ajan/logfx`).
2. What the directory is and how it fits the project, in a few sentences.
3. The commands or the import and usage example, each runnable as pasted.
4. The rules a change in the directory must keep, with links to the documents
   that own them.

---

## Descriptive Text Preservation

Scope: Rewriting any document or comment

Rule: Keep descriptive details that are still true when rewriting: a limit,
a reason, a case the original names ("without durable storage", "on Unix
sockets in temporary directories") exist for a reader. Shortening is fine;
losing a fact is not. Text that became false is deleted or rewritten, not
preserved.

Correct:

```markdown
The PostgreSQL tests start their own clusters from these binaries on Unix
sockets in temporary directories.
```

Incorrect:

```markdown
The PostgreSQL tests start their own clusters. <!-- lost where they run -->
```

---

## Docs Change With the Code

Scope: Any change to a setting, environment variable, Make target, command,
route, persisted format, user-visible behavior, user-facing string, or install
and first-run step

Rule: The documentation diff ships in the same change as the code diff. Before
handing the task over, find each document that describes the changed surface
and make it describe the new state (the pull request template asks the same:
"My change requires changes to the documentation"):

- Setting or environment variable: the file that lists it
  (`apps/webclient/.env.example` for webclient variables,
  `apps/services/config.json`, the `environment` block of `compose.yml` and
  `compose.production.yml`).
- Make target: the `##` help comment on its recipe, which is the line
  `make help` prints, and Development Commands in AGENTS.md when it lists it.
- A new framework package under `pkg/ajan/`, or one with a new responsibility:
  its README.md; a change of layers or directories: AGENTS.md: Project
  Structure and `apps/services/AGENTS.md: Layout`.
- A design an ADR describes: the ADR's status line.
- A change a user or operator notices: the top entry of `CHANGELOG.md`.
- A string a person reads in the webclient: the entry in all 13 locale files
  (AGENTS.md: Internationalization & Locale Fallback).
- A tool, command or convention that AGENTS.md or a skill states: that file.
- Removed or changed behavior: delete or rewrite the sentences that describe
  the old behavior.

**Why:** a feature whose docs lag the code does not exist for the reader who
finds it later, and stale instructions cost more than missing ones.

---

## Documentation Boundary

Scope: Every Markdown file, code comment and string a person reads

Rule: A document describes what the code does and why. Task status, progress,
schedules and delivery plans ("Phase 2 done", "lands next sprint") do not
appear in docs, comments or copy; they live in GitHub issues and pull requests.
Linking an issue is fine where a format asks for it (the Technical Story line of
an ADR, `Fixes #123` in a pull request), and `CHANGELOG.md` records what a
release contains, not what is planned.

**Why:** status text is true for a week and wrong afterwards, and nobody
deletes it.

---

## Product Names in Documents

Scope: Every document and every string a person reads

Rule: Use the terms AGENTS.md and the code set, and one term per concept: AYA
for the platform, profile with its types individual, organization and product
(AGENTS.md: Profile System), story with its kinds (`MEMORY.md`), locale for a
language variant. Do not introduce a synonym for a domain term in a document or
a string. Identifiers keep their stable spelling (`locale_code`, `story_tx`)
and are written as code.
