# Errors

How code in this repository raises, codes, wraps, reports and matches errors,
in Go and TypeScript. The Go idioms (`%w`, sentinels, typed errors,
`errors.AsType`) are in go-practices: Coded Errors; what may reach a response
is in security-practices: Error Sanitization.

## Contents

- Error Handling
- Errors Carry a Code
- Error Wrapping
- Actionable Error Reasons
- Ignored Errors
- Handle Each Error Once
- Result Objects
- Assertions
- Errors the UI Shows

---

## Error Handling

Scope: All code

Rule: Handle every failure that can actually happen: I/O, network, parsing,
input, a dependency returning an error. Do not add handling for states the
types or an invariant already rule out; a check for the impossible hides the
real contract and never runs in a test. Catch narrowly: handle the failures you
expect and let the rest propagate. A normal outcome (a missing optional file,
a lookup that finds nothing) is a return value, not an error.

Correct:

```go
data, err := os.ReadFile(path)
if errors.Is(err, fs.ErrNotExist) {
    return Settings{}, false, nil // no file is a normal answer: defaults apply
}
if err != nil {
    return Settings{}, false, fmt.Errorf("read settings %s: %w", path, err)
}
```

```ts
async function fetchUser(id: string): Promise<User | null> {
  try {
    const response = await fetch(`/api/users/${id}`);
    if (!response.ok) return null;
    return await response.json();
  } catch (error) {
    console.error("Failed to fetch user", error);
    return null;
  }
}
```

Incorrect:

```go
data, err := os.ReadFile(path)
if err != nil {
    return Settings{}, false, nil // an outage, a permission error and a bug all look like "no file"
}
```

```ts
async function fetchUser(id: string): Promise<User> {
  const response = await fetch(`/api/users/${id}`); // a rejected fetch and a non-2xx answer are unhandled
  return await response.json();
}
```

**Why:** a catch-all turns every failure into the same silent answer, and the
caller cannot tell a missing value from a broken system.

---

## Errors Carry a Code

Scope: Every failure a caller, another process, a test or a log reader has to
tell apart

Rule: Give each failure a stable, machine-readable identity before writing its
message: a sentinel (`var ErrExportInProgress = errors.New(...)`) or a typed
error in Go, a `code` field next to the message in the JSON error body where a
client has to branch on it, an `Error` subclass or a code from that body in
TypeScript. Callers branch on the identity (`errors.Is`, `errors.As`, the
`code`) and tests assert it; nobody matches message text, English or
translated. A typed error carries the data its caller needs
(`privacy.CooldownError.NextRequestAt`, `privacy.OwnedProfilesError.Profiles`).
Name the cause precisely: the failing check, the record and the reason, never a
catch-all such as `operation_failed` or one message for several conditions. A
substring check on text is acceptable only to prove that something sensitive is
absent.

Correct:

```go
case errors.Is(err, privacy.ErrExportInProgress):
    return ctx.Results.Error(http.StatusConflict, httpfx.WithJSON(map[string]string{
        "error": "an export is already in progress", "code": "export_in_progress",
    }))
```

```ts
const details = readPrivacyError(error); // reads `code` from the cause of the thrown Error
if (details.code === "export_cooldown" && details.nextRequestAt !== null) { /* ... */ }
```

```ts
class ValidationError extends Error {
  constructor(message: string, public field: string) {
    super(message);
    this.name = "ValidationError";
  }
}
```

Incorrect:

```go
if !strings.Contains(err.Error(), "already in progress") { // breaks on rewording
```

```ts
if (error.message.includes("cooldown")) showCooldown(); // breaks on translation
throw new Error("Something went wrong"); // too generic to tell apart
```

**Why:** text and translations change; codes are the contract. An error that
names its check and record is diagnosed from the log, without reading the
source.

---

## Error Wrapping

Scope: Errors propagated through more than one function

Rule: When propagating, add the step that failed and keep the cause, so the
chain reads from the outermost operation to the root. In Go the repository
joins a sentinel, the context and the cause so both stay matchable:
`fmt.Errorf("%w(userID: %s): %w", ErrFailedToGetRecord, userID, err)`; for a
plain step, `fmt.Errorf("read policy %q: %w", path, err)`. In TypeScript pass
the cause: `new Error("Payment processing failed", { cause: error })`. Never
throw strings, never rebuild an error from its `.message` alone, never rethrow
the original without context, and never lose the identity of the cause when
wrapping. Put values in properties or the cause, not in a concatenated
message.

Correct:

```go
if err := s.repo.Save(ctx, record); err != nil {
    return fmt.Errorf("%w(id: %s): %w", ErrFailedToInsertRecord, record.ID, err) // errors.Is finds both
}
```

```ts
try {
  await processPayment(data);
} catch (error) {
  throw new PaymentError("Failed to process payment", txId, { cause: error });
}
```

Incorrect:

```go
if err := s.repo.Save(ctx, record); err != nil {
    return errors.New(err.Error()) // cause and identity gone
}
```

```ts
throw "Payment failed"; // not an Error object, no stack
throw originalError; // no context added
throw new Error(originalError.message); // loses the stack and the cause
throw new Error("Failed for user " + userId + " amount " + amount); // data glued into text
```

**Why:** the identity lets callers branch without parsing text, and the cause
chain keeps the root failure that a copied message throws away.

