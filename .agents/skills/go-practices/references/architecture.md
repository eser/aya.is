# Go Architecture and Package Structure

Where Go code lives in `apps/services`, how business logic meets its adapters
and how packages are constructed. The layer rules themselves are in
architecture-guidelines: Hexagonal Architecture and architecture-guidelines: Composition Roots; the scoped facts of this directory are in
`apps/services/AGENTS.md`.

## Contents

- Repository Layout
- Modules, Ports and Adapters
- Dependency Injection
- Adding a Port or Module Package
- File Naming
- Constructors and Config
- Configuration Loading
- No `init()`
- Compile-Time Interface Checks

---

## Repository Layout

Scope: All Go code

Rule: The backend is one Go module, `github.com/eser/aya.is/services`
(`apps/services/go.mod`, no `go.work`), with the `Makefile`, `sqlc.yaml` and
`.golangci.yaml` next to it.

| Path                           | Holds                                                                        |
| ------------------------------ | ---------------------------------------------------------------------------- |
| `cmd/<executable>/`            | One program: `serve` (HTTP server and workers), `migrate`, `manage` (CLI)    |
| `pkg/api/business/<module>/`   | Pure business logic and its ports (`privacy`, `stories`, `profiles`, ...)    |
| `pkg/api/adapters/<adapter>/`  | Implementations of the ports and every external integration                  |
| `pkg/api/adapters/appcontext/` | The composition root: builds config, logger, connections, repository, services |
| `pkg/ajan/<name>fx/`           | Framework utilities (`logfx`, `httpfx`, `configfx`, `connfx`, `processfx`)   |
| `pkg/lib/<name>/`              | Small shared helpers (`caching`, `cursors`, `vars`)                          |
| `etc/data/default/`            | `migrations/`, `queries/` (sqlc input) and `seed/`                           |

Do not create `internal/`, `api/`, `configs/` or `test/` trees. A new business
module is a directory under `pkg/api/business/`, a new integration one under
`pkg/api/adapters/`.

---

## Modules, Ports and Adapters

Scope: Packages under `pkg/api/`

Rule:

- A business package holds its logic and declares its ports as interfaces next
  to the code that uses them (`privacy.Repository`, `privacy.ObjectStorage`,
  `privacy.Enqueuer` in `pkg/api/business/privacy/types.go`).
- Business packages import other business packages, `pkg/ajan` and `pkg/lib`,
  never an adapter, a database driver or a vendor SDK. A business package that
  imports `github.com/redis/go-redis/v9` is a violation.
- Adapters do the I/O (PostgreSQL through sqlc, S3, HTTP providers) and
  implement the ports. The adapter packages import the business package for
  its types, never the reverse.
- A port is small. A large repository interface needs a
  `//nolint:interfacebloat // reason` that says why one port is right
  (`privacy.Repository`).

Correct:

```go
// pkg/api/business/user/service.go: the port, owned by the business layer
package user

type Repository interface {
    GetUserByID(ctx context.Context, id string) (*User, error)
}

type Service struct {
    repo Repository // an interface, not a concrete type
}

func NewService(repo Repository) *Service {
    return &Service{repo: repo}
}
```

Incorrect:

```go
package user

import "github.com/redis/go-redis/v9" // an external dependency in business logic

type Service struct {
    redis *redis.Client // a concrete type, not an interface
}
```

**Why:** logic that only knows its own ports runs in a test with an in-memory
fake (`fakeRepo` in `pkg/api/business/privacy/service_test.go`) and survives a
change of vendor.

---

## Dependency Injection

Scope: Services, adapters and the composition root

Rule: A type receives everything it depends on through its constructor. The
constructor returns a concrete type (a `*Service`, not an interface); its
parameters are the consumer's own port interfaces, the shared `*logfx.Logger`
and a `*Config`.

- No hidden dependencies: a constructor does not open connections, read the
  environment or reach for a package-level variable.
- No service locator and no registry. `appcontext.AppContext.Init`
  (`pkg/api/adapters/appcontext/appcontext.go`) is the only code that builds
  the config, the logger, the connections, the repository and the token
  services in order and hands them on; `cmd/serve` then constructs the
  services from it.
- Stateful resources (connection pools, the HTTP client, the logger) are created
  once there and shared.
- A bad setting or an unreachable dependency stops start-up (`Init` returns
  the error, `main` stops), so no request meets a half-built process.
- Wiring stays manual. A DI container or code generator is a new dependency
  and an architecture change: it needs tooling-standards: Adding a Dependency
  first.
- Time and ids that tests must control are fields with a production default,
  not calls inside the logic: `newID: DefaultIDGenerator`, `now: time.Now`,
  replaced through a setter such as `SetClock` (`pkg/api/business/privacy/service.go`).

Correct:

```go
func NewService(logger *logfx.Logger, config *Config, repo Repository, storage ObjectStorage /* nil allowed */) *Service {
    return &Service{logger: logger, config: config, repo: repo, storage: storage, newID: DefaultIDGenerator, now: time.Now}
}
```

Incorrect:

