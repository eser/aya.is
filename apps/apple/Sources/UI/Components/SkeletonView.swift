import SwiftUI

/// An animated shimmer placeholder used while content is loading.
public struct SkeletonView: View {
    var height: CGFloat

    @State private var shimmerOffset: CGFloat = -1

    /// Creates a skeleton view with the given height.
    public init(height: CGFloat = 20) {
        self.height = height
    }

    /// The shimmer rectangle that animates across the view.
    public var body: some View {
        RoundedRectangle(cornerRadius: AYACornerRadius.sm)
            .fill(AYAColors.contentBackground)
            .frame(height: height)
            .overlay {
                GeometryReader { geo in
                    LinearGradient(
                        colors: [.clear, AYAColors.surfacePrimary.opacity(0.5), .clear],
                        startPoint: .leading,
                        endPoint: .trailing
                    )
                    .frame(width: geo.size.width * 0.4)
                    .offset(x: shimmerOffset * geo.size.width)
                }
                .clipped()
            }
            .clipShape(RoundedRectangle(cornerRadius: AYACornerRadius.sm))
            .onAppear {
                withAnimation(AYAAnimation.shimmer) {
                    shimmerOffset = 1.5
                }
            }
    }
}

/// A composite skeleton placeholder that mimics the layout of a typical content card.
public struct SkeletonCardView: View {
    /// Creates a skeleton card view.
    public init() {}

    /// The skeleton card layout: image area, title bar, subtitle bar, and metadata row.
    public var body: some View {
        VStack(alignment: .leading, spacing: AYASpacing.sm) {
            SkeletonView(height: 180)
            SkeletonView(height: 20)
            SkeletonView(height: 14)
                .frame(maxWidth: 200)
            HStack {
                SkeletonView(height: 14)
                    .frame(width: 80)
                Spacer()
                SkeletonView(height: 14)
                    .frame(width: 60)
            }
        }
        .padding(AYASpacing.md)
        .background(AYAColors.surfacePrimary)
        .clipShape(RoundedRectangle(cornerRadius: AYACornerRadius.lg))
    }
}
