# Frontend Tests

How the Vitest tests of the webclient (`apps/webclient`) are written and run.
What to test and which tier runs it is in strategy.md.

## Contents

- Runner and Files
- Test Commands
- Snapshot Testing
- Translations Test

---

## Runner and Files

Scope: Tests under `apps/webclient/src`

Rule:

- Tests use Vitest, configured in `apps/webclient/vitest.config.ts`: it
  includes `src/**/*.test.ts`, runs in the `node` environment and resolves the
  `@/` alias through `vite-tsconfig-paths`. The config is separate from
  `vite.config.ts` on purpose, so a test never loads the TanStack Start, nitro
  or devtools plugins. Never add a second test runner or assertion library.
- A test file sits next to the code it tests, as `foo.test.ts`
  (kebab-case, like every webclient file), and imports `describe`, `expect` and
  `test` from `vitest` explicitly.
- The suite has no DOM environment and no component-rendering library: tests
  cover pure functions, so the logic a component needs is extracted into a
  `foo-utils.ts` file with no `import.meta.env`, framework or React context
  (javascript-practices: Pure Function Extraction, AGENTS.md: Pure Function
  Extraction for Testability). Adding `jsdom` or a rendering library is a new
  dependency and needs the owner's approval (workflow-practices: Ask First).
- A test that reads a data file does it with `node:fs` relative to
  `import.meta.url`, as `src/messages/translations.test.ts` does.

Correct:

```typescript
// apps/webclient/src/lib/privacy.test.ts
import { describe, expect, test } from "vitest";
import { isDeletionConfirmed } from "./privacy.ts";

describe("isDeletionConfirmed", () => {
  test("accepts the username regardless of case and surrounding spaces", () => {
    expect(isDeletionConfirmed("  JANE ", "jane")).toBe(true);
  });
});
```

Incorrect:

```typescript
import { siteConfig } from "@/config.ts"; // transitively reads import.meta.env
export function buildUrl(locale: string) {
  return `${siteConfig.host}/${locale}`; // logic that cannot be tested without the config
}
```

**Why:** a node-only suite is fast and has no browser to flake; the price is
that logic must be importable without the framework, which is also what makes
it reusable.

---

## Test Commands

Scope: Running the webclient tests

Rule: Run them from `apps/webclient` (or the repository root for `pnpm test`):

| Command                                  | Does                                                                                   |
| ---------------------------------------- | -------------------------------------------------------------------------------------- |
| `pnpm test`                              | `vitest run`: runs once; writes snapshots that do not exist yet; a mismatch fails      |
| `pnpm test:watch`                        | `vitest`: watch mode while writing                                                     |
| `pnpm test:update`                       | `vitest run -u`: rewrites snapshots; only after an intentional change                  |
| `pnpm test:ci`                           | `CI=true vitest run`: read-only, never writes, a missing snapshot fails; used by `make ok` |
| `pnpm exec vitest run <file> -t <name>`  | One file or one test by name, for the targeted tier (strategy.md: Test Tiers)          |

Correct:

```sh
pnpm exec vitest run src/lib/locale-utils.test.ts
```

Incorrect:

```sh
pnpm test:update   # to turn a red run green, without reading the .snap diff
```

---

## Snapshot Testing

Scope: Vitest test files (`*.test.ts`) that guard structured output

Rule: Use `expect(value).toMatchSnapshot()` for anti-regression tests of
structured output that is tedious to assert by hand: a locale table, a generated
tag, parsed props, a key set. Snapshots are written to
`__snapshots__/<test-file>.snap` next to the test and are committed. The test
name is the snapshot key, so give each test a descriptive, stable name. A
snapshot holds deterministic data only (no timestamps, random ids or absolute
paths).

Update a snapshot only after an intentional change: run `pnpm test:update`,
then read the `.snap` diff and check that it shows exactly the change you meant,
nothing else. Never delete a snapshot file to make a failure disappear. Run
`pnpm test:ci` before handing over: it never writes, so a snapshot you forgot to
commit or update fails there as it would in CI.

Correct:

```typescript
import { expect, test } from "vitest";
import { SUPPORTED_LOCALES, supportedLocales } from "./locale-utils.ts";

// Snapshot the supported locales array: catches accidental additions/removals
test("SUPPORTED_LOCALES - complete list", () => {
  expect(SUPPORTED_LOCALES).toMatchSnapshot();
});

test("supportedLocales - all entries", () => {
  expect(supportedLocales).toMatchSnapshot();
});
```

Incorrect:

```typescript
test("buildReport", () => {
  expect(buildReport(new Date())).toMatchSnapshot(); // the date changes on every run
});
```

**Why:** a snapshot makes an unintended change to a table or key set visible as
a diff in review, and keeping the update step explicit stops a regression from
being accepted by accident.

---

## Translations Test

Scope: `apps/webclient/src/messages/translations.test.ts`

Rule: This test reads the 13 locale files of `src/messages/` and guards the
translation keys (AGENTS.md: Internationalization & Locale Fallback):

- It snapshots the sorted `Section.English text` key set of `en.json`
  (`translations - English key set baseline`). Adding or removing an English key
  changes the snapshot, so run `pnpm test:update` and check that the diff shows
  only your keys.
- It fails when a locale has different sections than English.
- It only warns, with `console.warn`, about keys missing in a locale and about
  orphan keys that exist in a locale but not in English. A warning does not fail
  the run, so read the `Translation gaps found` output after adding keys and
  fill every gap in all 13 locales, in the language of each file.
- i18n lookup keys must not contain a colon, which `src/lib/privacy.test.ts`
  enforces.

Correct:

```
Added t("Privacy", "Delete account") -> pnpm test:update -> the .snap diff adds
"Privacy.Delete account" -> 12 locale files get a translated value -> no
"Translation gaps found" warning
```

Incorrect:

```
Added the key to en.json only, ran pnpm test:update, saw "Translation gaps found"
and committed anyway
```
