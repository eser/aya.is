# Modules and Imports

How the webclient (`apps/webclient`) and the scripts (`scripts/`) lay out,
import and export their modules, and how pure logic is split from
environment-dependent code so Vitest can import it.

## Contents

- Languages and Files
- Imports
- Exports
- Module Paths
- Route Files
- Pure Function Extraction
- Node Scripts
- Formatting

---

## Languages and Files

Scope: Every JavaScript or TypeScript file in the repository

Rule: Each place has one language, and a new file follows it:

| Place | Language | Files |
| --- | --- | --- |
| `apps/webclient/src` | TypeScript, ES modules (`"type": "module"`), `strict` | `.ts`; `.tsx` for modules that render; `*.module.css` beside the component; `*.test.ts` beside the code it tests |
| `apps/webclient` root | TypeScript | `vite.config.ts`, `vitest.config.ts` |
| `apps/webclient/src/routeTree.gen.ts` | Generated | never edited by hand; oxfmt and oxlint ignore `*.gen.ts` |
| `scripts/` | TypeScript, run by Node.js 24 directly | `node scripts/<name>.ts` (see the root `Makefile` `tag` and `release` targets); linted by `pnpm lint:scripts` |

Vitest runs only `src/**/*.test.ts` in a `node` environment
(`vitest.config.ts`), so a `.test.tsx` file would not run: test logic, not
rendered components. `pnpm typecheck` covers `src`, `vite.config.ts` and
`vitest.config.ts`; `scripts/` is linted but not type-checked, so keep script
types simple and run the script after editing it.

**Why:** the tools are configured per pattern; a file outside the pattern
silently skips the check that would have caught its mistake.

---

## Imports

Scope: Imports in the webclient and the scripts

Rule:

- Application code uses named imports (`import { useState } from "react"`)
  and `import type { ... }` for types only. `import * as React from "react"`
  is how shadcn-style components and many existing files read React; both
  forms are accepted.
- A namespace import is the right tool when a module's names are generic or
  would collide with local ones: scripts write `import * as path from "node:path"`
  and use `path.join(...)`, not a bare `join`.
- Inside `src`, import by the alias `@/` (mapped to `src/` by `tsconfig.json`
  `paths`, `vite-tsconfig-paths` and a Vite alias for CSS and assets); use
  `./` only for siblings, and avoid `../` chains that climb out of a feature
  folder.
- Node.js built-ins carry the `node:` prefix (`node:fs/promises`, `node:url`).
- Extensions: `tsconfig.json` allows `.ts` and `.tsx` in specifiers, and
  AGENTS.md examples write them (`@/lib/locale-utils.ts`); most existing app
  code leaves them off (`@/config`). Match the file you edit and never rewrite
  the specifiers of lines a change does not touch. Test files and the modules
  they import by relative path write the extension (`./privacy.ts`).
- Backend calls are imported through the one `backend` object, not from the
  individual request files (AGENTS.md: Backend Object Pattern).

Correct:

```typescript
import * as path from "node:path";
import { readFile } from "node:fs/promises";
import type { PrivacyExportState } from "@/modules/backend/privacy/types";
import { backend } from "@/modules/backend/backend";
import { isDeletionConfirmed } from "@/lib/privacy";
```

Incorrect:

```typescript
import { join, resolve } from "node:path"; // `join` and `resolve` collide with local names in scripts
import { getProfile } from "@/modules/backend/profiles/get-profile"; // skips the backend object
import { x } from "../../../lib/x"; // climbs out of the feature folder
```

**Why:** the alias keeps moves of a file from rewriting its imports, type-only
imports vanish at build time, and one `backend` object is the single place
that lists what the backend can do.

---

## Exports

Scope: Webclient and script modules

Rule: Export on the declaration (`export function`, `export const`,
`export type`). No default exports, except where a tool or the framework reads
the default: `vite.config.ts` and `vitest.config.ts`, a module loaded with
`React.lazy` (`pdf-viewer-inner.tsx`) and the shared i18next instance
(`src/modules/i18n/i18n.ts`). Route modules export `Route` by name. A barrel
`index.ts` may re-export what its folder owns.

Correct:

```typescript
export function buildCommand() {}
export const CONFIG_PATH = "./config";
export type BuildOptions = {};
```

Incorrect:

```typescript
function buildCommand() {}
const CONFIG_PATH = "./config";
export { buildCommand, CONFIG_PATH }; // indirect export, the declaration does not show it is public

export default buildCommand; // default export, every importer picks its own name
```

**Why:** a named export has one name across the code base, so search, rename
and auto-import work; a default export gives each importer a free choice.

---

## Module Paths

Scope: Files, templates and data that a module reads from disk

Rule: Separate the module's own location from the caller's. A path inside the
package comes from `import.meta.dirname` or `new URL(..., import.meta.url)`;
the directory the user works in comes from `process.cwd()` or a parameter. A
function that needs a root takes an optional `projectRoot?: string` and falls
back to the module path, so a test passes a temporary directory instead of
changing global state.

Correct:

```typescript
const rootDir = path.resolve(import.meta.dirname ?? ".", "..");
const filePath = new URL(`./${locale}.json`, import.meta.url);

export async function analyzeComponents(srcDir: string, projectRoot?: string): Promise<Component[]> {
  const root = projectRoot ?? import.meta.dirname;
}
```

Incorrect:

```typescript
const root = process.cwd(); // depends on where the command was started, untestable
const root = dirname(fileURLToPath(import.meta.url)); // `import.meta.dirname` already is this
```

