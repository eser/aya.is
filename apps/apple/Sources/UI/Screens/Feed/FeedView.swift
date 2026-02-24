import SwiftUI

/// The main feed screen displaying stories, activities, profiles, and search results.
public struct FeedView: View {
    @Bindable var viewModel: FeedViewModel
    @Environment(\.locale) private var appLocale

    /// Creates a feed view backed by the given view model.
    public init(viewModel: FeedViewModel) {
        self.viewModel = viewModel
    }

    public var body: some View {
        VStack(spacing: 0) {
            // Search bar
            AYASearchBar(
                text: Binding(
                    get: { viewModel.searchQuery },
                    set: { viewModel.searchQuery = $0; viewModel.onSearchChanged() }
                ),
                placeholder: String(localized: "feed.search.placeholder", defaultValue: "Search stories, people, products...", locale: appLocale)
            )
            .padding(.horizontal, AYASpacing.md)
            .padding(.vertical, AYASpacing.sm)

            // Filter chips
            FilterChipBar(
                chips: [
                    FilterChip(id: FeedFilter.stories.rawValue, label: String(localized: "feed.filter.stories", defaultValue: "Stories", locale: appLocale)),
                    FilterChip(id: FeedFilter.activities.rawValue, label: String(localized: "feed.filter.activities", defaultValue: "Activities", locale: appLocale)),
                    FilterChip(id: FeedFilter.people.rawValue, label: String(localized: "feed.filter.people", defaultValue: "People", locale: appLocale)),
                    FilterChip(id: FeedFilter.products.rawValue, label: String(localized: "feed.filter.products", defaultValue: "Products", locale: appLocale)),
                ],
                selectedID: Binding(
                    get: { viewModel.activeFilter?.rawValue },
                    set: { id in
                        Task {
                            await viewModel.setFilter(id.flatMap { FeedFilter(rawValue: $0) })
                        }
                    }
                )
            )
            .padding(.bottom, AYASpacing.sm)

            // Content
            ZStack {
                if !viewModel.searchQuery.trimmingCharacters(in: .whitespaces).isEmpty {
                    searchContent
                        .transition(.opacity)
                } else if viewModel.isLoading && viewModel.items.isEmpty {
                    AYALoadingView()
                        .transition(.opacity)
                } else if let error = viewModel.error, viewModel.items.isEmpty {
                    AYAErrorView(
                        message: error,
                        retryLabel: String(localized: "feed.error.retry", defaultValue: "Try Again", locale: appLocale)
                    ) {
                        Task { await viewModel.refresh() }
                    }
                    .transition(.opacity)
                } else if viewModel.items.isEmpty {
                    AYAEmptyView(
                        title: String(localized: "feed.empty", defaultValue: "No content found", locale: appLocale),
                        systemImage: "newspaper"
                    )
                    .transition(.opacity)
                } else {
                    feedContent
                        .transition(.opacity)
                }
            }
            .animation(.easeInOut(duration: 0.2), value: viewModel.searchQuery.isEmpty)
            .animation(.easeInOut(duration: 0.2), value: viewModel.isSearching)
        }
        .task { await viewModel.load() }
    }

    // MARK: - Feed Content

    private static let gridBreakpoint: CGFloat = 700

    private var feedContent: some View {
        GeometryReader { geo in
            let useGrid = geo.size.width >= Self.gridBreakpoint
            ScrollView {
                if useGrid {
                    gridFeedContent(width: geo.size.width)
                } else {
                    singleColumnFeedContent
                }
            }
            .refreshable { await viewModel.refresh() }
        }
    }

    private var singleColumnFeedContent: some View {
        LazyVStack(spacing: AYASpacing.md) {
            featuredHeroCard
            let displayItems = feedDisplayItems
            ForEach(displayItems) { item in
                feedItemView(item)
                    .onAppear {
                        if item == displayItems.last {
                            Task { await viewModel.loadMore() }
                        }
                    }
            }
            loadingIndicator
        }
        .padding(AYASpacing.md)
    }

