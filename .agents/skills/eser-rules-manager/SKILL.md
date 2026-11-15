---
name: eser-rules-manager
description: "Picks the skills a task needs and maintains the rules in .agents/skills. Use when starting any task to choose skills, when the owner states a preference or a rule to remember, or when creating, changing, splitting, validating or evaluating a skill."
---

# Skill Discovery and Rule Management

Chooses which skills apply to a task, and says how rules are added and kept in
order. Repository-wide rules live in AGENTS.md, not in skills.

## Always

- Before the first answer, load the skills whose description matches the task (a
  task that touches two areas needs two): a skill loaded late cannot change the
  plan. Read a reference file only when its "Read when" matches
- A new rule goes into the one skill that owns the topic; other skills point to
  it as `skill-name: Exact Section Title`, and to AGENTS.md as
  `AGENTS.md: Section Title`
- Write only what a strong model would not do by default here: project choices,
  real paths and commands, and the reason for each rule
- Every example follows the other skills' rules and names real aya.is code
- SKILL.md stays an overview of at most 80 lines that links every reference file
  with a "Read when"
- Run `node .agents/skills/eser-rules-manager/scripts/validate-skills.ts --fix`
  after every skill change
- A change meant to alter behavior is evaluated against a snapshot of the old
  skill

## Available Skills

| Skill                       | Covers                                                                                  |
| --------------------------- | --------------------------------------------------------------------------------------- |
| `enrichment`                | Round-based interview until no decision is left open                                    |
| `handoff`                   | `/handoff`: a document a fresh session continues from                                   |
| `workflow-practices`        | Clarifying, roles, approvals, task flow, commits, gates                                 |
| `qa-practices`              | Test levels, Go and webclient tests, snapshots, CI checks, flaky tests, reporting       |
| `architecture-guidelines`   | Layout, hexagonal layers and ports, ADRs                                                |
| `design-principles`         | Pure functions, immutability, composition, async                                        |
| `coding-practices`          | Errors, logging, validation, identity comparison, readability                           |
| `javascript-practices`      | Webclient TypeScript and React, modules, types, pure function extraction                |
| `go-practices`              | Go conventions, hexagonal backend, errors, context, data structures                     |
| `security-practices`        | Secrets, input validation, SSRF, injection, production hardening                        |
| `tooling-standards`         | Node and pnpm tooling, Make targets, dependencies, scripts, file naming                 |
| `documentation-conventions` | READMEs, docs and prose style                                                           |
| `release-and-ci`            | Versions, releases, GitHub Actions workflows, Coolify deploys                           |

## References

| File                                                  | Read when                                          |
| ----------------------------------------------------- | -------------------------------------------------- |
| [skill-discovery.md](references/skill-discovery.md)   | Unsure which skills a task needs                   |
| [skill-format.md](references/skill-format.md)         | Writing or restructuring a skill or a rule         |
| [skill-evaluation.md](references/skill-evaluation.md) | Testing that a skill changes behavior and triggers |
