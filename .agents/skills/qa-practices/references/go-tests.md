# Go Tests

How Go tests are written and run here (`apps/services`). What to test and which
tier runs it is in strategy.md.

## Contents

- Table-Driven and Parallel
- Files, Packages and Lint
- Gated Tests Skip, Never Fail
- Fakes at the Port
- Real Database Tests
- Time and Concurrency
- Benchmarks

---

## Table-Driven and Parallel

Scope: All Go tests

Rule: Write table-driven tests with `t.Parallel()` as the first statement of
the test and of every subtest. Since Go 1.22 each loop iteration has its own
variable, so the old `tc := tc` copy is not needed. Use `t.Context()` instead
of `context.Background()`. The `paralleltest` linter is on (`.golangci.yaml`
has `default: all`), so a test that must stay serial (shared process state such
as `t.Setenv`, or a shared database) says so in a `//nolint:paralleltest` with
the reason, as the privacy integration test does. A parent whose subtests are
parallel releases what they use with `t.Cleanup`, never `defer`.

Correct:

```go
func TestIsSessionStale(t *testing.T) {
    t.Parallel()

    now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
    cases := []struct {
        name       string
        lastActive time.Time
        want       bool
    }{
        {"active a minute ago", now.Add(-time.Minute), false},
        {"idle for two hours", now.Add(-2 * time.Hour), true},
        {"never active", time.Time{}, true},
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            t.Parallel()
            if got := IsSessionStale(tc.lastActive, now); got != tc.want {
                t.Errorf("IsSessionStale(%v) = %v, want %v", tc.lastActive, got, tc.want)
            }
        })
    }
}
```

Incorrect:

```go
func TestIsSessionStale(t *testing.T) {
    if !IsSessionStale(time.Time{}, time.Now()) { // one case, wall clock, serial
        t.Error("failed")
    }
}
```

---

## Files, Packages and Lint

Scope: Where a Go test lives and how it is named

Rule:

- A test sits next to the code in `foo_test.go`; an integration test in
  `foo_integration_test.go`.
- A test lives in the external package (`package privacy_test`) and uses only
  the exported API, because the `testpackage` linter is on. A test that needs
  package internals (an unexported repository field) uses the internal package
  and says why in the `//nolint:testpackage` line at the top of the file.
- Assertions use the standard library (`t.Errorf`, `t.Fatalf`) or `testify`
  (`assert`, `require`), which `.golangci.yaml` allows and `pkg/ajan` tests use.
  Never `testify/mock` or generated mocks: fakes are hand-written
  (Fakes at the Port).
- A `//nolint` for a test file names every linter and carries a reason after
  `//`, like the existing ones.

Correct:

```go
//nolint:exhaustruct,varnamelen,cyclop,funlen // table-free scenario tests with small in-memory fakes
package privacy_test
```

Incorrect:

```go
//nolint
package privacy
```

---

## Gated Tests Skip, Never Fail

Scope: Go tests that need an external service, a real database or a long run

Rule: A test that needs more than memory skips when its dependency is not
configured, before it starts anything, with a reason that names the variable:
`t.Skip("PRIVACY_TEST_DSN is not set")`. It never fails because the dependency is
missing, so `go test ./...` and the gate stay hermetic. A shared helper does
the check once (`privacyTestRepository` in
`pkg/api/adapters/storage/repository_privacy_integration_test.go`,
`integrationClient` in `pkg/api/adapters/s3client/client_integration_test.go`),
so every caller inherits it. Describe the variables and a sample command in the
comment at the top of the file.

Correct:

```go
func integrationClient(t *testing.T) *Client {
    t.Helper()

    endpoint := os.Getenv("S3_TEST_ENDPOINT")
    if endpoint == "" {
        t.Skip("S3_TEST_ENDPOINT is not set")
    }
    // ...
}
```

Incorrect:

```go
func TestPutObject(t *testing.T) {
    client, err := New(t.Context(), Config{Endpoint: os.Getenv("S3_TEST_ENDPOINT")})
    if err != nil {
        t.Fatal(err) // fails for everyone who has no S3 server
    }
}
```