    private func gridFeedContent(width: CGFloat) -> some View {
        let columns = [
            GridItem(.flexible(), spacing: AYASpacing.md),
            GridItem(.flexible(), spacing: AYASpacing.md),
        ]
        let displayItems = feedDisplayItems

        return LazyVStack(spacing: AYASpacing.md) {
            featuredHeroCard

            LazyVGrid(columns: columns, spacing: AYASpacing.md) {
                ForEach(displayItems) { item in
                    feedItemView(item)
                        .onAppear {
                            if item == displayItems.last {
                                Task { await viewModel.loadMore() }
                            }
                        }
                }
            }

            loadingIndicator
        }
        .padding(AYASpacing.md)
    }

    @ViewBuilder
    private var featuredHeroCard: some View {
        if viewModel.activeFilter == nil,
           case .story(let featured) = viewModel.items.first {
            NavigationLink(value: FeedDestination.story(featured)) {
                FeaturedStoryCard(
                    title: featured.title,
                    summary: featured.summary,
                    imageUrl: featured.storyPictureUri,
                    authorName: featured.authorProfile?.title
                )
            }
            .buttonStyle(.plain)
        }
    }

    private var feedDisplayItems: [FeedItem] {
        let startIndex = (viewModel.activeFilter == nil) ? 1 : 0
        return Array(viewModel.items.dropFirst(startIndex))
    }

    @ViewBuilder
    private var loadingIndicator: some View {
        if viewModel.isLoading {
            ProgressView()
                .padding()
        }
    }

    // MARK: - Search Content

    private var searchContent: some View {
        ScrollView {
            LazyVStack(spacing: AYASpacing.sm) {
                if viewModel.isSearching && viewModel.searchResults.isEmpty {
                    ProgressView()
                        .frame(maxWidth: .infinity)
                        .padding(.top, AYASpacing.xxl)
                } else if viewModel.searchResults.isEmpty && !viewModel.isSearching {
                    AYAEmptyView(
                        title: String(localized: "feed.empty", defaultValue: "No content found", locale: appLocale),
                        systemImage: "magnifyingglass"
                    )
                    .padding(.top, AYASpacing.xl)
                } else {
                    ForEach(viewModel.searchResults) { result in
                        NavigationLink(value: FeedDestination.search(result)) {
                            searchResultRow(result)
                        }
                        .buttonStyle(.plain)
                        .onAppear {
                            if result == viewModel.searchResults.last {
                                Task { await viewModel.loadMore() }
                            }
                        }
                    }

                    if viewModel.isSearching {
                        ProgressView()
                            .padding()
                    }
                }
            }
            .padding(AYASpacing.md)
        }
    }

    // MARK: - Item Views

    @ViewBuilder
    private func feedItemView(_ item: FeedItem) -> some View {
        switch item {
        case .story(let story):
            NavigationLink(value: FeedDestination.story(story)) {
                StoryCard(
                    title: story.title,
                    summary: story.summary,
                    imageUrl: story.storyPictureUri,
                    authorName: story.authorProfile?.title,
                    authorImageUrl: story.authorProfile?.profilePictureUri,
                    date: story.publishedAt,
                    kind: story.kind
                )
            }
            .buttonStyle(.plain)
            .accessibilityHint("Double tap to view details")

        case .activity(let activity):
            NavigationLink(value: FeedDestination.activity(activity)) {
                ActivityCard(
                    title: activity.title,
                    summary: activity.summary,
                    imageUrl: activity.storyPictureUri,
                    activityKind: activity.properties?.activityKind,
                    timeStart: activity.properties?.activityTimeStart,
                    timeEnd: activity.properties?.activityTimeEnd
                )
            }
            .buttonStyle(.plain)
            .accessibilityHint("Double tap to view details")

        case .profile(let profile):
            Button {
                viewModel.selectedProfile = profile
            } label: {
                ProfileCard(
                    title: profile.title,
                    description: profile.description,
                    imageUrl: profile.profilePictureUri,
                    kind: profile.kind,
                    points: profile.points
                )
            }
            .buttonStyle(.plain)
            .accessibilityHint("Double tap to view profile")

        case .product(let profile):
            Button {
                viewModel.selectedProfile = profile
            } label: {
                ProductCard(
                    title: profile.title,
                    description: profile.description,
                    imageUrl: profile.profilePictureUri
                )
            }
            .buttonStyle(.plain)
            .accessibilityHint("Double tap to view profile")
        }
    }

