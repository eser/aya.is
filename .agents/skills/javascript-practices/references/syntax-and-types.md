# Syntax, Values, Errors and Async Work

How webclient and script code checks values, reads data, names its numbers,
reports errors and runs asynchronous work. oxlint (`apps/webclient/.oxlintrc.json`:
correctness errors plus the `typescript`, `react` and `oxc` plugins) and
`tsc --noEmit` with `strict` enforce the basics; run `pnpm lint` and
`pnpm typecheck` rather than restating their rules. oxfmt owns layout
(double quotes, semicolons, trailing commas, 120 columns): run `pnpm fmt`,
never format by hand.

## Contents

- Variables, Equality and Strings
- Explicit Checks
- Defaults and Absent Values
- Null and Undefined
- Iteration Over Optional Values
- Type Checks
- Dangerous Features
- Named Bounds
- Comments Describe the Why
- Copy, Do Not Mutate
- Errors Are Matched by Code
- Async Handlers and Floating Promises
- Browser Storage
- String Prefix Extraction
- Identifiers and Untrusted Strings

---

## Variables, Equality and Strings

Scope: Every `.ts` and `.tsx` file

Rule:

- `const` by default, `let` only for a binding that is reassigned (oxlint
  `prefer-const` and `no-var` are errors).
- Strict equality `===` and `!==` only; `==` coerces (`0 == ""` is true).
- Semicolons end every statement; oxfmt writes them.
- Template strings, not concatenation, for built text.
- `slice()`, never `substr()` (deprecated) or `substring()` (swaps arguments).
- Rest parameters (`...numbers`), never the `arguments` object.
- `for..of` for array values (`for (const [i, item] of items.entries())` when
  the index matters); a classic `for` only when `break`, `continue` or index
  arithmetic needs it.
- `delete` removes an object property only; remove array elements with
  `splice()` or `filter()`.
- `Number(value)` converts, never the unary `+`.

Correct:

```typescript
const url = `/api/users/${userId}/posts/${postId}`;
const port = Number(process.env.PORT ?? "8000");
const last = text.slice(-5);

function sum(...numbers: number[]): number {
  return numbers.reduce((total, n) => total + n, 0);
}
```

Incorrect:

```typescript
let userName = "John"; // never reassigned
if (value == 0) {} // matches 0, "0", false and ""
const port = +userInput; // intent hidden in a symbol
const url = "/api/users/" + userId + "/posts/" + postId;
```

**Why:** each form removes one source of silent coercion or one reader
question; none of them costs anything.

---

## Explicit Checks

Scope: Conditions on strings, numbers, arrays, objects and nullable values

Rule: Truthiness is only for values typed `boolean`. Compare everything else
explicitly (AGENTS.md: Explicit Checks).

Correct:

```typescript
if (array.length === 0) {}
if (user !== null) {}
if (value !== 0) {}
if (isActive) {} // boolean
```

Incorrect:

```typescript
if (!array.length) {} // reads as "no array"
if (user) {} // false for null, undefined, 0 and ""
```

**Why:** `0`, `""` and `NaN` are valid values for most fields; a truthiness
check turns them into "missing" without anyone deciding that.

---

## Defaults and Absent Values

Scope: Fallbacks for values that may be missing

Rule: Use `??` when only a missing value (`null` or `undefined`) falls back,
which is the case for nearly every field of an API answer. Use `||` only when
an empty string should fall back too and the reason is visible. Read optional
data with optional chaining and a fallback of the right shape.

Correct:

```typescript
const port = config.port ?? 8000;
const label = user.name ?? user.email ?? t("Layout.Profile");
```

Incorrect:

```typescript
const port = config.port || 8000; // port 0 falls back
const value = userInput || "default"; // "" falls back
```

---

## Null and Undefined

Scope: Return types, state and fields the code defines itself

Rule: Use `null` for "no value" in the types and fields this code owns, so
absence is a decision a reader can see (`User | null`, `profile: Profile | null`).
`undefined` stays where the language produces it (optional properties and
parameters, `Map.get`, array lookups); convert at the boundary with `?? null`.
Both are compared explicitly (`=== null`, `=== undefined`).

Correct:

```typescript
function findUser(id: string): User | null {
  return users.get(id) ?? null;
}
```

---

## Iteration Over Optional Values

Scope: Loops over a value that may be missing

Rule: Check for the absent case first, then iterate. Do not write
`?? {}` or `?? []` only to loop over nothing.

Correct:

```typescript
if (props !== undefined) {
  for (const [key, value] of Object.entries(props)) {
    process(key, value);
  }
}
```

