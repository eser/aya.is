# versions
FROM node:25-trixie-slim AS base

# builder stage
FROM base AS builder

RUN apt-get update && apt-get -y install --no-install-recommends \
  ca-certificates \
  curl

# Install pnpm
RUN curl -fsSL https://get.pnpm.io/install.sh | \
  ENV="$HOME/.bashrc" SHELL=/bin/bash bash -

# Set pnpm in PATH
ENV PNPM_HOME="/root/.local/share/pnpm"
ENV PATH="$PNPM_HOME:$PATH"

WORKDIR /app

# Copy package files
COPY package.json pnpm-lock.yaml ./

# Install dependencies
RUN --mount=type=cache,id=pnpm,target=/root/.local/share/pnpm/store \
  pnpm install --frozen-lockfile

# Copy source code
COPY . .

# Build the application
RUN node --run build

# runner stage
FROM base AS runner

WORKDIR /app

# Copy built output
COPY --from=builder /app/.output ./.output
COPY --from=builder /app/package.json ./

# Expose port
EXPOSE 3000
ENV PORT=3000

# Start the server
CMD ["node", "--run", "start"]
