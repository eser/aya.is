---
name: release-and-ci
description: "The release flow (VERSION, CHANGELOG.md, make tag, make release, v* tags), the GitHub Actions workflows in .github/workflows and the Compose and Coolify deployment of AYA. Use when cutting a release, choosing a version, editing a workflow, debugging a red CI run or checking a deployment. Not for which tests a workflow runs (use qa-practices)."
---

# Releases and CI

How a version is cut and tagged, how the workflows in `.github/workflows/` are
written and what they do, and how the Compose files and Coolify run AYA. Which
tests each check runs is in qa-practices: Test Hooks and Workflows.

## Always

- The owner releases and commits. Run `make tag` or `make release`, push a `v*`
  tag or dispatch a workflow only when the owner asks for that exact action;
  agents prepare the change and propose a one-line commit message
  (workflow-practices: Git Commit Policy).
- A version is `x.y.z`. It lives in `VERSION`; the tag is `vx.y.z`; the release
  notes are the `## x.y.z - YYYY-MM-DD` section of `CHANGELOG.md`. The three
  name the same version or the release job fails.
- A pushed tag is never moved or deleted, and a published release is never
  replaced: fix the cause and cut the next version.
- Workflow changes need the owner's approval (workflow-practices: Ask First).
- Every new or changed job sets `timeout-minutes`, and no action floats on a
  branch (`@main`, `@master`).
- Workflows grant `contents: read`; a write permission is declared on the one
  job that needs it. Secrets come from `secrets.*` into step `env`, are never
  printed, and pull request jobs pass without them.
- The Go patch version in `integrity.yml` is pinned on purpose (govulncheck
  fails on older standard library patches); it moves with a reason, not with
  `latest`.
- Inspect before acting: `gh run view <id> --log-failed`, `gh pr checks`,
  `docker ps`, `coolify app logs`. Use `gh`, not raw API calls with a token.
- Redeploying, restarting or changing a volume on a running deployment is the
  owner's action. The repository holds no Kubernetes, kubectl or ArgoCD
  configuration, so none of them is a tool for this project.

## References

| File                                                | Read when                                                                          |
| --------------------------------------------------- | ---------------------------------------------------------------------------------- |
| [release-process.md](references/release-process.md) | Cutting or recovering a release, choosing a version, writing a changelog entry     |
| [github-actions.md](references/github-actions.md)   | Editing a workflow or job, bumping a toolchain pin, using `gh`, debugging a CI run |
| [deployments.md](references/deployments.md)         | Running the Compose stacks, checking or reasoning about a Coolify deployment       |
