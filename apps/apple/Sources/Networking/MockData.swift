import Foundation

/// Canned JSON responses for UI testing. Each property returns a valid JSON string
/// matching the shape of the corresponding API endpoint.
enum MockData {

    // MARK: - Stories

    static let storiesResponse = """
    {
        "data": [
            {
                "id": "s1",
                "slug": "scaling-community-platforms",
                "kind": "article",
                "title": "Scaling Community Platforms",
                "summary": "Best practices for building and scaling open-source community platforms to millions of users.",
                "content": null,
                "story_picture_uri": null,
                "published_at": "2025-12-01T10:00:00Z",
                "created_at": "2025-12-01T09:00:00Z",
                "author_profile": {
                    "id": "a1",
                    "slug": "eser",
                    "title": "Eser Ozvataf",
                    "profile_picture_uri": null
                }
            },
            {
                "id": "s2",
                "slug": "open-source-sustainability",
                "kind": "article",
                "title": "Open Source Sustainability",
                "summary": "How open-source projects can achieve long-term financial sustainability while keeping their community-first values.",
                "content": null,
                "story_picture_uri": null,
                "published_at": "2025-11-15T08:00:00Z",
                "created_at": "2025-11-15T07:00:00Z",
                "author_profile": {
                    "id": "a2",
                    "slug": "jane",
                    "title": "Jane Developer",
                    "profile_picture_uri": null
                }
            },
            {
                "id": "s3",
                "slug": "building-developer-tools",
                "kind": "article",
                "title": "Building Developer Tools",
                "summary": "A deep dive into the architecture and design patterns behind modern developer tooling.",
                "content": null,
                "story_picture_uri": null,
                "published_at": "2025-11-01T12:00:00Z",
                "created_at": "2025-11-01T11:00:00Z",
                "author_profile": {
                    "id": "a1",
                    "slug": "eser",
                    "title": "Eser Ozvataf",
                    "profile_picture_uri": null
                }
            },
            {
                "id": "s4",
                "slug": "swift-concurrency-patterns",
                "kind": "article",
                "title": "Swift Concurrency Patterns",
                "summary": "Modern approaches to concurrent programming in Swift using async/await, actors, and structured concurrency.",
                "content": null,
                "story_picture_uri": null,
                "published_at": "2025-10-20T14:00:00Z",
                "created_at": "2025-10-20T13:00:00Z",
                "author_profile": {
                    "id": "a3",
                    "slug": "alp",
                    "title": "Alp Ozcan",
                    "profile_picture_uri": null
                }
            },
            {
                "id": "s5",
                "slug": "future-of-decentralized-apps",
                "kind": "essay",
                "title": "The Future of Decentralized Apps",
                "summary": "Exploring how decentralized applications are reshaping the way we think about software ownership and data privacy.",
                "content": null,
                "story_picture_uri": null,
                "published_at": "2025-10-10T09:00:00Z",
                "created_at": "2025-10-10T08:00:00Z",
                "author_profile": {
                    "id": "a2",
                    "slug": "jane",
                    "title": "Jane Developer",
                    "profile_picture_uri": null
                }
            }
        ],
        "cursor": null
    }
    """

    static let storyDetail = """
    {
        "data": {
            "id": "s1",
            "slug": "scaling-community-platforms",
            "kind": "article",
            "title": "Scaling Community Platforms",
            "summary": "Best practices for building and scaling open-source community platforms to millions of users.",
            "content": "# Scaling Community Platforms\\n\\nBuilding a community platform that serves millions of users requires careful architectural decisions.\\n\\n## Key Principles\\n\\n1. **Start simple** — avoid premature optimization\\n2. **Measure everything** — data-driven decisions\\n3. **Cache aggressively** — reduce database load\\n4. **Scale horizontally** — stateless services\\n\\n## Architecture\\n\\nThe AYA platform uses a hexagonal architecture that separates business logic from infrastructure concerns. This makes it easy to swap out implementations as the platform grows.\\n\\n## Conclusion\\n\\nScaling is not just about technology — it is about building the right abstractions early and evolving them as your community grows.",
            "story_picture_uri": null,
            "published_at": "2025-12-01T10:00:00Z",
            "created_at": "2025-12-01T09:00:00Z",
            "author_profile": {
                "id": "a1",
                "slug": "eser",
                "title": "Eser Ozvataf",
                "profile_picture_uri": null
            }
        }
    }
    """

    // MARK: - Activities

    static let activitiesResponse = """
    {
        "data": [
            {
                "id": "act1",
                "slug": "swift-workshop-2025",
                "kind": "activity",
                "title": "Swift Workshop 2025",
                "summary": "Hands-on workshop covering the latest Swift features and best practices.",
                "content": null,
                "story_picture_uri": null,
                "published_at": null,
                "created_at": "2025-11-01T10:00:00Z",
                "author_profile": {
                    "id": "a1",
                    "slug": "eser",
                    "title": "Eser Ozvataf",
                    "profile_picture_uri": null
                },
                "properties": {
                    "activity_kind": "workshop",
                    "activity_time_start": "2026-03-15T14:00:00Z",
                    "activity_time_end": "2026-03-15T16:00:00Z",
                    "external_activity_uri": null,
                    "external_attendance_uri": null,
                    "rsvp_mode": "open"
                }
            },
            {
                "id": "act2",
                "slug": "community-meetup-istanbul",
                "kind": "activity",
                "title": "Community Meetup Istanbul",
                "summary": "Monthly community gathering to discuss open-source projects and network with fellow developers.",
                "content": null,
                "story_picture_uri": null,
                "published_at": null,
                "created_at": "2025-10-15T08:00:00Z",
                "author_profile": null,
                "properties": {
                    "activity_kind": "meetup",
                    "activity_time_start": "2026-04-01T18:30:00Z",
                    "activity_time_end": null,
                    "external_activity_uri": null,
                    "external_attendance_uri": null,
                    "rsvp_mode": "open"
                }
            }
        ],
        "cursor": null
    }
    """

