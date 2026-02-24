import SwiftUI

@Observable @MainActor
public final class ProfileDetailViewModel {
    var profile: Profile?
    var pages: [Page] = []
    var stories: [Story] = []
    var activities: [Activity] = []
    var isLoading = false
    var error: String?
    var selectedTab: ProfileTab = .pages
    let slug: String
    let client: APIClientProtocol
    private let locale: String

    public enum ProfileTab: String, CaseIterable, Sendable {
        case pages = "Pages"
        case stories = "Stories"
        case activities = "Activities"
    }

    public init(
        slug: String,
        client: APIClientProtocol,
        locale: String = LocaleHelper.currentLocale,
        initialProfile: Profile? = nil
    ) {
        self.slug = slug
        self.client = client
        self.locale = locale
        self.profile = initialProfile
    }

    func load() async {
        isLoading = true
        do {
            async let profileFetch = client.fetchProfile(locale: locale, slug: slug)
            async let pagesFetch = client.fetchProfilePages(locale: locale, slug: slug)
            async let storiesFetch = client.fetchProfileStories(locale: locale, slug: slug)

            profile = try await profileFetch
            pages = try await pagesFetch.data
            stories = try await storiesFetch.data
        } catch {
            self.error = error.localizedDescription
        }
        isLoading = false
    }
}

public struct ProfileDetailView: View {
    @Bindable var viewModel: ProfileDetailViewModel
    @Environment(\.locale) private var appLocale

    public init(viewModel: ProfileDetailViewModel) {
        self.viewModel = viewModel
    }

    public var body: some View {
        ScrollView {
            if let profile = viewModel.profile {
                VStack(alignment: .leading, spacing: AYASpacing.md) {
                    HStack(spacing: AYASpacing.md) {
                        RemoteImage(urlString: profile.profilePictureUri, cornerRadius: 32)
                            .frame(width: 72, height: 72)

                        VStack(alignment: .leading, spacing: AYASpacing.xs) {
                            Text(profile.title)
                                .font(AYATypography.title2)
                                .fontWeight(.bold)
                                .foregroundStyle(AYAColors.textPrimary)

                            Text(profile.kind.capitalized)
                                .font(AYATypography.caption)
                                .foregroundStyle(.white)
                                .padding(.horizontal, 8)
                                .padding(.vertical, 3)
                                .background(AYAColors.accent)
                                .clipShape(Capsule())

                            if profile.points > 0 {
                                Text("\(profile.points) \(String(localized: "profile.points", defaultValue: "points", locale: appLocale))")
                                    .font(AYATypography.caption)
                                    .foregroundStyle(AYAColors.textSecondary)
                            }
                        }
                    }

                    if !profile.description.isEmpty {
                        RichContentView(content: profile.description)
                    }

                    Picker("Section", selection: $viewModel.selectedTab) {
                        ForEach(ProfileDetailViewModel.ProfileTab.allCases, id: \.self) { tab in
                            Text(localizedTabLabel(tab, locale: appLocale)).tag(tab)
                        }
                    }
                    .pickerStyle(.segmented)

                    switch viewModel.selectedTab {
                    case .pages:
                        if viewModel.pages.isEmpty {
                            Text(String(localized: "feed.empty", defaultValue: "No content found", locale: appLocale))
                                .font(AYATypography.body)
                                .foregroundStyle(AYAColors.textTertiary)
                        } else {
                            ForEach(viewModel.pages) { page in
                                VStack(alignment: .leading, spacing: AYASpacing.sm) {
                                    Text(page.title)
                                        .font(AYATypography.headline)
                                        .fontWeight(.semibold)
                                    if let content = page.content {
                                        RichContentView(content: content)
                                    }
                                }
                            }
                        }

                    case .stories:
                        if viewModel.stories.isEmpty {
                            Text(String(localized: "feed.empty", defaultValue: "No content found", locale: appLocale))
                                .font(AYATypography.body)
                                .foregroundStyle(AYAColors.textTertiary)
                        } else {
                            ForEach(viewModel.stories) { story in
                                StoryCard(
                                    title: story.title,
                                    summary: story.summary,
                                    imageUrl: story.storyPictureUri,
                                    date: story.publishedAt,
                                    kind: story.kind
                                )
                            }
                        }

                    case .activities:
                        Text(String(localized: "feed.empty", defaultValue: "No content found", locale: appLocale))
                            .font(AYATypography.body)
                            .foregroundStyle(AYAColors.textTertiary)
                    }
                }
                .padding(AYASpacing.md)
            } else if viewModel.isLoading {
                AYALoadingView()
            } else if let error = viewModel.error {
                AYAErrorView(message: error) { Task { await viewModel.load() } }
            }
        }
        .navigationTitle(viewModel.profile?.title ?? String(localized: "detail.profile", defaultValue: "Profile", locale: appLocale))
        #if os(iOS)
        .navigationBarTitleDisplayMode(.inline)
        #endif
        .task { await viewModel.load() }
    }

    private func localizedTabLabel(_ tab: ProfileDetailViewModel.ProfileTab, locale: Locale) -> String {
        switch tab {
        case .pages: String(localized: "profile.pages", defaultValue: "Pages", locale: locale)
        case .stories: String(localized: "profile.stories", defaultValue: "Stories", locale: locale)
        case .activities: String(localized: "profile.activities", defaultValue: "Activities", locale: locale)
        }
    }
}
