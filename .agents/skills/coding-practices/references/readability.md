# Readability

Naming, comments, abstraction, checks and control flow, in Go, TypeScript and
shell scripts. Go formatting belongs to the formatters of `.golangci.yaml` (`gofmt`,
`gofumpt`, `goimports`, `golines`), TypeScript formatting to oxfmt
(`printWidth: 120`, `pnpm fmt:check` at the repository root) and oxlint checks
the webclient (`pnpm lint` in `apps/webclient`); what they enforce is not
repeated here.

## Contents

- Self-Documenting Code
- Comments
- Don't Repeat Yourself
- Magic Values
- Explicit Checks
- Exhaustive Branching
- Identity Comparison
- Early Returns
- Line Length and Breaking
- Prefix Extraction
- No Runtime Code Generation or Global Patching
- Circular Dependencies

---

## Self-Documenting Code

Scope: All languages

Rule: Name similar operations the same way within a package and across
packages: all readers are `read*`, not a mix of `read*`, `fetch*` and `get*`.
Use full words, not abbreviations (`userAuthenticationService`, not `uas`;
`items` and `taxRate`, not `arr` and `r`). A product concept uses one noun
compound in every layer: `ProfilePage` in Go, `profilePage` in TypeScript,
`profile_page` in JSON keys and tables. Casing follows the language: Go's
exported and unexported names, camelCase in TypeScript.

Correct:

```ts
function calculateTotalPrice(items: Item[], taxRate: number): number {
  const subtotal = items.reduce((sum, item) => sum + item.price, 0);
  return subtotal + subtotal * taxRate;
}
```

Incorrect:

```ts
function calc(arr: Item[], r: number): number {
  const s = arr.reduce((a, b) => a + b.price, 0); // what is s?
  return s + s * r;
}

class UAS {} // unclear abbreviation
```

**Why:** consistent verbs and one name per concept let a reader and a search
find every place a concept appears; a mix hides some of them.

---

## Comments

Scope: All languages

Rule: Comments explain why: a trade-off, a constraint, a non-obvious rule, the
reason for a value. Every deliberate shortcut or exception (a `//nolint`, an
ignored error, a skipped case, a test that stays serial) carries a comment
giving the reason. A Go doc comment states what the identifier is and does,
including its zero value and the errors it returns. Do not commit
commented-out code. A `TODO` or `FIXME` names its owner, as `TODO(@eser)` does
in `pkg/ajan`, so it can be found and answered.

Correct:

```go
&mcp.Implementation{ //nolint:exhaustruct // external SDK type
```

```ts
// Cache results for 5 minutes to reduce database load.
// Trade-off: slight staleness for better performance.
const cache = new Map<string, CachedValue>();
```

Incorrect:

```go
&mcp.Implementation{ //nolint:exhaustruct
```

```ts
// This function finds a user
function findUser(users: User[], id: string): User | null {} // restates the code

const x = 5 * 60 * 1000; // no explanation of why this value
```

**Why:** an unexplained exception looks like an oversight, so the next person
either copies it or removes it without knowing what breaks. Comment the
trade-off, the non-obvious algorithm and the constraint or assumption; the code
already shows the what.

---

## Don't Repeat Yourself

Scope: All languages

Rule: Extract repeated logic at its third real use, or at the second when the
logic is complex or must stay in sync (a validation rule, a path layout, a wire
format). A form that a writer and a reader both depend on has one
implementation, in a shared package, not a copy per side. Lines that only look
alike but encode different decisions are not repetition; keep them apart.
Prefer a plain function over a configurable abstraction; add parameters when a
caller needs them, not in advance.

- Used 3 or more times: abstract.
- Used twice: abstract when the logic is complex or must stay in sync.
- Used once: keep it inline unless extracting improves clarity.

Correct:

```ts
function formatUserName(user: User): string {
  return `${user.firstName} ${user.lastName}`;
}

console.log(`Welcome, ${formatUserName(user)}`);
email.send(`Hello ${formatUserName(user)}`);
```

