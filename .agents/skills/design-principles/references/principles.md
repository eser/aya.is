# Design Principles

How functions, modules, data and state are shaped in this repository, in Go
(`apps/services`) and in the TypeScript of `apps/webclient`. Architecture-level
structure (ports, adapters, composition roots, module boundaries) is in
architecture-guidelines: Hexagonal Architecture.

## Contents

- Pure Functions
- Immutability
- Small Interfaces
- Data as Plain Objects
- Single Responsibility
- One Clock Read per Operation
- Cancellation Flows With the Work
- Choose Data Structures by Access Pattern
- Async Patterns
- String Formatting
- Loop Control Flow
- Module-Level Side Effects

Early returns are in coding-practices: Early Returns.

---

## Pure Functions

Scope: All languages

Rule: Put logic in pure functions that take everything they need as parameters
and return a result without touching outside state: same input, same output.
Keep that logic away from the environment, flags, the file system, the network
and React context; a thin wrapper reads those and passes plain values in. A
function that must have an effect keeps it at the edge: compute first, then
perform one effect with the result. Do not write to a variable outside the
function to hand a result back; return it. In the webclient, put the pure part
in its own file (`foo-utils.ts`) so a test can import it without
`import.meta.env`, and let the environment-dependent file (`foo.ts`) re-export
it (AGENTS.md: Pure Function Extraction for Testability).

Correct:

```typescript
// locale-utils.ts: pure, no import.meta.env, tested with Vitest
export function isValidLocale(locale: string): boolean { /* ... */ }

// config.ts: re-exports the pure part and holds the environment-dependent config
export { isValidLocale } from "@/lib/locale-utils.ts";
export const siteConfig = { host: import.meta.env.VITE_SITE_URL };
```

```go
// storyAgeBucket is pure: tested with values. (illustrative)
func storyAgeBucket(publishedAt, now time.Time) string { /* ... */ }

// bucketOf is the thin wrapper: it reads state, then calls the pure part.
func (s *Service) bucketOf(ctx context.Context, id string) (string, error) {
	story, err := s.repo.GetStory(ctx, id)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedToGetRecord, err)
	}

	return storyAgeBucket(story.PublishedAt, s.now()), nil
}
```

Incorrect:

```typescript
// Pure logic mixed with an environment dependency: a test cannot import it
import { siteConfig } from "@/config.ts";
export function buildUrl(locale: string) { return `${siteConfig.host}/${locale}`; }

let total = 0;
function addToTotal(item: Item) { // modifies external state
  total += item.price;
}
```

```go
var lastBucket string // hidden state shared by every caller

func storyAgeBucket(publishedAt time.Time) string {
	now := time.Now()             // reads the clock
	lastBucket = os.Getenv("MODE") // reads the environment
	// ...
}
```

**Why:** a pure function is tested with plain values, and a change in the
environment cannot change its result.

---

## Immutability

Scope: All languages

Rule: Do not change inputs or shared data; return a new value. In the
webclient, build new objects and arrays (`{...row, name}`, `toSorted`, `map`,
`filter`) instead of changing the ones a query or a prop holds; React and the
React Query cache compare by reference. In Go, a function that receives a slice
or map it does not own does not write into it, and a type hands out copies of
its internal slices and maps (`slices.Clone`, `maps.Clone`). Changing a local
value that nobody else can see is fine.

Correct:

```typescript
const sorted = rows.toSorted((a, b) => a.label.localeCompare(b.label));
const updated = { ...user, age: 31 };
const withItem = [...items, 4];
```

```go
func withGroup(groups []string, id string) []string {
	return append(slices.Clone(groups), id)
}
```

Incorrect:

```typescript
rows.sort((a, b) => a.label.localeCompare(b.label)); // changes the cached query data
user.age = 31;                                        // the caller's object changes too
items.push(4);
```

**Why:** a caller that passed a value does not expect it to change; in the
webclient a changed object also hides the change from React.

