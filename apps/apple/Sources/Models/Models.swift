import Foundation

// MARK: - API Response Wrapper

/// Paginated API response wrapper.
public struct APIResponse<T: Decodable & Sendable>: Decodable, Sendable {
    /// The response payload items.
    public let data: [T]
    /// Cursor for fetching the next page, `nil` if no more pages.
    public let cursor: String?
}

// MARK: - Profile

/// A user or organization profile.
public struct Profile: Decodable, Sendable, Identifiable, Hashable {
    /// Unique identifier.
    public let id: String
    /// URL-friendly slug.
    public let slug: String
    /// Profile type (e.g. individual, organization, product).
    public let kind: String
    /// Display name.
    public let title: String
    /// Profile bio or description text.
    public let description: String
    /// URI of the profile picture, if available.
    public let profilePictureUri: String?
    /// Reputation points.
    public let points: Int
    /// Preferred locale for this profile.
    public let defaultLocale: String?
    /// ISO 8601 creation timestamp.
    public let createdAt: String?
    /// Whether a translation exists for the requested locale.
    public let hasTranslation: Bool?

    /// Hashes by `id` for `Hashable` conformance.
    public func hash(into hasher: inout Hasher) {
        hasher.combine(id)
    }

    /// Equality based on `id`.
    public static func == (lhs: Profile, rhs: Profile) -> Bool {
        lhs.id == rhs.id
    }
}

// MARK: - Story

/// A published story or article.
public struct Story: Decodable, Sendable, Identifiable, Hashable {
    /// Unique identifier.
    public let id: String
    /// URL-friendly slug.
    public let slug: String
    /// Story type classifier.
    public let kind: String
    /// Headline or title.
    public let title: String
    /// Short summary or excerpt.
    public let summary: String?
    /// Full body content (may contain HTML/Markdown).
    public let content: String?
    /// URI of the story's cover image.
    public let storyPictureUri: String?
    /// ISO 8601 publication timestamp.
    public let publishedAt: String?
    /// ISO 8601 creation timestamp.
    public let createdAt: String?
    /// The author's profile summary.
    public let authorProfile: AuthorProfile?

    /// Hashes by `id` for `Hashable` conformance.
    public func hash(into hasher: inout Hasher) {
        hasher.combine(id)
    }

    /// Equality based on `id`.
    public static func == (lhs: Story, rhs: Story) -> Bool {
        lhs.id == rhs.id
    }
}

// MARK: - AuthorProfile (nested in Story)

/// Lightweight author information embedded in a story or activity.
public struct AuthorProfile: Decodable, Sendable, Hashable {
    /// Unique identifier.
    public let id: String
    /// URL-friendly slug.
    public let slug: String
    /// Display name.
    public let title: String
    /// URI of the author's profile picture.
    public let profilePictureUri: String?
}

// MARK: - Activity (Story with properties)

/// A story enriched with activity-specific scheduling and RSVP properties.
public struct Activity: Decodable, Sendable, Identifiable, Hashable {
    /// Unique identifier.
    public let id: String
    /// URL-friendly slug.
    public let slug: String
    /// Activity type classifier.
    public let kind: String
    /// Headline or title.
    public let title: String
    /// Short summary or excerpt.
    public let summary: String?
    /// Full body content.
    public let content: String?
    /// URI of the activity's cover image.
    public let storyPictureUri: String?
    /// ISO 8601 publication timestamp.
    public let publishedAt: String?
    /// ISO 8601 creation timestamp.
    public let createdAt: String?
    /// The author's profile summary.
    public let authorProfile: AuthorProfile?
    /// Activity-specific metadata (schedule, links, RSVP).
    public let properties: ActivityProperties?

    /// Hashes by `id` for `Hashable` conformance.
    public func hash(into hasher: inout Hasher) {
        hasher.combine(id)
    }

    /// Equality based on `id`.
    public static func == (lhs: Activity, rhs: Activity) -> Bool {
        lhs.id == rhs.id
    }
}

/// Scheduling and RSVP metadata for an activity.
public struct ActivityProperties: Decodable, Sendable, Hashable {
    /// Sub-type of activity (e.g. event, workshop).
    public let activityKind: String?
    /// ISO 8601 start time.
    public let activityTimeStart: String?
    /// ISO 8601 end time.
    public let activityTimeEnd: String?
    /// External link to the activity.
    public let externalActivityUri: String?
    /// External link for attendance/registration.
    public let externalAttendanceUri: String?
    /// RSVP mode (e.g. open, closed).
    public let rsvpMode: String?
}

// MARK: - Search Result

/// A single item returned from the search API.
public struct SearchResult: Decodable, Sendable, Identifiable, Hashable {
    /// Unique identifier.
    public let id: String
    /// Entity type (e.g. story, profile).
    public let type: String
    /// URL-friendly slug.
    public let slug: String
    /// Display title.
    public let title: String
    /// Short summary or excerpt.
    public let summary: String?
    /// URI of a representative image.
    public let imageUri: String?
    /// Slug of the associated profile, if any.
    public let profileSlug: String?
    /// Title of the associated profile, if any.
    public let profileTitle: String?
    /// Entity kind classifier.
    public let kind: String?
    /// Search relevance score.
    public let rank: Double?

    /// Hashes by `id` for `Hashable` conformance.
    public func hash(into hasher: inout Hasher) {
        hasher.combine(id)
    }

    /// Equality based on `id`.
    public static func == (lhs: SearchResult, rhs: SearchResult) -> Bool {
        lhs.id == rhs.id
    }
}

// MARK: - Page

/// A static content page belonging to a profile.
public struct Page: Decodable, Sendable, Identifiable, Hashable {
    /// Unique identifier.
    public let id: String
    /// URL-friendly slug.
    public let slug: String
    /// Page title.
    public let title: String
    /// Full body content.
    public let content: String?
    /// ISO 8601 creation timestamp.
    public let createdAt: String?

    /// Hashes by `id` for `Hashable` conformance.
    public func hash(into hasher: inout Hasher) {
        hasher.combine(id)
    }

    /// Equality based on `id`.
    public static func == (lhs: Page, rhs: Page) -> Bool {
        lhs.id == rhs.id
    }
}

// MARK: - Spotlight

/// Featured content collections for the spotlight/home screen.
public struct Spotlight: Decodable, Sendable {
    /// Featured stories.
    public let stories: [Story]?
    /// Featured activities.
    public let activities: [Activity]?
    /// Featured profiles.
    public let profiles: [Profile]?
}
