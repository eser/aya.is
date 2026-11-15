---
name: documentation-conventions
description: "Where aya.is docs live (README, CHANGELOG.md, docs/adr, READMEs beside code), docs that change with the code, implemented versus proposed, domain terms, and plain prose without empty qualifiers, unbacked claims or em dashes. Use when writing or editing READMEs, docs, skill files, ADR text, issues, commit messages or code comments. Not for ADR structure (use architecture-guidelines)."
---

# Documentation Conventions

Rules for what the repository's documents say, where they live and how they
say it. What may never appear in a document (task status, progress, delivery
plans) is readmes.md: Documentation Boundary; follow it for every Markdown file.

## Always

- Put a document where its kind already lives (readmes.md: Where Documentation
  Lives); technical docs go under `docs/`, decisions under `docs/adr/`,
  release notes in `CHANGELOG.md`.
- Task status, progress and delivery plans go into GitHub issues and pull
  requests, never into docs, comments or copy.
- The documentation diff ships in the same change as the code it describes:
  the file that lists a setting, the `make help` comment of a target, a
  package README, the locale files and the changelog entry when they are
  affected.
- State what is implemented and what is proposed; never let a design read as an
  available capability.
- Use the domain terms of AGENTS.md (AYA, profile and its types, story and its
  kinds, locale) and one term per concept.
- Keep descriptive details that are still true when rewriting; delete
  sentences that became false.
- No empty qualifiers (comprehensive, seamless, robust), no claim, number or
  quote that is not in the repository, no em dashes in new or edited text.
- Before delivering, read the text as a stranger: what looks generated, and
  which fact has no source? Fix both.

## References

| File                                | Read when                                                                   |
| ----------------------------------- | --------------------------------------------------------------------------- |
| [readmes.md](references/readmes.md) | Adding or moving a document, writing a README, changing documented behavior |
| [prose.md](references/prose.md)     | Writing or reviewing any prose: docs, skills, commits, issues, comments     |
