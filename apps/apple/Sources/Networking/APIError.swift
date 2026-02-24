import Foundation

public enum APIError: LocalizedError, Sendable {
    case invalidURL(String)
    case invalidResponse
    case httpError(statusCode: Int)
    case decodingError(Error)

    public var errorDescription: String? {
        switch self {
        case .invalidURL(let path): "Invalid URL: \(path)"
        case .invalidResponse: "Invalid server response"
        case .httpError(let code): "HTTP error \(code)"
        case .decodingError(let error): "Failed to decode: \(error.localizedDescription)"
        }
    }
}
