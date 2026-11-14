# Architecture Guide

This guide explains how this full-stack codebase is organized.

## Overview

This is a Docker Compose-based monorepo that consists of three main services:

- **Frontend (webclient)**: Next.js application serving the user interface
- **Backend (services)**: Go-based REST API services
- **Database (postgres)**: PostgreSQL database for data persistence

## Directory Structure

```
aya.is/
├── webclient/          # Frontend application
│   ├── src/           # Source code
│   ├── public/        # Static assets
│   └── Dockerfile     # Multi-stage Docker build
├── services/          # Backend services
│   ├── cmd/           # Command-line applications
│   │   ├── serve/     # HTTP server
│   │   ├── migrate/   # Database migrations
│   │   └── manage/    # Management CLI
│   ├── pkg/           # Shared packages
│   │   ├── ajan/      # Core framework
│   │   └── api/       # API implementation
│   └── Dockerfile     # Multi-stage Docker build
├── compose.yml        # Development environment
├── compose.production.yml  # Production environment
└── Makefile          # Build and deployment commands
```

## Services Architecture

### Frontend Service (webclient)

- **Framework**: Next.js with TypeScript
- **UI Library**: Shadcn UI components
- **Styling**: Tailwind CSS
- **Port**: 3000
- **Features**:
  - Server-side rendering
  - Internationalization (i18n)
  - Theme switching
  - Responsive design

### Backend Service (services)

- **Language**: Go
- **Framework**: Custom HTTP framework (ajan)
- **Port**: 8080
- **Architecture**: Clean architecture with:
  - **cmd/**: Entry points for different commands
  - **pkg/api/**: API layer with adapters and business logic
  - **pkg/ajan/**: Core framework utilities

### Database

- **Engine**: PostgreSQL 16+
- **Port**: 5432
- **Key Tables**:
  - `profile`: User profiles (individuals, organizations, communities)
  - `story`: Content and articles
  - `user`: Authentication and user accounts
  - `session`: User sessions
  - `event`: Events and activities
  - `*_tx`: Translation tables

## Development Workflow

1. **Prerequisites**: Docker, Make, Git

2. **Getting Started**:

   ```bash
   $ git clone git@github.com:eser/aya.is.git
   $ cd aya.is
   $ make up     # Start services
   ```

3. **Development Commands**:
   ```bash
   $ make logs    # View logs
   $ make cli     # Access backend shell
   $ make ok      # Run linters, formatters and tests
   ```

## Data Flow

1. **Client Request** → Frontend (Next.js)
2. **API Call** → Backend (Go services)
3. **Data Query** → PostgreSQL
4. **Response** → Backend → Frontend → Client

## Security Considerations

- Environment-based configuration
- JWT-based authentication
- CORS middleware
- Rate limiting
- SQL injection prevention via prepared statements

## Deployment

The project uses multi-stage Docker builds for both development and production:

- Development: Hot-reload enabled, verbose logging
- Production: Optimized builds, minimal logging

Use `compose.production.yml` for production deployments.
