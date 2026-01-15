# ------------------------------------------------------------
# Versions
# ------------------------------------------------------------

FROM denoland/deno:2.5.1 AS upstream-builder
FROM denoland/deno:distroless-2.5.1 AS upstream-runner

# Create a minimal image base-debian12 or static-debian12
# (see: https://github.com/GoogleContainerTools/distroless#why-should-i-use-distroless-images)

## ------------------------------------------------------------
## Base stage for both development and production
## ------------------------------------------------------------
FROM upstream-builder AS base

# Install base system dependencies
RUN apt-get update && \
  apt-get -y install --no-install-recommends curl ca-certificates

# Install node.js 24
RUN curl -fsSL https://deb.nodesource.com/setup_24.x -o nodesource_setup.sh
RUN sh nodesource_setup.sh
RUN apt-get install -y nodejs

ENV NEXT_TELEMETRY_DISABLED=1
ENV DENO_DIR=./.deno_cache

# Set working directory
WORKDIR /srv/app

## ------------------------------------------------------------
## Development stages
## ------------------------------------------------------------
FROM base AS development-runner

EXPOSE 3000

ENV NODE_ENV=development
ENV HOSTNAME="0.0.0.0"
ENV PORT=3000

CMD ["task", "dev"]

## ------------------------------------------------------------
## Production stages
## ------------------------------------------------------------
FROM base AS production-builder

# Install dependencies first to improve caching
COPY package.json deno.lock ./
RUN deno install --allow-scripts

# Copy the rest of the application
COPY . ./

# Create necessary directories and set permissions in builder stage
RUN mkdir -p .next
RUN mkdir -p .deno_cache

RUN deno task build

# Production image, copy all the files and run next
FROM upstream-runner AS production-runner

# Set working directory
WORKDIR /srv/app

ENV NODE_ENV=production
ENV HOSTNAME="0.0.0.0"
ENV PORT=3000

# # Automatically leverage output traces to reduce image size
# # https://nextjs.org/docs/app/api-reference/next-config-js/output
# COPY --from=production-builder --chown=deno:deno /srv/app/.next/standalone ./
# COPY --from=production-builder --chown=deno:deno /srv/app/.next/static ./.next/static
# COPY --from=production-builder --chown=deno:deno /srv/app/public ./public
COPY --from=production-builder /srv/app/ ./

EXPOSE 3000

# server.js is created by next build from the standalone output
# https://nextjs.org/docs/pages/api-reference/next-config-js/output
# CMD ["run", "--allow-all", "server.js"]

CMD ["task", "start"]
