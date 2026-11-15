# Test Strategy

What gets tested, at which level and tier, and when a task counts as done.
Language mechanics are in go-tests.md and frontend-tests.md.

## Contents

- What to Test
- Test Tiers
- Table-Driven Tests
- Done and Reported

---

## What to Test

Scope: All code changes

Rule:

- Business logic (`apps/services/pkg/api/business/`) has unit tests that use
  fake implementations of its interfaces. Adapters
  (`apps/services/pkg/api/adapters/`) have integration tests against the real
  dependency (PostgreSQL, an S3 server). Dependency injection is what makes
  both possible: a service takes its interfaces in its constructor.
- Tests cover failure paths (invalid input, missing row, timeout, dependency
  error) as well as the happy path.
- Assert observable behavior: return values, emitted output, state changes. A
  test that only checks a function was called, or still passes when the
  behavior breaks, is not a test.
- New code ships with its tests. A change never lowers the coverage of the
  packages it touches; aim for 80% or more on critical paths
  (`make test-cov` in `apps/services` measures it).
- Test files sit next to the code: `foo_test.go` in Go, `foo.test.ts` in the
  webclient. A Go integration test is named `*_integration_test.go`.
- A disabled test (an unconditional `t.Skip`, `test.skip`) carries a comment
  with the reason. A gated skip (a missing database or server variable) names
  the need instead.

Correct:

```
processExport uses a fake privacy.Repository and a fake storage       -> unit test
storage.Repository against a real PostgreSQL, PRIVACY_TEST_DSN gated   -> integration test
isDeletionConfirmed in lib/privacy.ts                                  -> Vitest test
```

Incorrect:

```typescript
// No test for logic that decides whether an account may be deleted
export function isDeletionConfirmed(typed: string, username: string): boolean {
  return typed.trim().toLowerCase() === username.toLowerCase();
}
```

---

## Test Tiers

Scope: Choosing which tests to run, and writing a test that needs more than
memory and a few seconds

Rule: The suite runs in three tiers. Every test belongs to one:

- **Targeted:** the tests of the changed behavior by name, then the changed Go
  package or the changed Vitest file. Run it while writing code; it takes
  seconds.
- **Gate:** `make ok`, which is what the `default` workflow of `.manifest.yml`
  and the Integrity Pipeline run (`eser ok`): `go test -failfast -race -count 1
  ./...` over every Go package, `pnpm test:ci` for the webclient, and the lint,
  type, format and vulnerability checks around them. It is hermetic, needs no
  database and takes minutes. Run it once before handing a task over, not
  after every edit.
- **Opt-in:** tests gated by an environment variable that need a real
  dependency or change real data: `PRIVACY_TEST_DSN` (a PostgreSQL database
  with every migration applied) for
  `apps/services/pkg/api/adapters/storage`, and `S3_TEST_ENDPOINT`,
  `S3_TEST_BUCKET`, `S3_TEST_KEY`, `S3_TEST_SECRET` for
  `apps/services/pkg/api/adapters/s3client`. Neither the gate nor CI sets
  them. Run them when a change touches that surface (privacy SQL, redaction
  steps, export datasets, the S3 client), and say in the report when they did
  not run.

A change whose effect only the opt-in tier shows says so in the report. The
commands are in AGENTS.md: Development Commands.

Correct:

```sh
cd apps/services
go test -race -count 1 -run '^TestRequestExportQueuesWorkAndEnforcesCooldown$' ./pkg/api/business/privacy
go test -race -count 1 ./pkg/api/business/privacy
cd ../webclient && pnpm exec vitest run src/lib/privacy.test.ts
PRIVACY_TEST_DSN="postgres://...@127.0.0.1:55450/privacy_test_real?sslmode=disable" \
  go test ./pkg/api/adapters/storage/ -run TestPrivacy -count=1 -v
```

Incorrect:

```sh
# PRIVACY_TEST_DSN pointing at a shared or production database: the tests DESTROY data
make ok          # after every single edit
```

**Why:** the targeted tests catch almost every regression of a change in
seconds; the gate is the same run CI makes, so a failure it finds locally is
one the pipeline never sees; the opt-in tests delete accounts, so they must
only ever see a throwaway copy.

---

## Table-Driven Tests

Scope: Unit tests with several input cases

Rule: Put the cases in a table and loop over it. Each case has a name that
describes the behavior, not the implementation. The Go form is in go-tests.md:
Table-Driven and Parallel; in Vitest use `test.each` or a loop that registers
one `test` per case.

Correct:

```typescript
test.each([
  { name: "zero is its own class", input: 0, want: "zero" },
  { name: "negative numbers", input: -1, want: "negative" },
  { name: "positive numbers", input: 42, want: "positive" },
])("classifyNumber - $name", ({ input, want }) => {
  expect(classifyNumber(input)).toEqual(want);
});
```

Incorrect:

```typescript
test("test1", () => expect(classifyNumber(0)).toEqual("zero"));
test("test2", () => expect(classifyNumber(-1)).toEqual("negative"));
test("test3", () => expect(classifyNumber(42)).toEqual("positive"));
```

**Why:** a new case is one line, the inputs sit next to their expected
outputs, and a failure names the behavior that broke.

---

## Done and Reported

Scope: Ending a task

Rule: A task is done when `make ok` passed for what it changed, or when the
report says exactly which part could not run here and why (for example the
opt-in database tests without a database, or the Apple tests on a machine
without Xcode). The report says in one line:

- the commands that ran, with their results;
- what is left to CI, the opt-in tier or another machine;
- the tests that skipped, with their reasons;
- a known flaky test that failed and passed on its single rerun, named as
  flaky;
- files the Go lint autofix rewrote that the task did not change.

Never claim a tier that did not run, and never present a failing or skipped
check as passing.

**Why:** the person who takes the change over decides from the report what is
still unproven; a report that hides a gap moves the risk to them unannounced.
