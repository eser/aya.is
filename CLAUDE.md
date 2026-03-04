# AYA Development Guidelines

Project-specific configuration and quick reference. For detailed rules, see `.claude/skills/*/SKILL.md`.

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

## Project Structure

```
aya.is/
├── apps/
│   ├── webclient/                # TanStack Start frontend (Deno)
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
└── .claude/skills/               # Development rules
```

## Development Commands

```bash
# Frontend (in /apps/webclient)
deno task dev         # Start dev server
deno task build       # Build for production
deno lint && deno fmt # Lint and format
deno task test        # Run snapshot tests (writes snapshots)
deno task test:update # Regenerate snapshots after intentional changes
deno task test:ci     # Run tests read-only (CI mode, used by make ok)

# Backend (in /apps/services)
make restart          # Restart after changes
make lint             # Run golangci-lint
make test             # Run tests

# Root
make ok               # All quality checks (includes frontend tests)
```

## Key Conventions (CRITICAL)

### Explicit Checks

```typescript
// Correct
if (value === null) {}
if (array.length === 0) {}

// Wrong - never use truthy/falsy except for booleans
if (!value) {}
if (!array.length) {}
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
// Correct - pure logic in separate file, testable with Deno
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
- All `_tx` tables use `CHAR(12)` for `locale_code` — **always `strings.TrimRight(value, " ")` when mapping to Go business types**

#### 3-Tier Locale Fallback Pattern

All `_tx` table joins (story_tx, profile_tx, profile_page_tx, profile_link_tx) MUST use the 3-tier fallback subquery — never restrict with `AND (locale = X OR locale = Y)`:

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

For stories, `<default_locale_reference>` is `(SELECT p_loc.default_locale FROM "profile" p_loc WHERE p_loc.id = s.author_profile_id)`. For profiles/pages/links, it's `p.default_locale` directly.

### React Query with SSR (CRITICAL)

All route data fetching uses **TanStack Query (React Query)** with SSR dehydration/hydration via `@tanstack/react-router-ssr-query`.

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
const permissions = await context.queryClient.ensureQueryData(
  profilePermissionsQueryOptions(locale, slug),
).catch(() => null);

// Live streams — optional, auto-refetch in component
await context.queryClient.prefetchQuery(liveNowQueryOptions(locale));
```

### Profile System

- Users can have only ONE individual profile
- Supports: individual, organization, product profile types
- Slug-based routing: `/{locale}/{profile-slug}`

### Shadcn UI Components

- Location: `/apps/webclient/src/components/ui/`
- Generated code - follows shadcn patterns (props destructuring allowed)
- Don't modify inline Tailwind in these files

### Base UI Select (CRITICAL)

`<SelectValue />` renders the **raw value** by default. Always use a children render function:

```tsx
<SelectValue>
  {(value: string) => labelMap.get(value) ?? value}
</SelectValue>
```

For rich items with descriptions, use `SelectPrimitive.Item` directly — label in `ItemText`, description outside it.

## Remember

- **Run `make ok` before committing**
- **Business logic stays dependency-free** (hexagonal architecture)
- **Sentinel errors in Go** - no `fmt.Errorf("message")` without wrapping
- Use Chrome DevTools to debug and verify
