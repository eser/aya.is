---
name: qa-practices
description: "Testing in this repository: test levels, the targeted, gate and opt-in tiers, Go and Vitest tests, snapshot tests, real-database tests, CI test checks, flaky tests and what a task reports as verified. Use when writing a test, choosing which tests to run, deciding whether a task is done, updating a snapshot, reading a failing check or changing a test hook. Not for releases or deployment (use release-and-ci)."
---

# QA Practices

How the code is tested and how a change proves it works: test levels, test
tiers, the mechanics of Go and webclient tests, the test hooks and workflows,
and the report a task ends with. Repository commands live in AGENTS.md:
Development Commands.

## Always

- Every new test belongs to a tier when it is written; a test that needs an
  external service or changes real data skips when its variable is unset, with
  a reason that names it, so `go test ./...` stays hermetic
- While working run the tests of what changed by name; before handing a task
  over run `make ok`, the same hooks CI runs, then read what its Go lint step
  rewrote (workflow-practices: Quality Gates)
- Logic gets unit tests with hand-written fakes at its interfaces, adapters get
  integration tests against the real dependency, failure paths are tested,
  coverage of a touched package never drops
- Assert observable behavior; a test that still passes when the behavior
  breaks is not a test
- Go tests are table-driven, call `t.Parallel()` first (or carry a
  `//nolint:paralleltest` with its reason), live in `package x_test` and run
  with `-race`; time is injected, never waited for
- Webclient tests are Vitest tests of pure functions next to their source;
  logic that a component needs is extracted first (javascript-practices: Pure
  Function Extraction)
- Snapshots guard against regression: `pnpm test:update` only after an
  intentional change, with the `.snap` diff read; `pnpm test:ci` never writes
- A bug fix starts with a failing test that reproduces it
- The report names in one line what ran, what could not run here and the
  skipped tests with their reasons; a known flaky test that fails is rerun once
  and reported as flaky, never as a pass

## References

| File                                               | Read when                                                                    |
| -------------------------------------------------- | ---------------------------------------------------------------------------- |
| [strategy.md](references/strategy.md)              | Deciding what to test, at which level and tier, or whether a task is done    |
| [go-tests.md](references/go-tests.md)              | Writing Go tests, fakes for interfaces, gated real-database tests, time      |
| [frontend-tests.md](references/frontend-tests.md)  | Writing or running Vitest tests, snapshots or the translations test          |
| [ci-tests.md](references/ci-tests.md)              | Reading a failing test check, or adding or changing a test hook or workflow  |
