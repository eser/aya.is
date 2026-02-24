import SwiftUI

// MARK: - Loading View

/// A full-screen loading placeholder composed of skeleton cards.
public struct AYALoadingView: View {
    /// Creates a loading view.
    public init() {}

    /// Two stacked skeleton cards centered vertically.
    public var body: some View {
        VStack(spacing: AYASpacing.md) {
            Spacer()
            SkeletonCardView()
            SkeletonCardView()
            Spacer()
        }
        .padding(AYASpacing.md)
        .frame(maxWidth: .infinity)
        .accessibilityLabel("Loading content")
        .accessibilityAddTraits(.updatesFrequently)
    }
}

// MARK: - Error View

/// A centered error view with a warning icon, message, and retry button.
public struct AYAErrorView: View {
    let message: String
    let retryLabel: String
    let onRetry: () -> Void

    /// Creates an error view.
    public init(message: String, retryLabel: String = "Try Again", onRetry: @escaping () -> Void) {
        self.message = message
        self.retryLabel = retryLabel
        self.onRetry = onRetry
    }

    /// The error layout: icon, message text, and retry button.
    public var body: some View {
        VStack(spacing: AYASpacing.md) {
            Spacer()
            Image(systemName: "exclamationmark.triangle")
                .font(.system(size: 48))
                .foregroundStyle(AYAColors.warning)
                .accessibilityHidden(true)
            Text(message)
                .font(AYATypography.body)
                .foregroundStyle(AYAColors.textSecondary)
                .multilineTextAlignment(.center)
            Button(retryLabel, action: onRetry)
                .buttonStyle(.bordered)
            Spacer()
        }
        .frame(maxWidth: .infinity)
        .padding()
        .accessibilityElement(children: .combine)
        .accessibilityLabel("Error: \(message)")
    }
}

// MARK: - Empty View

/// A centered empty-state view with an icon and descriptive title.
public struct AYAEmptyView: View {
    let title: String
    let systemImage: String

    /// Creates an empty-state view.
    public init(title: String, systemImage: String = "tray") {
        self.title = title
        self.systemImage = systemImage
    }

    /// The empty-state layout: icon and title centered vertically.
    public var body: some View {
        VStack(spacing: AYASpacing.md) {
            Spacer()
            Image(systemName: systemImage)
                .font(.system(size: 48))
                .foregroundStyle(AYAColors.textTertiary)
                .accessibilityHidden(true)
            Text(title)
                .font(AYATypography.body)
                .foregroundStyle(AYAColors.textSecondary)
            Spacer()
        }
        .frame(maxWidth: .infinity)
        .accessibilityElement(children: .combine)
        .accessibilityLabel(title)
    }
}

// MARK: - Markdown View

/// A convenience wrapper that renders Markdown content using ``RichContentView``.
public struct MarkdownView: View {
    let content: String

    /// Creates a Markdown view from raw Markdown text.
    public init(content: String) {
        self.content = content
    }

    /// The rendered Markdown output.
    public var body: some View {
        RichContentView(content: content)
    }
}
