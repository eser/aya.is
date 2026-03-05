import Foundation

/// Errors that can occur during API communication.
public enum APIError: LocalizedError, Sendable {
    /// The constructed URL path was invalid.
    case invalidURL(String)
    /// The server returned a non-HTTP response.
    case invalidResponse
    /// The server returned a non-2xx status code.
    case httpError(statusCode: Int)
    /// JSON decoding failed.
    case decodingError(Error)

    /// Human-readable error description.
    public var errorDescription: String? {
        switch self {
        case .invalidURL(let path): "Invalid URL: \(path)"
        case .invalidResponse: "Invalid server response"
        case .httpError(let code): "HTTP error \(code)"
        case .decodingError(let error): "Failed to decode: \(error.localizedDescription)"
        }
    }
}
