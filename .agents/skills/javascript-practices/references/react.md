# React in the Webclient

How components of the webclient (`apps/webclient/src`) are built with React 19,
TanStack Start and Query, i18next, Base UI through the shadcn-style
`components/ui`, and Tailwind with CSS Modules. Conventions that the whole
repository shares stay in AGENTS.md and are pointed at: Single Props Object,
CSS Modules with @apply, React Query with SSR, Internationalization
& Locale Fallback, Accessibility, DropdownMenu Navigation (Base UI),
Shadcn UI Components and Base UI Select.

## Contents

- Components and Props
- Where State Lives
- React Compiler
- Every String Is Translated
- Loading, Empty, Failure and Pending
- Forms
- Dialogs and Disabled Actions
- Accessibility
- Styles

---

## Components and Props

Scope: Function components and hooks

Rule: A component is a named function export that takes one `props` object and
reads `props.name` (AGENTS.md: Single Props Object); the props type is a
`type` declared above it, and `Readonly<Props>` is accepted. A component that
only serves one screen stays in that screen's file below the exported one.
Shared building blocks live in `src/components/` (grouped by feature:
`privacy/`, `forms/`, `content-editor/`); check them and `src/components/ui`
(generated shadcn components, which keep their own destructuring and Tailwind)
before writing a new one. Icons come from `lucide-react`. Default exports are
only for framework files and `React.lazy` modules (javascript-practices:
Exports).

Correct:

```tsx
type UserProfileProps = {
  userId: string;
  showActions: boolean;
};

export function UserProfile(props: Readonly<UserProfileProps>) {
  return (
    <div>
      <h1>User: {props.userId}</h1>
      {props.showActions && <button type="button">Edit</button>}
    </div>
  );
}
```

Incorrect:

```tsx
function UserProfile({ userId, showActions }: UserProfileProps) { /* destructured in the signature */ }
export default UserProfile; // default export for a regular component
```

**Why:** `props.userId` shows where a value comes from at every use, adding a
prop does not touch the signature, and a name that exists in one place is the
one a refactor finds.

---

## Where State Lives

Scope: Any value a component reads or changes

Rule: Each kind of state has one home; do not copy it into another:

| State | Home |
| --- | --- |
| Data from the backend | TanStack Query, through the option factories in `src/modules/backend/queries.ts`, prefetched in the route loader; never copied into `useState` (AGENTS.md: React Query with SSR) |
| Route values (locale, slug, ids) | The router's params and loader data (`Route.useParams()`, `Route.useLoaderData()`) |
| Signed-in user and preferences | `useAuth()` from `src/lib/auth/auth-context.tsx` |
| Theme | `src/components/theme-provider.tsx` |
| Form fields, an open dialog, a busy flag | `useState` in the component, or `@tanstack/react-form` (`useForm` in `src/components/forms`) |

Writes call the `backend` object (AGENTS.md: Backend Object Pattern) and then
invalidate the affected queries (`queryClient.invalidateQueries`). A component
does not call `fetch` itself.

**Why:** one owner per value keeps refetching and navigation consistent; a copy
in `useState` goes stale on the next refetch.

---

## React Compiler

Scope: Memoization in `.tsx` files and hooks

Rule: Write idiomatic, compiler-friendly components and do not add `useMemo`,
`useCallback` or `memo` without a reason you can state (a measured cost, or an
identity a hook dependency or a schema needs). `babel-plugin-react-compiler`
is installed, but `vite.config.ts` calls `viteReact()` without a compiler
option, so nothing is memoized automatically today; the existing
`useMemo` and `useCallback` calls (theme provider, forms) are for stable
identities. oxlint warns on `react/purity`, `react/immutability`,
`react/set-state-in-effect`, `react/static-components` and
`react-hooks/exhaustive-deps`; fix those instead of silencing them.

Correct:

```tsx
export function ExpensiveList(props: { items: Item[] }) {
  const sorted = props.items.toSorted((a, b) => a.name.localeCompare(b.name));
  return (
    <ul>
      {sorted.map((item) => (
        <li key={item.id}>{item.name}</li>
      ))}
    </ul>
  );
}
```

Incorrect:

```tsx
const sorted = useMemo(() => props.items.toSorted(byName), [props.items]); // no profile shows a cost
```

---

## Every String Is Translated

Scope: Text a person sees, labels, `aria-label` values, toasts and error
sentences

Rule: Source copy is English and the key is one string,
`Section.English text`: `t("Error.Something went wrong")` reads
`messages/<locale>.json` as `{ "Error": { "Something went wrong": "..." } }`
(AGENTS.md: Internationalization & Locale Fallback). A new key goes
into all 13 locale files with a real translation; `translations.test.ts`
fails when a locale lacks a key. Components get `t` from
`useTranslation()` (`react-i18next`); route loaders and `head()` use
`i18next.getFixedT(locale)`. Values are interpolated with i18next
(`{{date}}`), never by concatenating translated pieces, so each language keeps
its word order; dates go through the named formatters (`{{value, date-long}}`).