---

## Route Files

Scope: `apps/webclient/src/routes`

Rule: Routes are file-based: the file path is the URL (`routes/$locale/index.tsx`
is `/$locale/`), and each route module exports `Route` made by
`createFileRoute` (the loader, `head` and component patterns are in AGENTS.md:
React Query with SSR). `routeTree.gen.ts` is generated by the
TanStack plugin; never edit it. A file or folder beside the routes that is not
a route starts with `-` (`-home.module.css`, `-markdown.ts`) or lives in an
`_components/` folder (`routeFileIgnorePattern` in `vite.config.ts`). Route
file names follow the router, so `src/routes` is exempt from the kebab-case
rule (tooling-standards: File Naming). The router itself is swapped by Vite
aliases (AGENTS.md: Router (speedy-router)); import from
`@tanstack/react-router` as before.

---

## Pure Function Extraction

Scope: A module that has both testable logic and a dependency on the
environment (`import.meta.env`, a Vite-only API, React context, the site
configuration)

Rule: Keep pure functions free of environment and framework dependencies so
Vitest can import them directly (AGENTS.md: Pure Function Extraction for
Testability). When one file needs both, put the logic in `foo-utils.ts` (pure,
explicit parameters) and keep `foo.ts` as the thin wrapper that re-exports it
and adds the environment-dependent values. Tests import `foo-utils.ts`; the
application imports `foo.ts`. The test for a pure helper sits beside it as
`foo-utils.test.ts` (qa-practices: Snapshot Testing for the snapshot form).

Correct:

```typescript
// locale-utils.ts: pure functions, no import.meta.env
export const SUPPORTED_LOCALES = ["en", "tr", "fr" /* ... all 13 locales */] as const;
export function isValidLocale(locale: string): boolean { /* ... */ }

// config.ts: re-exports the pure part and adds the env-dependent configuration
export { SUPPORTED_LOCALES, isValidLocale } from "@/lib/locale-utils.ts";
export const siteConfig = { host: import.meta.env.VITE_HOST ?? "https://aya.is" /* ... */ };
```

```typescript
// seo-utils.ts: pure, every dependency is a parameter
export function generateMetaTags(meta: SeoMeta, defaults: { host: string; name: string; defaultImage: string }): MetaTag[] {
  /* ... */
}

// seo.ts: thin wrapper that injects siteConfig
import { generateMetaTags as generateMetaTagsPure } from "@/lib/seo-utils.ts";
import { siteConfig } from "@/config";
const SITE_DEFAULTS = { host: siteConfig.host, name: siteConfig.name, defaultImage: `${siteConfig.host}/og-image.png` };
export function generateMetaTags(meta: SeoMeta) {
  return generateMetaTagsPure(meta, SITE_DEFAULTS);
}
```

Incorrect:

```typescript
// Pure logic mixed with the env dependency: a test has to load the site configuration to call it
import { siteConfig } from "@/config";
export function buildUrl(locale: string, ...segments: string[]) {
  return `${siteConfig.host}/${locale}/${segments.join("/")}`;
}
```

Naming: `foo-utils.ts` for the pure functions, `foo.ts` for the config-aware
wrapper that re-exports and delegates (`locale-utils.ts` and `config.ts`,
`seo-utils.ts` and `seo.ts` in `src/lib`). A helper kept free of React says so
in its header comment (`privacy.ts`: "kept free of React so they are
unit-testable").

**Why:** `vitest.config.ts` runs tests in a plain `node` environment and
deliberately loads none of the TanStack Start, nitro or devtools plugins; a
test that has to import the configuration drags the environment into every
case, while a pure file needs no setup.

---

## Node Scripts

Scope: Files under `scripts/`

Rule: A script is a TypeScript file that Node.js runs directly, started
through a `Makefile` target rather than from memory (`make tag`,
`make release`). It starts with the license header like
every other `.ts` file, finds the repository root from its own location
(`import.meta.dirname`), checks its inputs first, and fails with a thrown
`Error` whose message names the command or file that failed and a non-zero exit
status. Process spawning goes through `node:child_process`
(`spawnSync(cmd, args, { cwd, stdio })`) with an argument array, never a
shell string built from input. File names are snake_case
(`push_current_release_tag.ts`): `scripts/` is exempt from the kebab-case
rule in `.manifest.yml`.

Correct:

```typescript
function run(cmd: string, args: string[]): void {
  const result = spawnSync(cmd, args, { cwd: rootDir, stdio: "inherit" });
  if (result.status !== 0) {
    throw new Error(`${cmd} ${args.join(" ")} failed with exit code ${result.status}`);
  }
}
```

Incorrect:

```typescript
execSync(`git tag ${tag}`); // a shell string; the tag text becomes shell syntax
spawnSync("git", ["tag", tag]); // result ignored, a failure continues silently
```

---

## Formatting

Scope: Layout of source text

Rule: oxfmt is the formatter (`.oxfmtrc.json`: `printWidth: 120`, ignoring
generated files, `__snapshots__`, `.agents/skills` and a few directories).
Run `pnpm fmt` from the repository root to write, `pnpm fmt:check` to verify;
the git hook and CI run the check. oxfmt rewrites files in place by default, so
use `--check` whenever the tree must stay untouched. Do not hand-format, and
do not reformat lines a change does not touch. `.editorconfig` supplies
two-space indentation, LF endings and a final newline.

**Why:** a hand-formatted line is changed by the next `pnpm fmt`, and a
reformat buries the real change in the diff.