Incorrect:

```ts
console.log(`Welcome, ${user.firstName} ${user.lastName}`); // duplicated
email.send(`Hello ${user.firstName} ${user.lastName}`); // duplicated
```

**Why:** duplicated logic drifts apart when one copy is fixed; a premature
abstraction couples unrelated callers and is harder to undo than a copy.

---

## Magic Values

Scope: All languages

Rule: A number or string with meaning beyond its value (a limit, a timeout, a
status name, a path segment) is a named constant, and its unit is in the name
(`TOKEN_REFRESH_THRESHOLD_MS`, `summaryMaxTokens`, `secondsPerYear`). The
comment on it says why that value. Before choosing a bound for a field, look
for the rule the code already applies to it and reuse it; never replace a named
constant with a literal. Durations count too, in tests as well as in code: a
wait, a deadline, a poll interval or a fallback delay is a named constant,
never an inline `10 * time.Second` or `setTimeout(fn, 250)`. A value the code
under test needs shorter in tests is a config field whose default is the named
production value, which the test sets from its own named constant.

Correct:

```go
// providerRequestTimeout bounds one provider call, so a stalled
// provider cannot hold the caller.
const providerRequestTimeout = 10 * time.Second

ctx, cancel := context.WithTimeout(ctx, providerRequestTimeout)
```

```ts
const MAX_RETRIES = 3;
const API_TIMEOUT_MS = 5000;

if (retries >= MAX_RETRIES) {}
setTimeout(callback, API_TIMEOUT_MS);
```

Incorrect:

```go
ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
```

```ts
if (retries >= 3) {} // what is 3?
setTimeout(callback, 5000); // what is 5000?
```

**Why:** the name says what the value means and gives one place to change it;
a literal drifts from the rule it copied. An unnamed wait in a test hides
whether the test proves a deadline or only sleeps, which is what makes suites
slow and flaky.

---

## Explicit Checks

Scope: TypeScript and JavaScript (AGENTS.md: Explicit Checks)

Rule: Use implicit truthy/falsy checks only on booleans. Compare everything
else explicitly: `=== null`, `=== undefined`, `=== ""`, `.length === 0`,
`=== 0`, `!== -1`. `value == null` is the one loose comparison allowed: it
checks `null` and `undefined` together. Defaults use `??`, never `||`; equality
is otherwise always `===`.

Correct:

```ts
if (value === null) {}
if (value !== undefined) {}
if (string === "") {}
if (array.length === 0) {}
if (count === 0) {}
if (index !== -1) {}

if (!response.ok) {} // boolean
if (isValid) {}
if (user.isActive) {}

const result = value ?? defaultValue;
```

Incorrect:

```ts
if (!value) {} // also true for 0, "" and false
if (!string) {} // also true for ""
if (!array.length) {} // implicit number check
if (user) {} // ambiguous: null, undefined or a real object?

const result = value || defaultValue; // replaces a valid 0 or ""
```

**Why:** `0`, `""`, `false`, `null` and `undefined` are all falsy, and the
first three are often valid values; an implicit check treats them as missing and
the bug shows up only with those inputs. Explicit checks state the intent and
make review easier.

---

## Exhaustive Branching

Scope: `switch` or `if` chains over enums, codes and kinds

Rule: Handle every case. In Go, a `switch` over an enum-like type has a
`default` that returns an error or the documented fallback for an unknown
value. In TypeScript, a lookup table or a `switch` over a kind has an explicit
fallback for a value the code does not know, never a silent `undefined`; over a
union type, assign the value to `never` in the `default` so the compiler flags
the case a new member adds.

Correct:

```go
switch role {
case roleOwner, roleLead, roleMember:
default:
    return fmt.Errorf("%w: %q", ErrUnknownRole, role) // a role this build does not know is refused
}
```