---

## Small Interfaces

Scope: Interfaces and the functions that accept them

Rule: Define small, purpose-specific interfaces in the package that uses them
(a port), with only the methods that package calls. Constructors accept the
interface, not a concrete type, so the composition root passes the real adapter
and tests pass a fake. Prefer composition over inheritance: pass dependencies
in (a constructor parameter, a field, a prop) instead of extending a base class.
Do not add a wrapper that only forwards to another method; call the method
directly. A module's own repository interface may be wide when it is the one
thing that talks to storage; every external service (AI, mail, object storage)
gets its own narrow port.

Correct:

```go
// ContentTranslator defines the port for AI-powered content translation.
// Implementations live in the adapter layer.
type ContentTranslator interface {
	Translate(
		ctx context.Context,
		sourceLocale, targetLocale string,
		title, summary, content string,
	) (translatedTitle, translatedSummary, translatedContent string, err error)
}
```

```typescript
interface Logger {
  log(message: string): void;
}

class UserService {
  constructor(private logger: Logger) {}
}
```

Incorrect:

```go
type Services interface {
	Translate(ctx context.Context, /* ... */) error
	SendMail(ctx context.Context, /* ... */) error
	PutObject(ctx context.Context, /* ... */) error
	Migrate(ctx context.Context) error // every consumer depends on all of it
}
```

```typescript
class BaseService {
  log(message: string) {}
}

class UserService extends BaseService { // tight coupling
  createUser(data: UserData) {
    this.log("Creating user");
  }
}
```

**Why:** a small interface is cheap to fake in a test and leaves each consumer
coupled only to what it uses; inheritance couples a class to everything its
parent does.

---

## Data as Plain Objects

Scope: Types in Go and objects in TypeScript

Rule: Model data as plain structs and objects and put behavior in functions
that take them. In Go, a record type has exported fields and JSON tags and no
hidden state; a service type with unexported fields owns resources and caches.
In the webclient, data from the API stays plain objects (use `readonly` fields
or return copies); components and hooks hold state. Use a class only for a
stateful service (a cache, a connection). No getters or setters around records:
a read or a write with a side effect surprises every caller.

Correct:

```go
type SlugAvailabilityResult struct {
	Message   string `json:"message,omitempty"`
	Severity  string `json:"severity,omitempty"`
	Available bool   `json:"available"`
}
```

```typescript
interface User {
  readonly id: string;
  readonly name: string;
  readonly email: string;
}

function updateUserEmail(user: User, email: string): User {
  return { ...user, email };
}

class UserRepository { // a stateful service: owns a cache
  private cache = new Map<string, User>();
}
```

Incorrect:

```typescript
class User {
  get email(): string {
    this.logAccess(); // a read with a side effect
    return this._email;
  }

  set email(value: string) {
    this.validateEmail(value); // a write with a side effect
    this._email = value;
  }
}
```

**Why:** plain data serializes, compares and copies without surprises; state
belongs in the few services that manage it.

---

## Single Responsibility

Scope: All languages

Rule: A function does one task and has a name that says which. When a function
validates, saves, notifies and logs, split it into those steps and let a short
caller sequence them. The steps stay testable alone, and the caller reads as the
list of what happens.

Correct:

```typescript
function validateEmail(email: string): boolean {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
}

function sendWelcomeEmail(email: string): void {
  sendEmail(email, createWelcomeMessage());
}
```

Incorrect:

```typescript
function processUser(email: string, data: UserData): void {
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) return; // validates
  database.save(data);                                      // saves
  sendEmail(email, "Welcome!");                             // notifies
  console.log("User processed");                            // logs
}
```

**Why:** a function with four reasons to change breaks for four unrelated
causes and cannot be reused for any one of them.

---

## One Clock Read per Operation

Scope: Code that timestamps or compares times