```go
func NewService() *Service {
    db, _ := sql.Open("pgx", os.Getenv("DATABASE_URL")) // hidden, unchecked
    return &Service{db: db}
}
```

**Why:** when every dependency is a constructor input, a test swaps any of
them for a fake, and a reader sees what a service touches without reading its
body.

---

## Adding a Port or Module Package

Scope: New logic in `apps/services`

Rule:

1. Create `pkg/api/business/<module>/` with `types.go` (types, ports,
   constants), `errors.go` (sentinels) and `service.go`.
2. Declare a port for everything the logic needs from outside, in that package.
3. Implement the ports in `pkg/api/adapters/<adapter>/`; database access goes
   through a repository method over sqlc queries (SQL and Generated Code in
   [tooling.md](tooling.md)).
4. Construct the service in `cmd/serve` (or the command that needs it) from
   what `appcontext` provides, and register its HTTP routes in
   `pkg/api/adapters/http/`.
5. Run `make ok` in `apps/services` (go-practices: Quality Tools).

---

## File Naming

Scope: Go files

Rule: Go files are snake_case (enforced by the `validate-filenames` hook for
`apps/services/`, see `.manifest.yml`), named for their content, with the
`_test` suffix last. Generated files keep the generator's suffix: `*.sql_gen.go`,
`db_gen.go`, `querier_gen.go` (sqlc), `*_string.go` (stringer).

Correct:

```
user_service.go
payment_handler.go
routes_privacy.go
```

Incorrect:

```
UserService.go
paymentHandler.go
auth-middleware.go
```

---

## Constructors and Config

Scope: Constructors of long-lived types (services, stores, workers, clients)

Rule: A constructor takes its dependencies as parameters and returns the
concrete type: `NewService(logger, config, repo, ...) *Service`. Settings that
come from the environment live in `Config` structs (`AppConfig`,
`httpfx.Config`, `httpclient.Config`), filled by `configfx` from `conf:"..."`
and `default:"..."` tags, so a default is written once, next to the field:

```go
type Config struct {
    ReadHeaderTimeout time.Duration `conf:"read_header_timeout" default:"5s"`
    ProfilingEnabled  bool          `conf:"profiling"           default:"false"`
}
```

- A dependency that is optional says so in the parameter comment and the
  struct field (`storage ObjectStorage // nil when no private bucket is configured`);
  the code checks it for `nil` before use.
- A constructor that can fail returns `(*T, error)` and a sentinel for an
  invalid value, so a bad setting fails at construction instead of at first use.
- Struct literals set every field, with `nil` or a zero value where it is
  intended: `exhaustruct` runs in `.golangci.yaml`. A deliberately partial
  literal of a third-party type carries `//nolint:exhaustruct` (and a reason when
  it is not obvious).
- Do not grow a signature by one `bool` per feature: a group of related
  settings becomes a `Config` struct parameter.

Incorrect:

```go
func New(repo Repository, timeout time.Duration, verbose bool, cache bool) *Service
// every new setting changes the signature, and a call site reads "true, false"
```

---

## Configuration Loading

Scope: Loading settings in `appcontext`

Rule: Settings load in priority order, lowest first, so the later source wins:
`config.json`, then `.env`, then the system environment
(`ConfigManager.LoadDefaults` in `pkg/ajan/configfx/manager.go`). A tag
default applies when no source sets the key. Local overrides live in
git-ignored `.env.local`, `.env.*.local` and `config.local.json`
(`apps/services/AGENTS.md`); never commit secrets.

Correct:

```go
// pkg/api/adapters/appcontext/appcontext.go
cl := configfx.NewConfigManager()
a.Config = &AppConfig{} //nolint:exhaustruct

err := cl.LoadDefaults(a.Config)
if err != nil {
    return fmt.Errorf("%w: %w", ErrInitFailed, err)
}
```

Incorrect:

```go
port := os.Getenv("API_PORT") // read in the middle of the logic: no default, no validation
```

---

## No `init()`

Scope: All Go code

Rule: Avoid `init()`. It runs implicitly at import time, cannot return an
error and runs in every test that imports the package; the `gochecknoinits` and
`gochecknoglobals` linters are on. Do the work in a constructor that
`appcontext` or `main` calls and whose error it handles. A package-level `var`
with a constant value (a compiled regexp, a lookup table) is fine; see
design-principles: Module-Level Side Effects. Where a table really has to be
built at import time (`ApplicationPresets` in
`pkg/api/business/profiles/application_presets.go`), the `init()` and its
variable carry `//nolint:gochecknoinits` and `//nolint:gochecknoglobals`, and a
comment says why a constructor does not fit.

---

## Compile-Time Interface Checks

Scope: Types that must satisfy an interface they are not passed as directly

Rule: Assert the relationship next to the type, so a changed method signature
fails the build in the implementing package instead of at a distant call site:
`var _ ConfigLoader = (*ConfigManager)(nil)` (`pkg/ajan/configfx/manager.go`),
`var _ Querier = (*Queries)(nil)` (generated by sqlc).
