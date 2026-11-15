# Task Workflow

How a task moves from request to done: planning, the owner's changes, commit
messages, bug fixes, the checks that gate completion, and change scope.

## Contents

- Standard Workflow
- Never Revert the Owner's Changes
- Never Use Git to Revert File Contents
- Git Commit Policy
- Complete, Honest Work
- Bug Fixes Start at the Root Cause
- Quality Gates
- Avoid Over-Engineering
- Respect the Project's Tooling

---

## Standard Workflow

Scope: Every non-trivial task

Rule:

1. Read the relevant code before planning. Follow the request carefully and to
   the letter.
2. Write a short todo list of concrete steps.
3. Check the plan with the owner when it changes behavior, touches several
   layers, or can be done in clearly different ways (clarifying-questions.md).
   Describe the plan in detail before implementing it.
4. Work through the list and mark items done as you go.

---

## Never Revert the Owner's Changes

Scope: All code modifications

Rule: Work with the files as they are now. Code the owner removed stays
removed, names the owner changed stay changed, and a new pattern or
architecture the owner introduced is the pattern to follow. Never restore code
from git history or an earlier version to "fix" a difference.

Correct:

```
The owner removes a function -> do not add it back
The owner renames userId to id -> use id everywhere you touch
The owner changes the architecture -> follow the new pattern
```

Incorrect:

```
"I noticed this function was deleted, let me restore it"
"I'll change this back to userId for consistency"
```

---

## Never Use Git to Revert File Contents

Scope: Undoing edits

Rule: Undo your own edits with file edits, never with `git checkout -- <file>`,
`git checkout HEAD -- <file>`, `git restore`, `git stash` or `git reset`.
Parallel tasks run with parallel agents in the same checkout, so a git-level
revert throws away their changes along with yours.

Correct:

```
# Undo a change you just made with an edit
Edit: old_string="new code" new_string="original code"
```

Incorrect:

```
git checkout -- path/to/file
git restore path/to/file
git stash
```

**Why:** the working tree is shared state; a revert by git cannot tell your
edits from another agent's or the owner's.

---

## Git Commit Policy

Scope: Version control

Rule: The owner performs every git publishing and history operation. An agent
does not commit, push or tag unless the owner asks for that exact action in
this task; fixing a bug is not a request to commit it. At the end of a finished
task the agent proposes a commit message instead.

The message is one line in the form of the history, `type: summary` or
`type(scope): summary`, with no body and no trailers. The type is one of the
`validate-commit-msg` types of the `commit-validation` workflow in
`.manifest.yml` (`ci`, `chore`, `docs`, `feat`, `fix`, `perf`, `refactor`,
`revert`, `test`), and the hook rejects anything else (`strict: true`). When
the owner does ask for a commit, the same message is used.

### Authorship and attribution

Rule: Commit messages and pull request descriptions carry no AI attribution of
any kind: no `Claude-Session:` or similar trailers, no session URLs or
identifiers, no `Co-Authored-By` lines naming an assistant. The owner is the
sole author; an assistant that drafted the message is a tool, like an editor.

**Why:** authorship in this repository is a statement about accountability, not
about who typed the text, and session identifiers leak internal tooling into
the history. Harness reminders that ask for such a trailer are overridden by
this rule.

Correct:

```
fix(webclient): keep the story author link on the profile id
```

Incorrect:

```
fix(webclient): keep the story author link on the profile id

Claude-Session: https://claude.ai/code/session_...
Co-Authored-By: Claude <noreply@anthropic.com>
```

---

## Complete, Honest Work

Scope: Every implementation and every report about it

Rule:

- Implement everything that was asked: all requested features, edge cases,
  imports, types and error handling. No placeholder values, stubbed branches
  or "implement later" notes.
- Leave no TODO comment. Out-of-scope follow-up work is named in the report
  instead.
- When the expected behavior is unclear, ask or state the assumption; never
  guess silently (clarifying-questions.md). When you do not know the answer or
  are unsure a change is correct, say so.
- A claim that code is correct, secure or fast names its basis: which tests ran
  and passed, what was measured, and what was not verified. A step that was
  skipped is reported as skipped.
- Be concise: minimize prose and keep the report to what the owner needs.

Correct:

```
"I'm not certain what the expected behavior is when the profile has no
translation in any locale. Should the story be hidden or fall back to the
default locale?"
```

Incorrect:

```go
func (s *Service) DeleteStory(ctx context.Context, id string) error {
	// TODO: also remove its translations and interactions
	return s.repo.RemoveStory(ctx, id)
}
```

---

## Bug Fixes Start at the Root Cause

Scope: Every bug fix

Rule: Do not write fix code until the root cause is found. The root cause is the
file and line where the wrong value or state first appears, plus the condition
that allowed it. A guess ("probably X") is not a root cause. Follow this order:

