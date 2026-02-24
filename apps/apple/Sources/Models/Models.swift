import Foundation

// MARK: - API Response Wrapper

public struct APIResponse<T: Decodable & Sendable>: Decodable, Sendable {
    public let data: [T]
    public let cursor: String?
}

// MARK: - Profile

public struct Profile: Decodable, Sendable, Identifiable, Hashable {
    public let id: String
    public let slug: String
    public let kind: String
    public let title: String
    public let description: String
    public let profilePictureUri: String?
    public let points: Int
    public let defaultLocale: String?
    public let createdAt: String?
    public let hasTranslation: Bool?

    public func hash(into hasher: inout Hasher) {
        hasher.combine(id)
    }

    public static func == (lhs: Profile, rhs: Profile) -> Bool {
        lhs.id == rhs.id
    }
}

// MARK: - Story

public struct Story: Decodable, Sendable, Identifiable, Hashable {
    public let id: String
    public let slug: String
    public let kind: String
    public let title: String
    public let summary: String?
    public let content: String?
    public let storyPictureUri: String?
    public let publishedAt: String?
    public let createdAt: String?
    public let authorProfile: AuthorProfile?

    public func hash(into hasher: inout Hasher) {
        hasher.combine(id)
    }

    public static func == (lhs: Story, rhs: Story) -> Bool {
        lhs.id == rhs.id
    }
}

// MARK: - AuthorProfile (nested in Story)

public struct AuthorProfile: Decodable, Sendable, Hashable {
    public let id: String
    public let slug: String
    public let title: String
    public let profilePictureUri: String?
}

// MARK: - Activity (Story with properties)

public struct Activity: Decodable, Sendable, Identifiable, Hashable {
    public let id: String
    public let slug: String
    public let kind: String
    public let title: String
    public let summary: String?
    public let content: String?
    public let storyPictureUri: String?
    public let publishedAt: String?
    public let createdAt: String?
    public let authorProfile: AuthorProfile?
    public let properties: ActivityProperties?

    public func hash(into hasher: inout Hasher) {
        hasher.combine(id)
    }

    public static func == (lhs: Activity, rhs: Activity) -> Bool {
        lhs.id == rhs.id
    }
}

public struct ActivityProperties: Decodable, Sendable, Hashable {
    public let activityKind: String?
    public let activityTimeStart: String?
    public let activityTimeEnd: String?
    public let externalActivityUri: String?
    public let externalAttendanceUri: String?
    public let rsvpMode: String?
}

// MARK: - Search Result

public struct SearchResult: Decodable, Sendable, Identifiable, Hashable {
    public let id: String
    public let type: String
    public let slug: String
    public let title: String
    public let summary: String?
    public let imageUri: String?
    public let profileSlug: String?
    public let profileTitle: String?
    public let kind: String?
    public let rank: Double?

    public func hash(into hasher: inout Hasher) {
        hasher.combine(id)
    }

    public static func == (lhs: SearchResult, rhs: SearchResult) -> Bool {
        lhs.id == rhs.id
    }
}

// MARK: - Page

public struct Page: Decodable, Sendable, Identifiable, Hashable {
    public let id: String
    public let slug: String
    public let title: String
    public let content: String?
    public let createdAt: String?

    public func hash(into hasher: inout Hasher) {
        hasher.combine(id)
    }

    public static func == (lhs: Page, rhs: Page) -> Bool {
        lhs.id == rhs.id
    }
}

// MARK: - Spotlight

public struct Spotlight: Decodable, Sendable {
    public let stories: [Story]?
    public let activities: [Activity]?
    public let profiles: [Profile]?
}