Rule: Read the clock once per logical operation and pass the value down. Do not
call `time.Now()` or `Date.now()` inside a loop or in several steps that must
agree on the same time. A Go type whose behavior depends on time takes its clock
as a `now func() time.Time` field that defaults to `time.Now` and has a setter
for tests (`privacy.Service.SetClock`), or a `clockfx.Clock`; a TypeScript
function takes `now` as a parameter that defaults to `Date.now()`.

Correct:

```go
now := s.now()
for _, item := range items {
	if item.ExpiresAt.Before(now) {
		expired = append(expired, item.ID)
	}
}
```

```typescript
export function isExpired(expiresAt: number, now = Date.now()): boolean {
  return expiresAt <= now;
}
```

Incorrect:

```go
for _, item := range items {
	if item.ExpiresAt.Before(time.Now()) { // a different moment per item
		expired = append(expired, item.ID)
	}
}
```

**Why:** two reads in one operation can straddle a boundary and give the steps
different answers about the same moment, and an injected clock lets a test prove
behavior over time without waiting.

---

## Cancellation Flows With the Work

Scope: Any operation that can be abandoned: HTTP requests, background jobs,
worker tasks, webclient queries

Rule: An operation carries one cancellation handle from the point it starts to
every call it makes: a `context.Context` in Go, an `AbortSignal` in the
webclient (React Query hands one to each query function as `signal`). Pass the
caller's handle down instead of creating a new one, so when the caller gives up
(a client disconnects, a person leaves the page, a deadline passes) everything
the operation started stops. A step may shorten the remaining time for its own
part, never extend it. Work that must finish after the caller leaves is
detached on purpose, gets its own named time limit and an owner that reports its
failure.

Correct:

```go
func (s *Service) translateAll(ctx context.Context, items []Item) error {
	for i, item := range items {
		if err := s.translator.Translate(ctx, /* ... */); err != nil {
			return fmt.Errorf("translate item %d of %d: %w", i+1, len(items), err)
		}
	}

	return nil
}
```

```typescript
useQuery({
  queryKey: ["stories", locale],
  queryFn: ({ signal }) => fetch(`/api/${locale}/stories`, { signal }).then((r) => r.json()),
});
```

Incorrect:

```go
_ = s.translator.Translate(context.Background(), /* ... */) // detached from the caller's deadline
```

The Go form is in go-practices (`context.md`).

**Why:** a chain that drops the handle at one link keeps consuming CPU,
connections and provider quota for a caller that is no longer there.

---

## Choose Data Structures by Access Pattern

Scope: Collections in Go and TypeScript

Rule: Pick the structure from how the data is read and written:

- Keyed lookup or membership: a `map` in Go, a `Map` or `Set` in TypeScript, not
  a linear search through a slice or array.
- Ordered data read in sequence: a slice or array. Use linked structures only
  for frequent insertion in the middle, and only after measuring.
- Known or estimable size: allocate once up front instead of growing in a loop.
- Shared state handed to a caller: return a copy, so the caller cannot change
  the owner's data.

Correct:

```typescript
const profilesById = new Map(profiles.map((p) => [p.id, p]));
const title = profilesById.get(id)?.title ?? id;
```

Incorrect:

```typescript
const title = profiles.find((p) => p.id === id)?.title ?? id; // O(n) per row in a loop
```

The Go form is in go-practices (`data-structures.md`).

**Why:** the wrong structure turns a lookup loop quadratic, and a shared
internal collection lets any caller corrupt the owner's state.

---

## Async Patterns

Scope: TypeScript

Rule: Use promises with `async` / `await` instead of callbacks. A function that
does asynchronous work returns a `Promise`; do not pass a continuation.

Correct:

```typescript
async function fetchUser(id: string): Promise<User> {
  const response = await fetch(`/api/users/${id}`);
  return response.json();
}
```

Incorrect:

```typescript
function fetchUser(id: string, callback: (user: User) => void) {
  fetch(`/api/users/${id}`).then((r) => r.json()).then(callback);
}
```

