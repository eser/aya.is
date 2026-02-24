import SwiftUI

/// A placeholder avatar icon used when no profile image is available.
public struct DefaultAvatarView: View {
    var size: CGFloat = 48

    /// Creates a default avatar view with the given size.
    public init(size: CGFloat = 48) {
        self.size = size
    }

    /// The avatar icon, tinted with the accent color.
    public var body: some View {
        Image(systemName: "person.crop.circle.fill")
            .resizable()
            .aspectRatio(contentMode: .fit)
            .foregroundStyle(AYAColors.accent.opacity(0.5))
            .frame(width: size, height: size)
    }
}

/// An asynchronous image view that loads from a remote URL with loading and fallback states.
public struct RemoteImage: View {
    let url: URL?
    var cornerRadius: CGFloat = AYACornerRadius.md

    /// Creates a remote image from an optional URL string.
    public init(urlString: String?, cornerRadius: CGFloat = AYACornerRadius.md) {
        self.url = urlString.flatMap { URL(string: $0) }
        self.cornerRadius = cornerRadius
    }

    /// The image view with loading spinner, success, and fallback phases.
    public var body: some View {
        AsyncImage(url: url) { phase in
            switch phase {
            case .success(let image):
                image
                    .resizable()
                    .aspectRatio(contentMode: .fill)
            case .failure:
                fallbackView
            default:
                if url != nil {
                    placeholder
                        .overlay { ProgressView() }
                } else {
                    fallbackView
                }
            }
        }
        .clipShape(RoundedRectangle(cornerRadius: cornerRadius))
    }

    private var placeholder: some View {
        Rectangle()
            .fill(AYAColors.contentBackground)
    }

    private var fallbackView: some View {
        Rectangle()
            .fill(AYAColors.contentBackground)
            .overlay {
                Image(systemName: "person.crop.circle.fill")
                    .resizable()
                    .aspectRatio(contentMode: .fit)
                    .foregroundStyle(AYAColors.accent.opacity(0.4))
                    .padding(8)
            }
    }
}
