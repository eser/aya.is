# Release Process

What an AYA release is, which files and commands make one, what each pushed tag
triggers, and how a version stays compatible with the one before it. The scripts
are `scripts/push_current_release_tag.ts` and
`scripts/sync_release_notes_from_changelog.ts`; the workflows are in
github-actions.md.

## Contents

- What a Release Is
- Versions
- Release Commands
- Tags Trigger Three Workflows
- Release Notes Come from the Changelog
- Tag Patterns and Environments
- Failed Releases
- Compatibility Across Versions
- Commit Messages

---

## What a Release Is

Scope: Any request to release, tag or ship AYA

Rule: A release is a Git tag plus a GitHub Release whose notes are a section of
`CHANGELOG.md`. The repository builds and publishes no package or image from a
release: no workflow pushes an image, and the `Deployment Pipeline` only
classifies the tag (Tags Trigger Three Workflows). The Apple apps have their own
checks (`apple.yml`) and no release flow here; the local tag `apple-apps-v0.1.0`
matches none of the `v*` workflow triggers.

---

## Versions

Scope: Choosing and recording a version

Rule: The owner chooses the number. A release version is three numbers,
`x.y.z`, kept in three places that must agree:

| Place          | Holds                                                  |
| -------------- | ------------------------------------------------------ |
| `VERSION`      | `x.y.z` alone, one line (`0.1.0` today)                |
| Git tag        | `vx.y.z`, annotated, made by `make tag` from `VERSION` |
| `CHANGELOG.md` | A heading `## x.y.z - YYYY-MM-DD` above that release's notes |

The sync script matches a heading by `v` plus its version, so the tag `v0.2.0`
needs the heading `## 0.2.0 - 2026-..-..` and nothing else will do.

Correct:

```markdown
## 0.2.0 - 2026-12-01

### Added

- Privacy settings with data export, hide and account deletion.
```

Incorrect:

```markdown
## Version 0.2.0 (December)    <!-- heading the script cannot parse -->
```

**Why:** the heading, the file and the tag are read by different tools; a
mismatch shows up only after the tag is pushed, when it is public.

---

## Release Commands

Scope: Turning a committed `VERSION` and changelog entry into a tag. An agent
runs these only when the owner asks for that exact action.

Rule: Use the root Make targets:

```sh
make tag       # node scripts/push_current_release_tag.ts
make release   # node scripts/sync_release_notes_from_changelog.ts --repo eser/aya.is
```

`make tag` reads `VERSION`, creates the annotated tag `vX.Y.Z` when it does not
exist locally, refuses when that tag points at a commit other than `HEAD`
("create a new release commit before pushing this tag"), runs
`git push origin HEAD`, then pushes the tag unless `origin` already has it. It
therefore publishes the branch as well as the tag. The commit that sets
`VERSION` and the changelog entry is made and reviewed before the target runs.

`make release` syncs the newest changelog entry to an existing GitHub Release
with `gh release edit`; without `--create-if-missing` it creates nothing and
prints that it skipped. It needs an authenticated `gh` (`gh auth status`).

---

## Tags Trigger Three Workflows

Scope: What happens after a `v*` tag reaches `origin`

Rule: One tag push starts these runs, each with its own concurrency group:

| Workflow                | Trigger pattern              | Does                                                        |
| ----------------------- | ---------------------------- | ----------------------------------------------------------- |
| `Integrity Pipeline`    | `v[0-9]+.[0-9]+.[0-9]+*`     | The same checks as a branch push (github-actions.md)         |
| `Deployment Pipeline`   | `v[0-9]+.[0-9]+.[0-9]+*`     | Classifies the tag and prints the variables (next sections)  |
| `Release Notes Sync`    | `v*`                         | Creates or updates the GitHub Release from the changelog    |

Read all three runs (`gh run list --limit 10`) before calling a release done.
The `Deployment Pipeline` job has no build, push or deploy step; do not report a
green run of it as "deployed". How a running Coolify application picks up new
code is configured in Coolify, not in the repository (deployments.md).

