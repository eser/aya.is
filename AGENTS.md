# AYA Development Guidelines

Project-specific configuration and quick reference. For detailed rules, see `.agents/skills/*/SKILL.md`.

## Quick Reference

| Topic                | Skill                     |
| -------------------- | ------------------------- |
| Project structure    | `architecture-guidelines` |
| Code quality         | `coding-practices`        |
| JS/TS patterns       | `javascript-practices`    |
| Go patterns          | `go-practices`            |
| Security             | `security-practices`      |
| Development workflow | `workflow-practices`      |
| CI/CD                | `ci-cd-practices`         |
| Tooling              | `tooling-standards`       |

## gstack

Use the `/browse` skill from gstack for all web browsing. **Never use `mcp__claude-in-chrome__*` tools.**

Available skills:

- `/browse` — Web browsing and page interaction (Playwright-based)
- `/plan-ceo-review` — CEO-level plan review
- `/plan-eng-review` — Engineering plan review
- `/review` — Code review
- `/ship` — Ship code (commit, PR, deploy workflow)
- `/qa` — Quality assurance testing
- `/setup-browser-cookies` — Import browser cookies for authenticated browsing
- `/retro` — Retrospective analysis

## Project Structure

```
aya.is/
├── apps/
│   ├── webclient/                # TanStack Start frontend (Node.js + pnpm)
│   │   └── src/
│   │       ├── routes/           # File-based routing
│   │       ├── components/       # UI components (shadcn in ui/)
│   │       ├── modules/          # Feature modules (backend, auth, i18n)
│   │       └── lib/              # Utilities
│   └── services/                 # Go backend (Hexagonal Architecture)
│       └── pkg/
│           ├── api/business/     # Pure business logic
│           ├── api/adapters/     # External implementations
│           └── ajan/             # Framework utilities
└── .agents/skills/               # Development rules
```

## Development Commands

```bash
# Frontend (in /apps/webclient)
pnpm dev               # Start dev server
pnpm build             # Build for production
pnpm lint              # Lint (oxlint, type-aware)
pnpm typecheck         # Type check (TypeScript 7 native tsc)
pnpm test              # Run tests (Vitest; writes new snapshots)
pnpm test:update       # Regenerate snapshots after intentional changes
pnpm test:ci           # Run tests read-only (CI mode, used by make ok)

# Backend (in /apps/services)
make restart          # Restart after changes
make lint             # Run golangci-lint
make test             # Run tests

# Root
pnpm fmt               # Format everything (oxfmt, from the repo root)
pnpm fmt:check         # Check formatting (used by make ok)
make ok               # All quality checks (includes frontend tests)
```

### Tooling override (this project)

This project runs on **Node.js (24+) and pnpm** (workspace: root + `apps/webclient`). It uses **oxlint**, **oxfmt**, the
native **TypeScript 7** `tsc` and **Vitest** instead of Deno's built-in tooling. Deno is only used to install and run the
`eser` CLI (git hook runner, published on JSR), nothing else. Where a skill under `.agents/skills` says to use Deno tooling
(`deno install/fmt/lint/test/task`), this section takes precedence.

- Lint: `pnpm lint` in `apps/webclient` (oxlint, type-aware; config in `apps/webclient/.oxlintrc.json`)
- Type check: `pnpm typecheck` in `apps/webclient` (TypeScript 7 `tsc --noEmit`, test files included)
- Tests: `pnpm test` in `apps/webclient` (Vitest; snapshots live in `__snapshots__/*.snap`; `pnpm test:ci` never writes)
- Format: `pnpm fmt` / `pnpm fmt:check` from the repo root (oxfmt, `printWidth: 120`; config in `.oxfmtrc.json`)
- oxfmt writes files in place by default; always use `--check` (or `fmt:check`) for verification
- `oxlint` and `oxfmt` are pinned to exact versions on purpose (oxfmt is pre-1.0), bump them deliberately

