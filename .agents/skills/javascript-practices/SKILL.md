---
name: javascript-practices
description: "TypeScript and React in the aya.is webclient (apps/webclient: TanStack Start, React 19, Tailwind with CSS Modules, i18next, Vitest) and the Node scripts under scripts/. Use when writing or reviewing .ts or .tsx files, a route, component, form, translated string, an oxlint finding or a pure helper extracted for tests. Not for dependencies or commands (use tooling-standards), test design (use qa-practices) or Go (use go-practices)."
---

# JavaScript Practices

How TypeScript is written in this repository: the webclient (`apps/webclient`,
Node.js 24, strict TypeScript, React 19) and the scripts in `scripts/`. The
repo-wide conventions live in AGENTS.md and are pointed at, not repeated;
tests are qa-practices, tools and commands are tooling-standards.

## Always

- TypeScript only, `.tsx` for modules that render; no `any`, `@ts-ignore` or
  `@ts-expect-error` (oxlint errors on them), no lint suppression comments
- Explicit checks, never truthiness on a non-boolean (AGENTS.md: Explicit
  Checks); `===` always; `??` for a missing value, `||` only when an empty
  string should fall back too
- Named exports on the declaration; a default export only where a tool or
  `React.lazy` reads it
- Components take one `props` object (AGENTS.md: Single Props Object); shadcn
  files in `src/components/ui` keep their generated destructuring
- Server data comes from the `backend` object through query options (AGENTS.md:
  Backend Object Pattern, React Query with SSR); no direct `fetch`
  in a component
- Every visible string goes through `t("Section.English text")` with a real
  translation in all 13 locale files (AGENTS.md: Internationalization & Locale
  Fallback)
- Component styles are a co-located `*.module.css` with `@apply` (AGENTS.md:
  CSS Modules with @apply); inline `style` only for a computed value
- Pure logic that tests need lives in `foo-utils.ts`, the env-dependent wrapper
  in `foo.ts` (AGENTS.md: Pure Function Extraction for Testability)
- `return await` inside `try`; a promise nobody awaits is marked `void`
- Every number with a meaning is a named constant with a comment on its source
- Browser APIs go through `globalThis` and an environment guard, because the
  same module runs in the SSR server
- Every `.ts` and `.tsx` file starts with the one-line license header
  (tooling-standards: Source Files and Headers)
- Before finishing: `pnpm lint`, `pnpm typecheck` and `pnpm test:ci` in
  `apps/webclient`, `pnpm fmt:check` at the root (tooling-standards: Tool per
  Place)

## References

| File                                                  | Read when                                                                                  |
| ----------------------------------------------------- | ------------------------------------------------------------------------------------------ |
| [modules.md](references/modules.md)                   | Adding a file, imports or exports, a route, a script, a pure helper next to its wrapper    |
| [syntax-and-types.md](references/syntax-and-types.md) | Checks, defaults, null, strings, copying data, API errors, async work, storage, prefix cuts |
| [react.md](references/react.md)                       | Building a component, state, translations, loading and error states, dialogs, styles       |
