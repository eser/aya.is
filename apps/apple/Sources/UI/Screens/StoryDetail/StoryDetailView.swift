import SwiftUI

/// View model that loads and holds the full content of a single story.
@Observable @MainActor
public final class StoryDetailViewModel {
    var story: Story?
    var isLoading = false
    var error: String?
    let slug: String
    let client: APIClientProtocol
    private let locale: String

    /// Creates a story detail view model.
    /// - Parameters:
    ///   - slug: The story's URL slug.
    ///   - client: The API client for fetching story content.
    ///   - locale: The locale for content localization.
    ///   - initialStory: An optional pre-loaded story to display immediately.
    public init(
        slug: String,
        client: APIClientProtocol,
        locale: String = LocaleHelper.currentLocale,
        initialStory: Story? = nil
    ) {
        self.slug = slug
        self.client = client
        self.locale = locale
        self.story = initialStory
    }

    func loadStory() async {
        guard story?.content == nil || story?.content?.isEmpty == true else { return }
        isLoading = true
        do {
            story = try await client.fetchStory(locale: locale, slug: slug)
        } catch {
            self.error = error.localizedDescription
        }
        isLoading = false
    }
}

/// Detail screen displaying a story's hero image, author info, and rich content.
public struct StoryDetailView: View {
    @Bindable var viewModel: StoryDetailViewModel
    @AppStorage("preferredLocale") private var preferredLocale: String = LocaleHelper.currentLocale

    /// Creates a story detail view backed by the given view model.
    public init(viewModel: StoryDetailViewModel) {
        self.viewModel = viewModel
    }

    public var body: some View {
        ScrollView {
            if let story = viewModel.story {
                VStack(alignment: .leading, spacing: 0) {
                    if let imageUrl = story.storyPictureUri, let url = URL(string: imageUrl) {
                        AsyncImage(url: url) { phase in
                            switch phase {
                            case .success(let image):
                                image
                                    .resizable()
                                    .aspectRatio(contentMode: .fit)
                                    .frame(maxWidth: .infinity)
                            default:
                                Color.clear.frame(height: 0)
                            }
                        }
                    }

                    VStack(alignment: .leading, spacing: AYASpacing.md) {
                        Text(story.title)
                            .font(AYATypography.largeTitle)
                            .fontWeight(.bold)
                            .foregroundStyle(AYAColors.textPrimary)

                        if let author = story.authorProfile {
                            HStack(spacing: AYASpacing.sm) {
                                RemoteImage(urlString: author.profilePictureUri, cornerRadius: 16)
                                    .frame(width: 32, height: 32)
                                VStack(alignment: .leading, spacing: 2) {
                                    Text(author.title)
                                        .font(AYATypography.subheadline)
                                        .fontWeight(.medium)
                                        .foregroundStyle(AYAColors.textPrimary)
                                    if let date = story.publishedAt {
                                        Text(formattedDate(date))
                                            .font(AYATypography.caption)
                                            .foregroundStyle(AYAColors.textTertiary)
                                    }
                                }
                            }
                            .accessibilityElement(children: .combine)
                            .accessibilityLabel("Author: \(author.title)")
                        }

                        Divider()

                        if let content = story.content, !content.isEmpty {
                            RichContentView(content: content)
                        }
                    }
                    .padding(.horizontal, AYASpacing.lg)
                    .padding(.top, AYASpacing.md)
                    .padding(.bottom, AYASpacing.lg)
                }
            } else if viewModel.isLoading {
                AYALoadingView()
            } else if let error = viewModel.error {
                AYAErrorView(message: error) { Task { await viewModel.loadStory() } }
            }
        }
        .navigationTitle(viewModel.story?.title ?? LocaleHelper.localized("detail.story", defaultValue: "Story", locale: preferredLocale))
        #if os(iOS)
        .navigationBarTitleDisplayMode(.inline)
        #endif
        .task { await viewModel.loadStory() }
    }

    private func formattedDate(_ iso: String) -> String {
        let formatter = ISO8601DateFormatter()
        formatter.formatOptions = [.withInternetDateTime, .withFractionalSeconds]
        guard let date = formatter.date(from: iso) else { return iso }
        let display = DateFormatter()
        display.dateStyle = .medium
        return display.string(from: date)
    }
}
