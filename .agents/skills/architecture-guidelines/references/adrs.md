# Architectural Decision Records

When to write an ADR in `docs/adr/`, what it contains and how its status
changes. `docs/adr/` holds `README.md` and `template.md` (MADR); it has no
records yet.

---

## Architectural Decision Records

Scope: Decisions that change a process or module boundary, the contract between
the webclient and the Go API, a persisted or wire format, the data model, a
dependency with lasting cost (a framework, a router, a database), or a product
rule

Rule: Record the decision in `docs/adr/NNNN-kebab-title.md`, numbered after the
highest existing file (`ls docs/adr`) with four digits, starting from `0001`
(`0001-example-decision-title.md`). Start from `docs/adr/template.md`
(or `pnpm dlx adr-tool create <decision>`, which `docs/adr/README.md`
documents). Write it before or with the change, in English, and keep the
template's section order, dropping the optional ones that have nothing to say:

- Title line `# Short title of the problem and its solution`, then `Status`
  (`proposed`, `accepted`, `deprecated` or `superseded by [ADR-NNNN](...)`),
  `Deciders` and `Date` (when the decision was last updated).
- `## Context and Problem Statement`: the forces and the question.
- `## Decision Drivers`: what the decision optimizes for.
- `## Considered Options`, then `## Pros and Cons of the Options` with the
  trade-offs of each.
- `## Decision Outcome`: the chosen option and why it beat the others, with
  `### Positive Consequences` and `### Negative Consequences`: what becomes
  easier and what becomes harder.
- `## Links`: related records and the change that implements it.

The owner decides. An agent drafts a record as `proposed` with a recommended
choice for every open question and never marks it accepted on its own; a point
the owner did not decide explicitly stays marked proposed. An agent never
commits the record.

A later record that changes an accepted one supersedes or refines it: the new
record's status or links name the older one, and the older record's status
points to the new record. The original reasoning stays. State in the record
which parts exist in the code and which are still design, so no reader takes a
proposal for a capability. An ADR carries no task IDs, milestones or progress
notes; those belong in the tracker.

**Why:** the reasoning behind a structure is what a later change needs, and it
is the part the code cannot show.
