# Security Rules

Secrets, output hygiene, input handling, authorization, injection, outbound
requests, errors, production settings and model trust in the Go services and
the webclient.

## Contents

- Secrets Come From the Environment
- Secrets and Personal Data Never Reach Output
- Secure Randomness
- Input Validation
- Authorization in the Handler
- No Injection Through Strings
- Outbound Requests to User-Supplied URLs
- Error Sanitization
- Safe Defaults and Explicit Outcomes
- Production Deployment Checklist
- LLM Trust Boundary

---

## Secrets Come From the Environment

Scope: Keys, tokens, passwords, database configuration and provider
credentials, in the services and the webclient

Rule: A secret never lives in code, committed configuration or fixtures. It
reaches a process through the environment: a field of the config struct
(`conf` tag, no `default`) that `configfx` fills from the system environment,
the git-ignored `.env.local`, `.env.*.local` and `config.local.json` for local
development (apps/services/AGENTS.md: Local configuration), or the deployment
platform's secret variables. Examples are `AUTH__JWT_SECRET`, an AI provider's
`api_key`, the S3 credentials of the export bucket (`S3_EXPORT__ACCESS_KEY_ID`,
`S3_EXPORT__SECRET_ACCESS_KEY`), `PPROF_TOKEN` and the database DSN. A missing
required secret stops start-up with an error that names the setting, never the
value (`ErrJWTSecretMissing`). In the webclient, a `VITE_` variable is bundled
into the browser and so is public; a secret has no `VITE_` prefix and is read
only in server code (`apps/webclient/.env.example`).

Correct:

```go
type Config struct {
    JwtSecret string `conf:"jwt_secret"` // required: no default, checked at start-up
}

if a.Config.Auth.JwtSecret == "" {
    return fmt.Errorf("%w: %w", ErrInitFailed, ErrJWTSecretMissing)
}
```

```ts
// src/env.ts declares server and client variables with zod (@t3-oss/env-core)
export const env = createEnv({
  server: { PROVIDER_API_KEY: z.string().min(1) }, // illustrative: required, no default, never VITE_
  clientPrefix: "VITE_",
  client: {},
});
```

Incorrect:

```go
jwtSecret := "my-super-secret-key" // hardcoded secret
```

```json
{ "jwtSecret": "my-super-secret-key" }
```

**Why:** a secret in the repository is a secret in every clone, fork and CI log
for as long as the history lives; the environment keeps it out of both and
lets each deployment have its own.

---

## Secrets and Personal Data Never Reach Output

Scope: Logs, audit records, error answers, telemetry, CLI output and the
browser

Rule: Credentials, keys, tokens and passwords never appear in any output, and
neither does content or personal data: request and response bodies, headers,
personal text, e-mail addresses. A session ID is a credential (it is the value
of the session cookie), so it is not logged either. Logs and audit records
carry identifiers, codes, counts and fixed outcomes. Put input in `slog`
attributes, never in the message string, and strip query strings from URLs
before they are logged. A new sink (an OpenTelemetry attribute, a webhook, an
analytics script) gets the same review before it is wired in.

- Personal data has one lifecycle: a new table or column that holds personal
  data is added to both the export datasets and the deletion steps, and a
  secret column to `secretColumns` (AGENTS.md: Privacy (export, hide, delete)).
- A type that holds a secret is never logged or formatted with `%v`; log its
  ID only.

A change that adds a log line or event on a path that handles content adds a
test that puts a sentinel string in the content and asserts it is absent from
the output.

Correct:

```go
logger.WarnContext(ctx, "provider call failed",
    slog.String("provider", provider), slog.String("user_id", userID), slog.Any("error", err))
```

Incorrect:

```go
slog.Warn("failed to update session activity", slog.String("session_id", sessionID)) // the cookie value
logger.Debug("request", slog.String("body", string(body)))                           // content in a log
```

**Why:** logs are copied to more places and read by more people than the
database is, and a leaked credential or personal text in them cannot be taken
back.

---

## Secure Randomness

Scope: Keys, tokens, OAuth state, PKCE verifiers, one-time codes, session
references, request IDs

Rule: Anything an attacker must not guess comes from `crypto/rand` in Go and
`crypto.getRandomValues` or `crypto.randomUUID` in TypeScript. A failing random
source fails the operation. `math/rand/v2` and `Math.random` are only for
choices that are not secret, such as jitter in a retry schedule
(`pkg/ajan/httpclient/retry_strategy.go`) or a skeleton's width.

