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

# Backend (in /apps/services)
make restart          # Restart after changes
make lint             # Run golangci-lint
make test             # Run tests

# Root
make ok               # All quality checks
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

## Project-Specific Notes

### Internationalization

- 13 locales: ar, de, en, es, fr, it, ja, ko, nl, pt-PT, ru, tr, zh-CN
- Translation keys use English text: `t("Section", "English text")`
- Messages in `/apps/webclient/src/messages/[locale].json`

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
