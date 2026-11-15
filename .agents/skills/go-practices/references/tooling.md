# Go Tooling, Tests and Generated Code

The commands that check `apps/services`, what the linter settings mean for
the code you write, the shape of a Go test, and how SQL and generated files are
handled. Scoped facts of the directory are in `apps/services/AGENTS.md`; test
mechanics are in qa-practices.

## Contents

- Quality Tools
- Linter Rules That Shape the Code
- Go Tests
- SQL and Generated Code

---

## Quality Tools

Scope: Go changes in `apps/services`

Rule: Run the checks from `apps/services` through the Makefile, which uses the
tools pinned in the `tool` block of `go.mod` (`go tool <name>`):

| Command         | Runs                                                                      |
| --------------- | ------------------------------------------------------------------------- |
| `make lint`     | `go tool golangci-lint run ./...` (config in `.golangci.yaml`)            |
| `make check`    | `modernize -test`, `govulncheck`, `betteralign`, `go vet ./...`           |
| `make test`     | `go test -failfast -race -count 1 ./...`                                  |
| `make ok`       | `lint`, `check` and `test`                                                |
| `make fix`      | `modernize -fix`, `betteralign -apply`, `go fmt ./...`                    |
| `make generate` | `go generate ./...`: sqlc, stringer (SQL and Generated Code)              |

The same steps run as hooks from `.manifest.yml` (`go-vet`, `go-build`,
`go-lint`, `go-vulncheck`, `go-betteralign`, `go-modernize`, `go-tests`), so a
commit and CI see what `make ok` sees; the root `make ok` adds the frontend.
Conventional one-line commit messages are validated by the same manifest.

`.golangci.yaml` sets `fix: true`, so linting rewrites files in place. When the
tree must stay untouched (a review, a comparison), run
`go tool golangci-lint run --fix=false ./...`. To reproduce the Linux CI run
locally, see `apps/services/AGENTS.md`, Linting and CI.

Correct:

```bash
cd apps/services
make ok                                        # before handing a change over
go tool golangci-lint run --fix=false ./...    # lint without rewriting files
```

Incorrect:

```bash
golangci-lint run ./...   # a global install: not the pinned version
go test ./...             # skips -race, which is what make test runs
```

---

## Linter Rules That Shape the Code

Scope: All Go code

Rule: `.golangci.yaml` enables every linter (`default: all`) except `godox`,
`tagliatelle`, `ireturn`, `wsl` and `funcorder`, and formats with `gci`,
`gofmt`, `gofumpt`, `goimports` and `golines`. Write code that
passes them instead of silencing them:

- `err113`, `wrapcheck`, `nilnil`: sentinels, wrapped errors, no `return nil, nil`
  without a reason (errors-and-logging.md: Error Handling).
- `exhaustruct`: set every field of a struct literal, with `nil` or a zero value
  where that is intended (architecture.md: Constructors and Config).
- `mnd`: a number with a meaning is a named constant with a comment saying why
  (coding-practices: Magic Values).
- `gochecknoglobals`, `gochecknoinits`: no mutable package state and no `init()`
  (architecture.md: No `init()`).
- `varnamelen`, `funlen`, `cyclop`, `contextcheck`, `depguard`: short names stay
  in short scopes (`err`, `ctx`, `i`, `n`, `t`, `r`, `w` are allowed), long
  functions are split, a derived context is passed on, and an import outside the
  `depguard` allow list in `.golangci.yaml` fails the build. A new third-party
  module is added to that list in the same change as `go.mod` (tooling-standards: Adding a Dependency).
- `io/ioutil` is denied; use `io` and `os`.

A suppression names the linter and says why when the reason is not obvious from
the line: `//nolint:nilnil // not found is a normal answer`. Never a bare
`//nolint`, and never a `nolint` that hides a real defect to get `make ok`
green. A file-level `//nolint:...` above `package` is for test files with
in-memory fakes, with the reason after `//` (`service_test.go` in
`pkg/api/business/privacy`).

---

## Go Tests

Scope: `*_test.go` files in `apps/services`

Rule: Tests are plain `testing` tests (the module also allows
`github.com/stretchr/testify`), table-driven where one behavior has several
inputs, parallel with `t.Parallel()`, and run with the race detector
(`make test`). Business logic is tested through its ports with small in-memory
fakes (`fakeRepo` in `pkg/api/business/privacy/service_test.go`), not through a
database. Each iteration has its own variable since Go 1.22, so `tc := tc` is
not needed. Use `t.Context()` instead of `context.Background()` and
`for b.Loop()` in benchmarks (modernization.md). The test tiers and the table
shape are in qa-practices: Table-Driven and Parallel.

Correct:

```go
func TestCalculate(t *testing.T) {
    t.Parallel()

    cases := []struct {
        name     string
        input    int
        expected int
    }{
        {"zero", 0, 0},
        {"positive", 5, 25},
        {"negative", -3, 9},
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            t.Parallel()

            result := Calculate(tc.input)
            if result != tc.expected {
                t.Errorf("got %d, want %d", result, tc.expected)
            }
        })
    }
}
```

Incorrect:

```go
func TestCalculate(t *testing.T) {
    // no parallel execution, no table, no message with the values
    if Calculate(5) != 25 {
        t.Error("failed")
    }
}
```

---

## SQL and Generated Code

Scope: Queries, migrations and generated files

Rule:

- Queries live in `etc/data/default/queries/*.sql`, migrations in
  `etc/data/default/migrations/` (goose, numbered `0001_initial.sql`, ...).
  Every migration has both `-- +goose Up` and `-- +goose Down` sections; write
  the Down migration.
- `make generate` runs sqlc (`sqlc.yaml`: `database/sql`, `emit_interface`,
  snake_case JSON tags, output in `pkg/api/adapters/storage/`) and stringer.
  `make migrate-up` / `make migrate-down` apply migrations.
- Never hand-edit a generated file (`*.sql_gen.go`, `db_gen.go`, `querier_gen.go`,
  `types_gen.go`, `*_string.go`); change the SQL or the enum and regenerate.
  Repository methods in `pkg/api/adapters/storage/repository_*.go` wrap the
  generated `Queries` and map rows to business types.
- Every `_tx` table uses `CHAR(12)` for `locale_code`: apply
  `strings.TrimRight(value, " ")` when mapping to a business type
  (`strings.TrimRight(story.LocaleCode, " ")` in
  `pkg/api/business/stories/service.go`), and join translations with the 3-tier
  fallback (AGENTS.md: Internationalization & Locale Fallback).
- Business types do not import the generated `storage` package; the repository
  adapter converts.

Correct:

```sql
-- +goose Up
CREATE TABLE IF NOT EXISTS "profile_note" ("id" CHAR(26) NOT NULL PRIMARY KEY);

-- +goose Down
DROP TABLE IF EXISTS "profile_note";
```

Incorrect:

```go
// editing pkg/api/adapters/storage/profiles.sql_gen.go by hand: the next
// `make generate` overwrites it
```