---

## Actionable Error Reasons

Scope: The reason of an error, CLI output, and the message of an operator-facing
log entry

Rule: The reason states what went wrong, why, and what would fix it. Include
the value, path or identifier that failed and, for a validation error, what a
valid value looks like. "Something went wrong", "failed" and a bare "invalid"
are defects. Start with the operation (`read policy ...`) where the wrap adds
one. A start-up refusal names the setting to change and never the value
(`ErrJWTSecretMissing` for the JWT secret).

Words for people in the user interface are not reasons: they are chosen from
the code (Errors the UI Shows).

Correct:

```go
return fmt.Errorf("%w: %w", ErrInitFailed, ErrJWTSecretMissing)
return fmt.Errorf("%w: unsupported scheme %q", ErrInvalidURL, parsed.Scheme)
```

Incorrect:

```go
return errors.New("invalid")
```

**Why:** the person reading the log usually has only the entry; one that names
the value and the fix saves a trip into the code.

---

## Ignored Errors

Scope: All languages

Rule: Never ignore a returned error or a rejected promise. An ignored error is
allowed only on a cleanup path where a second failure must not hide the first,
or a best-effort write whose failure changes nothing, and it carries a comment
saying why. A `catch` block is never empty.

Correct:

```go
defer func() { _ = resp.Body.Close() }() // best effort: the real error was already returned
```

```ts
try {
  await saveToDatabase(data);
} catch (error) {
  console.error("Failed to save", error);
  throw new Error("Database operation failed", { cause: error });
}
```

Incorrect:

```go
result, _ := operation() // error discarded
```

```ts
try {
  await saveToDatabase(data);
} catch (error) {
  // ignored
}
save().catch(() => {}); // the person never learns the save failed
```

**Why:** an ignored error surfaces later as corrupted state, far from the call
that could have explained it.

---

## Handle Each Error Once

Scope: All error handling

Rule: An error is either handled or returned, never both. Code that returns an
error adds context and passes it up without logging it. The code that handles
it (an HTTP handler, a worker loop, a cleanup path that decides to continue)
logs it once and turns it into an answer, a retry or a fallback.

Correct:

```go
func (s *Service) load(ctx context.Context, id string) (Story, error) {
    story, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return Story{}, fmt.Errorf("%w(id: %s): %w", ErrFailedToGetRecord, id, err) // returned, not logged
    }
    return story, nil
}

// the handler maps the identity to a status and logs only what it cannot classify
if errors.Is(err, stories.ErrStoryNotFound) {
    return ctx.Results.NotFound(httpfx.WithErrorMessage("not found"))
}
return ctx.Results.Error(http.StatusInternalServerError, httpfx.WithSanitizedError(err)) // logs once, answers generically
```

Incorrect:

```go
if err != nil {
    s.logger.WarnContext(ctx, "store read failed", slog.Any("error", err))
    return Story{}, err // logged here and again by every caller
}
```

**Why:** logging and then returning writes the same failure several times,
each entry with less context than the last.

---

## Result Objects

Scope: TypeScript functions whose failure is an expected outcome

Rule: Where a failure is an expected outcome of the call (parsing input, a
lookup a user can get wrong), consider returning a result object instead of
throwing; the caller then has to look at both cases. Throw for failures the
caller cannot reasonably handle.

Correct:

```ts
type Result<T, E> = { ok: true; value: T } | { ok: false; error: E };

function parseData(input: string): Result<Data, string> {
  if (input === "") return { ok: false, error: "Empty input" };
  try {
    return { ok: true, value: JSON.parse(input) };
  } catch {
    return { ok: false, error: "Invalid JSON" };
  }
}
```

**Why:** a thrown error is invisible in the signature; a result type makes the
failure path part of the contract.

---

## Assertions

Scope: Pre-conditions, post-conditions and invariants

Rule: Use an assertion to document an assumption the code relies on and that
only a programming mistake can break. Input from outside is validated and
refused with an error (inputs-and-config.md: Input Validation), not asserted.
`console.assert` only logs; an invariant that must stop execution throws.

Correct:

```ts
function divide(dividend: number, divisor: number): number {
  if (divisor === 0) throw new Error("divide: divisor must not be zero");
  return dividend / divisor;
}
```

**Why:** a documented invariant fails at the line that broke it, instead of as
a wrong number three calls later.

---

## Errors the UI Shows

Scope: The webclient and any text a person reads about a failure

Rule: A person sees words chosen from the error's code or status, translated
through `t("Section.English text")` with an entry in every locale
(AGENTS.md: Internationalization & Locale Fallback); never the
server's message, an identifier, a stack or an error code. A screen branches on
the code and falls back to a generic translated sentence for a code it does not
know. The words say what happened and what the person can do next.

Correct:

```ts
const details = readPrivacyError(error);

if (details.code === "export_unavailable") {
  toast.error(t("Privacy.Downloading your data is not available on this server."));
} else if (details.code !== "export_in_progress") {
  toast.error(t("Privacy.Could not complete the request. Please try again."));
}
```

Incorrect:

```tsx
<p>{error.message}</p> // the server's English text, untranslated, possibly internal
```

**Why:** the message on the wire is for operators and may change; a person
needs words in their language that tell them what to do, and the code is the
stable key for choosing them.
