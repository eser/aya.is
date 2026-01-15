# ------------------------------------------------------------
# Monorepo Services Dockerfile
# Build context: repository root (not apps/services/)
# Usage: docker build -f apps/services/Dockerfile .
# ------------------------------------------------------------

# ------------------------------------------------------------
# Versions
# ------------------------------------------------------------

FROM golang:1-bookworm AS upstream-builder
# FROM gcr.io/distroless/base-debian12 AS upstream-runner
FROM debian:bookworm-slim AS upstream-runner

# Create a minimal image base-debian12 or static-debian12
# (see: https://github.com/GoogleContainerTools/distroless#why-should-i-use-distroless-images)

## ------------------------------------------------------------
## Base stage for both development and production
## ------------------------------------------------------------
FROM upstream-builder AS base

# Configuration
ARG GH_LOGIN
ENV GH_LOGIN=$GH_USER

ARG GH_ACCESS_TOKEN
ENV GH_ACCESS_TOKEN=$GH_ACCESS_TOKEN

ARG GH_PATH
ENV GH_PATH=$GH_PATH

# Setup for private Go Modules where available
RUN echo "machine github.com login ${GH_LOGIN} password ${GH_ACCESS_TOKEN}" > ~/.netrc
ENV GOPRIVATE="${GH_PATH}/*"
RUN go env -w GOPRIVATE="${GH_PATH}/*"
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux

# Set working directory to monorepo root
WORKDIR /srv/monorepo

## ------------------------------------------------------------
## Development stages
## ------------------------------------------------------------
FROM base AS development-runner

WORKDIR /srv/monorepo/apps/services

EXPOSE 8080

ENTRYPOINT ["go", "run", "./cmd/serve/"]

## ------------------------------------------------------------
## Production stages
## ------------------------------------------------------------
FROM base AS production-builder

# Copy package dependencies first (for better layer caching)
# Future: If using shared packages, copy packages/ here first
# COPY packages/ ./packages/

# Copy services module files for dependency download
COPY apps/services/go.mod apps/services/go.sum ./apps/services/

# Download dependencies
WORKDIR /srv/monorepo/apps/services
RUN go mod download && go mod verify

# Copy the rest of the services application
WORKDIR /srv/monorepo
COPY apps/services/ ./apps/services/

# Build the application
WORKDIR /srv/monorepo/apps/services
RUN go build ./cmd/serve/


# Production image, copy all the files and run next
FROM upstream-runner AS production-runner

# Install base system dependencies
RUN apt-get update && \
  apt-get -y install --no-install-recommends curl ca-certificates && \
  groupadd -r nonroot && \
  useradd -r -g nonroot nonroot

# Set working directory
WORKDIR /srv/app

# Copy the binary from the production-builder container
COPY --from=production-builder --chown=nonroot:nonroot /srv/monorepo/apps/services/serve ./
COPY --from=production-builder --chown=nonroot:nonroot /srv/monorepo/apps/services/.env ./
COPY --from=production-builder --chown=nonroot:nonroot /srv/monorepo/apps/services/config.json ./

COPY apps/services/etc/data/default/migrations ./etc/data/default/migrations
COPY apps/services/etc/data/default/seed/seed.sql ./etc/data/default/seed/seed.sql

# Run as a non-root user
USER nonroot

EXPOSE 8080

ENTRYPOINT ["./serve"]
