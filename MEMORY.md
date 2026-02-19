# Story Architecture & Decisions

## Unified Content Model

Stories are the **central content primitive** in AYA. All content types are modeled as stories with different `kind` values.

### Story Kinds
- `status` — short updates
- `announcement` — important announcements
- `article` — long-form articles
- `news` — news items
- `content` — generic content
- `presentation` — SpeakerDeck presentations (managed, synced via RSS)
- `activity` — activities/events (meetup, workshop, conference, broadcast, meeting)

### Key Tables

| Table | Purpose |
|-------|---------|
| `story` | Core content (id, slug, kind, author_profile_id, properties JSONB, is_managed, remote_id, series_id) |
| `story_tx` | Locale-specific content (title, summary, content per locale) |
| `story_publication` | Multi-profile publishing (published_at, is_featured per profile) |
| `story_series` | Groups stories (article series, activity series, etc.) |
| `story_interaction` | Profile-to-story interactions (RSVP, likes, bookmarks) |

### Properties JSONB Pattern

Activity-specific fields live in `story.properties` rather than dedicated columns:
```json
{
  "activity_kind": "meetup",
  "activity_time_start": "2026-03-15T18:00:00Z",
  "activity_time_end": "2026-03-15T20:00:00Z",
  "external_activity_uri": "https://meetup.com/...",
  "external_attendance_uri": "https://meetup.com/.../rsvp",
  "rsvp_mode": "enabled"
}
```

**Promotion pattern:** If a JSONB field needs frequent querying/indexing, promote it to a first-class column. Examples: `remote_id` started in properties, promoted in migration 0033. `activity_time_start` has a functional JSONB index for sorting.

### RSVP Mode (for activities)
- `enabled` — in-platform RSVP buttons (attending/interested/not_attending)
- `managed_externally` — show link to `external_attendance_uri`
- `disabled` — no RSVP functionality

### Story Interactions
- `story_interaction` table: UNIQUE(story_id, profile_id, kind)
- Allows multiple interaction types per user (e.g., both "like" AND "bookmark")
- RSVP kinds (attending/interested/not_attending) are **mutually exclusive** — enforced at application layer, not DB
- Application removes existing RSVP interactions before setting new one

### Story Series
- `story_series` → `story.series_id` (nullable FK)
- Any story kind can belong to a series
- Flat structure (no nested series)

### Managed Stories
- `is_managed = true` for externally synced content (YouTube, SpeakerDeck)
- `remote_id` stores external identifier (unique per author_profile_id)
- Workers periodically sync and create/update managed stories

### Telegram Integration

Profiles can connect their Telegram accounts via the `profile_link` system (`kind='telegram'`, `is_managed=true`).

**Bot:** "ayabot" — handles account linking, status checks, and unlinking via private messages.

**Deep Link Token Flow:**
1. User clicks "Connect Telegram" in AYA web UI
2. API generates a crypto-random token (32 bytes → 64 hex chars, 10min TTL), stores in `telegram_link_token` table
3. Returns deep link: `https://t.me/ayabot?start=<token>`
4. User opens link → sends `/start <token>` to bot
5. Bot validates token → creates `profile_link` with `kind='telegram'` → confirms

**Bot Commands:** `/start`, `/start <token>`, `/help`, `/status`, `/unlink`

**Two Modes:**
- **Webhooks** (production): Telegram pushes to `POST /telegram/webhook`, verified via `X-Telegram-Bot-Api-Secret-Token`
- **Long-polling** (dev): Polling worker fetches updates — no public URL needed

**Environment Variables:**
| Variable | Default | Description |
|---|---|---|
| `TELEGRAM__ENABLED` | `false` | Enable Telegram integration |
| `TELEGRAM__BOT_TOKEN` | *(required)* | Token from @BotFather |
| `TELEGRAM__BOT_USERNAME` | `ayabot` | Bot username for deep links |
| `TELEGRAM__WEBHOOK_URL` | — | Webhook URL (production) |
| `TELEGRAM__WEBHOOK_SECRET` | — | Webhook secret header |
| `TELEGRAM__USE_POLLING` | `false` | Use polling instead of webhooks |
| `WORKERS__TELEGRAM_BOT__ENABLED` | `false` | Enable polling worker |

**Key Tables:** `telegram_link_token` (migration 0035), `profile_link` (existing, `kind='telegram'`)

### Publishing Model
- Stories have no direct `published_at` — it's on `story_publication`
- A story without active publications is a draft
- `published_at` is computed as `MIN(story_publication.published_at)`
- Stories can be published to multiple profiles

### Authorization
- Admin users can always edit
- User's individual_profile_id == author_profile_id → can edit
- User has owner/lead/maintainer membership on author profile → can edit

### Existing Activity Tables (Legacy)
- `activity`, `activity_series`, `activity_attendance` tables exist (renamed from event_* in migration 0020)
- **NOT USED** — activities are modeled as stories with kind='activity'
- Can be dropped in a future cleanup migration
