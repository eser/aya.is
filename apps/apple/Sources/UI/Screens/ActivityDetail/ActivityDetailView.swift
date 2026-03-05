import SwiftUI

/// View model that loads and holds the full content of a single activity.
@Observable @MainActor
public final class ActivityDetailViewModel {
    var activity: Activity?
    var isLoading = false
    var error: String?
    let slug: String
    let client: APIClientProtocol
    private let locale: String

    /// Creates an activity detail view model.
    /// - Parameters:
    ///   - slug: The activity's URL slug.
    ///   - client: The API client for fetching activity content.
    ///   - locale: The locale for content localization.
    ///   - initialActivity: An optional pre-loaded activity to display immediately.
    public init(
        slug: String,
        client: APIClientProtocol,
        locale: String = LocaleHelper.currentLocale,
        initialActivity: Activity? = nil
    ) {
        self.slug = slug
        self.client = client
        self.locale = locale
        self.activity = initialActivity
    }

    func load() async {
        guard activity?.content == nil || activity?.content?.isEmpty == true else { return }
        isLoading = true
        do {
            activity = try await client.fetchActivity(locale: locale, slug: slug)
        } catch {
            self.error = error.localizedDescription
        }
        isLoading = false
    }
}

public struct ActivityDetailView: View {
    @Bindable var viewModel: ActivityDetailViewModel
    @AppStorage("preferredLocale") private var preferredLocale: String = LocaleHelper.currentLocale

    public init(viewModel: ActivityDetailViewModel) {
        self.viewModel = viewModel
    }

    public var body: some View {
        ScrollView {
            if let activity = viewModel.activity {
                VStack(alignment: .leading, spacing: AYASpacing.md) {
                    Rectangle()
                        .fill(kindColor(activity.properties?.activityKind))
                        .frame(height: 6)

                    VStack(alignment: .leading, spacing: AYASpacing.md) {
                        if let imageUrl = activity.storyPictureUri, let url = URL(string: imageUrl) {
                            AsyncImage(url: url) { phase in
                                switch phase {
                                case .success(let image):
                                    image
                                        .resizable()
                                        .aspectRatio(contentMode: .fit)
                                        .frame(maxWidth: .infinity)
                                        .clipShape(RoundedRectangle(cornerRadius: AYACornerRadius.md))
                                default:
                                    Color.clear.frame(height: 0)
                                }
                            }
                        }

                        Text(activity.title)
                            .font(AYATypography.largeTitle)
                            .fontWeight(.bold)
                            .foregroundStyle(AYAColors.textPrimary)

                        if let props = activity.properties {
                            HStack(spacing: AYASpacing.md) {
                                if let kind = props.activityKind {
                                    Label(kind.capitalized, systemImage: "tag")
                                        .font(AYATypography.subheadline)
                                        .foregroundStyle(kindColor(kind))
                                }
                                if let start = props.activityTimeStart {
                                    Label(formatTime(start), systemImage: "clock")
                                        .font(AYATypography.subheadline)
                                        .foregroundStyle(AYAColors.textSecondary)
                                }
                            }

                            if let uri = props.externalActivityUri, let url = URL(string: uri) {
                                Link(destination: url) {
                                    Label(LocaleHelper.localized("activity.join", defaultValue: "Join Activity", locale: preferredLocale), systemImage: "arrow.up.right.square")
                                }
                                .buttonStyle(.bordered)
                                .accessibilityHint("Opens external link to join the activity")
                            }

                            if let uri = props.externalAttendanceUri, let url = URL(string: uri) {
                                Link(destination: url) {
                                    Label(LocaleHelper.localized("activity.rsvp", defaultValue: "RSVP", locale: preferredLocale), systemImage: "person.badge.plus")
                                }
                                .buttonStyle(.bordered)
                                .accessibilityHint("Opens external link to RSVP")
                            }
                        }

                        if let content = activity.content, !content.isEmpty {
                            RichContentView(content: content)
                        }
                    }
                    .padding(.horizontal, AYASpacing.md)
                    .padding(.bottom, AYASpacing.lg)
                }
            } else if viewModel.isLoading {
                AYALoadingView()
            } else if let error = viewModel.error {
                AYAErrorView(message: error) { Task { await viewModel.load() } }
            }
        }
        .navigationTitle(viewModel.activity?.title ?? LocaleHelper.localized("detail.activity", defaultValue: "Activity", locale: preferredLocale))
        #if os(iOS)
        .navigationBarTitleDisplayMode(.inline)
        #endif
        .task { await viewModel.load() }
    }

    private func kindColor(_ kind: String?) -> Color {
        switch kind?.lowercased() {
        case "workshop": .purple
        case "meetup": .blue
        case "conference": .orange
        case "webinar": .teal
        default: AYAColors.accent
        }
    }

    private func formatTime(_ iso: String) -> String {
        let formatter = ISO8601DateFormatter()
        formatter.formatOptions = [.withInternetDateTime, .withFractionalSeconds]
        var date = formatter.date(from: iso)
        if date == nil {
            let short = DateFormatter()
            short.dateFormat = "yyyy-MM-dd'T'HH:mm"
            short.timeZone = .current
            date = short.date(from: iso)
        }
        guard let date else { return iso }
        let display = DateFormatter()
        display.dateStyle = .medium
        display.timeStyle = .short
        return display.string(from: date)
    }
}
