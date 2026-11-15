# Project Structure and Naming

How the monorepo, the webclient and file names are laid out. Go package
structure is in [hexagonal-architecture.md](hexagonal-architecture.md). Tests
are covered by qa-practices; where documentation lives by
documentation-conventions: Where Documentation Lives.

## Contents

- Module System
- Monorepo Structure
- Naming Conventions
- Centralized Backend Object Pattern

---

## Module System

Scope: JS/TS code in `apps/webclient`

Rule: Use ES Modules (`import` / `export`), the standard every modern browser
and server runtime supports. Do not use CommonJS or AMD. Import with the
`@/` alias for code under `src/` (`@/modules/backend/backend.ts`), as the
examples in AGENTS.md do, and follow the extension style of the file you are
editing.

Correct:

```typescript
import { isValidLocale } from "@/lib/locale-utils.ts";

export function processLocale(locale: string) {}
export const CONFIG = {};
```

Incorrect:

```typescript
const { isValidLocale } = require("@/lib/locale-utils"); // CommonJS

module.exports = { processLocale }; // CommonJS
exports.CONFIG = {}; // CommonJS
```

**Why:** ES Modules are statically analyzable, so bundlers tree-shake them and
the type checker follows them; CommonJS defeats both.

---

## Monorepo Structure

Scope: The repository as a whole

Rule: `apps/` holds deployable applications, each with its own build
configuration and toolchain; a shared package gets its own directory with clear
dependencies. The frontend is `apps/webclient` (pnpm workspace member) and the
backend is `apps/services` (Go module `github.com/eser/aya.is/services`).

```
aya.is/
├── apps/
│   ├── webclient/                # TanStack Start frontend (Node.js + pnpm)
│   │   └── src/
│   │       ├── routes/           # File-based routing
│   │       ├── components/       # UI components (shadcn-style primitives in ui/)
│   │       ├── modules/          # Feature modules (backend, i18n, navigation)
│   │       ├── lib/              # Utilities, pure logic and its tests
│   │       ├── hooks/            # React hooks
│   │       └── messages/         # i18n catalogs, one file per locale
│   └── services/                 # Go backend (hexagonal)
│       ├── cmd/                  # serve, manage, migrate
│       └── pkg/
│           ├── api/business/     # Pure business logic
│           ├── api/adapters/     # External implementations
│           ├── ajan/             # Framework utilities
│           └── lib/              # Shared helpers
├── docs/adr/                     # Architectural decision records
└── .agents/skills/               # Development rules
```

- The frontend uses file-based routing under `src/routes`; the router tree is
  generated (`routeTree.gen.ts`), never edited by hand.
- Backend services follow the hexagonal layout in
  [hexagonal-architecture.md](hexagonal-architecture.md).
- Run a quality check where the app lives: `pnpm lint`, `pnpm typecheck` and
  `pnpm test` in `apps/webclient`, `make lint` and `make test` in
  `apps/services`, `make ok` from the root before a commit.

---

## Naming Conventions

Scope: Files and directories

Rule:

- Webclient directories and files are `kebab-case`: `profile-envelopes/`,
  `app-links.ts`, `query-error.tsx`, `get-profile.ts`. Exported React
  components and types are `PascalCase`, functions and variables `camelCase`.
- A pure helper next to an environment-dependent module is `foo-utils.ts`
  (`locale-utils.ts` beside `config.ts`).
- Go files and business module directories are `snake_case`
  (`profile_points/`, `story_interactions/`, `repository_stories.go`); Go
  package names are lowercase (`apps/services/AGENTS.md: Layout`).
- Tests sit next to the code: `*.test.ts` / `*.test.tsx` (Vitest) and `*_test.go`
  (qa-practices).
- One API function per file in `modules/backend/<area>/`, named for the action
  (`get-profile.ts`, `update-profile-link.ts`).

Correct:

```
src/modules/backend/profiles/get-profile.ts
src/components/profile-picture-upload.tsx
src/lib/locale-utils.test.ts
pkg/api/business/profile_points/service.go
```

Incorrect:

```
src/modules/backend/profiles/getProfile.ts   // camelCase file
src/components/ProfilePictureUpload.tsx      // PascalCase file
pkg/api/business/profilePoints/service.go    // camelCase Go directory
```

---

## Centralized Backend Object Pattern

Scope: Webclient code that talks to the Go API

Rule: `src/modules/backend/backend.ts` exports one `backend` object that is the
single entry point for API calls. Import the object, not the individual
function files; type imports from submodules are allowed. Route loaders and
components do not call `backend` directly for page data: they use the query
option factories in `src/modules/backend/queries.ts`, which call `backend`
(AGENTS.md: React Query with SSR).

```
src/modules/backend/
├── backend.ts          # exports the backend object
├── queries.ts          # React Query option factories
├── fetcher.ts          # HTTP client utilities
├── types.ts            # shared types
├── profiles/
│   ├── get-profile.ts  # one function per file
│   └── ...
└── stories/
```

Correct:

```typescript
import { backend } from "@/modules/backend/backend.ts";
import type { Profile } from "@/modules/backend/types.ts";

const profile = await backend.getProfile("en", props.profileId);
```

```typescript
// backend.ts
import { getProfile } from "./profiles/get-profile.ts";
import { getStory } from "./stories/get-story.ts";

export const backend = {
  // Profiles
  getProfile,

  // Stories
  getStory,
} as const;
```

Incorrect:

```typescript
// Direct function import bypasses the facade
import { getProfile } from "@/modules/backend/profiles/get-profile.ts";

// Inline fetch in a component
const response = await fetch(`/api/profiles/${props.profileId}`);
const profile = await response.json();
```

**Why:** one entry point is the single source of truth for the contract with
the Go API, is easy to mock in a test, and keeps a refactor of an endpoint to
one file.
