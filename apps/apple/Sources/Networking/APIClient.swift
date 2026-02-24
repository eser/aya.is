import Foundation
import Models

// MARK: - API Client Protocol

public protocol APIClientProtocol: Sendable {
    func fetchStories(locale: String, cursor: String?) async throws -> APIResponse<Story>
    func fetchStory(locale: String, slug: String) async throws -> Story
    func fetchActivities(locale: String, cursor: String?) async throws -> APIResponse<Activity>
    func fetchActivity(locale: String, slug: String) async throws -> Activity
    func fetchProfiles(locale: String, filterKind: String?, cursor: String?) async throws -> APIResponse<Profile>
    func fetchProfile(locale: String, slug: String) async throws -> Profile
    func fetchProfilePages(locale: String, slug: String) async throws -> APIResponse<Page>
    func fetchProfileStories(locale: String, slug: String) async throws -> APIResponse<Story>
    func search(locale: String, query: String, cursor: String?) async throws -> APIResponse<SearchResult>
}

// MARK: - API Client

public struct APIClient: APIClientProtocol {
    private let baseURL: URL
    private let session: URLSession
    private let decoder: JSONDecoder

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

    public func fetchStories(locale: String, cursor: String? = nil) async throws -> APIResponse<Story> {
        var components = urlComponents(path: "/\(locale)/stories")
        if let cursor { components.queryItems = [URLQueryItem(name: "cursor", value: cursor)] }
        return try await request(components)
    }

    public func fetchStory(locale: String, slug: String) async throws -> Story {
        let components = urlComponents(path: "/\(locale)/stories/\(slug)")
        return try await request(components)
    }

    // MARK: - Activities

    public func fetchActivities(locale: String, cursor: String? = nil) async throws -> APIResponse<Activity> {
        var components = urlComponents(path: "/\(locale)/activities")
        if let cursor { components.queryItems = [URLQueryItem(name: "cursor", value: cursor)] }
        return try await request(components)
    }

    public func fetchActivity(locale: String, slug: String) async throws -> Activity {
        let components = urlComponents(path: "/\(locale)/activities/\(slug)")
        return try await request(components)
    }

    // MARK: - Profiles

    public func fetchProfiles(locale: String, filterKind: String? = nil, cursor: String? = nil) async throws -> APIResponse<Profile> {
        var components = urlComponents(path: "/\(locale)/profiles")
        var queryItems: [URLQueryItem] = []
        if let filterKind { queryItems.append(URLQueryItem(name: "filter_kind", value: filterKind)) }
        if let cursor { queryItems.append(URLQueryItem(name: "cursor", value: cursor)) }
        if !queryItems.isEmpty { components.queryItems = queryItems }
        return try await request(components)
    }

    public func fetchProfile(locale: String, slug: String) async throws -> Profile {
        let components = urlComponents(path: "/\(locale)/profiles/\(slug)")
        return try await request(components)
    }

    public func fetchProfilePages(locale: String, slug: String) async throws -> APIResponse<Page> {
        let components = urlComponents(path: "/\(locale)/profiles/\(slug)/pages")
        return try await request(components)
    }

    public func fetchProfileStories(locale: String, slug: String) async throws -> APIResponse<Story> {
        let components = urlComponents(path: "/\(locale)/profiles/\(slug)/stories")
        return try await request(components)
    }

    // MARK: - Search

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
