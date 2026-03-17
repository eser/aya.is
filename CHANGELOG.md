# Changelog

## 0.1.0 - 2026-03-08

### Added

- Multi-profile system with individual, organization, and product profile types, slug-based routing, and custom domain
  support.
- Full internationalization across 13 locales (ar, de, en, es, fr, it, ja, ko, nl, pt-PT, ru, tr, zh-CN) with 3-tier
  locale fallback.
- Story management system supporting articles, announcements, and other content kinds with Markdown/MDX editing.
- Story series for grouping related stories with ordering and navigation.
- Story discussions with threaded comments.
- Story translations for multi-locale content authoring.
- Date proposals system for collaborative scheduling.
- Share wizard for sharing content across social platforms.
- Full-text search across profiles, stories, and pages with profile scoping and spotlight UI.
- Profile memberships, teams, and points system.
- Profile links with OAuth-based GitHub and YouTube integration, sync workers, and content folder support.
- Profile custom pages with drag-and-drop ordering.
- Cover editor with background image support for profile pages.
- Referral and invitation system.
- Bulletin system for community announcements.
- Mailbox for in-app messaging.
- Cookie consent overlay with configurable preferences.
- YouTube broadcast detection for live stream awareness.
- Authentication via GitHub OAuth and Apple login with session management and PoW challenges.
- MCP (Model Context Protocol) adapter with tools for profiles, stories, and news.
- WebMCP support for browser-based MCP integration.
- Auto-summarization of story content.
- Native iOS and macOS apps.
- Admin area for profile management.
- Security headers and vulnerability reporting.
- SEO improvements with metadata, reading time, and Open Graph support.

### Improved

- Migrated frontend from Next.js to TanStack Start with file-based routing and SSR via React Query
  dehydration/hydration.
- Migrated backend to Go hexagonal architecture with connfx/aifx patterns and goose migrations.
- Adopted CSS Modules with Tailwind @apply for component styling.
- Implemented 88 Deno snapshot tests for pure utility functions.
- Enhanced accessibility with skip-to-content links, ARIA labels, and semantic data-slot attributes.
- Localization improvements across all 13 supported locales.
- Client-server call optimizations and cross-domain routing enhancements.
- Multi-stage Docker builds with distroless images for production deployment.

### Fixed

- Locale switching and locale fallback reliability.
- Story card aspect ratio rendering.
- Cross-domain navigation and custom domain link handling.
- OAuth redirect and token exchange flows.
- Dropdown menu interactions and form validation.

### Notes

- This is the initial changelog entry summarizing the project's full development history. Future entries will document
  changes incrementally per release.
