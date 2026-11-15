# GitHub Actions

Conventions for the workflows in `.github/workflows/`, what each one does, and
how to read their runs with `gh`. The release flow is in release-process.md;
which tests each check runs is qa-practices: Test Hooks and Workflows.

## Contents

- Workflows in This Repository
- GitHub CLI Usage
- Inspect Before Acting
- Job Settings
- Permissions
- Secrets in CI
- Integrity Pipeline
- Toolchain Pins
- Dependabot

---

## Workflows in This Repository

Scope: Knowing what runs where

| Workflow (file)                          | Trigger                                                                                      | Purpose                                                                  |
| ---------------------------------------- | -------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------ |
| `Integrity Pipeline` (`integrity.yml`)   | Push to any branch and `v*.*.*` tags; pull requests (opened, edited, synchronize); manual    | `eser ok` plus a Go build on Ubuntu; the repository's quality gate       |
| `Deployment Pipeline` (`deployment.yml`) | Push of a `v[0-9]+.[0-9]+.[0-9]+*` tag                                                       | Derives the environment and image name from the tag (release-process.md: Tag Patterns and Environments) |
| `Release Notes Sync` (`release_notes_sync.yml`) | `v*` tags; `main` pushes that touch `CHANGELOG.md`; manual                            | Creates or updates the GitHub Release body from the changelog            |
| `Apple Apps CI` (`apple.yml`)            | `main` pushes and pull requests that touch `apps/apple/**` or the workflow file              | XcodeGen, iOS and macOS builds, unit and UI tests, an import lint        |

`apple.yml` lists its own file among its `paths`, so a change to the workflow
runs it. A new input directory that a path-filtered workflow depends on is added
to its `paths`. `integrity.yml` has no path filter and runs on every push and
pull request. Changing any workflow is an Ask First action (workflow-practices:
Ask First).

---

## GitHub CLI Usage

Scope: GitHub operations (pull requests, runs, issues)

Rule: Use the `gh` CLI for GitHub tasks; do not call the API with `curl` and a
token when `gh` can do it. Read-only commands (`gh pr view`, `gh pr checks`,
`gh run list`, `gh run view`, `gh issue view`) are the everyday set; commands
that write (`gh pr create`, `gh pr merge`, `gh issue create`, `gh release`) are
run only when the owner asks for that action.

```bash
gh pr view 123            # description and reviewers
gh pr diff 123
gh pr checks 123          # CI state; add --watch to wait for it
gh run list --workflow integrity.yml --limit 5
gh run view <run-id> --log-failed    # only the failing steps
gh issue list
gh issue view 123
```

Correct:

```bash
# Look at CI before anything else
gh pr checks 123 --watch
gh run view 12345 --log-failed
```

Incorrect:

```bash
# Raw API call with a token when gh is available
curl -H "Authorization: token $TOKEN" https://api.github.com/repos/...
```

---

## Inspect Before Acting

Scope: Debugging CI, merging, rerunning jobs, and any change to a deployment

Rule: Read the current state before changing anything. Check a pull request
before a merge, and a run before a rerun:

```bash
gh pr checks <number> && gh pr view <number>   # before any merge
gh run view <run-id> --log-failed              # before any rerun
```

Incorrect:

```bash
gh pr merge 123 --squash    # merging without checking CI status
```

A red check is read as qa-practices: Reading a Failing Check says. A red
`Integrity Pipeline` is a real signal, not inherited noise: reproduce the failing
hook with its CI command (Integrity Pipeline) before calling it flaky. For a deployment, look first with read-only
commands (deployments.md: Inspect Before Acting).

---

## Job Settings

Scope: Every job in every workflow

Rule:

- **Timeouts:** every job sets `timeout-minutes`, sized from a green run, so a
  hung job does not hold a runner for hours. `integrity.yml` allows 15 and
  `deployment.yml` 30; a new job starts at 20 and is tuned. `apple.yml` and
  `release_notes_sync.yml` set none yet; add one when touching them. Hook-level
  timeouts in `.manifest.yml` (for example 300 seconds for `go-tests`) are
  separate; none may exceed the job's.
- **Runners:** `ubuntu-latest` for Linux work, `macos-26` for the Apple jobs
  (Xcode 26.2, from the `XCODE_VERSION` env of `apple.yml`). A new runner type
  is a workflow change and needs approval.
- **Actions** are pinned to a version, never to a branch: the workflows use
  major tags (`actions/checkout@v7`, `actions/setup-go@v7`,
  `actions/setup-node@v7`, `pnpm/action-setup@v6`, `denoland/setup-deno@v2`)
  and Dependabot moves them (Dependabot). A full commit SHA with the version in
  a comment is the stronger pin, because a tag can move; use it where the owner
  asks for it.
- **Checkout** uses `fetch-depth: 2` in `integrity.yml` and `deployment.yml`,
  full history (`0`) only where the job reads tags, and `lfs: true` where the
  Apple tests need assets. A job that never pushes sets
  `persist-credentials: false` (the existing jobs do not yet).
- **Toolchains** come from setup actions with versions in the file: Go through
  `actions/setup-go` (`cache: true`, `cache-dependency-path:
  ./apps/services/go.sum`), Node through `actions/setup-node` (`cache: pnpm`),
  pnpm through `pnpm/action-setup` with no version, so it follows
  `packageManager` in `package.json`.
- **Install frozen:** `pnpm install --frozen-lockfile`; a job never rewrites the
  lockfile.
- **Concurrency:** every workflow sets a group. Checks cancel superseded runs
  (`integrity.yml`: `${{ github.ref }}-integrity`; `apple.yml`: workflow plus
  ref); `release_notes_sync.yml` sets `cancel-in-progress: false` and groups by
  the dispatched tag or the ref, so two syncs for one release queue and never cut
  each other off.
