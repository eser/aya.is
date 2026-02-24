import SwiftUI

public struct ProfileCard: View {
    let title: String
    let description: String?
    let imageUrl: String?
    let kind: String?
    let points: Int?

    public init(
        title: String,
        description: String? = nil,
        imageUrl: String? = nil,
        kind: String? = nil,
        points: Int? = nil
    ) {
        self.title = title
        self.description = description
        self.imageUrl = imageUrl
        self.kind = kind
        self.points = points
    }

    public var body: some View {
        HStack(spacing: AYASpacing.md) {
            RemoteImage(urlString: imageUrl, cornerRadius: 22)
                .frame(width: 44, height: 44)

            VStack(alignment: .leading, spacing: 3) {
                Text(title)
                    .font(AYATypography.headline)
                    .foregroundStyle(AYAColors.textPrimary)
                    .lineLimit(1)

                if let description, !description.isEmpty {
                    Text(description)
                        .font(AYATypography.caption)
                        .foregroundStyle(AYAColors.textSecondary)
                        .lineLimit(1)
                }
            }

            Spacer()

            VStack(alignment: .trailing, spacing: 3) {
                if let kind {
                    Text(kind.capitalized)
                        .font(AYATypography.caption2)
                        .fontWeight(.medium)
                        .foregroundStyle(kindColor(kind))
                        .padding(.horizontal, 7)
                        .padding(.vertical, 2)
                        .background(kindColor(kind).opacity(0.12))
                        .clipShape(Capsule())
                }
                if let points, points > 0 {
                    HStack(spacing: 2) {
                        Image(systemName: "star.fill")
                            .font(.system(size: 9))
                        Text("\(points)")
                            .font(AYATypography.caption2)
                            .fontWeight(.bold)
                    }
                    .foregroundStyle(AYAColors.accent)
                }
            }
        }
        .padding(AYASpacing.md)
        .background(AYAColors.surfacePrimary)
        .clipShape(RoundedRectangle(cornerRadius: AYACornerRadius.xl))
        .shadow(color: .black.opacity(0.06), radius: 8, x: 0, y: 3)
        .shadow(color: .black.opacity(0.03), radius: 1, x: 0, y: 1)
    }

    private func kindColor(_ kind: String) -> Color {
        switch kind.lowercased() {
        case "individual": .blue
        case "organization": .green
        case "product": .purple
        default: AYAColors.accent
        }
    }
}
