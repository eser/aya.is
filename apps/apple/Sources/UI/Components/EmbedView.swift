import SwiftUI
import YouTubePlayerKit

// MARK: - Embed Dispatcher

/// Parses a URL and renders the appropriate embed: YouTube player, Spotify deep-link card, or a generic link card.
struct EmbedView: View {
    let url: String

    var body: some View {
        if let youtubeID = Self.youtubeVideoID(from: url) {
            YouTubeEmbedView(videoID: youtubeID)
        } else if let spotifyInfo = Self.spotifyInfo(from: url) {
            SpotifyLinkCard(kind: spotifyInfo.kind, title: spotifyInfo.kind.capitalized, webURL: url)
        } else {
            LinkCardView(urlString: url)
        }
    }

    // MARK: - URL Matching

    /// Extracts a YouTube video ID from watch, short, or embed URLs.
    static func youtubeVideoID(from urlString: String) -> String? {
        guard let url = URL(string: urlString) else { return nil }
        let host = url.host?.lowercased() ?? ""

        // youtube.com/watch?v=ID
        if host.contains("youtube.com"),
           let components = URLComponents(url: url, resolvingAgainstBaseURL: false),
           let videoID = components.queryItems?.first(where: { $0.name == "v" })?.value,
           !videoID.isEmpty {
            return videoID
        }

        // youtu.be/ID
        if host.contains("youtu.be") {
            let id = url.lastPathComponent
            if !id.isEmpty { return id }
        }

        // youtube.com/embed/ID
        if host.contains("youtube.com"), url.pathComponents.contains("embed") {
            if let idx = url.pathComponents.firstIndex(of: "embed"),
               idx + 1 < url.pathComponents.count {
                return url.pathComponents[idx + 1]
            }
        }

        return nil
    }

    /// Extracts Spotify content info from open.spotify.com URLs.
    static func spotifyInfo(from urlString: String) -> (kind: String, id: String)? {
        guard let url = URL(string: urlString) else { return nil }
        let host = url.host?.lowercased() ?? ""
        guard host.contains("spotify.com") else { return nil }

        let components = url.pathComponents.filter { $0 != "/" }
        guard components.count >= 2 else { return nil }
        let kind = components[0]
        let id = components[1]
        let validKinds: Set<String> = ["track", "album", "playlist", "episode", "show"]
        guard validKinds.contains(kind) else { return nil }
        return (kind, id)
    }
}

// MARK: - YouTube Embed (via YouTubePlayerKit SDK)

/// Renders a YouTube video using the YouTubePlayerKit SDK with inline playback.
struct YouTubeEmbedView: View {
    let videoID: String

    @State private var player: YouTubePlayer

    init(videoID: String) {
        self.videoID = videoID
        self._player = State(initialValue: YouTubePlayer(
            source: .video(id: videoID),
            configuration: .init(
                fullscreenMode: .system,
                autoPlay: false,
                showControls: true,
                loopEnabled: false
            )
        ))
    }

    var body: some View {
        YouTubePlayerView(player) { state in
            switch state {
            case .idle:
                ZStack {
                    RoundedRectangle(cornerRadius: AYACornerRadius.lg)
                        .fill(AYAColors.contentBackground)
                    ProgressView()
                }
            case .ready:
                EmptyView()
            case .error:
                LinkCardView(urlString: "https://www.youtube.com/watch?v=\(videoID)")
            }
        }
        .aspectRatio(16.0 / 9.0, contentMode: .fit)
        .clipShape(RoundedRectangle(cornerRadius: AYACornerRadius.lg))
    }
}

// MARK: - Spotify Link Card

/// A styled card that deep-links to the Spotify app or falls back to the web URL.
struct SpotifyLinkCard: View {
    let kind: String
    let title: String
    let webURL: String

    private var spotifyAppURL: URL? {
        // spotify: URI scheme opens the Spotify app directly
        guard let url = URL(string: webURL) else { return nil }
        let components = url.pathComponents.filter { $0 != "/" }
        guard components.count >= 2 else { return nil }
        return URL(string: "spotify:\(components[0]):\(components[1])")
    }

    private var iconName: String {
        switch kind {
        case "track": "music.note"
        case "album": "square.stack"
        case "playlist": "music.note.list"
        case "episode": "mic"
        case "show": "antenna.radiowaves.left.and.right"
        default: "music.note"
        }
    }

    var body: some View {
        if let url = URL(string: webURL) {
            Link(destination: spotifyAppURL ?? url) {
                HStack(spacing: AYASpacing.md) {
                    ZStack {
                        RoundedRectangle(cornerRadius: AYACornerRadius.md)
                            .fill(Color.green.opacity(0.12))
                            .frame(width: 44, height: 44)

                        Image(systemName: iconName)
                            .font(.system(size: 18, weight: .medium))
                            .foregroundStyle(.green)
                    }

                    VStack(alignment: .leading, spacing: 2) {
                        HStack(spacing: AYASpacing.xs) {
                            Text("Spotify")
                                .font(AYATypography.subheadline)
                                .fontWeight(.semibold)
                                .foregroundStyle(AYAColors.textPrimary)

                            Text(kind.capitalized)
                                .font(AYATypography.caption2)
                                .fontWeight(.medium)
                                .foregroundStyle(.green)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(.green.opacity(0.12))
                                .clipShape(Capsule())
                        }

                        Text(webURL)
                            .font(AYATypography.caption)
                            .foregroundStyle(AYAColors.textTertiary)
                            .lineLimit(1)
                    }

                    Spacer()

                    Image(systemName: "arrow.up.right")
                        .font(.system(size: 12, weight: .semibold))
                        .foregroundStyle(AYAColors.textTertiary)
                }
                .padding(AYASpacing.md)
                .background(AYAColors.contentBackground)
                .clipShape(RoundedRectangle(cornerRadius: AYACornerRadius.lg))
            }
            .buttonStyle(.plain)
        }
    }
}

// MARK: - Generic Link Card

/// A compact card that displays a URL with its domain and a link icon.
struct LinkCardView: View {
    let urlString: String

    var body: some View {
        if let url = URL(string: urlString) {
            Link(destination: url) {
                HStack(spacing: AYASpacing.md) {
                    ZStack {
                        RoundedRectangle(cornerRadius: AYACornerRadius.md)
                            .fill(AYAColors.accent.opacity(0.1))
                            .frame(width: 44, height: 44)

                        Image(systemName: "link")
                            .font(.system(size: 18, weight: .medium))
                            .foregroundStyle(AYAColors.accent)
                    }

                    VStack(alignment: .leading, spacing: 2) {
                        Text(url.host ?? urlString)
                            .font(AYATypography.subheadline)
                            .fontWeight(.medium)
                            .foregroundStyle(AYAColors.textPrimary)
                            .lineLimit(1)

                        Text(urlString)
                            .font(AYATypography.caption)
                            .foregroundStyle(AYAColors.textTertiary)
                            .lineLimit(1)
                    }

                    Spacer()

                    Image(systemName: "arrow.up.right")
                        .font(.system(size: 12, weight: .semibold))
                        .foregroundStyle(AYAColors.textTertiary)
                }
                .padding(AYASpacing.md)
                .background(AYAColors.contentBackground)
                .clipShape(RoundedRectangle(cornerRadius: AYACornerRadius.lg))
            }
            .buttonStyle(.plain)
        }
    }
}
