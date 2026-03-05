import Foundation
import SwiftUI

/// Resolves the user's preferred locale to a supported API locale.
public enum LocaleHelper: Sendable {
    /// All locales supported by the app.
    public static let supportedLocales: [(code: String, name: String, flag: String)] = [
        ("en", "English", "\u{1F1FA}\u{1F1F8}"),
        ("tr", "Türkçe", "\u{1F1F9}\u{1F1F7}"),
        ("fr", "Français", "\u{1F1EB}\u{1F1F7}"),
        ("de", "Deutsch", "\u{1F1E9}\u{1F1EA}"),
        ("es", "Español", "\u{1F1EA}\u{1F1F8}"),
        ("pt-PT", "Português (Portugal)", "\u{1F1F5}\u{1F1F9}"),
        ("it", "Italiano", "\u{1F1EE}\u{1F1F9}"),
        ("nl", "Nederlands", "\u{1F1F3}\u{1F1F1}"),
        ("ja", "日本語", "\u{1F1EF}\u{1F1F5}"),
        ("ko", "한국어", "\u{1F1F0}\u{1F1F7}"),
        ("ru", "Русский", "\u{1F1F7}\u{1F1FA}"),
        ("zh-CN", "简体中文", "\u{1F1E8}\u{1F1F3}"),
        ("ar", "العربية", "\u{1F1F8}\u{1F1E6}"),
    ]

    /// RTL locale codes.
    public static let rtlLocales: Set<String> = ["ar"]

    /// Returns true if the given locale code is right-to-left.
    public static func isRTL(_ code: String) -> Bool {
        rtlLocales.contains(code)
    }

    /// Returns the layout direction for a locale code.
    public static func layoutDirection(for code: String) -> LayoutDirection {
        isRTL(code) ? .rightToLeft : .leftToRight
    }

    /// Returns the user's preferred language if supported, otherwise `"en"`.
    public static var currentLocale: String {
        let preferred = Locale.preferredLanguages.first ?? "en"
        for locale in supportedLocales {
            let prefix = locale.code.split(separator: "-").first.map(String.init) ?? locale.code
            if preferred.hasPrefix(prefix) {
                return locale.code
            }
        }
        return "en"
    }

    /// Returns the display name for a locale code.
    public static func displayName(for code: String) -> String {
        supportedLocales.first { $0.code == code }?.name ?? code
    }

    /// Returns the flag emoji for a locale code.
    public static func flag(for code: String) -> String {
        supportedLocales.first { $0.code == code }?.flag ?? "🏳️"
    }

    /// Maps API locale codes to Apple `.lproj` bundle identifiers where they differ.
    private static let bundleLocaleMap: [String: String] = [
        "zh-CN": "zh-Hans",
    ]

    /// Returns the Apple bundle identifier for a given API locale code.
    static func bundleLocale(for code: String) -> String {
        bundleLocaleMap[code] ?? code
    }

    /// Returns a localized string from the correct `.lproj` bundle for runtime language switching.
    public static func localized(_ key: String, defaultValue: String, locale: String) -> String {
        let resource = bundleLocale(for: locale)
        guard let path = Bundle.main.path(forResource: resource, ofType: "lproj"),
              let bundle = Bundle(path: path) else {
            return defaultValue
        }
        let result = bundle.localizedString(forKey: key, value: defaultValue, table: nil)
        return result
    }
}
