# AYA Web Client

The frontend application for the [AYA platform](https://aya.is), built with
TanStack Start on the Deno runtime.

## Tech Stack

- **Deno** - Runtime (uses `nodeModulesDir: "auto"` for npm compatibility)
- **TanStack Start** - Full-stack React framework with file-based routing
- **Vite** - Bundler (with Nitro as the server engine)
- **React** - UI library
- **CSS Modules** - Component-scoped styling with Tailwind via `@apply`
- **shadcn/ui** - Component library
- **i18next** - Internationalization (13 locales)

## Getting Started

```bash
# Install dependencies
deno install

# Start development server (port 3000)
deno task dev

# Build for production
deno task build

# Preview production build
deno task preview

# Lint and format
deno lint && deno fmt
```

## Project Structure

```
src/
  components/             # React components
    ui/                   # shadcn/ui primitives (generated, don't modify)
    userland/             # MDX content components (cards, embeds, profile-card, etc.)
    page-layouts/         # Page layout wrappers (default layout with header/footer)
    forms/                # Form components (create/edit profile)
    elements/             # Reusable display components (filter bar, profile list)
    widgets/              # Composite widgets (astronaut, story-content, text-content)
    profile/              # Profile-specific components (picture upload)
    content-editor/       # Rich content editor
    cover-generator/      # Cover image generator
    icons.tsx             # Icon components
    profile-sidebar-layout.tsx  # Profile page sidebar layout
  hooks/                  # Custom React hooks (mobile detection, session prefs, etc.)
  lib/                    # Core utilities
    auth/                 # Auth context and helpers
    schemas/              # Validation schemas (profile, settings)
    hooks/                # Shared hook utilities
    cover-generator/      # Cover generation logic
    mdx.tsx               # MDX compilation and component mapping
    seo.ts                # SEO/meta tag utilities
    date.ts               # Date formatting
    url.ts                # URL utilities
  messages/               # i18n JSON files (ar, de, en, es, fr, it, ja, ko, nl, pt-PT, ru, tr, zh-CN)
  modules/                # Feature modules
    backend/              # API client (centralized via backend.ts)
    i18n/                 # i18n configuration
    navigation/           # Navigation context
  routes/                 # File-based routing (TanStack Router)
    $locale/              # Locale-prefixed routes
      $slug/              # Profile routes (by slug)
        qa/               # Q&A section
        settings/         # Profile settings
        stories/          # Profile stories
  server/                 # Server-side middleware and handlers
  workers/                # Web Workers (proof-of-work solver)
  config.ts               # App configuration
  env.ts                  # Environment variable definitions
  router.tsx              # Router setup
  styles.css              # Global styles
  theme.css               # Tailwind theme tokens
```

## Key Conventions

### CSS Modules with Tailwind

Every `.module.css` file **must** start with `@reference "@/theme.css";` as its
first line. Without it, `@apply` with Tailwind utility classes will crash the
production build.

```css
@reference "@/theme.css";

.card {
  @apply border rounded-lg p-4;
}
```

### Component Props

Use a single `props` object without destructuring:

```tsx
type CardProps = { title: string; description: string };

function Card(props: CardProps) {
  return <h2>{props.title}</h2>;
}
```

### Backend API Calls

Always use the centralized `backend` object:

```tsx
import { backend } from "@/modules/backend/backend.ts";

const profile = await backend.getProfile("en", slug);
```

### Routing

Routes use TanStack Router's file-based routing. The `-` prefix excludes files
from route generation (e.g., `-components/` directories for co-located
components).

All user-facing routes are prefixed with `$locale` (e.g., `/tr/eser`,
`/en/eser/qa`).

### Internationalization

- 13 locales: ar, de, en, es, fr, it, ja, ko, nl, pt-PT, ru, tr, zh-CN
- Translation files live in `src/messages/{locale}.json`
- Use `useTranslation()` hook: `t("Section.Key text")`

### Strict Equality

Always use explicit checks:

```tsx
// Correct
if (value === null) {}
if (array.length === 0) {}

// Wrong
if (!value) {}
if (!array.length) {}
```