Correct:

```go
secret := make([]byte, 32)
if _, err := rand.Read(secret); err != nil { // crypto/rand
    return "", fmt.Errorf("generate token: %w", err)
}
```

```ts
generateId: () => crypto.randomUUID(),
```

Incorrect:

```go
code := fmt.Sprint(mathrand.IntN(1_000_000)) // a guessable one-time code
```

**Why:** a predictable value is a guessable value; the standard random sources
cost nothing extra.

---

## Input Validation

Scope: Everything that enters a process: HTTP bodies, headers, path and query
parameters, settings, files, webhook and provider answers, and model output

Rule: Validate at the boundary where the input enters, before business logic,
and fail fast with an error that names the check. Bound the input first
(`http.MaxBytesReader`, or the `RequestSizeLimitMiddleware` of
`pkg/ajan/httpfx/middlewares`), then decode into a concrete type, never `any`,
and refuse unknown fields and a second document
(`decoder.DisallowUnknownFields()`). Names, IDs, slugs, locales and
enumerations are checked against an allowlist and a length bound; reuse the
validator the code already has (`profiles.IsValidLocale`). In the webclient,
parse `unknown` input with a zod schema. When adding an entry point, name its
inputs and the check for each in the change.

Correct:

```go
body := http.MaxBytesReader(ctx.ResponseWriter, ctx.Request.Body, maxBodyBytes)
decoder := json.NewDecoder(body)
decoder.DisallowUnknownFields()

var input CreateUserRequest
if err := decoder.Decode(&input); err != nil {
    return ctx.Results.BadRequest(httpfx.WithErrorMessage("invalid request body"))
}
if !isValidEmail(input.Email) {
    return ctx.Results.BadRequest(httpfx.WithErrorMessage("invalid email"))
}
```

```ts
const CreateUserSchema = z.object({
  email: z.email(),
  name: z.string().min(1).max(100),
  age: z.number().int().min(0).max(150),
});

const validated = CreateUserSchema.parse(input); // throws on invalid input
```

Incorrect:

```go
body, _ := io.ReadAll(ctx.Request.Body) // unbounded
_ = json.Unmarshal(body, &input)        // unknown keys pass, errors ignored
```

```ts
function createUser(input: any) {
  return userRepository.create(input); // no validation
}
```

**Why:** the boundary is the one place where every input can be seen; a check
placed later is skipped by the first path that does not go through it.

---

## Authorization in the Handler

Scope: Every route, MCP tool and action taken on a person's behalf

Rule: A protected route runs behind `AuthMiddleware(authService, userService)`,
which resolves the session from the cookie or the `Authorization` bearer token,
loads it and requires it to be active and logged in. The handler then passes
the acting user's ID (`*session.LoggedInUserID`) and the resource to a service
method, and that method checks the permission for this resource on every
request (`profiles.ErrInsufficientAccess`, answered as 403). Deny by default:
only an explicit allowance allows; a missing session answers 401, a failed
permission lookup answers an error, and neither falls through to the action.
Decide on purpose whether a refusal may reveal that a resource exists. A hidden
button or route is presentation, not access control. Tests cover the denial
paths: no session, another person's resource and a missing permission.

Correct:

```go
teams, err := profileService.ListTeams(ctx.Request.Context(), *session.LoggedInUserID, slugParam)
if err != nil {
    statusCode := http.StatusInternalServerError
    if errors.Is(err, profiles.ErrInsufficientAccess) {
        statusCode = http.StatusForbidden
    }

    return ctx.Results.Error(statusCode, httpfx.WithSanitizedError(err))
}
```

Incorrect:

```go
// the UI hides the button, so any caller that gets here is allowed
return profileService.DeleteTeam(ctx.Request.Context(), teamID)
```

**Why:** the client is not under our control; only a check next to the data
holds when someone calls the API directly.

---

## No Injection Through Strings

Scope: SQL, subprocesses, file paths and markup built from input

Rule: Never build a query, command, path or markup by concatenating input.

- SQL: queries live in `etc/data/default/queries/` and sqlc generates the
  access code with placeholders; never format a value into SQL with
  `fmt.Sprintf` or string concatenation. An identifier comes from a package
  constant, never from input.
