---
name: workflow-practices
description: "How agents work in this repository: clarifying requests, roles, approvals, forbidden git and release actions, root-cause fixes, quality gates. Use for every request to add, build, change, refactor, migrate or fix something, beside the domain skill, above all open-ended ones (add caching, make X faster), for reviews, and before any commit, push or tag."
---

# Workflow Practices

The flows and roles of agent work in this repository: clarifying the request,
what an agent may do alone, how a task runs, and when it is done. Repository
rules that hold for every task live in AGENTS.md; this skill points to them.

## Always

- Ask before building when the answer changes what gets built (architecture,
  business rules, formats, security or data loss, unbounded scope); otherwise
  state the assumption and proceed. Ask everything in one message, 2 to 4
  options, recommended first; never ask what the code, skills or checks decide
- Read the code, write a todo list, check a behavior-changing plan with the
  owner, then execute; follow the request to the letter and say so when unsure
- Stay inside the area the task names and change only what it needs; a change
  across layers follows the protocol; a refactor is its own change
- The owner commits and pushes: never commit, push, tag, reset, stash, check
  out, restore or run release commands unless asked for that exact action;
  propose a one-line conventional commit message (the types of
  `.manifest.yml` `commit-validation`) with no AI attribution, trailers or
  session ids
- Never revert what the owner changed; undo your own edits with edits, never
  with `git checkout`, `git restore` or `git stash`: parallel agents share the
  checkout
- Ask first before new dependencies (`go.mod`, `package.json`), workflow,
  deployment or release changes, replacing a tool, or anything a release
  cannot undo
- Never edit `*_gen.go`, `routeTree.gen.ts`, `go.sum` or `pnpm-lock.yaml` by
  hand; change the source and regenerate
- Bug fixes: root cause with file:line, failing test first, fix at the origin
- Implement everything asked, with no placeholder or TODO comment; follow-up
  work outside the request is named in the report
- Run the tests of what you changed by name while working, and `make ok`
  before handing over; mind that its Go lint step rewrites files
  (`--fix=false`); done means it passed or the report names what could not run
  (qa-practices: Done and Reported, AGENTS.md: Remember)
- Reviewers change nothing, cite file:line and mark each finding blocking or
  informational
- Code, comments, commit messages and skills are English whatever the
  conversation language; product copy follows the i18n rules (AGENTS.md:
  Internationalization & Locale Fallback)

## References

| File                                                          | Read when                                                                                    |
| ------------------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| [clarifying-questions.md](references/clarifying-questions.md) | Deciding whether to ask, or phrasing questions and assumptions                               |
| [roles-and-boundaries.md](references/roles-and-boundaries.md) | Checking whether an action is allowed, who owns an area, or ordering a change across layers  |
| [task-workflow.md](references/task-workflow.md)               | Planning a task, fixing a bug, proposing a commit message, or checking which gates must pass |