    static let activityDetail = """
    {
        "data": {
            "id": "act1",
            "slug": "swift-workshop-2025",
            "kind": "activity",
            "title": "Swift Workshop 2025",
            "summary": "Hands-on workshop covering the latest Swift features and best practices.",
            "content": "Join us for a comprehensive workshop on Swift development.",
            "story_picture_uri": null,
            "published_at": null,
            "created_at": "2025-11-01T10:00:00Z",
            "author_profile": {
                "id": "a1",
                "slug": "eser",
                "title": "Eser Ozvataf",
                "profile_picture_uri": null
            },
            "properties": {
                "activity_kind": "workshop",
                "activity_time_start": "2026-03-15T14:00:00Z",
                "activity_time_end": "2026-03-15T16:00:00Z",
                "external_activity_uri": null,
                "external_attendance_uri": null,
                "rsvp_mode": "open"
            }
        }
    }
    """

    // MARK: - Profiles

    static let profilesResponse = """
    {
        "data": [
            {
                "id": "p1",
                "slug": "eser",
                "kind": "individual",
                "title": "Eser Ozvataf",
                "description": "Software architect and open-source advocate. Building AYA.",
                "profile_picture_uri": null,
                "points": 2500,
                "default_locale": "en",
                "created_at": "2024-01-01T00:00:00Z",
                "has_translation": true
            },
            {
                "id": "p2",
                "slug": "aya-foundation",
                "kind": "organization",
                "title": "AYA Foundation",
                "description": "Building open-source tools for community empowerment.",
                "profile_picture_uri": null,
                "points": 5000,
                "default_locale": "en",
                "created_at": "2024-06-01T00:00:00Z",
                "has_translation": true
            },
            {
                "id": "p3",
                "slug": "alp",
                "kind": "individual",
                "title": "Alp Ozcan",
                "description": "Mobile developer and Swift enthusiast.",
                "profile_picture_uri": null,
                "points": 350,
                "default_locale": "en",
                "created_at": "2025-01-15T00:00:00Z",
                "has_translation": false
            }
        ],
        "cursor": null
    }
    """

    static let profileDetail = """
    {
        "data": {
            "id": "p1",
            "slug": "eser",
            "kind": "individual",
            "title": "Eser Ozvataf",
            "description": "Software architect and open-source advocate. Building AYA — an open software network for collaboration and knowledge sharing.",
            "profile_picture_uri": null,
            "points": 2500,
            "default_locale": "en",
            "created_at": "2024-01-01T00:00:00Z",
            "has_translation": true
        }
    }
    """

    // MARK: - Profile Sub-resources

    static let profilePages = """
    {
        "data": [
            {
                "id": "pg1",
                "slug": "about",
                "title": "About",
                "content": "Welcome to my profile. I build open-source software for developer communities.",
                "created_at": "2024-01-15T12:00:00Z"
            }
        ],
        "cursor": null
    }
    """

    static let profileStories = """
    {
        "data": [
            {
                "id": "s1",
                "slug": "scaling-community-platforms",
                "kind": "article",
                "title": "Scaling Community Platforms",
                "summary": "Best practices for building and scaling open-source community platforms.",
                "content": null,
                "story_picture_uri": null,
                "published_at": "2025-12-01T10:00:00Z",
                "created_at": "2025-12-01T09:00:00Z",
                "author_profile": {
                    "id": "a1",
                    "slug": "eser",
                    "title": "Eser Ozvataf",
                    "profile_picture_uri": null
                }
            }
        ],
        "cursor": null
    }
    """

    // MARK: - Search

    static let searchResponse = """
    {
        "data": [
            {
                "id": "sr1",
                "type": "story",
                "slug": "scaling-community-platforms",
                "title": "Scaling Community Platforms",
                "summary": "Best practices for building and scaling open-source community platforms.",
                "image_uri": null,
                "profile_slug": "eser",
                "profile_title": "Eser Ozvataf",
                "kind": "article",
                "rank": 0.95
            },
            {
                "id": "sr2",
                "type": "story",
                "slug": "open-source-sustainability",
                "title": "Open Source Sustainability",
                "summary": "How open-source projects can achieve long-term financial sustainability.",
                "image_uri": null,
                "profile_slug": "jane",
                "profile_title": "Jane Developer",
                "kind": "article",
                "rank": 0.85
            },
            {
                "id": "sr3",
                "type": "profile",
                "slug": "eser",
                "title": "Eser Ozvataf",
                "summary": "Software architect and open-source advocate.",
                "image_uri": null,
                "profile_slug": null,
                "profile_title": null,
                "kind": "individual",
                "rank": 0.80
            },
            {
                "id": "sr4",
                "type": "story",
                "slug": "building-developer-tools",
                "title": "Building Developer Tools",
                "summary": "A deep dive into the architecture behind modern developer tooling.",
                "image_uri": null,
                "profile_slug": "eser",
                "profile_title": "Eser Ozvataf",
                "kind": "article",
                "rank": 0.75
            },
            {
                "id": "sr5",
                "type": "story",
                "slug": "swift-concurrency-patterns",
                "title": "Swift Concurrency Patterns",
                "summary": "Modern approaches to concurrent programming in Swift.",
                "image_uri": null,
                "profile_slug": "alp",
                "profile_title": "Alp Ozcan",
                "kind": "article",
                "rank": 0.70
            }
        ],
        "cursor": null
    }
    """
}
