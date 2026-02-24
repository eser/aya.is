import Foundation
import Models

public enum FeedItem: Identifiable, Hashable, Sendable {
    case story(Story)
    case activity(Activity)
    case profile(Profile)
    case product(Profile)

    public var id: String {
        switch self {
        case .story(let s): "story-\(s.id)"
        case .activity(let a): "activity-\(a.id)"
        case .profile(let p): "profile-\(p.id)"
        case .product(let p): "product-\(p.id)"
        }
    }

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

    public var title: String {
        switch self {
        case .story(let s): s.title
        case .activity(let a): a.title
        case .profile(let p): p.title
        case .product(let p): p.title
        }
    }

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
        short.timeZone = TimeZone(identifier: "Europe/Istanbul")
        return short.date(from: iso)
    }
}

public enum FeedFilter: String, CaseIterable, Sendable {
    case stories
    case activities
    case people
    case products
}