- Subprocesses: `exec.CommandContext` with an argument array and an explicit
  environment; a secret goes through stdin, never argv; no `sh -c` with input.
- Paths: confine a path from input to a root with `os.OpenRoot` and its methods
  instead of joining and cleaning.
- Markup: React escapes text. `dangerouslySetInnerHTML` receives only markup
  that was escaped first or generated from a fixed string
  (`renderInlineMarkdown` in `src/lib/inline-markdown.tsx` escapes HTML before
  formatting and links only `http` and `https`); a string with user input goes
  through the same escaping or a sanitizer.

Correct:

```go
cmd := exec.CommandContext(ctx, "/usr/bin/convert", "-resize", size, inputPath, outputPath)
```

```ts
function escapeHtml(text: string): string {
  return text.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;").replace(/"/g, "&quot;");
}
```

Incorrect:

```go
cmd := exec.CommandContext(ctx, "/bin/sh", "-c", "convert -resize "+size+" "+inputPath)
row := db.QueryRow("SELECT * FROM story WHERE slug = '" + slug + "'")
```

```tsx
<div dangerouslySetInnerHTML={{ __html: comment.body }} /> // user text as markup
```

**Why:** input that becomes code, a path or markup is attacker-controlled
behavior; placeholders, argument arrays and escaping keep data as data.

---

## Outbound Requests to User-Supplied URLs

Scope: Webhooks, callbacks, link and image imports, provider `base_url`
settings and any request whose address comes from settings or input

Rule: Validate the address before the request and check what the connection
actually dials. `lib.ValidateExternalURL(rawURL, requireHTTPS)` in
`pkg/ajan/lib/net.go` parses the URL, requires `http` or `https` (only
`https` when `requireHTTPS` is true), resolves the host and refuses it if any
resolved IP is private or reserved. Use it for every such request, with
`requireHTTPS` true outside local development. Because the check resolves
once and the client resolves again, a DNS answer can change in between;
for a caller that fetches on behalf of users, also check in the dialer
(`net.Dialer.ControlContext`) or dial only the checked IPs. Apply the check to
every redirect or refuse redirects.

Blocked ranges (`privateNetworks` in `pkg/ajan/lib/net.go`):

| Range | Meaning |
| --- | --- |
| `10.0.0.0/8`, `172.16.0.0/12`, `192.168.0.0/16` | RFC 1918 private networks |
| `127.0.0.0/8`, `::1/128` | loopback |
| `169.254.0.0/16`, `fe80::/10` | link-local, including cloud metadata addresses |
| `0.0.0.0/8` | unspecified |
| `100.64.0.0/10` | shared address space (RFC 6598) |
| `fc00::/7` | IPv6 unique local |

Correct:

```go
if err := lib.ValidateExternalURL(webhookURL, true); err != nil {
    return fmt.Errorf("webhook url: %w", err)
}

resp, err := client.Do(req) // client with a timeout and redirects refused or re-validated
```

Incorrect:

```ts
async function sendWebhook(url: string, payload: unknown) {
  await fetch(url, { method: "POST", body: JSON.stringify(payload) }); // no validation
}
```

**Why:** a server that fetches an address it was given can be pointed at its own
network or the cloud metadata service; the response of an internal service then
leaves through the public one.

---

## Error Sanitization

Scope: Error answers of the HTTP API, the MCP server and server-side rendering

Rule: An answer carries a stable code or a fixed message chosen by the kind of
failure, plus what helps the caller fix the request. It never carries a Go
error string, a stack, SQL, an internal path or echoed input. An unclassified
error is logged with its cause and answered generically:
`httpfx.WithSanitizedError(err)` logs the full error server-side and returns
`{"error": "an error occurred"}`. The only switch is `features.disclose_errors`
(`FEATURES__DISCLOSE_ERRORS`, default false), which returns the real message,
and it stays off in production. Callers match a sentinel with `errors.Is` or the
`code` of the body, never the text (coding-practices: Errors Carry a Code). A
refusal that would reveal whether something exists (an unknown account versus a
wrong credential) answers the same for every cause.

Correct:

```go
default:
    return ctx.Results.Error(http.StatusInternalServerError, httpfx.WithSanitizedError(err))
```

Incorrect:

