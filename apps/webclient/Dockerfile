# ------------------------------------------------------------
# Monorepo Webclient Dockerfile (TanStack Start + Vite + Nitro)
# Runtime: Node.js (pnpm workspace)
# Build context: repository root (not apps/webclient/)
# Usage: docker build -f apps/webclient/Dockerfile .
# ------------------------------------------------------------

# ------------------------------------------------------------
# Base image
# ------------------------------------------------------------
# speedy-router requires Node >= 24.
FROM node:24-bookworm-slim AS base

# pnpm is pinned through the root package.json "packageManager" field (activated by corepack)
ENV PNPM_HOME=/pnpm
ENV PATH=$PNPM_HOME:$PATH
RUN corepack enable

# ------------------------------------------------------------
# Dependencies stage (shared by development and builder)
# ------------------------------------------------------------
FROM base AS deps

RUN apt-get update && apt-get -y install --no-install-recommends \
  ca-certificates \
  curl \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /monorepo

# Workspace manifests first (better layer caching): lockfile, workspace file, registry config and package manifests
COPY ./package.json ./pnpm-lock.yaml ./pnpm-workspace.yaml ./.npmrc ./
COPY ./patches ./patches
COPY ./apps/webclient/package.json ./apps/webclient/package.json

# Install the webclient workspace package and its dependencies (store cached across builds)
RUN --mount=type=cache,id=pnpm-store,target=/pnpm/store \
    pnpm install --frozen-lockfile --filter @aya/webclient...

WORKDIR /monorepo/apps/webclient

# ------------------------------------------------------------
# Development stage (with hot reload)
# ------------------------------------------------------------
FROM deps AS development

# Copy source code
COPY ./apps/webclient/ ./

EXPOSE 3000

# Development server with hot reload (uses package.json scripts)
CMD ["pnpm", "dev"]

# ------------------------------------------------------------
# Builder stage
# ------------------------------------------------------------
FROM deps AS builder

# Copy source code
COPY ./apps/webclient/ ./

# Give Vite/Rolldown enough heap for SSR transform on memory-constrained builders
ENV NODE_OPTIONS=--max-old-space-size=4096

# Build the application
RUN pnpm build

# ------------------------------------------------------------
# Production runner (Node.js)
# ------------------------------------------------------------
FROM node:24-bookworm-slim AS runner

# curl is required by the container healthcheck (must be installed as root,
# before switching to the 'node' user). The base image ships neither curl nor wget.
RUN apt-get update && apt-get -y install --no-install-recommends \
  curl \
  && rm -rf /var/lib/apt/lists/*

# Run as non-root user (the Node image provides 'node')
USER node

WORKDIR /app

# Copy built output from builder (self-contained Nitro server)
COPY --from=builder --chown=node:node /monorepo/apps/webclient/.output ./

EXPOSE 3000
ENV PORT=3000
ENV HOST=::

# Start the Nitro server directly
CMD ["node", "server/index.mjs"]
