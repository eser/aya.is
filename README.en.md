# AYA - Open Source Network

[![Discord](https://img.shields.io/discord/1072074800622739476?color=7289da&logo=discord&logoColor=white)](https://discord.gg/itdepremyardim)
[![GitHub issues](https://img.shields.io/github/issues/eser/aya.is)](https://github.com/eser/aya.is/issues)

(For Turkish please click [here](README.md))

We are a close-knit and passionate community of individuals who share a common interest in open source software and open
source data. United by our dedication to leveraging technology for the betterment of society, we joined forces following
the devastating [February 2023 earthquake](https://en.wikipedia.org/wiki/2023_Turkey%E2%80%93Syria_earthquake) in
Turkey.

In the aftermath of the earthquake, we recognized the pressing need to support those affected and assist in the process
of reuniting families and loved ones. With this purpose in mind, we set out to develop innovative tools and solutions
aimed at facilitating the search and connection between individuals.

Through tireless collaboration and a collective effort, our community has grown to an inspiring size of 24,000
individuals. Together, we continue to refine and expand our tools, harnessing the power of technology to help reunite
families, restore hope, and bring solace to those impacted by the earthquake.

As we forge ahead, our commitment remains unwavering. We strive to enhance our tools, actively contribute to open source
projects, and foster a supportive environment where knowledge sharing and collaboration thrive. We are proud to be a
part of the open source community and look forward to continuing our journey with you.

## Our Mission

We are actively engaged in utilizing open-source solutions, applying information systems, and implementing engineering
practices to contribute to the betterment of the society we reside in. Our primary focus lies in fulfilling social
responsibility needs and addressing various societal challenges through these means. By leveraging our expertise and
resources, we strive to make a positive impact and foster sustainable development within our community.

## Technology

This is a Docker Compose-based monorepo project. It includes the following components:

- **Frontend (webclient)**: Built with Next.js and Shadcn UI
- **Backend (services)**: REST API services written in Go
- **Database**: PostgreSQL

Prerequisites:

- [Docker](https://docker.com) (Orbstack is recommended)
- [Make](https://www.gnu.org/software/make/) (usually pre-installed on Unix/macOS systems)
- [Git](https://git-scm.com/)

## Setup and Getting Started

Clone the GitHub repository:

```bash
$ git clone git@github.com:eser/aya.is.git
$ cd aya.is
```

Set up the required environment variables for the backend service:

```bash
# Create .env.local file (this file is gitignored)
$ cat > apps/services/.env.local << EOF
AUTH__JWT_SECRET=any-secret-key-here
EOF
```

> **Note:** `AUTH__JWT_SECRET` is required for authentication. You can use any value for local development.

Use Make commands to start the project:

```bash
# Build and start all services
$ make up
```

Other useful Make commands:

```bash
$ make help      # Show all commands
$ make logs      # Show container logs
$ make stop      # Stop services
$ make restart   # Restart services
$ make down      # Remove containers completely
```

## Development with Nix (Optional)

The project includes a [Nix](https://nixos.org/download/) flake that provides all development tools. With this approach,
only PostgreSQL runs in Docker while the backend and frontend run natively.

Activate the Nix development environment:

```bash
$ nix develop
```

> **Tip:** If you use [direnv](https://direnv.net/), run `direnv allow` and the shell will activate automatically.

The Nix shell provides Go 1.24, Node.js 20, Deno, pnpm, sqlc, air, golangci-lint, and all other required tools.

Set up the environment variables for the backend service (the default database connection to `localhost:5432` is
already defined in `config.json`):

```bash
$ cat > apps/services/.env.local << EOF
AUTH__JWT_SECRET=any-secret-key-here
EOF
```

Start the services in three separate terminals:

```bash
# 1. PostgreSQL
$ docker compose up -d postgres

# 2. Backend (hot-reload)
$ cd apps/services && make dev

# 3. Frontend
$ cd apps/webclient && deno task dev
```

The API will be available at `http://localhost:8080`, the frontend at `http://localhost:3000`.

## Project Management and CLI

To connect to the service container:

```bash
$ make cli
```

This command connects you to the backend service's bash shell. From here you can perform database management and other
operations.

### Example

Getting a profile:

```js
await backend.getProfile("en", "eser");
```

## Data Model

The project uses the following main data structures:

- **Profile**: User profiles (individuals, organizations, communities)
- **Story**: Content and articles (blog posts, news, events)
- **User**: System users and authentication information

## How to Contribute

We welcome contributions from everyone. To start please [read our contributing guide](CONTRIBUTING.en.md). If you want
to help you can check out our [issues](https://github.com/eser/aya.is/issues). If you have any questions, feel free to
join our [Discord server](https://discord.gg/itdepremyardim). If you are stuck at any point, feel free to ask for help
on GitHub Issues or Discord.

## License

Apache 2.0, see [LICENSE](LICENSE) file for details.