Correct:

```tsx
const { t } = useTranslation();
<h2>{t("Error.Something went wrong")}</h2>;
toast.error(t("Privacy.You can request a new copy after {{date}}.", { date: formatDate(nextRequestAt) }));
```

Incorrect:

```tsx
t("Error.something_went_wrong"); // snake_case key, not English text
t("Profile", "Failed to connect"); // a second string argument is i18next's default value, not a key part
t("Layout.Profile") || "Profile"; // a missing key is a test failure, not a runtime fallback
<small>{t("Tier") + " " + tier + ", " + t("weight") + " " + weight}</small>; // word order fixed to English
```

**Why:** one key shape lets a reviewer search for any sentence, and the test
that compares every locale to English catches a missing translation before a
reader sees a key.

---

## Loading, Empty, Failure and Pending

Scope: Every view that reads or writes data

Rule: A view that reads data accounts for four states and does not render a
blank page for any of them:

- pending: a route loader awaits its data (`ensureQueryData`) before the route
  renders, so the component that reads it with `useSuspenseQuery` has no
  pending state of its own; a component that fetches by itself (`useQuery`)
  renders a `Skeleton` or `Spinner` from `components/ui` while `isPending`;
- failed: a route sets `errorComponent: QueryError`
  (`src/components/query-error.tsx`); data that is only nice to have uses
  `.catch(() => null)` or `prefetchQuery` (AGENTS.md: React Query with SSR); an action shows its failure in a toast or an `Alert`;
- empty: an explicit sentence or the `Empty` component, never nothing;
- data: the list or form.

While a write is pending its button is disabled and shows progress
(`isBusy`, `Loader2`); what the person entered stays.

Incorrect:

```tsx
if (data === undefined) return null; // a failure looks like an empty page
```

---

## Forms

Scope: Forms that create or change records

Rule: Forms use `@tanstack/react-form` with a zod schema (`create-profile-form.tsx`,
`edit-profile-form.tsx`) and keep native constraints (`required`, `maxLength`).
A schema with translated messages is built inside the component from the
current `t` (`localizedSchema` in `create-profile-form.tsx`). A field that
fails validation gets `aria-invalid`; a refusal from the backend is chosen by
its code (javascript-practices: Errors Are Matched by Code). Limits are named
constants that mirror the backend's.

---

## Dialogs and Disabled Actions

Scope: Confirmations, destructive actions and buttons that cannot act

Rule: Ask in a dialog (`AlertDialog` from `components/ui/alert-dialog`, as
`privacy/delete-section.tsx` does), not with `globalThis.confirm()`; the
three existing `confirm` calls (`content-editor.tsx`, `date-poll.tsx`) are not
a pattern to copy. A destructive dialog names what is removed and what stays,
and a very destructive one asks to type the identifier
(`isDeletionConfirmed`). A button that cannot act now stays visible and
disabled with its reason next to it; do not hide an action because it is not
allowed right now. Menu items that navigate use the `render` prop with
`LocaleLink` (AGENTS.md: DropdownMenu Navigation (Base UI)).

---

## Accessibility

Scope: Interactive elements

Rule: The rules are in AGENTS.md: Accessibility (`data-slot="card"`, an
`aria-label` on every icon-only `Button`, `role="menu"` and `role="menuitem"`
on custom menus) and AGENTS.md: Base UI Select. `.oxlintrc.json`
enables the `typescript`, `react` and `oxc` plugins and no `jsx-a11y`, so no
tool checks these: review them by hand. A field sits inside its `label`, an
icon next to a text label is `aria-hidden`, and focus returns to the control
that opened a dialog.

---

## Styles

Scope: Classes, Tailwind and inline styles

Rule: Component styling is a co-located `*.module.css` that begins with
`@reference "@/theme.css";` (Tailwind 4 needs it for `@apply`) and uses
`@apply` with semantic camelCase class names, imported as `styles` and used as
`styles.headerDot` (AGENTS.md: CSS Modules with @apply). Direct Tailwind
classes are for a small, non-reusable adjustment or a top-level layout
container, combined with `cn()` from `@/lib/utils` when mixed with a module
class. Inline `style` only for a value computed at runtime
(`style={{ height }}`); never a literal color or size. Much existing code
still writes long Tailwind class lists inline; do not convert files a change
does not touch, but new components follow the rule.

Correct:

```css
/* live-now.module.css */
@reference "@/theme.css";

.card {
  @apply border rounded-lg p-4 shadow-md flex flex-col gap-2;
}
```

```tsx
import styles from "./live-now.module.css";

<div className={cn(styles.imageContainer, "w-full md:w-[250px]")} />;
```

Incorrect:

```tsx
<div className="border rounded-lg p-4 shadow-md flex flex-col gap-2"> {/* utility soup on a reusable component */}
<div style={{ border: "1px solid", padding: "16px" }}> {/* literal values inline */}
```
