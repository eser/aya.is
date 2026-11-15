# Inputs and Settings

Where values come from: validated input, settings and paths. Validation at
trust boundaries (HTTP, files, webhooks, provider and model output) is owned by
security-practices: Input Validation.

## Contents

- Input Validation
- No Machine-Specific Assumptions
- Settings Precedence

---

## Input Validation

Scope: Functions that receive data they did not produce

Rule: Validate data where it enters the program, following security-practices:
Input Validation, and turn it into a typed value there. Code past that boundary
trusts the types and does not re-validate. Reuse the validator the codebase
already applies to that value (`profiles.IsValidLocale` for a locale path
parameter, the zod schema of the form for a field) instead of writing a new
bound. Refuse a value with an error that names the value and says what a valid
one looks like (errors.md: Errors Carry a Code). In TypeScript, parse
unknown input with a schema; do not cast it.

Correct:

```go
locale := ctx.Request.PathValue("locale")
if !profiles.IsValidLocale(locale) {
    return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
}
```

```ts
import { z } from "zod";

const CreateUserSchema = z.object({
  email: z.email(),
  name: z.string().min(1).max(100),
  age: z.number().int().min(0).max(150),
});

function createUser(input: unknown) {
  const validated = CreateUserSchema.parse(input); // throws on invalid input
  return userRepository.create(validated);
}
```

Incorrect:

```go
func readFile(root, name string) ([]byte, error) {
    return os.ReadFile(filepath.Join(root, name)) // "../" reaches any file
}
```

```ts
function createUser(input: any) {
  return userRepository.create(input); // no validation
}

function createUser(email: string, age: number): User {
  return { email, age }; // nothing checks either value
}
```

**Why:** one check at the boundary is easy to audit; checks scattered through
the code get skipped on some path.

---

## No Machine-Specific Assumptions

Scope: Source, scripts, tests and configuration committed to the repository

Rule: No absolute paths, user names, home directories or local ports from one
machine. Derive every path from one of these anchors:

- Files that ship with a package (seeds, locales, fixtures) are embedded with
  `//go:embed` or read from the package directory in its tests (`testdata/`).
  Default locations are relative to the working directory and configurable
  (`migration_path` defaults to `etc/data/default/migrations`).
- Files that belong to a running component resolve against a root the caller
  passes in (`t.TempDir()` in tests), never a hard-coded location inside the
  function.
- Scripts run from the repository root, as the Makefile and the package
  scripts call them, and take machine-specific values from variables that name
  them. A test that needs a real service takes its address from a variable and
  is skipped when it is unset (`PRIVACY_TEST_DSN`, `S3_TEST_ENDPOINT`).
  Machine-only overrides live in the git-ignored `.env.local`,
  `.env.*.local` and `config.local.json`, never in a committed file.

Correct:

```go
func Open(stateDirectory string) (*Store, error) {
    return openStore(filepath.Join(stateDirectory, "cache"))
}
```

Incorrect:

```go
const stateDirectory = "/Users/me/aya-state" // one machine only
```

**Why:** taking the root as a parameter lets a test point the code at a
temporary directory without changing global state, and the code runs on any
machine and in CI.

---

## Settings Precedence

Scope: Settings of the Go services and the webclient server

Rule: A Go setting is one field of a config struct with a `conf` tag and, where
a safe value exists, a `default` tag (`pkg/ajan/config.go`,
`pkg/api/adapters/appcontext/config.go`). `configfx` loads in this order, each
source overriding the one before it: `config.json`, the `.env` file, then the
system environment; nested keys use `__` in a variable name (`LOG__LEVEL`,
`HTTP__CORS_ORIGIN`). A setting is declared once and read from the loaded
struct; code deep in a package does not call `os.Getenv` for a second name of
the same setting. A required setting has no default and is checked at
start-up (`AUTH__JWT_SECRET`, `ErrJWTSecretMissing`). An invalid value is an error at
start-up, never a silent fall back to the default, and the error names the
setting without repeating its value, which may be a secret. In the webclient,
`VITE_`-prefixed variables are bundled into the browser and carry nothing
secret; server-only values (`BACKEND_URI`) have no prefix
(`apps/webclient/.env.example`).

Correct:

```go
type FeatureFlags struct {
    DiscloseErrors bool `conf:"disclose_errors" default:"false"` // safe default
}

if a.Config.Auth.JwtSecret == "" {
    return fmt.Errorf("%w: %w", ErrInitFailed, ErrJWTSecretMissing) // names the setting, not a value
}
```

Incorrect:

```go
listen := os.Getenv("LISTEN") // a second name for the setting, read deep in a package
if listen == "" {
    listen = ":8080"
}
```

**Why:** one declaration per setting gives an operator one place to look and
one name to set, and a start-up error stops the server from running with a
value nobody chose.