- **Manual entry:** a new workflow offers `workflow_dispatch` so it can be run
  while debugging. `integrity.yml` and `release_notes_sync.yml` have it;
  `integrity.yml` declares a `run_checks` input that no step reads.

Correct:

```yaml
name: Example
on:
  push:
    branches: [main]
  pull_request:
  workflow_dispatch:

concurrency:
  group: ${{ github.ref }}-example
  cancel-in-progress: true

jobs:
  check:
    runs-on: ubuntu-latest
    timeout-minutes: 20
    permissions:
      contents: read
    steps:
      - uses: actions/checkout@v7
      - run: pnpm install --frozen-lockfile
```

Incorrect:

```yaml
jobs:
  check:
    runs-on: ubuntu-latest      # no timeout, no permissions
    steps:
      - uses: actions/checkout@main    # a branch, not a version
      - run: pnpm install              # may rewrite the lockfile
```

**Why:** a job without a timeout holds a runner when it hangs, a branch-pinned
action changes under the same name, and an unfrozen install makes a run
irreproducible.

---

## Permissions

Scope: `permissions` in every workflow

Rule: Declare `permissions` explicitly, `contents: read` by default (at the top
of `apple.yml`, on the job in `integrity.yml` and `deployment.yml`). A job that
needs more declares it on that job alone: `Release Notes Sync` has
`contents: write` because `gh release create` and `gh release edit` write to the
repository's releases through `GITHUB_TOKEN`. No workflow commits, pushes or
creates tags; tags are the owner's (release-process.md: Release Commands). A new
write permission is a workflow change for the owner to approve.

---

## Secrets in CI

Scope: Any job that reads a secret, and every job that runs on `pull_request`

Rule: Pull request jobs read no secrets and pass without them: the integrity
checks need none. Secrets reach a job as `secrets.*` or `GITHUB_TOKEN` through
step `env` or an expression, never hard-coded. The ones in use are
`GITHUB_TOKEN` (release notes) and `GCLOUD_PROJECT` and `APP_NAME`
(`deployment.yml`, tag runs only). Never print, upload or cache a secret or a
credential file, and do not echo a variable built from one into a log that
others read. `validate-secrets` in `.manifest.yml` scans commits for
credentials; it is a second line, not the first. Security-sensitive handling is
in security-practices.

Correct:

```yaml
- name: Sync release body from changelog
  env:
    GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
  run: node scripts/sync_release_notes_from_changelog.ts --repo "$REPO"
```

Incorrect:

```yaml
- run: echo "token is ${{ secrets.GITHUB_TOKEN }}"
```

---

## Integrity Pipeline

Scope: `.github/workflows/integrity.yml`

Rule: One job, `Integration`, on `ubuntu-latest` (a matrix whose only row is
Ubuntu, Node.js 24, Go 1.27.4 and Deno 2.9.5), in this order: checkout, Deno
(only to host the `eser` CLI, published on JSR), Go, pnpm, Node.js,
`pnpm install --frozen-lockfile`, `deno install -gAf --name eser jsr:@eser/cli`,
`eser ok`, then `go build -v ./...` in `apps/services`.

`eser ok` is the `ok` script of `.manifest.yml` (`eser workflows run -e
precommit`), the same `default` workflow that the pre-commit and pre-push hooks
run: file validators, `oxfmt`, `oxlint`, the TypeScript 7 type check, Vitest
(`pnpm test:ci`), then the Go checks (`go vet`, `go build`, golangci-lint,
govulncheck, betteralign, modernize, `go test -race`). A hook that fails in CI
fails on the owner's machine the same way, so reproduce it with the command
and working directory named in `.manifest.yml`. Which hook checks what is in
qa-practices: Test Hooks and Workflows.

Reproduce without rewriting files: `.golangci.yaml` sets `fix: true` and the
`oxfmt` command writes by default, so run `go tool golangci-lint run --fix=false
./...` and `pnpm fmt:check`, not their fixing forms, in a tree that must stay
clean.

---

## Toolchain Pins

Scope: Version numbers in `integrity.yml` and the files it must agree with

Rule: These versions are written down on purpose; change one with a reason and
move its neighbours in the same change.

| Pin                             | Where                                             | Moves when                                       |
| ------------------------------- | ------------------------------------------------- | ------------------------------------------------ |
| Go patch `1.27.4`               | `golang-version` in `integrity.yml`               | govulncheck reports a new standard library finding |
| Go minor `1.27`                 | `go` line of `apps/services/go.mod`, `golang:1.27-bookworm` in `apps/services/Dockerfile` | A new Go minor; check all three |
| Node.js `24`                    | `node-version` in `integrity.yml`, `engines` in the root `package.json`, `node:24-bookworm-slim` in `apps/webclient/Dockerfile` | A new Node major; check all three |
| pnpm `12.18.0`                  | `packageManager` in the root `package.json`       | A deliberate pnpm bump                           |
| Deno `2.9.5`                    | `deno-version` in `integrity.yml`                 | Only the `eser` CLI needs it                     |
| `oxlint`, `oxfmt` exact versions | Root `package.json`                              | A deliberate bump (oxfmt is pre-1.0; AGENTS.md: Tooling override) |

A bump is run through the same checks as any change; a red hook after a bump is
read from its log, not retried.

---

## Dependabot

Scope: `.github/dependabot.yml`

Rule: Dependabot opens weekly update pull requests for `github-actions`, `npm`,
`docker` and `gomod` at the repository root, labelled `dependencies`, with the
commit prefix `ci` for actions and `chore` for the rest. Review them like any
change: the Integrity Pipeline runs on them, and a toolchain bump in Toolchain
Pins is also a change to the files listed there. Adding an ecosystem is a
configuration change for the owner to approve.
