import SwiftUI

public struct DefaultAvatarView: View {
    var size: CGFloat = 48

    public init(size: CGFloat = 48) {
        self.size = size
    }

    public var body: some View {
        Image(systemName: "person.crop.circle.fill")
            .resizable()
            .aspectRatio(contentMode: .fit)
            .foregroundStyle(AYAColors.accent.opacity(0.5))
            .frame(width: size, height: size)
    }
}

public struct RemoteImage: View {
    let url: URL?
    var cornerRadius: CGFloat = AYACornerRadius.md

    public init(urlString: String?, cornerRadius: CGFloat = AYACornerRadius.md) {
        self.url = urlString.flatMap { URL(string: $0) }
        self.cornerRadius = cornerRadius
    }

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
