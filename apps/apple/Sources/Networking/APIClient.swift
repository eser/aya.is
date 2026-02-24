import Foundation

// MARK: - API Client Protocol

/// Contract for all AYA API operations, enabling testable dependency injection.
public protocol APIClientProtocol: Sendable {
    /// Fetches a paginated list of stories.
    func fetchStories(locale: String, cursor: String?) async throws -> APIResponse<Story>
    /// Fetches a single story by its slug.
    func fetchStory(locale: String, slug: String) async throws -> Story
    /// Fetches a paginated list of activities.
    func fetchActivities(locale: String, cursor: String?) async throws -> APIResponse<Activity>
    /// Fetches a single activity by its slug.
    func fetchActivity(locale: String, slug: String) async throws -> Activity
    /// Fetches a paginated list of profiles, optionally filtered by kind.
    func fetchProfiles(locale: String, filterKind: String?, cursor: String?) async throws -> APIResponse<Profile>
    /// Fetches a single profile by its slug.
    func fetchProfile(locale: String, slug: String) async throws -> Profile
    /// Fetches pages belonging to a profile.
    func fetchProfilePages(locale: String, slug: String) async throws -> APIResponse<Page>
    /// Fetches stories authored by a profile.
    func fetchProfileStories(locale: String, slug: String) async throws -> APIResponse<Story>
    /// Performs a full-text search across all content types.
    func search(locale: String, query: String, cursor: String?) async throws -> APIResponse<SearchResult>
}

// MARK: - API Client

/// Concrete API client that communicates with the AYA REST API over HTTPS.
public struct APIClient: APIClientProtocol {
    private let baseURL: URL
    private let session: URLSession
    private let decoder: JSONDecoder

    /// Creates an API client.
    /// - Parameters:
    ///   - baseURL: Root URL of the API. Defaults to `https://api.aya.is`.
    ///   - session: URL session used for network requests. Defaults to `.shared`.
    public init(
        baseURL: URL = URL(string: "https://api.aya.is")!,
        session: URLSession = .shared
    ) {
        self.baseURL = baseURL
        self.session = session

        let decoder = JSONDecoder()
        decoder.keyDecodingStrategy = .convertFromSnakeCase
        self.decoder = decoder
    }

    // MARK: - Stories

    /// Fetches a paginated list of stories for the given locale.
    public func fetchStories(locale: String, cursor: String? = nil) async throws -> APIResponse<Story> {
        var components = urlComponents(path: "/\(locale)/stories")
        if let cursor { components.queryItems = [URLQueryItem(name: "cursor", value: cursor)] }
        return try await request(components)
    }

    /// Fetches a single story by its slug.
    public func fetchStory(locale: String, slug: String) async throws -> Story {
        let components = urlComponents(path: "/\(locale)/stories/\(slug)")
        return try await request(components)
    }

    // MARK: - Activities

    /// Fetches a paginated list of activities for the given locale.
    public func fetchActivities(locale: String, cursor: String? = nil) async throws -> APIResponse<Activity> {
        var components = urlComponents(path: "/\(locale)/activities")
        if let cursor { components.queryItems = [URLQueryItem(name: "cursor", value: cursor)] }
        return try await request(components)
    }

    /// Fetches a single activity by its slug.
    public func fetchActivity(locale: String, slug: String) async throws -> Activity {
        let components = urlComponents(path: "/\(locale)/activities/\(slug)")
        return try await request(components)
    }

    // MARK: - Profiles

    /// Fetches a paginated list of profiles, optionally filtered by kind.
    public func fetchProfiles(locale: String, filterKind: String? = nil, cursor: String? = nil) async throws -> APIResponse<Profile> {
        var components = urlComponents(path: "/\(locale)/profiles")
        var queryItems: [URLQueryItem] = []
        if let filterKind { queryItems.append(URLQueryItem(name: "filter_kind", value: filterKind)) }
        if let cursor { queryItems.append(URLQueryItem(name: "cursor", value: cursor)) }
        if !queryItems.isEmpty { components.queryItems = queryItems }
        return try await request(components)
    }

    /// Fetches a single profile by its slug.
    public func fetchProfile(locale: String, slug: String) async throws -> Profile {
        let components = urlComponents(path: "/\(locale)/profiles/\(slug)")
        return try await request(components)
    }

    /// Fetches pages belonging to a profile.
    public func fetchProfilePages(locale: String, slug: String) async throws -> APIResponse<Page> {
        let components = urlComponents(path: "/\(locale)/profiles/\(slug)/pages")
        return try await request(components)
    }

    /// Fetches stories authored by a profile.
    public func fetchProfileStories(locale: String, slug: String) async throws -> APIResponse<Story> {
        let components = urlComponents(path: "/\(locale)/profiles/\(slug)/stories")
        return try await request(components)
    }

    // MARK: - Search

    /// Performs a full-text search across all content types.
    public func search(locale: String, query: String, cursor: String? = nil) async throws -> APIResponse<SearchResult> {
        var components = urlComponents(path: "/\(locale)/search")
        var queryItems = [URLQueryItem(name: "q", value: query)]
        if let cursor { queryItems.append(URLQueryItem(name: "cursor", value: cursor)) }
        components.queryItems = queryItems
        return try await request(components)
    }

    // MARK: - Private Helpers

    private func urlComponents(path: String) -> URLComponents {
        var components = URLComponents(url: baseURL, resolvingAgainstBaseURL: false)!
        components.path = path
        return components
    }

    private func request<T: Decodable & Sendable>(_ components: URLComponents) async throws -> T {
        guard let url = components.url else {
            throw APIError.invalidURL(components.path)
        }

        var request = URLRequest(url: url)
        request.setValue("application/json", forHTTPHeaderField: "Accept")
        request.timeoutInterval = 30

        let (data, response) = try await session.data(for: request)

        guard let httpResponse = response as? HTTPURLResponse else {
            throw APIError.invalidResponse
        }

        guard (200...299).contains(httpResponse.statusCode) else {
            throw APIError.httpError(statusCode: httpResponse.statusCode)
        }

        do {
            return try decoder.decode(T.self, from: data)
        } catch {
            throw APIError.decodingError(error)
        }
    }
}