```ts
const label = labels[kind] ?? t("Unknown");

switch (status) {
  case "draft": return t("Draft");
  case "published": return t("Published");
  default: {
    const unreachable: never = status; // a new status is a compile error here
    throw new Error(`unhandled status: ${unreachable}`);
  }
}
```

Incorrect:

```ts
switch (row.status) {
  case "refused": return t(row.message);
} // any other status renders nothing
```

**Why:** an explicit fallback makes the new case visible where it lands,
instead of as an empty cell or a nil value far away.

---

## Identity Comparison

Scope: All languages

Rule: Always compare entities by their stable unique identifiers (IDs), never
by display strings such as slugs, usernames, email addresses, host labels or
titles. Strings are mutable, can be duplicated and may differ in whitespace,
casing and encoding; IDs are immutable and unique. A slug from a URL is a
lookup key: resolve it to the record once, then compare IDs. A slug really does
change (deleting an account renames it to `deleted-<id>`, AGENTS.md: Privacy
(export, hide, delete)).

Correct:

```ts
// Compare by ID: stable and unambiguous
const isAuthorProfile = story.author_profile?.id === profile.id;
const isSameUser = comment.user_id === currentUser.id;
const isOwnProfile = membership.profile_id === userProfile.id;
```

```go
if story.AuthorProfileID == profile.ID {
    // author's own profile
}
```

Incorrect:

```ts
const isAuthorProfile = story.author_profile?.slug === slug; // a slug can change
const isSameUser = comment.username === currentUser.username; // a username can change, casing differs
const isOwner = member.email === user.email; // an address is renamed or normalized
```

**Why:** slugs, usernames and emails are renamed and reformatted, and string
comparison is vulnerable to whitespace, casing and encoding differences. IDs
are immutable, unique and indexed, so using them makes refactoring safer:
renaming a slug cannot break the logic.

---

## Early Returns

Scope: All languages

Rule: Handle errors and edge cases first and return early, so the happy path
stays at the lowest indentation and reads straight down. Each guard checks one
condition and leaves; the main work comes after all of them, not in an `else`.
New validation is then one more guard at the top.

Correct:

```go
func readStory(store Store, id string) (Story, error) {
    if id == "" {
        return Story{}, ErrInvalidStoryID
    }

    story, err := store.Story(id)
    if err != nil {
        return Story{}, fmt.Errorf("read story %q: %w", id, err)
    }

    return story, nil
}
```

```ts
function processUser(user: User | null): Result {
  if (user === null) return { error: "User not found" };
  if (!user.isActive) return { error: "User is inactive" };

  return { data: loadProfile(user.id) }; // happy path, no nesting
}
```

Incorrect:

```ts
function processUser(user: User | null): Result {
  if (user !== null) {
    if (user.isActive) {
      return { data: loadProfile(user.id) };
    } else {
      return { error: "User is inactive" };
    }
  } else {
    return { error: "User not found" };
  }
}
```

**Why:** nested branches make the reader hold every condition at once; guards
retire them one by one, and the error handling is visible at the top.

---

## Line Length and Breaking

Scope: Go, and shell scripts

Rule: `gofmt` does not wrap lines; `golines` (a formatter in
`.golangci.yaml`) wraps what is too long. When a call or condition no longer reads at a glance,
break it where the structure already has a seam (between arguments, after an
operator, between key-value pairs of a log call), not at whatever column the
line reaches. A signature that is too long to read usually has too many
parameters; the fix is a `Config` struct (go-practices: Constructors and
Config), not cleverer wrapping. In TypeScript, oxfmt decides; keep the layout
of the file you edit.

Correct:

```go
p.logger.WarnContext(ctx, "Failed to create request for PDF URL",
    slog.String("url", presentationURL),
    slog.Any("error", err))

func New(config Config) (*Service, error)
```

Incorrect:

```go
func New(store Store, models Models, inspector Inspector, logger *logfx.Logger, now func() time.Time, workers int, retry time.Duration) (*Service, error)
```

**Why:** a line broken at its seams can be read and diffed argument by
argument.

---