### Agent tooling layout

`.agents/` is the only real agent directory; everything tool-specific lives inside it. `.claude`, `.codex` and
`.mcp.json` are symlinks that exist only because Claude Code and Codex read those fixed paths. Never put real files
behind a symlink, edit them in `.agents/`.

- `AGENTS.md` — project instructions. A nested `AGENTS.md` scopes guidance to a directory, the closest one wins (see
  `apps/services/AGENTS.md`).
- `.agents/skills/<name>/SKILL.md` — skills in the open Agent Skills format.
- `.agents/hooks/` — hook scripts shared by every tool.
- `.agents/settings.json` — Claude Code permissions, env and hook registration (via `.claude`).
- `.agents/config.toml`, `.agents/hooks.json` — Codex MCP servers, env and hook registration (via `.codex`).
- `.agents/mcp.json` — MCP servers for Claude Code (via `.mcp.json`).
- `.agents/settings.local.json` and `.agents/worktrees/` are local state and git-ignored.

There is no `.cursor/` directory: Cursor reads `AGENTS.md` and `.agents/skills` natively.

## Key Conventions (CRITICAL)

### Explicit Checks

```typescript
// Correct
if (value === null) {
}
if (array.length === 0) {
}

// Wrong - never use truthy/falsy except for booleans
if (!value) {
}
if (!array.length) {
}
```

### Backend Object Pattern

```typescript
// Correct - use centralized backend object
import { backend } from "@/modules/backend/backend.ts";
const profile = await backend.getProfile("en", id);

// Wrong - direct imports
import { getProfile } from "@/modules/backend/profiles/get-profile.ts";
```

### CSS Modules with @apply

```css
/* product-card.module.css */
.card {
  @apply border rounded-lg p-4;
}
```

```tsx
import styles from "./product-card.module.css";
<div className={styles.card}>...</div>;
```

### Single Props Object

```tsx
// Correct
function Component(props: Props) {
  return <div>{props.title}</div>;
}

// Wrong - destructured
function Component({ title }: Props) {}
```

### Pure Function Extraction for Testability

```typescript
// Correct - pure logic in separate file, testable with Vitest
// locale-utils.ts (pure, no import.meta.env)
export function isValidLocale(locale: string) { ... }

// config.ts (re-exports + env-dependent config)
export { isValidLocale } from "@/lib/locale-utils.ts";
export const siteConfig = { host: import.meta.env.VITE_SITE_URL };

// Wrong - pure logic mixed with env dependency
import { siteConfig } from "@/config.ts";
export function buildUrl(locale: string) { return `${siteConfig.host}/${locale}`; }
```

## Project-Specific Notes

### Internationalization & Locale Fallback (CRITICAL)

- 13 locales: ar, de, en, es, fr, it, ja, ko, nl, pt-PT, ru, tr, zh-CN
- Translation keys use English text: `t("Section", "English text")`
- Messages in `/apps/webclient/src/messages/[locale].json`
- **NEVER put English text in non-English locale files** — every value in a locale file MUST be translated to that
  locale's language. When adding new i18n keys, always provide proper translations for ALL 13 locales.
- All `_tx` tables use `CHAR(12)` for `locale_code` — **always `strings.TrimRight(value, " ")` when mapping to Go
  business types**

#### 3-Tier Locale Fallback Pattern

All `_tx` table joins (story_tx, profile_tx, profile_page_tx, profile_link_tx) MUST use the 3-tier fallback subquery —
never restrict with `AND (locale = X OR locale = Y)`:

