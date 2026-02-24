import Foundation

public enum LocaleHelper: Sendable {
    public static var currentLocale: String {
        let preferred = Locale.preferredLanguages.first ?? "en"
        if preferred.hasPrefix("tr") { return "tr" }
        return "en"
    }
}
