# CI Tests

Which hook and workflow runs which tests, and how a test check is read and
changed. Workflow conventions (timeouts, permissions, runners, deployment) are
in release-and-ci.

---

## Test Hooks and Workflows

Scope: Knowing which check runs which tests

| Where                                  | Runs on                                                 | Tests                                                                                                          |
| -------------------------------------- | ------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------- |
| `.manifest.yml` workflow `default`     | `precommit` and `prepush` hooks; `make ok`; CI          | `go-tests`: `go test -failfast -race -count 1 ./...` (timeout 300 s); `webclient-tests`: `pnpm test:ci` (120 s) |
| `.github/workflows/integrity.yml`      | Every push, every pull request, `v*` tags; manual       | `eser ok`, so the whole `default` workflow, then `go build -v ./...`, on Go 1.27.4 and Node 24                  |
| `.github/workflows/apple.yml`          | Push and pull request to `main` that touch `apps/apple` | Unit tests, builds and UI tests of the Apple apps; the snapshot test classes are skipped in CI            |
| `.github/workflows/deployment.yml`     | `v*` tags                                               | No tests; it only classifies the tag and builds nothing (release-and-ci)                                                                          |

The hooks and the Integrity Pipeline run the same commands, so a task that
passed `make ok` has run what CI runs. The opt-in tests (strategy.md: Test
Tiers) run in neither. A new test must finish inside the timeout of its hook.
A new test check joins the `default` workflow of `.manifest.yml` as a `shell`
step with a name and a timeout, so it runs in the hooks and in CI alike, never
as a separate workflow that would run beside it. Changing a hook or a workflow
needs approval (workflow-practices: Ask First).

---

## Reading a Failing Check

Scope: A red test check on a pull request or `main`

Rule: Read the failing steps before rerunning anything:

```bash
gh pr checks <number>
gh run view <run-id> --log-failed
```

The log of the Integrity Pipeline names the failing hook (`go-tests`,
`webclient-tests`, `go-lint`, ...). Reproduce it locally with that hook's
command from `.manifest.yml`, in its `workingDirectory`, and for a single test
by name (strategy.md: Test Tiers). A Go lint failure that appears only in CI
is reproduced with `GOOS=linux GOARCH=amd64 $(go tool -n golangci-lint) run
--fix=false ./...` (apps/services/AGENTS.md: Linting and CI).

A test that fails once and passes on a single rerun is reported as flaky, with
its name; it is not reported as a pass, and a second failure is a failure. A
new test that fails only under load is a bug in the test (real time, shared
state, an unnamed wait, test order), not a reason to rerun: fix the cause. A
snapshot failure in `webclient-tests` is a regression or an update that was not
made, never a reason to run `pnpm test:update` blindly (frontend-tests.md:
Snapshot Testing).

**Why:** a rerun that turns a check green without a reason hides the failure
the next change runs into.