Incorrect:

```typescript
for (const [key, value] of Object.entries(props ?? {})) {
  process(key, value);
}
```

**Why:** no throw-away object is allocated, the conditional nature of the loop
is visible, and a log line or breakpoint has an obvious place to go.

---

## Type Checks

Scope: Narrowing `unknown` values and feature detection

Rule: `typeof` is the tool for primitives and for "is this API present"
(`typeof value === "string"`, `typeof globalThis.localStorage === "undefined"`);
use it to write type guards for data that arrives as `unknown`. Use `instanceof`
for classes (`error instanceof Error`) and `Array.isArray` for arrays. Do not
use `typeof value === "object"` alone, which is true for `null` and arrays.

Correct:

```typescript
function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null;
}

if (Array.isArray(profiles)) {
  details.profiles = profiles.filter(isOwnedProfile);
}
```

Incorrect:

```typescript
if (typeof value === "object") {} // true for null and for arrays
```

---

## Dangerous Features

Scope: Every `.ts` and `.tsx` file

Rule: No `eval` or `new Function`, no prototype mutation
(`Array.prototype.x = ...`), no `Object.defineProperty` on shared objects, no
`@ts-ignore`, `@ts-expect-error` or `oxlint-disable` comment, no `any`. Fix the
type or the code; `ban-ts-comment` and `no-explicit-any` are errors, and the
tree has no suppression today.

---

## Named Bounds

Scope: Limits, delays, sizes and counts

Rule: Every number with a meaning is a named constant with a comment that says
where it comes from or what unit it has. A bound the backend enforces mirrors
the backend's value, so a form refuses what the backend would refuse.

Correct:

```typescript
const TOKEN_REFRESH_THRESHOLD_MS = 5 * 60 * 1000; // 5 minutes before expiry
const HTTP_STATUS_UNAUTHORIZED = 401;
const MAX_RECENT = 5;
```

Incorrect:

```typescript
if (Date.now() > expiresAt - 300000) refresh(); // what is 300000?
if (response.status === 401) logout();
```

---

## Comments Describe the Why

Scope: Exported functions, components and module-level constants

Rule: A comment above an exported declaration says what it is for and which
rule or decision shapes it, in full sentences. A comment inside a function
explains a decision the code cannot show (why a step waits, why a value is
kept), never what the next line does. A pure helper that is kept free of React
or `import.meta.env` says so in its header.

Correct:

```typescript
/**
 * Deleting an account requires typing the username. The comparison ignores case and surrounding whitespace, the
 * same way the backend does, so a user is never blocked by the UI for something the server would accept.
 */
export function isDeletionConfirmed(typed: string, slug: string): boolean {
```

Incorrect:

```typescript
// check if the strings are equal
return typed.trim().toLowerCase() === slug.toLowerCase();
```

---

## Copy, Do Not Mutate

Scope: State, props, query data and API answers

Rule: Build a new object or array instead of changing one in place: spread,
`map`, `filter`, `toSorted`, `structuredClone` for a deep copy. Query data from
TanStack Query and props are never mutated (oxlint `react/immutability` warns).

Correct:

```typescript
setRows(rows.map((row) => (row.id === id ? { ...row, ...change } : row)));
const sorted = items.toSorted((a, b) => a.name.localeCompare(b.name));
```

Incorrect:

```typescript
Object.assign(rows.find((row) => row.id === id), change);
setRows(rows); // same identity, nothing re-renders
items.sort((a, b) => a.name.localeCompare(b.name)); // sorts the prop in place
```

**Why:** React compares by identity; a mutated object renders nothing new, and
a mutated query cache shows the change on every screen that reads it before the
backend accepted it.

---

## Errors Are Matched by Code

Scope: Failed backend requests

Rule: The fetcher (`src/modules/backend/fetcher.ts`) throws an `Error` whose
`cause.response` holds the backend's JSON answer. Decide what to show from the
machine-readable `code` (and the HTTP status), read through a small guard that
accepts only known codes; show the sentence from a translation, never the
backend's `message` text. A module with its own codes keeps the reader beside
its types (`readPrivacyError` in `src/modules/backend/privacy/errors.ts`).

Correct:

```typescript
const details = readPrivacyError(error);
if (details.code === "export_cooldown" && details.nextRequestAt !== null) {
  toast.error(t("Privacy.You can request a new copy after {{date}}.", { date: formatDate(details.nextRequestAt) }));
}
```

Incorrect:

```typescript
toast.error(error.message); // backend text, one language
if (error.message.includes("cooldown")) {} // matched by wording
```

