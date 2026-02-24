import Foundation

/// A feed item representing one of the content types displayed in the main feed.
public enum FeedItem: Identifiable, Hashable, Sendable {
    /// A story entry authored by a profile.
    case story(Story)
    /// An activity such as a workshop, meetup, or conference.
    case activity(Activity)
    /// An individual profile (person).
    case profile(Profile)
    /// A product profile.
    case product(Profile)

    /// A unique identifier combining the content type prefix and the item's own ID.
    public var id: String {
        switch self {
        case .story(let s): "story-\(s.id)"
        case .activity(let a): "activity-\(a.id)"
        case .profile(let p): "profile-\(p.id)"
        case .product(let p): "product-\(p.id)"
        }
    }

    /// The associated date for sorting, parsed from the item's date strings.
    public var date: Date? {
        let dateString: String? = switch self {
        case .story(let s): s.publishedAt ?? s.createdAt
        case .activity(let a): a.properties?.activityTimeStart ?? a.publishedAt ?? a.createdAt
        case .profile(let p): p.createdAt
        case .product(let p): p.createdAt
        }
        guard let dateString else { return nil }
        return Self.parseDate(dateString)
    }

    /// The display title of the feed item.
    public var title: String {
        switch self {
        case .story(let s): s.title
        case .activity(let a): a.title
        case .profile(let p): p.title
        case .product(let p): p.title
        }
    }

    /// The URL-friendly slug used for navigation and deep linking.
    public var slug: String {
        switch self {
        case .story(let s): s.slug
        case .activity(let a): a.slug
        case .profile(let p): p.slug
        case .product(let p): p.slug
        }
    }

    private static func parseDate(_ iso: String) -> Date? {
        let formatter = ISO8601DateFormatter()
        formatter.formatOptions = [.withInternetDateTime, .withFractionalSeconds]
        if let date = formatter.date(from: iso) { return date }
        formatter.formatOptions = [.withInternetDateTime]
        if let date = formatter.date(from: iso) { return date }
        let short = DateFormatter()
        short.dateFormat = "yyyy-MM-dd'T'HH:mm"
        short.timeZone = .current
        return short.date(from: iso)
    }
}

/// Filter options for narrowing the feed to a specific content type.
public enum FeedFilter: String, CaseIterable, Sendable {
    case stories
    case activities
    case people
    case products
}