## Prefix Extraction

Scope: Taking a value after a known prefix: HTTP headers, URIs, file paths,
protocol strings, setting values

Rule: Check that the value starts with the prefix, then cut exactly the
prefix. In Go that is `strings.CutPrefix`; in TypeScript `startsWith()` to
validate and `slice(prefix.length)` to extract. Never use a search-and-replace
for this: it matches the first occurrence anywhere in the string, so
`"xBearer token"` or `"s3://file://bucket"` are silently rewritten. Never slice
a hard-coded length (`slice(7)`): it does not check that the prefix is there
and breaks when the prefix changes. Return "no value" when the prefix is
missing or nothing follows it.

Correct:

```go
token, ok := strings.CutPrefix(authHeader, "Bearer ")
if !ok || token == "" {
    return "", false
}
```

```ts
// Generic helper for any prefix extraction
function extractAfterPrefix(value: string | null | undefined, prefix: string): string | null {
  if (value === null || value === undefined) return null;
  if (!value.startsWith(prefix)) return null;

  const extracted = value.slice(prefix.length);
  return extracted.length > 0 ? extracted : null;
}

const BEARER_PREFIX = "Bearer ";
const token = extractAfterPrefix(authHeader, BEARER_PREFIX);
const filePath = extractAfterPrefix(url, "file://");
```

Incorrect:

```go
token := strings.Replace(authHeader, "Bearer ", "", 1) // matches anywhere
```

```ts
const token = authHeader?.replace("Bearer ", "") ?? null;
// "xBearer token" becomes "xtoken": the prefix position is not validated

const path = uri?.replace("file://", "") ?? null;
// "s3://file://bucket" becomes "s3://bucket"

const token = authHeader?.slice(7) ?? null;
// does not check that "Bearer " is actually there
```

**Why:** a replace that matches mid-string accepts malformed input as valid,
which in a header parser is an authentication bug. `startsWith` with `slice`
(or `CutPrefix`) is explicit about the expected format and handles an empty
remainder on purpose. `getBearerToken` in
`pkg/ajan/httpfx/middlewares/auth_middleware.go` is the Go form.

---

## No Runtime Code Generation or Global Patching

Scope: All languages

Rule: Do not build and run code at runtime, and do not change shared built-ins
or other packages' internals. In JavaScript that means no `eval`, no
`new Function`, no prototype changes and no `Object.defineProperty` on shared
objects. `dangerouslySetInnerHTML` receives only markup that was escaped first
or generated from a fixed string (`renderInlineMarkdown` in
`src/lib/inline-markdown.tsx` escapes HTML before it applies formatting and
only allows `http`/`https` links); never a string that carries user input
(security-practices: No Injection Through Strings). In Go it means no `unsafe`
or `reflect` writes into another package's unexported state.

Correct:

```ts
const handlers = { add, remove };
const handler = handlers[action];
```

Incorrect:

```ts
eval(`${action}(input)`); // any action string becomes code
Array.prototype.last = function () {}; // changes every array in the page
```

**Why:** both break static analysis, and any input that reaches them becomes
code.

---

## Circular Dependencies

Scope: Imports between Go packages and TypeScript modules

Rule: Imports form a graph without cycles. Go refuses an import cycle at build
time; ES modules do not, so check the import graph after moving code between
files in `apps/webclient/src`. Break a cycle by moving the shared part into a
package or module both sides import, never by importing a module's internals:
between business packages the shared part is a port, an interface owned by the
consumer (`stories.ContentSummarizer`, implemented in `pkg/api/adapters/ai`;
go-practices: Modules, Ports and Adapters).

Incorrect:

```ts
// file-a.ts
import { funcB } from "./file-b.ts";
export function funcA() { funcB(); }

// file-b.ts
import { funcA } from "./file-a.ts"; // circular
export function funcB() { funcA(); }
```

**Why:** a cycle makes evaluation order decide which side sees `undefined`
exports, and the failure depends on which module is imported first.
