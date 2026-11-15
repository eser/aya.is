---
name: tooling-standards
description: "aya.is tooling: Node.js 24 with a pnpm workspace, oxlint, oxfmt, TypeScript 7 tsc and Vitest for the webclient; Go with make and go tool for the services; the eser CLI and .manifest.yml for git hooks and CI. Use when running a build or check, adding a script, make target, dependency, patch or package, editing package.json, go.mod or pnpm-workspace.yaml, or naming files. Not for code style (use go-practices or javascript-practices)."
---

# Tooling Standards

The tools of this repository and how they are run: pnpm and Node.js for the
webclient, make and `go tool` for the Go services, the `eser` hooks, the
dependencies, file naming and layout. Which tests to write is qa-practices;
workflows and deployment are release-and-ci.

## Always

- Webclient: Node.js 24 and pnpm 12 in a workspace (the root and
  `apps/webclient`), oxlint (type-aware), oxfmt, the native TypeScript 7 `tsc`
  and Vitest; install once at the root with `pnpm install`
  (`--frozen-lockfile` in CI and Docker)
- Never `npm`, `npx` or `yarn`; never `deno fmt`, `deno lint`, `deno test` or
  `deno task`: Deno only installs and hosts the `eser` CLI
- Go: `make` inside `apps/services` (`make ok`, `make lint`, `make test`,
  `make generate`); linters run as `go tool`, and `.golangci.yaml` has
  `fix: true`, so use `--fix=false` to keep the tree untouched
- `.manifest.yml` defines the checks; git hooks and CI both run them through
  `eser` (`eser ok`), and `make ok` runs the same ones by hand; never skip the
  hooks
- oxfmt writes in place by default: verify with `pnpm fmt:check`, write with
  `pnpm fmt`
- A command people or CI run is a `pnpm` script, a make target or an
  `.manifest.yml` script, with real logic in `scripts/`
- Never edit `go.sum`, `pnpm-lock.yaml` or a generated file (`*.gen.ts`,
  `*_gen.go`) by hand
- Pin exactly the tools and the router family that must move together;
  every override or patch carries a comment with its exit condition
- Ask the owner before adding any dependency, with purpose, license,
  advisories, health and exit cost written down
- File names: kebab-case in the webclient and the root, snake_case in the Go
  services, as `validate-filenames` enforces; every `.ts` and `.tsx` file
  starts with the license header
- Do not swap or re-wire an existing tool without the owner's approval

## References

| File                                                          | Read when                                                              |
| ------------------------------------------------------------- | ---------------------------------------------------------------------- |
| [commands-and-scripts.md](references/commands-and-scripts.md) | Running a check, git hooks, adding a make target or a script           |
| [dependencies.md](references/dependencies.md)                 | Adding or upgrading a dependency, a patch, an override or a pin        |
| [repo-conventions.md](references/repo-conventions.md)         | Naming files, license headers, placing a package, module or document   |