    private func searchResultRow(_ result: SearchResult) -> some View {
        HStack(spacing: AYASpacing.md) {
            RemoteImage(urlString: result.imageUri, cornerRadius: AYACornerRadius.md)
                .frame(width: 48, height: 48)

            VStack(alignment: .leading, spacing: AYASpacing.xs) {
                HStack(spacing: AYASpacing.sm) {
                    Text(result.type.capitalized)
                        .font(AYATypography.caption2)
                        .foregroundStyle(.white)
                        .padding(.horizontal, 6)
                        .padding(.vertical, 2)
                        .background(colorForType(result.type))
                        .clipShape(Capsule())

                    if let kind = result.kind {
                        Text(kind.capitalized)
                            .font(AYATypography.caption2)
                            .foregroundStyle(AYAColors.textTertiary)
                    }
                }

                Text(result.title)
                    .font(AYATypography.headline)
                    .foregroundStyle(AYAColors.textPrimary)
                    .lineLimit(1)

                if let summary = result.summary, !summary.isEmpty {
                    Text(summary)
                        .font(AYATypography.caption)
                        .foregroundStyle(AYAColors.textSecondary)
                        .lineLimit(2)
                }
            }

            Spacer()
        }
        .padding(AYASpacing.md)
        .background(AYAColors.surfacePrimary)
        .clipShape(RoundedRectangle(cornerRadius: AYACornerRadius.lg))
        .shadow(color: .black.opacity(0.04), radius: 6, x: 0, y: 2)
    }

    private func colorForType(_ type: String) -> Color {
        switch type {
        case "story": .blue
        case "profile": .green
        case "page": .orange
        default: .gray
        }
    }
}

// MARK: - Navigation Destination

/// Navigation destinations reachable from the feed.
public enum FeedDestination: Hashable {
    /// Navigate to a story detail screen.
    case story(Story)
    /// Navigate to an activity detail screen.
    case activity(Activity)
    /// Navigate to a profile detail screen.
    case profile(Profile)
    /// Navigate to the detail screen for a search result.
    case search(SearchResult)
}

// MARK: - Navigation Container

/// Root navigation container that wraps `FeedView` in a `NavigationStack` with toolbar, locale, and theme controls.
public struct FeedNavigationView: View {
    @Bindable var viewModel: FeedViewModel
    @AppStorage("preferredLocale") private var preferredLocale: String = LocaleHelper.currentLocale
    @AppStorage("appColorScheme") private var appColorScheme: String = "system"
    @Environment(\.colorScheme) private var systemColorScheme

    /// Creates the navigation view with the given feed view model.
    public init(viewModel: FeedViewModel) {
        self.viewModel = viewModel
    }

