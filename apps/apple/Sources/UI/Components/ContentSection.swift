import SwiftUI

/// A titled section with an optional "See all" action, wrapping arbitrary content.
public struct ContentSection<Content: View>: View {
    let title: String
    var seeAllAction: (() -> Void)?
    var seeAllLabel: String = "See all"
    @ViewBuilder let content: () -> Content

    /// Creates a content section.
    public init(
        title: String,
        seeAllAction: (() -> Void)? = nil,
        seeAllLabel: String = "See all",
        @ViewBuilder content: @escaping () -> Content
    ) {
        self.title = title
        self.seeAllAction = seeAllAction
        self.seeAllLabel = seeAllLabel
        self.content = content
    }

    /// The section layout: header row with title and optional action, followed by the content.
    public var body: some View {
        VStack(alignment: .leading, spacing: AYASpacing.sm) {
            HStack {
                Text(title)
                    .font(AYATypography.title3)
                    .fontWeight(.bold)
                    .foregroundStyle(AYAColors.textPrimary)

                Spacer()

                if let seeAllAction {
                    Button(seeAllLabel, action: seeAllAction)
                        .font(AYATypography.subheadline)
                        .foregroundStyle(AYAColors.accent)
                }
            }

            content()
        }
    }
}
