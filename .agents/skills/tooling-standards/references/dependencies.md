# Dependencies

Where each dependency is declared, how versions are pinned and what a new
dependency needs before it is added.

## Contents

- Manifests and Lockfiles
- Exact Versions
- Overrides and Patches
- Registry and JSR
- Adding a Dependency

---

## Manifests and Lockfiles

Scope: Every dependency of the repository

Rule: Each toolchain has one manifest and one lockfile:

| Place | Manifest | Lockfile | Install |
| --- | --- | --- | --- |
| Go module `github.com/eser/aya.is/services` | `apps/services/go.mod` | `apps/services/go.sum` | `go get <module>@<version>`, then `go mod tidy` (`make dep` downloads and tidies); Go tools: `go get -tool <module>@<version>` (`make init-tools`, `make init-generators`, `make init-checkers`) |
| Webclient | `apps/webclient/package.json` | `pnpm-lock.yaml` at the root | `pnpm add <package> --filter @aya/webclient`; CI and Docker run `pnpm install --frozen-lockfile` |
| Repository tooling (oxlint, oxfmt) | the root `package.json` | `pnpm-lock.yaml` | `pnpm add -wD <package>@<version>` |
| Workspace settings | `pnpm-workspace.yaml` (`packages`, `overrides`, `patchedDependencies`), `.npmrc` | | edited by hand, with a comment for every entry |

Never edit `go.sum` or `pnpm-lock.yaml` by hand; the tool writes them. The
`eser` CLI is not a project dependency: it is installed globally with Deno
(tooling-standards: Tool per Place).

---

## Exact Versions

Scope: JavaScript dependencies

Rule: Most webclient dependencies use a caret range, and these are pinned to an
exact version, with the reason in a comment or in AGENTS.md:

- tools whose output changes between releases: `oxlint`, `oxlint-tsgolint` and
  `oxfmt` (pre-1.0, so a bump is a deliberate change that is formatted and
  linted on its own);
- packages that must move together: the TanStack router and start family and
  `speedy-router*`, pinned both in `package.json` and in the `overrides` of
  `pnpm-workspace.yaml` (AGENTS.md: Router (speedy-router));
- `react` and `react-dom`, which the router set requires at one version;
- pre-release packages (`nitro` is a dated beta).

A new package is caret-ranged unless it falls into one of these groups.
Upgrades are deliberate changes that update the lockfile, run `make ok` and
say what changed.

Correct:

```json
{"dependencies": {"react": "19.2.8", "speedy-router": "0.1.25", "zod": "^4.4.3"}}
```

Incorrect:

```json
{"dependencies": {"speedy-router": "^0.1.25"}} // a 0.x router moves without review
```

**Why:** the router set was verified as a whole; one package that moves alone
crashes server rendering, as the notes in `pnpm-workspace.yaml` record.

---

## Overrides and Patches

Scope: `overrides` and `patchedDependencies` in `pnpm-workspace.yaml`, and
`patches/`

Rule: An override or a patch needs a comment that says what breaks without
it, which upstream version fixes it and when to drop it. A patch file is named
`@scope__package@version.patch` in `patches/` and is registered under
`patchedDependencies` with the exact version it applies to. When the pinned
version moves, the patch is re-checked or removed in the same change. Do not
patch a package to avoid asking the owner about a dependency.

Correct:

```yaml
# router-ssr-query-core 1.169.1 calls hydrate(queryClient, undefined) when the query stream ends ... Drop this patch
# together with the pin once both move up.
patchedDependencies:
  "@tanstack/router-ssr-query-core@1.169.1": patches/@tanstack__router-ssr-query-core@1.169.1.patch
```

Incorrect:

```yaml
patchedDependencies:
  "some-package@1.2.3": patches/some-package.patch # no reason, no exit condition
```

**Why:** an uncommented patch becomes permanent, and the next upgrade fails
for a reason nobody remembers.

---

## Registry and JSR

Scope: Where a JavaScript package comes from

Rule: `.npmrc` maps the `@jsr` scope to `https://npm.jsr.io`, so a package
published on jsr.io can be a normal dependency of `package.json`
(`"@std/path": "npm:@jsr/std__path@^1.0.0"`) and is installed by pnpm, not by
Deno. Prefer the JSR package when the same package exists on both registries
and the JSR one is maintained; use npm when there is no JSR alternative. No
JSR package is a dependency today.

Incorrect:

```json
{"dependencies": {"@std/path": "^1.0.0"}}
```

---

## Adding a Dependency

Scope: Every new Go module or npm package, direct or a new indirect one a
change pulls in

Rule: A new dependency is untrusted code and a long-term liability, and
workflow-practices: Ask First requires asking the owner. Before proposing one,
write down:

- Purpose, and why the standard library, `golang.org/x/*`, the browser or Node
  APIs, a package already in `go.mod` or `package.json`, or about 50 lines of
  our own code cannot do it.
- License: permissive (MIT, BSD, Apache-2.0, ISC) for anything that ships; the
  repository is Apache-2.0. GPL, LGPL, AGPL or an unclear license needs the
  owner's explicit approval.
- Advisories: no open advisory in the Go vulnerability database
  (`go tool govulncheck ./...` in `apps/services`) or in the npm advisory data
  (`pnpm audit`).
- Health: an identifiable maintainer, recent releases, and an entry point and
  dependency tree that contain nothing unexpected.
- Weight: for a webclient package, what it adds to the browser bundle and
  whether it runs in the SSR server too.
- Exit cost: what breaks if the package disappears, and how much code touches
  it.

"It looked fine" is not evidence. Prefer established, widely used packages
over new ones.
