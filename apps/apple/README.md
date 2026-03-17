# AYA - Apple Apps (iOS & macOS)

Native iOS and macOS client for the [AYA Open Software Network](https://aya.is), built with SwiftUI and Swift 6.0.

## System Design

```
┌─────────────────────────────────────────────────────────────┐
│                        App Layer                            │
│  ┌──────────────┐              ┌──────────────┐             │
│  │  AYA-iOS     │              │  AYA-macOS   │             │
│  │  @main App   │              │  @main App   │             │
│  │  RootView    │              │  ContentView │             │
│  └──────┬───────┘              └──────┬───────┘             │
│         └──────────────┬──────────────┘                     │
│                        ▼                                    │
├─────────────────────────────────────────────────────────────┤
│                    AYAKit (SPM Library)                     │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐    │
│  │                  UI Layer                            │    │
│  │  ┌────────────┐  ┌─────────────┐  ┌──────────────┐ │    │
│  │  │  Screens   │  │ Components  │  │   Tokens     │ │    │
│  │  │            │  │             │  │              │ │    │
│  │  │ Feed       │  │ StoryCard   │  │ Colors       │ │    │
│  │  │ StoryDetail│  │ ActivityCard│  │ Typography   │ │    │
│  │  │ Activity   │  │ ProfileCard │  │ Spacing      │ │    │
│  │  │ Profile    │  │ SearchBar   │  │ CornerRadius │ │    │
│  │  │            │  │ FilterChips │  │ Animation    │ │    │
│  │  └─────┬──────┘  │ Skeleton   │  └──────────────┘ │    │
│  │        │         │ RemoteImage│                    │    │
│  │        │         └────────────┘                    │    │
│  └────────┼───────────────────────────────────────────┘    │
│           ▼                                                 │
│  ┌────────────────────┐    ┌──────────────────────────┐    │
│  │  Models             │    │  Networking              │    │
│  │                     │    │                          │    │
│  │  Profile, Story,    │◄───│  APIClient               │    │
│  │  Activity, Page,    │    │  APIClientProtocol       │    │
│  │  SearchResult,      │    │  APIError                │    │
│  │  Spotlight          │    │  LocaleHelper            │    │
│  └─────────────────────┘    └──────────┬───────────────┘    │
│                                        ▼                    │
├─────────────────────────────────────────────────────────────┤
│                  External Services                          │
│                  https://api.aya.is                         │
└─────────────────────────────────────────────────────────────┘
```

## Architecture

**MVVM** (Model-View-ViewModel) with a single SPM package (`AYAKit`):

```
Sources/
├── Models/         Domain models (Profile, Story, Activity, etc.)
├── Networking/     API client, error types, locale detection
└── UI/
    ├── Tokens/     Design system (colors, typography, spacing, etc.)
    ├── Components/ Reusable UI components (cards, search bar, etc.)
    └── Screens/    Feature screens with View + ViewModel pairs
        ├── Feed/
        ├── StoryDetail/
        ├── ActivityDetail/
        └── ProfileDetail/
```

### Key Patterns

- **Direct instantiation** — ViewModels receive concrete dependencies via init
- **@Observable ViewModels** — Swift Observation framework with `@MainActor` isolation
- **Protocol-based networking** — `APIClientProtocol` for testability, `APIClient` as concrete implementation
- **Responsive layout** — single-column on mobile, 2-column grid on wide screens (700pt+)
- **13-language localization** — with `String(localized:)` and `.strings` files
- **Full accessibility** — VoiceOver labels, hints, traits, Dynamic Type, and RTL support

### Targets

| Target            | Platform | Bundle ID            | Min OS     |
| ----------------- | -------- | -------------------- | ---------- |
| AYA-iOS           | iPhone   | `is.aya.ios`         | iOS 18.0   |
| AYA-macOS         | Mac      | `is.aya.mac`         | macOS 15.0 |
| AYA-iOS-UITests   | iPhone   | `is.aya.ios.uitests` | iOS 18.0   |
| AYA-macOS-UITests | Mac      | `is.aya.mac.uitests` | macOS 15.0 |

## Prerequisites

- Xcode 16.0+
- [XcodeGen](https://github.com/yonaskolb/XcodeGen): `brew install xcodegen`
- (Optional) [Fastlane](https://fastlane.tools): `bundle install`

## Quick Start

```bash
# Generate Xcode project and build
make project
make build-ios
make build-macos

# Run tests
make test           # SPM unit + snapshot tests
make test-ios       # Xcode iOS simulator tests
make test-macos     # Xcode macOS tests
make uitest-ios     # UI tests on iOS simulator
make uitest-macos   # UI tests on macOS

# Lint
make lint           # Check code conventions
make check          # Lint + tests
```

## Make Commands

| Command             | Description                            |
| ------------------- | -------------------------------------- |
| `make project`      | Generate `.xcodeproj` via XcodeGen     |
| `make build-ios`    | Build iOS target for simulator         |
| `make build-macos`  | Build macOS target                     |
| `make build`        | Build both targets                     |
| `make test`         | Run SPM unit + snapshot tests          |
| `make test-ios`     | Run tests via Xcode (iOS simulator)    |
| `make test-macos`   | Run tests via Xcode (macOS)            |
| `make uitest-ios`   | Run UI tests (iOS simulator)           |
| `make uitest-macos` | Run UI tests (macOS)                   |
| `make lint`         | Verify code conventions                |
| `make check`        | Lint + test                            |
| `make format`       | Format Swift code                      |
| `make clean`        | Clean build artifacts and derived data |

### Fastlane

| Command                       | Description                                 |
| ----------------------------- | ------------------------------------------- |
| `make fastlane-setup`         | Install Ruby dependencies                   |
| `make fastlane-ios-build`     | Build iOS via Fastlane                      |
| `make fastlane-macos-build`   | Build macOS via Fastlane                    |
| `make fastlane-ios-test`      | Test iOS via Fastlane                       |
| `make fastlane-macos-test`    | Test macOS via Fastlane                     |
| `make fastlane-ios-archive`   | Archive iOS for App Store                   |
| `make fastlane-macos-archive` | Archive macOS for App Store                 |
| `make fastlane-beta`          | Upload to TestFlight (requires credentials) |

## API

All data is fetched from `https://api.aya.is`. Endpoints:

- `GET /{locale}/stories` — Paginated stories
- `GET /{locale}/stories/{slug}` — Story detail
- `GET /{locale}/activities` — Paginated activities
- `GET /{locale}/activities/{slug}` — Activity detail
- `GET /{locale}/profiles?filter_kind=` — Paginated profiles
- `GET /{locale}/profiles/{slug}` — Profile detail
- `GET /{locale}/profiles/{slug}/pages` — Profile pages
- `GET /{locale}/profiles/{slug}/stories` — Profile stories
- `GET /{locale}/search?q=` — Full-text search

## Features

- **Feed** — Browse stories, activities, people, and products with filter chips
- **Search** — Full-text search with debounced input (400ms)
- **Story Detail** — Rich markdown content with author info
- **Activity Detail** — Event details with join/RSVP links
- **Profile Detail** — Tabbed view (pages, stories, activities)
- **Dark/Light mode** — Toggle with AppStorage persistence
- **Language switching** — 13 languages with live reload
- **Infinite scroll** — Cursor-based pagination
- **Pull to refresh** — On all feed content
- **Responsive grid** — Adaptive 1/2 column layout
- **Accessibility** — VoiceOver, Dynamic Type, color blind friendly, RTL

## Localization

Strings files are in `Resources/{platform}/{locale}.lproj/Localizable.strings`. Supported languages:

| Language             | Code      | RTL |
| -------------------- | --------- | --- |
| English              | `en`      | No  |
| Türkçe               | `tr`      | No  |
| Français             | `fr`      | No  |
| Deutsch              | `de`      | No  |
| Español              | `es`      | No  |
| Português (Portugal) | `pt-PT`   | No  |
| Italiano             | `it`      | No  |
| Nederlands           | `nl`      | No  |
| 日本語               | `ja`      | No  |
| 한국어               | `ko`      | No  |
| Русский              | `ru`      | No  |
| 简体中文             | `zh-Hans` | No  |
| العربية              | `ar`      | Yes |

## Testing

- **Unit Tests** — Model decoding, API client, locale helper (Swift Testing)
- **Snapshot Tests** — Visual regression for all card components and utility views (swift-snapshot-testing)
- **RTL Snapshot Tests** — Visual regression for Arabic/RTL layouts
- **UI Tests** — End-to-end tests for feed, navigation, search, filters, detail screens, and accessibility

## License

MIT - See [LICENSE](../../LICENSE) for details.