```go
return ctx.Results.Error(http.StatusInternalServerError, httpfx.WithErrorMessage(err.Error())) // internal text to the client
```

```ts
res.status(500).json({ error: error.message, stack: error.stack, query: req.query });
```

**Why:** an error string names tables, paths and library versions; handing it
to a client does half of an attacker's reconnaissance.

---

## Safe Defaults and Explicit Outcomes

Scope: Settings, environment-specific behavior and anything that relaxes a check

Rule: A setting that weakens security defaults to the safe value and is turned
on by name: `disclose_errors`, `expose_internal_errors` and `profiling` default
to `false` (`pkg/api/adapters/appcontext/config.go`,
`pkg/ajan/httpfx/config.go`). A check that cannot run fails closed. A
relaxation is an explicit setting, never code that skips a check. Behavior
differs between development and production through named settings and the
config files selected by `env` (`config.json`, `config.<env>.json`), not
through `if env == "production"` scattered in code. Development may enable
profiling, plain HTTP webhooks and detailed errors; production has strict TLS,
no profiling endpoint unless `PPROF_TOKEN` protects it, sanitized errors and
minimal logging. TLS verification is never turned off to make something work,
and TLS 1.2 is the floor for serving (`MinVersion` in
`pkg/ajan/httpfx/http_service.go`).

Correct:

```go
ProfilingEnabled bool `conf:"profiling" default:"false"`

pprofToken := os.Getenv("PPROF_TOKEN")
if pprofToken == "" {
    slog.Warn("PPROF_TOKEN not set - pprof endpoints are unprotected")
}
```

Incorrect:

```go
tlsConfig := &tls.Config{InsecureSkipVerify: true} // "just for now"
```

**Why:** a default is what runs when nobody decided; making it the safe value
means a forgotten setting costs a feature, not a breach.

---

## Production Deployment Checklist

Scope: Production deployments (Coolify)

Rule: Before a deployment goes live, check each line:

- [ ] `AUTH__JWT_SECRET` and every provider key come from the platform's secret
      variables, not from a config file
- [ ] No `config.json` or `.env` committed with real credentials
- [ ] `FEATURES__DISCLOSE_ERRORS` is false (the default)
- [ ] Profiling is off, or `PPROF_TOKEN` is set
- [ ] Rate limiting and the request size limit are applied to the listener
      (`HTTP__RATE_LIMIT_REQUESTS`, `HTTP__MAX_REQUEST_SIZE_MB`)
- [ ] TLS for every external connection, including the database
- [ ] The export bucket is private and has no public domain
      (AGENTS.md: Privacy (export, hide, delete))
- [ ] `govulncheck` passes (`make ok`; apps/services/AGENTS.md: Linting and CI)

Incorrect:

```bash
FEATURES__DISCLOSE_ERRORS=true  # leaks internal details
```

**Why:** each line is a default someone can forget to carry into production;
one list is cheaper than finding out from an incident.

---

## LLM Trust Boundary

Scope: Summaries, translations and generated profile text (`pkg/ajan/aifx`,
`pkg/api/adapters/ai`), and the MCP server (`pkg/api/adapters/mcp`)

Rule: Model output is untrusted input. Read it as bounded text (set
`MaxTokens`, as `summaryMaxTokens` does), validate any structure it claims
before use, and store and render it exactly like user-written content, with the
same escaping. Model output never reaches a shell, a file path, SQL or a URL
without the checks any other input gets. Text from stories, profiles or users
that goes into a prompt is content the model reads, never an instruction the
service follows: keep the system prompt separate from the user message, and
give the model no tool that acts on its own.

The MCP server at `/mcp` is anonymous and its tools are read-only
(`list_*`, `get_*`, `search`). A tool returns only what an anonymous visitor
may read, validates its arguments like any other input, and a tool that writes
data or returns non-public data authenticates and authorizes like a route
(Authorization in the Handler).

Correct:

```go
aifx.NewTextMessage(aifx.RoleUser, buildSummaryPrompt(story, language)), // content as the user message
// System: summarySystemPrompt, MaxTokens: summaryMaxTokens
```

Incorrect:

```go
prompt := systemPrompt + "\n" + story.Content // content mixed into the instructions
html := model.Output                          // rendered as markup without escaping
```

**Why:** a model repeats what it read, including instructions an attacker put
in a story or a profile; treating its answer as data keeps those instructions
from becoming actions.
