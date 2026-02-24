import Foundation

/// Resolves the user's preferred locale to a supported API locale.
public enum LocaleHelper: Sendable {
    /// Returns `"tr"` if the user's preferred language is Turkish, otherwise `"en"`.
    public static var currentLocale: String {
        let preferred = Locale.preferredLanguages.first ?? "en"
        if preferred.hasPrefix("tr") { return "tr" }
        return "en"
    }
}
