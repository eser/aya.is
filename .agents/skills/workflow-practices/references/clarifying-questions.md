# Clarifying Questions

When to ask before implementing, when to proceed on a stated assumption, and how
to ask.

## Contents

- Ask, Assume or Proceed
- How to Ask
- Question Templates
- Stating Assumptions
- Anti-Patterns
- Context-Aware Clarification
- Using AskUserQuestion

---

## Ask, Assume or Proceed

Scope: Before implementing any non-trivial change

Rule:

```
Is the requirement clear?
├── yes → proceed
└── no → are there several valid approaches with different consequences?
    ├── yes → present them with trade-offs and a recommendation, ask
    └── no → is information missing that changes the result?
        ├── yes → ask specific questions
        └── no → state the assumption and proceed unless corrected
```

**Ask** when the answer changes what gets built:

- architecturally different implementations exist;
- business rules have unclear edge cases (soft or hard delete, who may do it);
- existing behavior, a public API or a persisted format would change;
- a security, data-loss or performance trade-off has to be chosen;
- the scope has no boundary ("add a feature" without details);
- the action is on the Ask First list (roles-and-boundaries.md: Ask First).

**Ask, or state a default, when** these are open and matter for the result: the
preferred library, the testing strategy, the error handling approach, the data
validation rules.

**Do not ask** about naming, formatting or anything the codebase, the skills,
AGENTS.md or the checks already decide, standard practice that applies clearly,
or when the owner said to use your judgment.

---

## How to Ask

Scope: Every clarifying question

Rule:

- Ask before writing code, not after.
- Ask everything related in one message; never drip questions one at a time.
- Make each question specific and answerable: "Should hiding a profile also
  hide its stories, or only the profile page?", not "How should this work?".
- For a choice between options, list 2 to 4, each with its consequence, and put
  the recommended one first.
- Offer the default you would take, so the owner can answer with "yes".

Correct:

```
Before implementing the profile feature, I have a few questions:

1. Scope: should it include social links, or just the basic info?
2. Privacy: should the page be public by default or private?
3. Images: should people upload an avatar, or use the GitHub one?

I can proceed with sensible defaults if you prefer.
```

Incorrect:

```
Should it include social links?
[waits for response]
Should it be public by default?
[waits for response]
Should people upload an avatar?
```

---

## Question Templates

Scope: Phrasing a question for a known kind of gap

Rule: Use the shape that fits the gap, and keep it short.

Scope gap:

```
I want to clarify the scope before implementing:
- Should this handle [edge case A]?
- Is [boundary condition C] in or out of scope?
```

Several approaches:

```
There are several ways to implement this:

1. [Approach A]: [one line]. Pros: [benefit]. Cons: [drawback].
2. [Approach B]: [one line]. Pros: [benefit]. Cons: [drawback].

I recommend A because [reason]. Which do you prefer?
```

Missing information:

```
To implement this correctly I need to know [specific question]. Without it I
will assume [default behavior]; tell me if that does not work.
```

Trade-off:

```
There is a trade-off: option A is [faster or simpler] but [limitation], option B
is [more robust] but [cost]. Which priority should I optimize for?
```

---

## Stating Assumptions

Scope: Proceeding without an explicit answer

Rule: When proceeding on an assumption, list it where the owner will see it,
before or together with the work, and name how to change it.

Correct:

```
Assumptions:
- Hiding a profile keeps its stories (they are history)
- The handler returns 404 for an unknown profile id, like the other profile
  handlers
- Validation matches the existing forms

Say if any should be different.
```

Incorrect:

```
I'll just implement it how I think is best.
```

---

## Anti-Patterns

Scope: Common mistakes in clarifying

Rule: Avoid these four.

- **Over-asking:** questions about naming, comments, indentation or anything the
  code style, the linters and existing patterns decide.
- **Under-asking:** assuming critical business logic without verification.
- **Vague questions:** "How should this work?", "Is this okay?".
- **Asking after implementation:** the question comes before the first line of
  code, not after 200 lines of it.

Incorrect:

```
Should I use camelCase or snake_case for variables?        (the linters decide)
*Implements delete as a soft delete without asking*        (a business rule)
*Writes 200 lines* By the way, should it have been B?      (too late)
```

Correct:

```
Should the delete action ask for confirmation or delete immediately?
When a user is deactivated, should their posts stay visible or be hidden?
Before I implement this, should I take approach A (faster) or B (more flexible)?
```

---

## Context-Aware Clarification

Scope: Adjusting the amount of asking to the task

Rule: Match the depth of clarification to the kind of work.

- **Existing code (the usual case):** follow the established pattern and ask
  only when deviating from it; name the pattern: "I'll follow the pattern of
  the existing route".
- **Quick fix:** minimal questions, state the assumption: "I'll fix this the
  simplest way unless you need something specific".
- **Major feature:** more thorough clarification; break a large scope into
  phases: "Let's align on the approach before I implement all of this".
- **New area with no pattern:** ask more about architecture and conventions
  before building.

Correct:

```
User: "Add a delete button to the dashboard"

Before adding it, I want to clarify:
1. What gets deleted: the dashboard itself, or items within it?
2. Should it show a confirmation dialog?
3. Who may delete: any user, or only owners?

I'll add a delete button for items with a confirmation dialog and owner-only
access, following the existing button patterns. Tell me if you need something
different.
```

---

## Using AskUserQuestion

Scope: When to use the structured question tool

Rule: Use AskUserQuestion for decisions with discrete options, and a plain
message for open-ended questions, questions that need code in the answer and
trivial confirmations.

Correct:

```
AskUserQuestion:
- header: "Auth method"
- question: "Which authentication approach should I use?"
- options: "JWT tokens (Recommended)" - stateless;
           "Session cookies" - better for SSR;
           "OAuth only" - delegate to providers
```