---

## Release Notes Come from the Changelog

Scope: `CHANGELOG.md` and `.github/workflows/release_notes_sync.yml`

Rule: The sync job runs on every `v*` tag, on a push to `main` that touches
`CHANGELOG.md`, and by hand (`workflow_dispatch` with `tag` and
`create_if_missing`). The script takes the section from its heading to the next
heading and writes it as the release body.

- On a tag push it passes that tag and `--create-if-missing`, so it creates the
  release `AYA vX.Y.Z` (`gh release create --verify-tag`); if a concurrent run
  created it first, it edits instead.
- On a `main` push with no tag it targets the top changelog entry and only edits
  a release that already exists; an entry without a release is skipped.
- A tag without a matching heading fails with "No matching changelog section
  found". A pre-release tag such as `v0.2.0-rc.1` needs a heading `## 0.2.0-rc.1`
  of its own, or its sync run fails even though the tag is valid.

The changelog is the source of the body. To fix the body of the newest release,
change `CHANGELOG.md` on `main`; for an older version, dispatch the workflow with
that `tag`, because a `main` push only targets the top entry.

---

## Tag Patterns and Environments

Scope: `.github/workflows/deployment.yml`, choosing a tag form

Rule: The deployment workflow derives an environment from the tag and fails on
any other form:

```text
v1.0.0          production
v1.0.0-rc.1     stage        (also v1.0.0-rc)
v1.0.0-dev.1    test         (also v1.0.0-dev)
```

It then builds `IMAGE_NAME` as `gcr.io/<GCLOUD_PROJECT>/<APP_NAME>:<tag>` from
the repository secrets `GCLOUD_PROJECT` and `APP_NAME` and echoes the four
variables (`TAG`, `DEPLOY_ENV`, `IMAGE_NAME`, `IMAGE_TAG`). Nothing consumes
them yet. A tag of any other shape, such as `v1.0.0-beta`, passes the trigger
glob and then fails the job.

**Why:** a stable, parseable tag grammar makes the environment a function of the
tag, so a rollback is a new tag and never a moved one.

---

## Failed Releases

Scope: A release run or sync that stopped

Rule: Read the failure before retrying: `gh run view <run-id> --log-failed`.

- A failed `Integrity Pipeline` on a tag means the tagged commit is red. The tag
  stays; fix `main` and cut the next version.
- A failed `Release Notes Sync` is repaired by correcting the changelog heading
  and re-running it with `workflow_dispatch` (`tag` and `create_if_missing`),
  not by retagging.
- A tag that points at the wrong commit is not deleted and re-pushed once others
  may have fetched it: fix forward with a new version. `make tag` refuses a
  local tag that differs from `HEAD` for the same reason.

---

## Compatibility Across Versions

Scope: A change that alters a stored shape, an API contract, or a setting that a
running deployment needs

Rule: A new build meets data and clients that an older build left behind.

- Database changes ship as goose migrations with both `-- +goose Up` and
  `-- +goose Down` sections in `apps/services/etc/data/default/migrations/`
  (`apps/services/AGENTS.md`). The production image copies that directory next
  to the binary.
- A new required setting (an environment variable in the Compose files or in
  Coolify) is named in the changelog entry, so whoever deploys can set it first.
- A change to a response or contract that the webclient or the Apple apps read
  keeps the old shape working until the clients that need the new one are out.

---

## Commit Messages

Scope: Messages for the owner's commits

Rule: Commits are conventional one-liners (`feat:`, `fix:`, `chore:`, `ci:`,
`docs:`, `perf:`, `refactor:`, `revert:`, `test:`), enforced by the
`commit-validation` workflow in `.manifest.yml`. The changelog is the human
record of a release, written by the owner per version; a commit subject records
one change. workflow-practices: Git Commit Policy holds the format and the rule
that agents never commit.
