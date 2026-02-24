import SwiftUI

/// A card view that displays an activity or event with date, time, and an optional image.
public struct ActivityCard: View {
    @Environment(\.referenceDate) private var referenceDate
    @Environment(\.locale) private var locale

    let title: String
    let summary: String?
    let imageUrl: String?
    let activityKind: String?
    let timeStart: String?
    let timeEnd: String?

    /// Creates an activity card.
    public init(
        title: String,
        summary: String? = nil,
        imageUrl: String? = nil,
        activityKind: String? = nil,
        timeStart: String? = nil,
        timeEnd: String? = nil
    ) {
        self.title = title
        self.summary = summary
        self.imageUrl = imageUrl
        self.activityKind = activityKind
        self.timeStart = timeStart
        self.timeEnd = timeEnd
    }

    /// The card layout: image or calendar placeholder, kind badge, time, title, and summary.
    public var body: some View {
        VStack(alignment: .leading, spacing: 0) {
            if imageUrl != nil {
                RemoteImage(urlString: imageUrl, cornerRadius: 0)
                    .frame(height: 200)
                    .frame(maxWidth: .infinity)
                    .clipped()
            } else {
                ZStack {
                    Rectangle()
                        .fill(kindColor.opacity(0.12))

                    VStack(spacing: AYASpacing.sm) {
                        Image(systemName: "calendar.badge.clock")
                            .font(.system(size: 36))
                            .foregroundStyle(kindColor.opacity(0.5))

                        if let timeStart {
                            VStack(spacing: 2) {
                                Text(dayString(timeStart))
                                    .font(.system(size: 28, weight: .bold, design: .rounded))
                                    .foregroundStyle(kindColor)
                                Text(monthString(timeStart))
                                    .font(AYATypography.caption)
                                    .fontWeight(.semibold)
                                    .foregroundStyle(AYAColors.textSecondary)
                                    .textCase(.uppercase)
                            }
                        }
                    }
                }
                .frame(height: 200)
                .frame(maxWidth: .infinity)
            }

            VStack(alignment: .leading, spacing: AYASpacing.sm) {
                HStack(spacing: AYASpacing.sm) {
                    if let activityKind {
                        Text(activityKind.capitalized)
                            .font(AYATypography.caption2)
                            .fontWeight(.bold)
                            .foregroundStyle(.white)
                            .padding(.horizontal, 8)
                            .padding(.vertical, 3)
                            .background(kindColor)
                            .clipShape(Capsule())
                    }

                    Spacer()

                    if let timeStart {
                        HStack(spacing: 3) {
                            Image(systemName: "clock")
                                .font(.system(size: 10))
                            Text(timeString(timeStart))
                                .font(AYATypography.caption)
                        }
                        .foregroundStyle(AYAColors.textTertiary)
                    }
                }

                Text(title)
                    .font(AYATypography.title3)
                    .fontWeight(.semibold)
                    .foregroundStyle(AYAColors.textPrimary)
                    .lineLimit(2)
                    .fixedSize(horizontal: false, vertical: true)

                if let summary, !summary.isEmpty {
                    Text(summary)
                        .font(AYATypography.subheadline)
                        .foregroundStyle(AYAColors.textSecondary)
                        .lineLimit(2)
                }
            }
            .padding(AYASpacing.md)
            .frame(maxHeight: .infinity, alignment: .top)
            .background(AYAColors.surfaceSecondary)
        }
        .accessibilityElement(children: .combine)
        .accessibilityLabel(activityAccessibilityLabel)
        .background(AYAColors.surfacePrimary)
        .clipShape(RoundedRectangle(cornerRadius: AYACornerRadius.xl))
        .shadow(color: .black.opacity(0.08), radius: 12, x: 0, y: 4)
        .shadow(color: .black.opacity(0.04), radius: 2, x: 0, y: 1)
    }

    private var activityAccessibilityLabel: String {
        var parts: [String] = []
        if let activityKind { parts.append(activityKind.capitalized) }
        parts.append(title)
        if let summary { parts.append(summary) }
        if let timeStart { parts.append("Starts \(timeString(timeStart))") }
        return parts.joined(separator: ", ")
    }

    private var kindColor: Color {
        switch activityKind?.lowercased() {
        case "workshop": .purple
        case "meetup": .blue
        case "conference": .orange
        case "webinar": .teal
        default: AYAColors.accent
        }
    }

    private func parseDate(_ iso: String) -> Date? {
        let formatter = ISO8601DateFormatter()
        formatter.formatOptions = [.withInternetDateTime, .withFractionalSeconds]
        if let date = formatter.date(from: iso) { return date }
        formatter.formatOptions = [.withInternetDateTime]
        if let date = formatter.date(from: iso) { return date }
        let short = DateFormatter()
        short.dateFormat = "yyyy-MM-dd'T'HH:mm"
        short.timeZone = .current
        return short.date(from: iso)
    }

    private func dayString(_ iso: String) -> String {
        guard let date = parseDate(iso) else { return "" }
        let f = DateFormatter()
        f.locale = locale
        f.dateFormat = "d"
        return f.string(from: date)
    }

    private func monthString(_ iso: String) -> String {
        guard let date = parseDate(iso) else { return "" }
        let f = DateFormatter()
        f.locale = locale
        f.dateFormat = "MMM"
        return f.string(from: date)
    }

    private func timeString(_ iso: String) -> String {
        guard let date = parseDate(iso) else { return "" }
        let f = DateFormatter()
        f.locale = locale
        f.timeStyle = .short
        return f.string(from: date)
    }
}
