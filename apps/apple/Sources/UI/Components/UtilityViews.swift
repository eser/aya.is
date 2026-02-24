import SwiftUI

// MARK: - Loading View

public struct AYALoadingView: View {
    public init() {}

    public var body: some View {
        VStack(spacing: AYASpacing.md) {
            Spacer()
            SkeletonCardView()
            SkeletonCardView()
            Spacer()
        }
        .padding(AYASpacing.md)
        .frame(maxWidth: .infinity)
    }
}

// MARK: - Error View

public struct AYAErrorView: View {
    let message: String
    let retryLabel: String
    let onRetry: () -> Void

    public init(message: String, retryLabel: String = "Try Again", onRetry: @escaping () -> Void) {
        self.message = message
        self.retryLabel = retryLabel
        self.onRetry = onRetry
    }

    public var body: some View {
        VStack(spacing: AYASpacing.md) {
            Spacer()
            Image(systemName: "exclamationmark.triangle")
                .font(.system(size: 48))
                .foregroundStyle(AYAColors.warning)
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
    }
}

// MARK: - Empty View

public struct AYAEmptyView: View {
    let title: String
    let systemImage: String

    public init(title: String, systemImage: String = "tray") {
        self.title = title
        self.systemImage = systemImage
    }

    public var body: some View {
        VStack(spacing: AYASpacing.md) {
            Spacer()
            Image(systemName: systemImage)
                .font(.system(size: 48))
                .foregroundStyle(AYAColors.textTertiary)
            Text(title)
                .font(AYATypography.body)
                .foregroundStyle(AYAColors.textSecondary)
            Spacer()
        }
        .frame(maxWidth: .infinity)
    }
}

// MARK: - Markdown View

public struct MarkdownView: View {
    let content: String

    public init(content: String) {
        self.content = content
    }

    public var body: some View {
        RichContentView(content: content)
    }
}