```sql
-- CORRECT: 3-tier fallback — always finds a translation if one exists
AND st.locale_code = (
  SELECT stx.locale_code FROM "story_tx" stx
  WHERE stx.story_id = s.id
  ORDER BY CASE
    WHEN stx.locale_code = sqlc.arg(locale_code) THEN 0          -- requested locale
    WHEN stx.locale_code = <default_locale_reference> THEN 1     -- entity's default
    ELSE 2                                                        -- any available
  END
  LIMIT 1
)

-- WRONG: 2-tier — drops entities when default_locale doesn't match any translation
AND (stx.locale_code = sqlc.arg(locale_code) OR stx.locale_code = <default>)
ORDER BY CASE WHEN stx.locale_code = sqlc.arg(locale_code) THEN 0 ELSE 1 END
```

For stories, `<default_locale_reference>` is
`(SELECT p_loc.default_locale FROM "profile" p_loc WHERE p_loc.id = s.author_profile_id)`. For profiles/pages/links,
it's `p.default_locale` directly.

### React Query with SSR (CRITICAL)

All route data fetching uses **TanStack Query (React Query)** with SSR dehydration/hydration via
`@tanstack/react-router-ssr-query`.

**Infrastructure:**

- `QueryClient` created in `src/router.tsx` with `staleTime: 60_000` and `retry: 2`
- `setupRouterSsrQueryIntegration({ router, queryClient })` handles dehydrate/hydrate/provider automatically
- Query option factories centralized in `src/modules/backend/queries.ts`
- Error handling via `QueryError` component (`src/components/query-error.tsx`)

**Route patterns:**

```typescript
// Loader: prefetch into cache for SSR
loader: async ({ params, context }) => {
  const { locale } = params;
  await context.queryClient.ensureQueryData(storiesByKindsQueryOptions(locale, ["article"]));
  return { locale };
},
errorComponent: QueryError,

// Component: read from hydrated cache, auto-refetch when stale
function Page() {
  const { locale } = Route.useLoaderData();
  const { data } = useSuspenseQuery(storiesByKindsQueryOptions(locale, ["article"]));
}
```

**Query key convention:** `[entity, locale, ...identifiers, { filters }]` — designed for future pagination.

```typescript
// Correct - use query options from queries.ts
import { storiesByKindsQueryOptions } from "@/modules/backend/queries";
await context.queryClient.ensureQueryData(storiesByKindsQueryOptions(locale, ["article"]));

// Wrong - direct backend calls in route loaders
const articles = await backend.getStoriesByKinds(locale, ["article"]);
```

**When head() needs backend data** (Pattern B), `ensureQueryData` returns the data AND caches it:

```typescript
const profile = await context.queryClient.ensureQueryData(profileQueryOptions(locale, slug));
// profile is available for head() via loaderData AND cached for components
```

**Optional data** (e.g., permissions, live streams) uses `.catch(() => null)` or `prefetchQuery`:

```typescript
// Permissions — optional, don't crash on failure
const permissions = await context.queryClient
  .ensureQueryData(profilePermissionsQueryOptions(locale, slug))
  .catch(() => null);

// Live streams — optional, auto-refetch in component
await context.queryClient.prefetchQuery(liveNowQueryOptions(locale));
```

### Router (speedy-router)

The webclient runs on [`speedy-router`](https://github.com/anonrig/router) (a TanStack Router API clone, pre-1.0, pinned to
exact versions) while source code still imports `@tanstack/react-router`. The swap happens only in `apps/webclient/vite.config.ts`
(`resolve.alias` for `@tanstack/react-router`, `@tanstack/router-core`, `@tanstack/history`, plus `ssr.noExternal` for
`@tanstack/*`, so TanStack Start does not load its own copy of the real router on the server).

- TanStack Start, `@tanstack/react-router-ssr-query` and the route-tree generator are still TanStack; types also still come
  from `@tanstack/react-router`.
- To switch back to TanStack Router: delete those three alias lines and the `ssr` block from `vite.config.ts`, remove the
  `speedy-router*` dependencies from `apps/webclient/package.json`, run `pnpm install`.