    public var body: some View {
        NavigationStack {
            FeedView(viewModel: viewModel)
                .navigationTitle(String(localized: "app.title", defaultValue: "AYA", locale: currentLocale))
                #if os(iOS)
                .navigationBarTitleDisplayMode(.large)
                #endif
                #if os(macOS)
                .toolbarTitleDisplayMode(.inline)
                #endif
                #if os(macOS)
                .navigationSubtitle(String(localized: "app.subtitle", defaultValue: "Open Software Network", locale: currentLocale))
                #endif
                .toolbar {
                    ToolbarItem(placement: .automatic) {
                        HStack(spacing: AYASpacing.sm) {
                            Button {
                                let isEffectivelyDark = (resolvedColorScheme ?? systemColorScheme) == .dark
                                appColorScheme = isEffectivelyDark ? "light" : "dark"
                            } label: {
                                Image(systemName: appearanceIcon)
                                    .imageScale(.medium)
                                    .contentTransition(.symbolEffect(.replace))
                            }
                            .accessibilityLabel("Toggle appearance")
                            .accessibilityHint("Switches between light and dark mode")

                            Menu {
                                ForEach(LocaleHelper.supportedLocales, id: \.code) { locale in
                                    Button {
                                        setLocale(locale.code)
                                    } label: {
                                        HStack {
                                            Text("\(locale.flag) \(locale.name)")
                                            if preferredLocale == locale.code {
                                                Image(systemName: "checkmark")
                                            }
                                        }
                                    }
                                }
                            } label: {
                                Text(LocaleHelper.flag(for: preferredLocale))
                                    .font(.body)
                                    .padding(.horizontal, 6)
                                    .padding(.vertical, 2)
                                    .background(AYAColors.accentSubtle)
                                    .clipShape(Capsule())
                            }
                            .accessibilityLabel("Change language, current: \(LocaleHelper.displayName(for: preferredLocale))")
                        }
                    }
                }
                .navigationDestination(for: FeedDestination.self) { destination in
                    destinationView(destination)
                }
                .sheet(item: $viewModel.selectedProfile) { profile in
                    NavigationStack {
                        ProfileDetailView(viewModel: ProfileDetailViewModel(
                            slug: profile.slug,
                            client: viewModel.client,
                            locale: viewModel.locale,
                            initialProfile: profile
                        ))
                        .toolbar {
                            ToolbarItem(placement: .confirmationAction) {
                                Button {
                                    viewModel.selectedProfile = nil
                                } label: {
                                    Image(systemName: "xmark.circle.fill")
                                        .foregroundStyle(AYAColors.textTertiary)
                                }
                            }
                        }
                    }
                    .presentationDetents([.medium, .large])
                    .presentationDragIndicator(.visible)
                }
        }
        .environment(\.locale, currentLocale)
        .preferredColorScheme(resolvedColorScheme)
        .onAppear {
            if viewModel.locale != preferredLocale {
                Task { await viewModel.switchLocale(preferredLocale) }
            }
        }
    }

    private var currentLocale: Locale {
        Locale(identifier: preferredLocale)
    }

    private var appearanceIcon: String {
        let isEffectivelyDark = (resolvedColorScheme ?? systemColorScheme) == .dark
        return isEffectivelyDark ? "moon.fill" : "sun.max.fill"
    }

    private var resolvedColorScheme: ColorScheme? {
        switch appColorScheme {
        case "light": .light
        case "dark": .dark
        default: nil
        }
    }

    private func setLocale(_ locale: String) {
        guard locale != preferredLocale else { return }
        preferredLocale = locale
        Task { await viewModel.switchLocale(locale) }
    }

    @ViewBuilder
    private func destinationView(_ destination: FeedDestination) -> some View {
        switch destination {
        case .story(let story):
            StoryDetailView(viewModel: StoryDetailViewModel(
                slug: story.slug,
                client: viewModel.client,
                locale: viewModel.locale,
                initialStory: story
            ))

        case .activity(let activity):
            ActivityDetailView(viewModel: ActivityDetailViewModel(
                slug: activity.slug,
                client: viewModel.client,
                locale: viewModel.locale,
                initialActivity: activity
            ))

        case .profile:
            EmptyView()

        case .search(let result):
            searchResultDestination(result)
        }
    }

    @ViewBuilder
    private func searchResultDestination(_ result: SearchResult) -> some View {
        switch result.type {
        case "story":
            StoryDetailView(viewModel: StoryDetailViewModel(
                slug: result.slug,
                client: viewModel.client,
                locale: viewModel.locale
            ))
        case "profile":
            ProfileDetailView(viewModel: ProfileDetailViewModel(
                slug: result.slug,
                client: viewModel.client,
                locale: viewModel.locale
            ))
        default:
            if let profileSlug = result.profileSlug {
                ProfileDetailView(viewModel: ProfileDetailViewModel(
                    slug: profileSlug,
                    client: viewModel.client,
                    locale: viewModel.locale
                ))
            } else {
                Text(result.title)
                    .font(.title)
                    .padding()
            }
        }
    }
}
