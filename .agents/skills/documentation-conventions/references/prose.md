# Plain Prose

Writing rules for every sentence a person reads in this repository.

## Contents

- Plain Prose
- Empty Qualifiers
- Claims Need Evidence
- No Em Dashes
- Borrowed Rhythm
- Filler and Hedging
- Name the Actor
- Neutral Terms
- Formatting That Shouts
- Audit Before Delivering

---

## Plain Prose

Scope: Every sentence a person reads: READMEs, docs, skill files, ADRs,
issues, commit messages, Go doc comments, code comments, interface copy

Rule: Documentation says what the code does, in the words a maintainer would
use out loud. Text that reads as generated (empty qualifiers, borrowed rhythm,
claims with nothing behind them) is a defect even when every fact in it is
true, because readers stop trusting the page. The rules below name the
patterns that give it away. When a sentence trips one, fix the sentence; do
not pad it back up with a different pattern.

---

## Empty Qualifiers

Rule: Cut words that sound like praise and carry no information:
comprehensive, seamless, robust, powerful, elegant, cutting-edge, next-level,
world-class, effortless, streamline, leverage, harness (as a verb), unlock,
elevate, empower, delve, showcase, testament, game-changer. Replace the word
with what the component actually does, or with nothing.

Correct:

```markdown
The story query returns the translation in the requested locale, then the
profile's default locale, then any available translation.
```

Incorrect:

```markdown
The story query provides robust, comprehensive localization that seamlessly
serves every reader. <!-- three claims, no information -->
```

---

## Claims Need Evidence

Rule: A statement about quality, scale or support is either backed by
something in the repository (a test, a CI run, a measured number with its
command and setup) or it is not written. This covers "production-ready",
"battle-tested", "fast", "supports every client", and any number, quote or name
that did not come from a source the reader can open. When you do not know a
fact, leave it out; never fill the gap with a plausible guess.

Correct:

```markdown
The migrations were applied with `make migrate-up` against the
`pgvector/pgvector:pg18` image that `compose.yml` starts.
```

Incorrect:

```markdown
The backend runs on any PostgreSQL version. <!-- which versions, tried how -->
The webclient supports every browser. <!-- no source, and false -->
```

---

## No Em Dashes

Rule: New or edited prose uses no em dash (`—`), spaced or not, and no double
hyphen standing in for one. Replace each with, in order of preference, a
period, a comma, a colon, or parentheses; restructure when none fits. Verbatim
quotes of program output keep their punctuation. Convert the sentences you are
already editing; do not reformat a file you were not asked to change.

Correct:

```markdown
`make down` removes the containers and keeps the `postgres-data` volume: the
stories and profiles in the database survive.
```

Incorrect:

```markdown
`make down` removes the containers — but keeps `postgres-data` — so the
stories and profiles in the database survive.
```

---

## Borrowed Rhythm

Rule: Avoid the sentence shapes that models reach for by default:

- "It's not just X, it's Y" and "not only X but also Y"
- A run of clipped fragments for drama ("No config. No setup. No surprises.")
- Every list forced to three items when the content has two or five
- "From X to Y, and everything in between" where X and Y are not a scale
- Aphorisms ("Simplicity is the language of trust")
- Signposting ("Let's dive in", "Here's what you need to know", "In this
  section we'll explore")
- Chatbot closers ("I hope this helps", "Let me know if you have questions")
- Upbeat send-offs ("Exciting times ahead", "The future looks bright")

Say the thing directly, in as many items as there are.

Correct:

```markdown
The profile page lists the stories of the profile.
```

Incorrect:

```markdown
It's not just a profile page, it's your home on the web. No clutter. No noise.
```

---

## Filler and Hedging

Rule: "In order to" is "to"; "due to the fact that" is "because"; "it is
important to note that" and "it is worth noting" are nothing. One qualifier
per claim: "may" or "usually", never "could potentially possibly".

Correct:

```markdown
The server closes the connection to stop reuse of a stale session.
```

Incorrect:

```markdown
It is worth noting that the server may potentially close the connection in
order to prevent reuse of a stale session.
```

---

## Name the Actor

Rule: When the actor is known, the sentence names it. "The API server signs
the session token", not "the session token is signed". Passive is fine when
the actor is unknown, irrelevant or deliberately withheld ("the server was
restarted at 03:00"). Do not give abstractions a mind: a resolver resolves, a
check reports, a relay forwards; none of them "understands", "decides" or
"wants".

---

## Neutral Terms

Rule: Use neutral technical terms in new code, flags and docs: allowlist and
denylist, primary and replica, main. Do not assume the reader's gender,
operating system or access to a paid tool.

---

## Formatting That Shouts

Rule: Emphasis is written into the sentence, not into the markup.

- No inline-header bullets that restate the item
  (`- **Performance:** it is faster`). A bold lead that names a distinct
  subject in a list of parallel cases (`- **Update.**`, `- **Backup.**`) is
  fine.
- No bold on every key term; bold at most the one phrase the reader must not
  miss.
- No quotation marks as emphasis or hedging; quotes are for real quotes,
  titles and program output.
- No all-caps clauses inside paragraphs.
- No emojis in documents, comments or copy.

---

## Audit Before Delivering

Rule: Before handing over any documentation, read it once as a stranger and
answer two questions: what in this text makes it look generated, and does it
state any fact, number, name, date or quote that is not in the source or in
the repository. Fix both. Then search the diff for `—`; a hit means the text
is not done.
