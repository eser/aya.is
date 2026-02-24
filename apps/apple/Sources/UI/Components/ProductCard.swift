import SwiftUI

/// A compact horizontal card that displays a product with its image, title, and description.
public struct ProductCard: View {
    let title: String
    let description: String?
    let imageUrl: String?

    /// Creates a product card.
    public init(title: String, description: String? = nil, imageUrl: String? = nil) {
        self.title = title
        self.description = description
        self.imageUrl = imageUrl
    }

    @AppStorage("preferredLocale") private var preferredLocale: String = LocaleHelper.currentLocale

    /// The card layout: product image, title, description, and a "Product" badge.
    public var body: some View {
        HStack(spacing: AYASpacing.md) {
            RemoteImage(urlString: imageUrl, cornerRadius: AYACornerRadius.lg)
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

            Text(LocaleHelper.localized("product.badge", defaultValue: "Product", locale: preferredLocale))
                .font(AYATypography.caption2)
                .fontWeight(.medium)
                .foregroundStyle(.purple)
                .padding(.horizontal, 7)
                .padding(.vertical, 2)
                .background(.purple.opacity(0.12))
                .clipShape(Capsule())
        }
        .padding(AYASpacing.md)
        .accessibilityElement(children: .combine)
        .accessibilityLabel(productAccessibilityLabel)
        .background(AYAColors.surfacePrimary)
        .clipShape(RoundedRectangle(cornerRadius: AYACornerRadius.xl))
        .shadow(color: .black.opacity(0.06), radius: 8, x: 0, y: 3)
        .shadow(color: .black.opacity(0.03), radius: 1, x: 0, y: 1)
    }

    private var productAccessibilityLabel: String {
        var parts: [String] = ["Product", title]
        if let description { parts.append(description) }
        return parts.joined(separator: ", ")
    }
}