- Bugs found in speedy-router are worth reporting or fixing upstream (https://github.com/anonrig/router/issues).

### Profile System

- Users can have only ONE individual profile
- Supports: individual, organization, product profile types
- Slug-based routing: `/{locale}/{profile-slug}`

### Privacy (export, hide, delete)

Self-service tab at `/{locale}/{slug}/settings/privacy` (own individual profile only). Backend: `pkg/api/business/privacy`,
routes in `adapters/http/routes_privacy.go`, redaction/export SQL in `adapters/storage/repository_privacy*.go` and
`privacy_datasets.go`.

- **Hide** reuses `profile.deleted_at` plus `hidden_at`; a login hook restores only rows with `hidden_at` set (admin
  deletions stay deleted). Hiding ends every session.
- **Delete** runs ordered `redactionSteps` in one transaction: secrets are hard-deleted, personal text becomes
  `**redacted**`, kept content moves to the `deleted-user` placeholder (`TombstoneUserID`), the slug becomes
  `deleted-<id>` so the username is free again. Sole owners of organizations/products are blocked until they transfer
  ownership. Steps and datasets are skipped when a table does not exist (production still has `legacy_question*`).
- **New table or personal column? Update both** `privacy_datasets.go` (export) and `redactionSteps` (delete), otherwise
  the data is missing from the archive or survives deletion. Secrets belong in `secretColumns`.
- **Export** archives are zipped JSON, stored in a **private** bucket (`S3_EXPORT__BUCKET_NAME`, plus
  `S3_EXPORT__ACCESS_KEY_ID` / `S3_EXPORT__SECRET_ACCESS_KEY`; endpoint and region fall back to the `S3__*` values). The
  bucket must have no public domain; downloads use 10-minute presigned URLs. Kept 7 days, one request per 30 days, user
  is e-mailed through Resend. Without the bucket the feature reports `export_unavailable`.
- i18n lookup keys must **not contain a colon** (i18next treats it as a namespace separator); `lib/privacy.test.ts`
  enforces this.
- Real-DB tests are gated: `PRIVACY_TEST_DSN` (storage) and `S3_TEST_ENDPOINT/BUCKET/KEY/SECRET` (use RustFS, not moto:
  moto does not verify signatures).

### Accessibility

- Use `data-slot="card"` (NOT `role="card"`) as CSS hook for card links — `role="card"` is not a valid ARIA role
- Icon-only `<Button>` elements MUST have `aria-label`
- Custom dropdown menus need `role="menu"` on container, `role="menuitem"` on items

### DropdownMenu Navigation (Base UI)

For menu items that navigate, use `render` prop with `LocaleLink` (renders a real `<a>` tag):

```tsx
// Correct — accessible link menu item
<DropdownMenuItem render={<LocaleLink to="/path" className="no-underline" />}>
  <Icon size="14" />
  {t("Label")}
</DropdownMenuItem>

// For action items (logout, dialogs) — use onClick
<DropdownMenuItem onClick={handleAction}>
  {t("Label")}
</DropdownMenuItem>
```

### Shadcn UI Components

- Location: `/apps/webclient/src/components/ui/`
- Generated code - follows shadcn patterns (props destructuring allowed)
- Don't modify inline Tailwind in these files

### Base UI Select (CRITICAL)

`<SelectValue />` renders the **raw value** by default. Always use a children render function:

```tsx
<SelectValue>{(value: string) => labelMap.get(value) ?? value}</SelectValue>
```

For rich items with descriptions, use `SelectPrimitive.Item` directly — label in `ItemText`, description outside it.

## Remember

- **Run `make ok` before committing**
- **Business logic stays dependency-free** (hexagonal architecture)
- **Sentinel errors in Go** - no `fmt.Errorf("message")` without wrapping
- Use Chrome DevTools to debug and verify
- **NEVER claim "not deployed" without evidence** — always verify the live site first (curl, browser). When the user
  reports something broken, investigate the actual cause instead of assuming deployment state