**Why:** a promise composes, propagates errors to one `catch` and can be awaited
by the caller; a callback does none of that.

---

## String Formatting

Scope: TypeScript

Rule: Build strings with template literals, not concatenation. Embedded
expressions and multi-line text read as written.

Correct:

```typescript
const greeting = `Hello, ${user.name}!`;
const url = `/api/users/${userId}/posts/${postId}`;
```

Incorrect:

```typescript
const greeting = "Hello, " + user.name + "!";
const url = "/api/users/" + userId + "/posts/" + postId;
```

**Why:** concatenation hides the shape of the result between quotes and plus
signs, and a missing space or `+` is easy to write and hard to see.

---

## Loop Control Flow

Scope: All languages

Rule: Use a plain `for` / `for...of` loop when the body needs `break`,
`continue` or an early `return`. Use array methods (`filter`, `map`, `reduce`)
when it needs none of them.

Correct:

```typescript
for (const item of items) {
  if (item.skip) continue;
  if (item.error) break;
  process(item);
}

const processed = items.filter((item) => !item.skip).map((item) => process(item));
```

Incorrect:

```typescript
items.forEach((item) => {
  if (item.error) return; // returns from the callback, not the loop: it does not stop
  process(item);
});
```

**Why:** `break` and `continue` do not exist in a callback, and a `return` there
only ends one iteration.

---

## Module-Level Side Effects

Scope: Go packages and TypeScript modules

Rule: Importing a package or module runs no I/O, reads no environment, starts no
goroutines or timers and does no heavy computation. Package and module scope
holds constants, types, compiled constant patterns and function declarations; no
global mutable variable (`let` at module scope, a package-level `var` that
changes). A value that is expensive or environment-dependent is computed on
first use and cached in a getter function, with a promise for async work. In Go,
`init()` and package-level mutable state are flagged by `gochecknoinits` and
`gochecknoglobals`; the few that exist fill a constant lookup table and carry a
`//nolint` with the reason. A lazily computed Go value uses `sync.OnceValue`.

Correct:

```go
// namePattern is the form of a slug.
var namePattern = regexp.MustCompile(`^[a-z][a-z0-9-]{0,62}$`)

var defaultCatalog = sync.OnceValues(func() (Catalog, error) {
	return parseCatalog(catalogJSON) // parsed on first use, not on import
})
```

```typescript
let cachedSettings: Settings | null = null;
let initPromise: Promise<Settings> | null = null;

export const getSettings = async (): Promise<Settings> => {
  if (cachedSettings !== null) return cachedSettings;

  initPromise ??= loadConfig().then((config) => {
    cachedSettings = processConfig(config);
    return cachedSettings;
  });

  return initPromise;
};
```

```typescript
// Acceptable at module scope: pure constants, types, declarations, plain literals
const MAX_RETRIES = 3;
interface Config { apiUrl: string }
const defaultOptions = { timeout: 5000, retries: 3 };
```

Incorrect:

```go
var settings = mustLoad(os.Getenv("SETTINGS")) // environment read on import

func init() {
	go refreshForever() // a goroutine nobody owns starts on import
}
```

```typescript
// Runs on import: can fail, blocks module loading
const config = await loadConfig();
export const settings = processConfig(config);

const apiUrl = process.env.API_URL;                 // environment read at module scope
const raw = readFileSync("./data.json", "utf8");    // I/O at module scope
const processedData = heavyComputation(rawData);    // heavy work on import

let currentUser = null; // global mutable state
export function setCurrentUser(user: User) { currentUser = user; }
```

Refactoring approach: when you meet module-level side effects, wrap them in a
getter function that computes on first call and caches the result.

**Why:** code at package or module scope runs on every import, in import order,
before a test can intercept it; one failure there breaks every importer, a
top-level `await` blocks the module graph, and the work cannot be tree-shaken.
