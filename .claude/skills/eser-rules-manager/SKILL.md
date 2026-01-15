---
name: eser-rules-manager
description: Skill discovery and rule management. Use when starting any conversation to identify applicable skills, and when user states preferences or asks to add/modify rules.
---

# eser-rules: Skill Discovery & Rule Management

Two functions: (1) Discover and invoke relevant skills before ANY response,
(2) Manage development rules across skills.

## Quick Start

1. **Every conversation**: Scan message → Invoke relevant skills → Announce → Respond
2. **Rule changes**: Identify scope → Update skill → Validate → Test

## Skill Discovery (Mandatory)

**Before ANY response** (including clarifying questions):

1. Scan user message for applicable skills
2. Invoke relevant skills using Skill tool
3. Announce: "Applying skills: [list]"
4. Follow skill instructions, then respond

Even 1% probability requires checking skills first. This is not optional.

## Rule Management

1. Identify scope → choose skill (or create new)
2. Add/update rule in `.claude/skills/<name>/references/rules.md`
3. Validate: `npx -y claude-skills-cli validate .claude/skills/<name>`

## Available Skills

| Skill                       | Triggers                                  |
| --------------------------- | ----------------------------------------- |
| `javascript-practices`      | JS/TS code, modules, types                |
| `go-practices`              | Go code, hexagonal architecture           |
| `workflow-practices`        | Task execution, git commits               |
| `requirement-clarification` | Unclear scope, multiple approaches        |
| `security-practices`        | Auth, secrets, validation                 |
| `ci-cd-practices`           | GitHub Actions, Kubernetes, deployments   |

## References

- [skill-discovery.md](references/skill-discovery.md) - Mandatory invocation rules
- [skill-format.md](references/skill-format.md) - Creating/updating skills
- [skill-testing.md](references/skill-testing.md) - TDD for skills