**Why:** a test that fails without its dependency turns every gate run red for
a reason that is not a regression, and people learn to ignore the gate.

---

## Fakes at the Port

Scope: Tests of business services

Rule: Test a service by passing hand-written fakes of its interfaces to the real
constructor. A fake is a small struct in the `_test.go` file that implements
the interface with an in-memory map or a recorded call list (`fakeRepo` in
`pkg/api/business/privacy/service_test.go`); no mock framework. Pass a real or
no-op value for every parameter, never `nil` for a dependency the code might
call: a nil logger or clock turns a test failure into a panic far from the
cause. Keep the real adapters for integration tests, and fake only the boundary
under test. A failure of a dependency is a field of the fake (`failCreate`), so
the failure path is a test too.

```go
type fakeStore struct {
    sessions map[string]*Session
}

func (f *fakeStore) Get(_ context.Context, id string) (*Session, error) {
    s, ok := f.sessions[id]
    if !ok {
        return nil, ErrSessionNotFound
    }
    return s, nil
}

func TestServiceLoad(t *testing.T) {
    t.Parallel()

    store := &fakeStore{sessions: map[string]*Session{"s1": {ID: "s1"}}}
    now := func() time.Time { return time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC) }
    svc := NewService(store, now, slog.New(slog.DiscardHandler))

    got, err := svc.Load(t.Context(), "s1")
    if err != nil {
        t.Fatalf("Load(%q): %v", "s1", err)
    }
    if got.ID != "s1" {
        t.Errorf("Load(%q).ID = %q, want %q", "s1", got.ID, "s1")
    }
}
```

**Why:** a fake that implements the interface breaks at compile time when the
interface changes, and the test exercises the same constructor the entry point
uses.

---

## Real Database Tests

Scope: Tests of SQL, redaction steps, export datasets and other storage code

Rule: A storage test that needs PostgreSQL reads its connection string from an
environment variable and skips without it (Gated Tests Skip, Never Fail). A test
that deletes or rewrites data also refuses a database it should not touch:
`PRIVACY_TEST_DSN` must name a database that starts with `privacy_test`, any
other name aborts the run. Run such tests against a throwaway database that has
every migration applied, ideally filled with a copy of realistic data, and
restore it between runs. A new table or personal column is covered by both
the export datasets and the redaction steps (AGENTS.md: Privacy (export, hide,
delete)); `TestPrivacyDatasetsCoverEveryTable` and
`TestPrivacySecretColumnsAreWithheld` fail when one is missing, but only in the
opt-in tier, so a migration that adds a table needs that run or a report that it
did not happen. Test an S3 client
against a server that verifies signatures (RustFS, MinIO, R2), never a
permissive fake such as moto: it accepts any signature and proves nothing.

Correct:

```sh
PRIVACY_TEST_DSN="postgres://postgres:...@127.0.0.1:55450/privacy_test_real?sslmode=disable" \
  go test ./pkg/api/adapters/storage/ -run TestPrivacy -count=1 -v
```

Incorrect:

```go
db, _ := sql.Open("pgx", os.Getenv("DATABASE_URL")) // the development database, no name guard
```

**Why:** these tests exercise deleting accounts; the name guard is the last
line between a typo in a DSN and lost data.

---

## Time and Concurrency

Scope: Tests of code with clocks, goroutines, timers or tickers

Rule: Prove behavior over time with an injected clock or an event, never by
waiting for real time (`time.Sleep`). A production interval the test shortens is
a field whose zero value is the named production default. Every duration in a
test is a named constant whose comment says why that value; a wait for an event
has a named ceiling, and the ceiling is never what the test proves. For
goroutines and timers use `testing/synctest.Test` with `synctest.Wait`; time
inside the bubble is virtual. `-race` is part of every run, and a pass does not
prove the absence of races on paths the tests do not run.

**Why:** a test bound by real time fails under the load of a parallel run and
makes the whole gate slow; the hook has a 300-second timeout for all Go tests.

---

## Benchmarks

Scope: `Benchmark*` functions

Rule: Use `for b.Loop()` and `b.ReportAllocs()`, and compare runs rather than
reading one number.
