# Deployments

How the Compose files and Coolify run AYA: the development and production-like
stacks in the repository, the Docker images they build, and what is known about
the Coolify deployment. The repository holds no Kubernetes manifests, kubectl
usage or ArgoCD configuration; the old habit of reaching for `kubectl` does not
apply here.

## Contents

- Compose Stacks
- Coolify Runs the Production Application
- Images and Pins
- Inspect Before Acting
- Volumes and the PostgreSQL Major

---

## Compose Stacks

Scope: Running the stack locally or checking what a deploy builds

Rule: Two Compose files at the repository root, three services each:
`webclient` (port 3000), `services` (port 8080) and `postgres` (port 5432).

| File                     | Project name           | Targets built                                          | Notes                                                    |
| ------------------------ | ---------------------- | ------------------------------------------------------ | -------------------------------------------------------- |
| `compose.yml`            | `aya-is-development`   | `webclient`: `development`; `services`: `development-runner` | `develop.watch` syncs sources; `docker compose watch` |
| `compose.production.yml` | `aya-is-production`    | `webclient`: `runner`; `services`: `production-runner` | `ENV=production`, `LOG__LEVEL=WARN`; built from the same Dockerfiles |

The root Makefile wraps the development stack: `make up`, `make down`,
`make build`, `make restart`, `make logs`, `make cli` (a shell in `services`).
Run them only for the local stack; they are not a deployment. The database
settings and password in both files (`s3cr3t`, `sslmode=disable`) are local
values and never copied to a real environment. The `VITE_*` variables are build
arguments baked into the webclient bundle; `BACKEND_URI` is a runtime variable
of its server side (the comments in the compose files say so).

---

## Coolify Runs the Production Application

Scope: Reasoning about, or checking, the live deployment

Rule: Production runs on Coolify. Its application settings (source, build
target, environment variables, domains, deploy triggers) live in Coolify, not in
the repository, so the repository cannot tell you whether a push or a tag
redeploys. Do not claim a version is or is not deployed without evidence (AGENTS.md:
Remember): look at the live site and at Coolify first.

What the repository does hold:

- `apps/services/pkg/api/adapters/coolify` is a client of the Coolify HTTP API,
  configured under the `coolify` key of the services configuration
  (`api_url`, `api_token`, `webclient_app_uuid`, `request_timeout`); it reads
  and updates the webclient application's domain list and can restart that
  application. The token is a secret and stays in the environment.
- The `postgres` comment in both compose files: the image is the same as
  production's.
- `.agents/settings.json` allows these read-only commands, which are the way to
  look: `coolify app list`, `coolify app deployments list`, `coolify app logs`,
  `coolify deploy list`.

A deploy, restart, rollback or settings change in Coolify is the owner's action.
Report what you found instead of working around it.

---

## Images and Pins

Scope: `apps/services/Dockerfile`, `apps/webclient/Dockerfile`, the compose files,
`flake.nix`

Rule: The base images are pinned by tag, not by digest, and Dependabot's weekly
`docker` run moves them (github-actions.md: Dependabot):

| Image                       | Where                                          |
| --------------------------- | ---------------------------------------------- |
| `golang:1.27-bookworm`      | Builder stage of `apps/services/Dockerfile`    |
| `debian:bookworm-slim`      | Runner stage of `apps/services/Dockerfile`     |
| `node:24-bookworm-slim`     | Base and runner stages of `apps/webclient/Dockerfile` |
| `pgvector/pgvector:pg18`    | `postgres` in both compose files               |

Move a pin by changing every place that names the same thing in one change: the
Go minor with `go.mod` and `integrity.yml`, Node with `engines` and
`integrity.yml` (github-actions.md: Toolchain Pins), and PostgreSQL 18 in both
compose files and in `flake.nix` (`postgresql_18`), so a developer's shell and
the stacks run the same major. Production runner stages run as a non-root user
(`nonroot` for services, `node` for the webclient) and add only `curl` (and
`ca-certificates` for services); keep it that way.

---

## Inspect Before Acting

Scope: Debugging or checking a deployment or a local stack

Rule: Look before changing anything, with read-only commands:

```sh
docker compose ps
docker ps --all
docker logs <container> --tail 100
coolify app list
coolify app deployments list <app-uuid>
coolify app logs <app-uuid>
```

`services` must reach `postgres` through its healthcheck
(`pg_isready -U postgres -d postgres`) before it starts, because it
`depends_on: condition: service_healthy`; a `services` container that exits at
start is read from its logs, not restarted in a loop. A restart, recreate or
volume change on a running deployment is an action the owner confirms first.

Correct:

```sh
docker logs aya-is-development-services-1 --tail 100   # read what failed first
```

Incorrect:

```sh
docker compose down -v && docker compose up -d   # destroys the database volume to "fix" a startup error
```

---

## Volumes and the PostgreSQL Major

Scope: The `postgres-data` volume and any change to the PostgreSQL image

Rule: PostgreSQL 18 images keep data under `/var/lib/postgresql/<major>/`, so the
compose files mount `postgres-data` at `/var/lib/postgresql`, one level up. A
volume created by an older major (the earlier `postgres:16` setup) is refused at
start. The two ways out are `docker compose down -v`, which deletes the local
data, or a dump and restore before the switch (the comment in the compose files
says both). Choosing either on a stack that holds data worth keeping is the
owner's decision; dump first. A change of the PostgreSQL major is announced in
the changelog entry of its release (release-process.md: Compatibility Across
Versions).
