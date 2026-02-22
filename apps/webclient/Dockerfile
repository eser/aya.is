# ------------------------------------------------------------
# Monorepo Webclient Dockerfile (TanStack Start + Vite + Nitro)
# Runtime: Deno
# Build context: repository root (not apps/webclient/)
# Usage: docker build -f apps/webclient/Dockerfile .
# ------------------------------------------------------------

# ------------------------------------------------------------
# Base image
# ------------------------------------------------------------
FROM denoland/deno:2.6.5 AS base

# ------------------------------------------------------------
# Dependencies stage (shared by development and builder)
# ------------------------------------------------------------
FROM base AS deps

RUN apt-get update && apt-get -y install --no-install-recommends \
  ca-certificates \
  curl \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /monorepo

# Copy root deno.json (shared lint/fmt settings)
COPY ./deno.json ./

# Install root dependencies (if any)
RUN --mount=type=cache,target=/root/.cache/deno \
    deno install

WORKDIR /monorepo/apps/webclient

# Copy webclient package files (no lock file - let Deno resolve all transitive deps)
COPY ./apps/webclient/package.json ./apps/webclient/deno.json ./

# Install webclient dependencies (cached across builds)
RUN --mount=type=cache,target=/root/.cache/deno \
    deno install

# ------------------------------------------------------------
# Development stage (with hot reload)
# ------------------------------------------------------------
FROM deps AS development

# Copy source code
COPY ./apps/webclient/ ./

EXPOSE 3000

# Development server with hot reload (uses package.json scripts)
CMD ["deno", "task", "dev"]

# ------------------------------------------------------------
# Builder stage
# ------------------------------------------------------------
FROM deps AS builder

# Copy source code
COPY ./apps/webclient/ ./

# Build the application (with Deno cache)
RUN --mount=type=cache,target=/root/.cache/deno \
    deno task build

# ------------------------------------------------------------
# Production runner (Deno)
# ------------------------------------------------------------
FROM base AS runner

# Run as non-root user (Deno image provides 'deno' user)
USER deno

WORKDIR /app

# Copy built output from builder (self-contained Nitro server)
COPY --from=builder /monorepo/apps/webclient/.output ./

EXPOSE 3000
ENV PORT=3000
ENV HOST=::

# Start the Nitro server directly
CMD ["deno", "run", "-A", "server/index.mjs"]
