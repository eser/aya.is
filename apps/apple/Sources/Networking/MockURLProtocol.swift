import Foundation

/// A `URLProtocol` subclass that intercepts all HTTP requests and returns
/// canned JSON responses. Activated by passing `--uitesting` as a launch argument.
///
/// This enables deterministic UI tests that never hit the network.
public final class MockURLProtocol: URLProtocol {

    /// Route handler: maps a URL path pattern to a JSON response body.
    nonisolated(unsafe) static var routes: [String: Data] = [:]

    override public class func canInit(with request: URLRequest) -> Bool {
        true
    }

    override public class func canonicalRequest(for request: URLRequest) -> URLRequest {
        request
    }

    override public func startLoading() {
        guard let url = request.url else {
            client?.urlProtocol(self, didFailWithError: URLError(.badURL))
            return
        }

        let path = url.path
        let data = Self.matchRoute(path: path, query: url.query)

        let response = HTTPURLResponse(
            url: url,
            statusCode: 200,
            httpVersion: "HTTP/1.1",
            headerFields: ["Content-Type": "application/json"]
        )!

        client?.urlProtocol(self, didReceive: response, cacheStoragePolicy: .notAllowed)
        client?.urlProtocol(self, didLoad: data)
        client?.urlProtocolDidFinishLoading(self)
    }

    override public func stopLoading() {}

    // MARK: - Route Matching

    private static func matchRoute(path: String, query: String?) -> Data {
        // Try exact path match first
        if let data = routes[path] {
            return data
        }

        // Try path with query
        if let query, let data = routes["\(path)?\(query)"] {
            return data
        }

        // Try pattern matching (e.g. /en/stories/{slug})
        for (pattern, data) in routes {
            if pathMatches(pattern: pattern, path: path) {
                return data
            }
        }

        // Fallback: empty response
        return Data("{\"data\":[],\"cursor\":null}".utf8)
    }

    private static func pathMatches(pattern: String, path: String) -> Bool {
        let patternParts = pattern.split(separator: "/")
        let pathParts = path.split(separator: "/")
        guard patternParts.count == pathParts.count else { return false }
        return zip(patternParts, pathParts).allSatisfy { p, v in
            p.hasPrefix("{") && p.hasSuffix("}") || p == v
        }
    }

    // MARK: - Registration

    /// Configures a `URLSession` that uses this mock protocol with all routes pre-loaded.
    public static func mockSession() -> URLSession {
        registerRoutes()
        let config = URLSessionConfiguration.ephemeral
        config.protocolClasses = [MockURLProtocol.self]
        return URLSession(configuration: config)
    }

    /// Registers all canned responses for every endpoint the app uses.
    private static func registerRoutes() {
        // Stories list
        routes["/en/stories"] = Data(MockData.storiesResponse.utf8)

        // Single story detail
        routes["/{locale}/stories/{slug}"] = Data(MockData.storyDetail.utf8)
        routes["/en/stories/scaling-community-platforms"] = Data(MockData.storyDetail.utf8)
        routes["/en/stories/open-source-sustainability"] = Data(MockData.storyDetail.utf8)
        routes["/en/stories/building-developer-tools"] = Data(MockData.storyDetail.utf8)

        // Activities list
        routes["/en/activities"] = Data(MockData.activitiesResponse.utf8)

        // Single activity detail
        routes["/{locale}/activities/{slug}"] = Data(MockData.activityDetail.utf8)

        // Profiles list
        routes["/en/profiles"] = Data(MockData.profilesResponse.utf8)

        // Single profile detail
        routes["/{locale}/profiles/{slug}"] = Data(MockData.profileDetail.utf8)
        routes["/en/profiles/eser"] = Data(MockData.profileDetail.utf8)

        // Profile pages
        routes["/{locale}/profiles/{slug}/pages"] = Data(MockData.profilePages.utf8)
        routes["/en/profiles/eser/pages"] = Data(MockData.profilePages.utf8)

        // Profile stories
        routes["/{locale}/profiles/{slug}/stories"] = Data(MockData.profileStories.utf8)
        routes["/en/profiles/eser/stories"] = Data(MockData.profileStories.utf8)

        // Search
        routes["/en/search"] = Data(MockData.searchResponse.utf8)
    }
}