**Why:** codes are the contract with the backend; message text changes without
notice and has no translation.

---

## Async Handlers and Floating Promises

Scope: Event handlers, effects, loaders and requests

Rule:

- Use `return await` in async functions, always, and especially inside `try`
  so a rejection reaches the `catch`.
- A handler that sends a request is `async` and awaits it; it closes the form
  or dialog only after success.
- A promise that is deliberately not awaited is marked `void`
  (`void queryClient.invalidateQueries(...)`), so a reader sees it was not
  forgotten. oxlint warns on `no-floating-promises` and `no-misused-promises`,
  and errors on `require-await`.
- Awaiting inside a loop needs a sequential reason, stated in a comment.

Correct:

```typescript
async function fetchUser(id: string): Promise<User> {
  return await userRepository.findById(id);
}

async function processData(): Promise<Result> {
  try {
    return await riskyOperation();
  } catch (error) {
    return await fallbackOperation();
  }
}
```

Incorrect:

```typescript
async function processData(): Promise<Result> {
  try {
    return riskyOperation(); // BUG: a rejection bypasses the catch block
  } catch (error) {
    return fallbackOperation();
  }
}
```

**Why:** `return await` keeps the function in the stack trace, makes `try`/`catch`
work, and costs no performance in current engines; the old `no-return-await`
lint rule is deprecated.

---

## Browser Storage

Scope: `localStorage`, `sessionStorage`, `document` and `window`

Rule: The same module runs in the browser and in the SSR server, and Node.js
25 and later expose a `globalThis.localStorage` on the server, so testing for
`localStorage` alone proves nothing: guard with the browser check
(`isBrowserEnvironment()` in `fetcher.ts` tests `globalThis.document`). Reach
browser APIs through `globalThis`, never `window`. Keys are named constants. A
read of stored JSON sits in `try`/`catch` with a fallback, because a private
window, blocked storage or a corrupt value throws.

Correct:

```typescript
const STORAGE_KEY = "aya-recent-blocks";

export function getRecentBlockIds(): string[] {
  if (typeof globalThis.localStorage === "undefined") {
    return [];
  }
  const stored = globalThis.localStorage.getItem(STORAGE_KEY);
  if (stored === null) {
    return [];
  }
  try {
    const parsed = JSON.parse(stored);
    return Array.isArray(parsed) ? parsed.filter((id): id is string => typeof id === "string") : [];
  } catch {
    return [];
  }
}
```

Incorrect:

```typescript
const isDefined = typeof window.localStorage !== "undefined"; // `window` does not exist on the server
const recent = JSON.parse(localStorage.getItem("recent")!); // unguarded, unvalidated, magic key
```

---

## String Prefix Extraction

Scope: Values after a known prefix (HTTP headers, URIs, file paths, protocol
strings)

Rule: Validate with `startsWith()`, extract with `slice(prefix.length)`. Never
`replace(prefix, "")`, which also edits the middle of the string, and never
`slice(7)` with a hard-coded length, which does not check that the prefix is
there and breaks when the prefix changes.

Correct:

```typescript
function extractAfterPrefix(value: string | null | undefined, prefix: string): string | null {
  if (value === null || value === undefined || !value.startsWith(prefix)) {
    return null;
  }

  const extracted = value.slice(prefix.length);
  return extracted.length > 0 ? extracted : null;
}

const BEARER_PREFIX = "Bearer ";
const token = extractAfterPrefix(authHeader, BEARER_PREFIX);
```

Incorrect:

```typescript
const token = authHeader?.replace("Bearer ", "") ?? null; // "xBearer abc" becomes "xabc"
const path = uri?.replace("file://", "") ?? null; // "s3://file://bucket" becomes "s3://bucket"
const token = authHeader?.slice(7) ?? null; // never checks that "Bearer " is there
```

**Why:** `replace` rewrites the first match anywhere; `startsWith` plus `slice`
states the expected format, rejects what does not have it and handles the empty
remainder on purpose.

---

## Identifiers and Untrusted Strings

Scope: Generated keys and values from the address bar, the backend or users

Rule: Generated identifiers come from `crypto.randomUUID()`, never
`Math.random` or an array index used as a stable key. A value read from the
query string or a form is checked against the allowed values or a bound before
use; an unknown value falls back to the default. Render untrusted text as
text: React escapes it. `dangerouslySetInnerHTML` only receives HTML that the
code produced from escaped input (`renderInlineMarkdown` in
`src/lib/inline-markdown.tsx` escapes first, then adds its own tags);
validation and sanitizing rules are in security-practices.
