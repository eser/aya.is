import Foundation

/// View model that manages feed data loading, filtering, searching, and pagination.
@Observable @MainActor
public final class FeedViewModel {
    // MARK: - Published State

    var items: [FeedItem] = []
    var isLoading = false
    var error: String?
    var searchQuery = ""
    var activeFilter: FeedFilter?
    var searchResults: [SearchResult] = []
    var isSearching = false
    var selectedProfile: Profile?

    // MARK: - Internal State

    private var stories: [Story] = []
    private var activities: [Activity] = []
    private var people: [Profile] = []
    private var products: [Profile] = []

    private var storyCursor: String?
    private var activityCursor: String?
    private var peopleCursor: String?
    private var productCursor: String?
    private var searchCursor: String?

    nonisolated(unsafe) private var searchTask: Task<Void, Never>?
    nonisolated(unsafe) private var paginationTask: Task<Void, Never>?

    let client: APIClientProtocol
    var locale: String

    // MARK: - Init

    /// Creates a new feed view model.
    /// - Parameters:
    ///   - client: The API client used to fetch feed data.
    ///   - locale: The locale for content localization.
    public init(client: APIClientProtocol, locale: String = LocaleHelper.currentLocale) {
        self.client = client
        self.locale = locale
    }

    deinit {
        searchTask?.cancel()
        paginationTask?.cancel()
    }

    // MARK: - Load

    func load() async {
        guard !isLoading else { return }
        isLoading = true
        error = nil
        do {
            async let storiesFetch = client.fetchStories(locale: locale, cursor: nil)
            async let activitiesFetch = client.fetchActivities(locale: locale, cursor: nil)

            let (storiesResp, activitiesResp) = try await (storiesFetch, activitiesFetch)
            stories = storiesResp.data
            storyCursor = storiesResp.cursor
            activities = activitiesResp.data
            activityCursor = activitiesResp.cursor
            rebuildItems()
        } catch {
            self.error = error.localizedDescription
        }
        isLoading = false
    }

    // MARK: - Filter

    func setFilter(_ filter: FeedFilter?) async {
        activeFilter = filter
        searchQuery = ""
        searchResults = []
        isSearching = false

        switch filter {
        case .people where people.isEmpty:
            await loadPeople()
        case .products where products.isEmpty:
            await loadProducts()
        default:
            break
        }
        rebuildItems()
    }

    // MARK: - Search

    func onSearchChanged() {
        searchTask?.cancel()
        let trimmed = searchQuery.trimmingCharacters(in: .whitespaces)
        guard !trimmed.isEmpty else {
            searchResults = []
            isSearching = false
            rebuildItems()
            return
        }
        isSearching = true
        searchTask = Task {
            try? await Task.sleep(for: .milliseconds(400))
            guard !Task.isCancelled else { return }
            await performSearch()
        }
    }

    func performSearch() async {
        let trimmed = searchQuery.trimmingCharacters(in: .whitespaces)
        guard !trimmed.isEmpty else { return }
        isSearching = true
        do {
            let response = try await client.search(locale: locale, query: trimmed, cursor: nil)
            searchResults = response.data
            searchCursor = response.cursor
        } catch {
            self.error = error.localizedDescription
        }
        isSearching = false
    }

    // MARK: - Pagination

    func loadMore() async {
        if isSearching || !searchQuery.trimmingCharacters(in: .whitespaces).isEmpty {
            await loadMoreSearch()
            return
        }
        switch activeFilter {
        case .stories: await loadMoreStories()
        case .activities: await loadMoreActivities()
        case .people: await loadMorePeople()
        case .products: await loadMoreProducts()
        case nil: await loadMoreStories()
        }
    }

    // MARK: - Refresh

    func refresh() async {
        paginationTask?.cancel()
        paginationTask = nil
        stories = []
        activities = []
        people = []
        products = []
        storyCursor = nil
        activityCursor = nil
        peopleCursor = nil
        productCursor = nil
        searchCursor = nil
        items = []
        await load()

        if activeFilter == .people { await loadPeople() }
        if activeFilter == .products { await loadProducts() }
        rebuildItems()
    }

    // MARK: - Locale

    func switchLocale(_ newLocale: String) async {
        locale = newLocale
        await refresh()
    }

    // MARK: - Private Helpers

    private func rebuildItems() {
        var result: [FeedItem] = []

        switch activeFilter {
        case .stories:
            result = stories.map { .story($0) }
        case .activities:
            result = activities.map { .activity($0) }
        case .people:
            result = people.map { .profile($0) }
        case .products:
            result = products.map { .product($0) }
        case nil:
            result = stories.map { .story($0) } + activities.map { .activity($0) }
            result.sort { ($0.date ?? .distantPast) > ($1.date ?? .distantPast) }
        }

        items = result
    }

    private func loadPeople() async {
        do {
            let response = try await client.fetchProfiles(locale: locale, filterKind: "individual", cursor: nil)
            people = response.data
            peopleCursor = response.cursor
        } catch {
            self.error = error.localizedDescription
        }
    }

    private func loadProducts() async {
        do {
            let response = try await client.fetchProfiles(locale: locale, filterKind: "product", cursor: nil)
            products = response.data
            productCursor = response.cursor
        } catch {
            self.error = error.localizedDescription
        }
    }

    private func loadMoreStories() async {
        guard let cursor = storyCursor else { return }
        do {
            let response = try await client.fetchStories(locale: locale, cursor: cursor)
            stories.append(contentsOf: response.data)
            storyCursor = response.cursor
            rebuildItems()
        } catch {
            self.error = error.localizedDescription
        }
    }

    private func loadMoreActivities() async {
        guard let cursor = activityCursor else { return }
        do {
            let response = try await client.fetchActivities(locale: locale, cursor: cursor)
            activities.append(contentsOf: response.data)
            activityCursor = response.cursor
            rebuildItems()
        } catch {
            self.error = error.localizedDescription
        }
    }

    private func loadMorePeople() async {
        guard let cursor = peopleCursor else { return }
        do {
            let response = try await client.fetchProfiles(locale: locale, filterKind: "individual", cursor: cursor)
            people.append(contentsOf: response.data)
            peopleCursor = response.cursor
            rebuildItems()
        } catch {
            self.error = error.localizedDescription
        }
    }

    private func loadMoreProducts() async {
        guard let cursor = productCursor else { return }
        do {
            let response = try await client.fetchProfiles(locale: locale, filterKind: "product", cursor: cursor)
            products.append(contentsOf: response.data)
            productCursor = response.cursor
            rebuildItems()
        } catch {
            self.error = error.localizedDescription
        }
    }

    private func loadMoreSearch() async {
        guard let cursor = searchCursor else { return }
        do {
            let response = try await client.search(locale: locale, query: searchQuery, cursor: cursor)
            searchResults.append(contentsOf: response.data)
            searchCursor = response.cursor
        } catch {
            self.error = error.localizedDescription
        }
    }
}
