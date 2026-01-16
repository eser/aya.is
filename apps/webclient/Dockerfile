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
# Development stage (with hot reload)
# ------------------------------------------------------------
FROM base AS development

RUN apt-get update && apt-get -y install --no-install-recommends \
  ca-certificates \
  curl \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /monorepo

# Copy root deno.json (shared lint/fmt settings)
COPY ./deno.json ./

# Install root dependencies (if any)
RUN deno install

WORKDIR /monorepo/apps/webclient

# Copy webclient package files (no lock file - let Deno resolve all transitive deps)
COPY ./apps/webclient/package.json ./apps/webclient/deno.json ./

# Install webclient dependencies
RUN deno install

# Copy source code
COPY ./apps/webclient/ ./

EXPOSE 3000

# Development server with hot reload (uses package.json scripts)
CMD ["deno", "task", "dev"]

# ------------------------------------------------------------
# Builder stage
# ------------------------------------------------------------
FROM base AS builder

RUN apt-get update && apt-get -y install --no-install-recommends \
  ca-certificates \
  curl \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /monorepo

# Copy root deno.json (shared lint/fmt settings)
COPY ./deno.json ./

# Install root dependencies (if any)
RUN deno install

WORKDIR /monorepo/apps/webclient

# Copy webclient package files
COPY ./apps/webclient/package.json ./apps/webclient/deno.json ./

# Install webclient dependencies
RUN deno install

# Copy source code
COPY ./apps/webclient/ ./

# Build the application
RUN deno task build

# ------------------------------------------------------------
# Production runner (Deno)
# ------------------------------------------------------------
FROM base AS runner

# Run as non-root user (Deno image provides 'deno' user)
USER deno

WORKDIR /app

# Copy built output from builder (self-contained, no package.json needed)
COPY --from=builder /monorepo/apps/webclient/.output ./

EXPOSE 3000
ENV PORT=3000
ENV HOST=::

# Start the Nitro server directly
CMD ["deno", "task", "start"]
