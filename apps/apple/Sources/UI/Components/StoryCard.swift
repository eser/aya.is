import SwiftUI

public struct StoryCard: View {
    let title: String
    let summary: String?
    let imageUrl: String?
    let authorName: String?
    let authorImageUrl: String?
    let date: String?
    let kind: String?

    public init(
        title: String,
        summary: String? = nil,
        imageUrl: String? = nil,
        authorName: String? = nil,
        authorImageUrl: String? = nil,
        date: String? = nil,
        kind: String? = nil
    ) {
        self.title = title
        self.summary = summary
        self.imageUrl = imageUrl
        self.authorName = authorName
        self.authorImageUrl = authorImageUrl
        self.date = date
        self.kind = kind
    }

    public var body: some View {
        VStack(alignment: .leading, spacing: 0) {
            if imageUrl != nil {
                RemoteImage(urlString: imageUrl, cornerRadius: 0)
                    .frame(height: 200)
                    .frame(maxWidth: .infinity)
                    .clipped()
            }

            VStack(alignment: .leading, spacing: AYASpacing.sm) {
                HStack(spacing: AYASpacing.sm) {
                    if let kind {
                        Text(kind.capitalized)
                            .font(AYATypography.caption2)
                            .fontWeight(.bold)
                            .foregroundStyle(.white)
                            .padding(.horizontal, 8)
                            .padding(.vertical, 3)
                            .background(kindColor(kind))
                            .clipShape(Capsule())
                    }

                    Spacer()

                    if let date {
                        Text(formattedDate(date))
                            .font(AYATypography.caption)
                            .foregroundStyle(AYAColors.textTertiary)
                    }
                }

                Text(title)
                    .font(AYATypography.title3)
                    .fontWeight(.semibold)
                    .foregroundStyle(AYAColors.textPrimary)
                    .lineLimit(2)
                    .fixedSize(horizontal: false, vertical: true)

                if let summary, !summary.isEmpty {
                    Text(summary)
                        .font(AYATypography.subheadline)
                        .foregroundStyle(AYAColors.textSecondary)
                        .lineLimit(2)
                }

                if let authorName {
                    HStack(spacing: AYASpacing.sm) {
                        if let authorImageUrl {
                            RemoteImage(urlString: authorImageUrl, cornerRadius: 12)
                                .frame(width: 24, height: 24)
                        } else {
                            Image(systemName: "person.circle.fill")
                                .resizable()
                                .frame(width: 24, height: 24)
                                .foregroundStyle(AYAColors.textTertiary)
                        }

                        Text(authorName)
                            .font(AYATypography.caption)
                            .fontWeight(.medium)
                            .foregroundStyle(AYAColors.textSecondary)
                    }
                }
            }
            .padding(AYASpacing.md)
            .background(AYAColors.surfaceSecondary)
        }
        .background(AYAColors.surfacePrimary)
        .clipShape(RoundedRectangle(cornerRadius: AYACornerRadius.xl))
        .shadow(color: .black.opacity(0.08), radius: 12, x: 0, y: 4)
        .shadow(color: .black.opacity(0.04), radius: 2, x: 0, y: 1)
    }

    private func kindColor(_ kind: String) -> Color {
        switch kind.lowercased() {
        case "article": .blue
        case "story": .indigo
        case "announcement": .orange
        case "activity": .purple
        default: AYAColors.accent
        }
    }

    private func formattedDate(_ iso: String) -> String {
        let formatter = ISO8601DateFormatter()
        formatter.formatOptions = [.withInternetDateTime, .withFractionalSeconds]
        if let date = formatter.date(from: iso) {
            return Self.relativeDateFormatter.localizedString(for: date, relativeTo: Date())
        }
        formatter.formatOptions = [.withInternetDateTime]
        if let date = formatter.date(from: iso) {
            return Self.relativeDateFormatter.localizedString(for: date, relativeTo: Date())
        }
        return iso
    }

    private static let relativeDateFormatter: RelativeDateTimeFormatter = {
        let f = RelativeDateTimeFormatter()
        f.unitsStyle = .abbreviated
        return f
    }()
}
