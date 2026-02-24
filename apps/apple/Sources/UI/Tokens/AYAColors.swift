import SwiftUI

public enum AYAColors {
    // MARK: - Backgrounds
    #if os(macOS)
    public static let windowBackground = Color(nsColor: .windowBackgroundColor)
    public static let contentBackground = Color(nsColor: .controlBackgroundColor)
    public static let groupedBackground = Color(nsColor: .underPageBackgroundColor)
    #else
    public static let windowBackground = Color(uiColor: .systemBackground)
    public static let contentBackground = Color(uiColor: .secondarySystemBackground)
    public static let groupedBackground = Color(uiColor: .systemGroupedBackground)
    #endif

    // MARK: - Text
    public static let textPrimary = Color.primary
    public static let textSecondary = Color.secondary
    public static let textTertiary = Color.secondary.opacity(0.7)

    // MARK: - Accent
    public static let accent = Color.accentColor
    public static let accentSubtle = Color.accentColor.opacity(0.15)

    // MARK: - Surfaces
    #if os(macOS)
    public static let surfacePrimary = Color(nsColor: .controlBackgroundColor)
    public static let surfaceSecondary = Color(nsColor: .windowBackgroundColor)
    #else
    public static let surfacePrimary = Color(uiColor: .systemBackground)
    public static let surfaceSecondary = Color(uiColor: .secondarySystemBackground)
    #endif

    // MARK: - Borders
    #if os(macOS)
    public static let border = Color(nsColor: .separatorColor)
    public static let borderSubtle = Color(nsColor: .separatorColor).opacity(0.5)
    #else
    public static let border = Color(uiColor: .separator)
    public static let borderSubtle = Color(uiColor: .separator).opacity(0.5)
    #endif

    // MARK: - Status
    public static let success = Color.green
    public static let warning = Color.orange
    public static let error = Color.red
    public static let info = Color.blue
}
