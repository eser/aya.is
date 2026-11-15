# Skill Discovery

Which skills to load for a task, and in what order.

---

## Load Before Responding

Scope: Every task, including questions and clarifying questions

Rule: Before the first answer, question or file read, load every skill that
might apply (a small chance is enough) and say which ones you are applying:
`Applying skills: workflow-practices, go-practices, security-practices`. A skill
loaded after the plan is made cannot change the plan. Loading one that turns out
not to matter costs little.

Look past the obvious one. "Add an API endpoint for user preferences" needs
go-practices, and also architecture-guidelines (layers and ports),
security-practices (input and auth), qa-practices (tests) and
workflow-practices (steps and gates). The same holds for a question: before
asking which caching strategy to use, load workflow-practices, because it says
when to ask and when to assume.

Correct:

```
User: Add caching to the API
Agent: Applying skills: workflow-practices, architecture-guidelines, go-practices
       (then asks about the data and the invalidation, as
       workflow-practices: Ask, Assume or Proceed describes)
```

Incorrect:

```
User: Add caching to the API
Agent: What caching strategy do you prefer?   <!-- asked before loading any skill -->
```

**Why:** agents tend to under-load skills; the ones missed are the ones whose
rules were broken.

---

## Skill Trigger Keywords

| Skill                       | Trigger keywords and patterns                                                 |
| --------------------------- | ----------------------------------------------------------------------------- |
| `enrichment`                | enrich, grill, stress-test, interrogate a plan or idea                        |
| `handoff`                   | user-invoked only: `/handoff [next focus]`                                    |
| `workflow-practices`        | any file change, vague request, review, commit, push, tag, delete             |
| `qa-practices`              | test, go test, vitest, snapshot, flaky, coverage, failing check, done         |
| `architecture-guidelines`   | structure, layout, layer, port, adapter, hexagonal, ADR                       |
| `design-principles`         | pure function, immutability, interfaces, composition, async                   |
| `coding-practices`          | errors, logging, validation, naming, comments, comparing entities             |
| `javascript-practices`      | .ts, .tsx, React, webclient, route, module, import, types                     |
| `go-practices`              | .go, apps/services, pkg/, context, errors, log/slog, sqlc, go vet             |
| `security-practices`        | auth, secret, token, input, SSRF, injection, CORS                             |
| `tooling-standards`         | dependency, go.mod, package.json, pnpm, Make target, oxlint, oxfmt, file name |
| `documentation-conventions` | README, docs, CHANGELOG, doc comments, ADR prose, commit message              |
| `release-and-ci`            | release, version, tag, workflow, CI, deploy, compose, Coolify                 |
| `eser-rules-manager`        | new rule, preference, create or change a skill, evals                         |

---

## Order

When several skills apply, apply them in this order:

1. Process: workflow-practices (clarifying, approvals, steps and gates).
2. Design and implementation: architecture-guidelines, design-principles, the
   language skill, coding-practices, security-practices.
3. Verification: qa-practices (which tests, what the report says).
4. Delivery: release-and-ci.

Process and safety skills are followed as written (commit policy, quality gates,
release steps). Design and language skills state defaults that the code around
the change may refine. Repository-wide rules in AGENTS.md hold over every skill.