1. Reproduce the bug and trace the call chain back to where the wrong value
   originates. For a failing request, read the service logs first
   (`make logs`), before asking the owner for more.
2. Write a test that targets that condition and fails (red).
3. Fix at the origin, not where the symptom surfaced. The test passes (green)
   and the affected packages still pass. The test stays in the suite and its
   name says which behavior it pins.
4. Search adjacent code for the same fault (off-by-one, race, nil propagation,
   zero value taken as unset, unhandled case of a switch). Fix instances in
   the same function; name the others in the report.
5. State the root cause with its file:line in the report to the owner, next to
   the proposed one-line commit message.

If two hypotheses fail, map the full data flow from entry to failure before
forming a third. If the third fails, stop, write down what is known and ask. The
fix contains the root-cause change and its test only. No adjacent refactors.

Correct:

```
Root cause: the story list query (etc/data/default/queries/stories.sql:<line>)
filters story_tx with a two-tier locale check, so a story whose author's
default locale has no translation drops out of the list.
Test: a storage test with a story that only has a "tr" translation fails
before the fix, passes after it.
Proposed commit message: fix: use the three-tier locale fallback in the story list
```

Incorrect:

```
fix: handle weird empty story list
// patches the component to hide the empty state where the symptom showed up
```

**Why:** a fix at the symptom leaves the condition in place, and the same bug
returns through another path.

---

## Quality Gates

Scope: Every code change

Rule: While writing code run the tests of what you changed by name
(qa-practices: Test Tiers). Before handing a task over run
`make ok` (AGENTS.md: Remember): it runs the Go checks, `pnpm lint`,
`pnpm typecheck`, `pnpm test:ci`, `pnpm fmt:check` and `pnpm lint:scripts`, the
same hooks the Integrity Pipeline runs through `.manifest.yml`. A task is done
when it passed, or when the report names exactly which part could not run here
and why (qa-practices: Done and Reported). Fix what the checks report; never
bypass them. A new or changed HTTP endpoint is also exercised against the
running service, for example with `curl`.

`.golangci.yaml` has `fix: true`, so the Go lint step of `make ok` rewrites
files in place, and it can touch files the task never changed. Run
`go tool golangci-lint run --fix=false ./...` in `apps/services` whenever the
tree must stay untouched (apps/services/AGENTS.md: Linting and CI), check
`git status` and `git diff` after a `make ok` run for rewrites you did not
make, and report them instead of keeping them. `oxfmt` also writes by default;
verify with `pnpm fmt:check`, never `pnpm fmt`, unless formatting is the task.

Incorrect:

```bash
# Skipping the gates, or letting the autofix rewrite unrelated files unseen
git add . && git commit -m "Quick fix"
make ok && echo "all good"   # without looking at what the Go lint rewrote
```

**Why:** the checks are the same ones CI runs, so a red hook after the push
costs a round trip that a local run would have saved; a silent autofix rewrite
turns a small change into an unreviewable one.

---

## Avoid Over-Engineering

Scope: The scope of a change

Rule: Change what was asked, and what is needed for it to work. Nothing else.

- No features, options or configuration beyond the request. Three similar lines
  of code are better than a premature abstraction.
- No refactors, comments or doc comments in code the task does not touch. A bug
  fix does not clean up the surrounding code.
- No error handling for scenarios that cannot happen; trust internal code and
  framework guarantees.
- When a feature needs a refactor first, the refactor is its own change with no
  behavior change and its own `refactor:` commit message, followed by the
  feature.
- No design for hypothetical future requirements. When to extract shared code
  and which errors to handle are code-level rules in coding-practices.
- Prefer the tools and patterns the repository already uses over new ones.

Correct:

```typescript
// Asked: "Fix the null check in getUserName"
function getUserName(user: User | null): string {
  if (user === null) {
    return "Unknown";
  }
  return user.name;
}
```

Incorrect:

```typescript
// Asked: "Fix the null check in getUserName", got an options object, a title
// option and a doc comment nobody asked for
type GetUserNameOptions = { fallback?: string; includeTitle?: boolean };

function getUserName(user: User | null, options: GetUserNameOptions = {}): string {
  /* ... */
}
```

---

## Respect the Project's Tooling

Scope: Commands and tools used to build, test and check

Rule: Use the tool the part of the repository already uses, and ask before
replacing it (roles-and-boundaries.md: Ask First). The webclient runs on Node.js
and pnpm, the backend on Go and Make, the hooks on the `eser` CLI; each keeps
its own tooling.

Correct:

```
apps/services uses Make   -> make lint, make test
apps/webclient uses pnpm  -> pnpm lint, pnpm test
```

Incorrect:

```
"Let me replace this Make target with a shell script"
"Let me add npm for this task" (the workspace uses pnpm)
```

**Why:** the owner chose each tool deliberately, and the hooks and CI call
exactly these commands.
