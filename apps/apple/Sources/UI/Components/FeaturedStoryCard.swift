import SwiftUI

/// A large, hero-style card with a full-bleed image and overlaid title, summary, and author.
public struct FeaturedStoryCard: View {
    let title: String
    let summary: String?
    let imageUrl: String?
    let authorName: String?
    let kind: String?

    /// Creates a featured story card.
    public init(
        title: String,
        summary: String? = nil,
        imageUrl: String? = nil,
        authorName: String? = nil,
        kind: String? = nil
    ) {
        self.title = title
        self.summary = summary
        self.imageUrl = imageUrl
        self.authorName = authorName
        self.kind = kind
    }

    /// The card layout: background image with gradient overlay, kind badge, title, summary, and author.
    public var body: some View {
        ZStack(alignment: .bottomLeading) {
            if let imageUrl {
                RemoteImage(urlString: imageUrl, cornerRadius: AYACornerRadius.xl)
                    .frame(height: 320)
                    .frame(maxWidth: .infinity)
            } else {
                RoundedRectangle(cornerRadius: AYACornerRadius.xl)
                    .fill(
                        LinearGradient(
                            colors: [AYAColors.accent, AYAColors.accent.opacity(0.6)],
                            startPoint: .topLeading,
                            endPoint: .bottomTrailing
                        )
                    )
                    .frame(height: 320)
            }

            LinearGradient(
                stops: [
                    .init(color: .black.opacity(0.85), location: 0),
                    .init(color: .black.opacity(0.5), location: 0.5),
                    .init(color: .clear, location: 1),
                ],
                startPoint: .bottom,
                endPoint: .top
            )
            .clipShape(RoundedRectangle(cornerRadius: AYACornerRadius.xl))

            VStack(alignment: .leading, spacing: AYASpacing.sm) {
                if let kind {
                    Text(kind.capitalized)
                        .font(AYATypography.caption2)
                        .fontWeight(.bold)
                        .foregroundStyle(.white)
                        .padding(.horizontal, 8)
                        .padding(.vertical, 3)
                        .background(.white.opacity(0.2))
                        .clipShape(Capsule())
                }

                Text(title)
                    .font(AYATypography.title2)
                    .fontWeight(.bold)
                    .foregroundStyle(.white)
                    .lineLimit(3)
                    .shadow(color: .black.opacity(0.3), radius: 4, x: 0, y: 2)

                if let summary, !summary.isEmpty {
                    Text(summary)
                        .font(AYATypography.subheadline)
                        .foregroundStyle(.white.opacity(0.9))
                        .lineLimit(2)
                }

                if let authorName {
                    HStack(spacing: AYASpacing.xs) {
                        Image(systemName: "person.circle.fill")
                            .font(.caption)
                        Text(authorName)
                            .font(AYATypography.caption)
                            .fontWeight(.medium)
                    }
                    .foregroundStyle(.white.opacity(0.8))
                }
            }
            .padding(AYASpacing.lg)
        }
        .clipShape(RoundedRectangle(cornerRadius: AYACornerRadius.xl))
        .shadow(color: .black.opacity(0.15), radius: 16, x: 0, y: 8)
    }
}
