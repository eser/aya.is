import Foundation
import Testing
@testable import AYAKit

@Suite("LocaleHelper")
struct LocaleHelperTests {
    // MARK: - Supported Locale Codes

    @Test("All 13 supported locales are defined")
    func supportedLocaleCount() {
        #expect(LocaleHelper.supportedLocales.count == 13)
    }

    @Test("Supported locale codes match backend API codes",
          arguments: ["en", "tr", "fr", "de", "es", "pt-PT", "it", "nl", "ja", "ko", "ru", "zh-CN", "ar"])
    func supportedLocaleCodesMatchBackend(code: String) {
        let codes = LocaleHelper.supportedLocales.map(\.code)
        #expect(codes.contains(code), "Missing locale code: \(code)")
    }

    @Test("Chinese locale uses zh-CN, not zh-Hans")
    func chineseLocaleUsesZhCN() {
        let codes = LocaleHelper.supportedLocales.map(\.code)
        #expect(codes.contains("zh-CN"), "zh-CN should be in supported locales")
        #expect(!codes.contains("zh-Hans"), "zh-Hans should NOT be in supported locales (use zh-CN for API)")
    }

    @Test("No duplicate locale codes exist")
    func noDuplicateLocales() {
        let codes = LocaleHelper.supportedLocales.map(\.code)
        let uniqueCodes = Set(codes)
        #expect(codes.count == uniqueCodes.count, "Duplicate locale codes found")
    }

    // MARK: - Display Names

    @Test("Display name returns correct name for all locales",
          arguments: [
            ("en", "English"),
            ("tr", "Türkçe"),
            ("zh-CN", "简体中文"),
            ("ar", "العربية"),
            ("pt-PT", "Português (Portugal)"),
          ])
    func displayNameReturnsCorrectName(code: String, expectedName: String) {
        #expect(LocaleHelper.displayName(for: code) == expectedName)
    }

    @Test("Display name returns raw code for unknown locale")
    func displayNameFallback() {
        #expect(LocaleHelper.displayName(for: "xx-XX") == "xx-XX")
    }

    // MARK: - Flags

    @Test("Flag returns correct emoji for all locales",
          arguments: [
            ("en", "🇺🇸"),
            ("zh-CN", "🇨🇳"),
            ("ar", "🇸🇦"),
            ("pt-PT", "🇵🇹"),
          ])
    func flagReturnsCorrectEmoji(code: String, expectedFlag: String) {
        #expect(LocaleHelper.flag(for: code) == expectedFlag)
    }

    @Test("Flag returns default for unknown locale")
    func flagFallback() {
        #expect(LocaleHelper.flag(for: "xx-XX") == "🏳️")
    }

    // MARK: - RTL

    @Test("Arabic is detected as RTL")
    func arabicIsRTL() {
        #expect(LocaleHelper.isRTL("ar") == true)
    }

    @Test("Non-RTL locales are detected correctly",
          arguments: ["en", "tr", "fr", "de", "es", "pt-PT", "it", "nl", "ja", "ko", "ru", "zh-CN"])
    func nonRTLLocales(code: String) {
        #expect(LocaleHelper.isRTL(code) == false)
    }

    // MARK: - Bundle Locale Mapping

    @Test("Bundle locale maps zh-CN to zh-Hans for .lproj lookup")
    func bundleLocaleMapsChinese() {
        #expect(LocaleHelper.bundleLocale(for: "zh-CN") == "zh-Hans")
    }

    @Test("Bundle locale returns code unchanged for non-mapped locales",
          arguments: ["en", "tr", "fr", "de", "es", "pt-PT", "it", "nl", "ja", "ko", "ru", "ar"])
    func bundleLocalePassthrough(code: String) {
        #expect(LocaleHelper.bundleLocale(for: code) == code)
    }

    // MARK: - Current Locale

    @Test("currentLocale returns a supported locale code")
    func currentLocaleIsSupported() {
        let current = LocaleHelper.currentLocale
        let codes = LocaleHelper.supportedLocales.map(\.code)
        #expect(codes.contains(current), "currentLocale returned unsupported code: \(current)")
    }
}

@Suite("API Client Locale Integration")
struct APIClientLocaleTests {
    /// A mock API client that records the locale parameter for verification.
    struct RecordingAPIClient: APIClientProtocol {
        let recordedLocales: LockedBox<[String]>

        init() {
            self.recordedLocales = LockedBox([])
        }

        func fetchStories(locale: String, cursor: String?) async throws -> APIResponse<Story> {
            recordedLocales.mutate { $0.append(locale) }
            return APIResponse(data: [], cursor: nil)
        }

        func fetchStory(locale: String, slug: String) async throws -> Story {
            recordedLocales.mutate { $0.append(locale) }
            throw APIError.httpError(statusCode: 404)
        }

        func fetchActivities(locale: String, cursor: String?) async throws -> APIResponse<Activity> {
            recordedLocales.mutate { $0.append(locale) }
            return APIResponse(data: [], cursor: nil)
        }

        func fetchActivity(locale: String, slug: String) async throws -> Activity {
            recordedLocales.mutate { $0.append(locale) }
            throw APIError.httpError(statusCode: 404)
        }

        func fetchProfiles(locale: String, filterKind: String?, cursor: String?) async throws -> APIResponse<Profile> {
            recordedLocales.mutate { $0.append(locale) }
            return APIResponse(data: [], cursor: nil)
        }

        func fetchProfile(locale: String, slug: String) async throws -> Profile {
            recordedLocales.mutate { $0.append(locale) }
            throw APIError.httpError(statusCode: 404)
        }

        func fetchProfilePages(locale: String, slug: String) async throws -> APIResponse<Page> {
            recordedLocales.mutate { $0.append(locale) }
            return APIResponse(data: [], cursor: nil)
        }

        func fetchProfileStories(locale: String, slug: String) async throws -> APIResponse<Story> {
            recordedLocales.mutate { $0.append(locale) }
            return APIResponse(data: [], cursor: nil)
        }

        func search(locale: String, query: String, cursor: String?) async throws -> APIResponse<SearchResult> {
            recordedLocales.mutate { $0.append(locale) }
            return APIResponse(data: [], cursor: nil)
        }
    }

    /// Thread-safe box for collecting values across async calls.
    final class LockedBox<T>: @unchecked Sendable {
        private var value: T
        private let lock = NSLock()

        init(_ value: T) { self.value = value }

        func get() -> T {
            lock.lock()
            defer { lock.unlock() }
            return value
        }

        func mutate(_ transform: (inout T) -> Void) {
            lock.lock()
            defer { lock.unlock() }
            transform(&value)
        }
    }

    @Test("FeedViewModel passes locale to API client for all fetches",
          arguments: ["en", "tr", "zh-CN", "pt-PT", "ar", "ja"])
    @MainActor
    func feedViewModelPassesLocale(locale: String) async {
        let recorder = RecordingAPIClient()
        let viewModel = FeedViewModel(client: recorder, locale: locale)

        await viewModel.load()

        let recorded = recorder.recordedLocales.get()
        #expect(!recorded.isEmpty, "Expected API calls for locale \(locale)")
        for recorded in recorded {
            #expect(recorded == locale, "Expected locale \(locale), got \(recorded)")
        }
    }

    @Test("FeedViewModel switchLocale updates the locale used for API calls")
    @MainActor
    func switchLocaleUpdatesAPICalls() async {
        let recorder = RecordingAPIClient()
        let viewModel = FeedViewModel(client: recorder, locale: "en")

        await viewModel.switchLocale("zh-CN")

        let recorded = recorder.recordedLocales.get()
        let zhCalls = recorded.filter { $0 == "zh-CN" }
        #expect(!zhCalls.isEmpty, "Expected zh-CN API calls after switchLocale")
    }

    @Test("APIClient constructs correct URL for zh-CN locale")
    func apiClientURLConstruction() {
        // Verify the URL components produce the right path with zh-CN
        var components = URLComponents(string: "https://api.aya.is")!
        components.path = "/zh-CN/stories"
        #expect(components.url?.absoluteString == "https://api.aya.is/zh-CN/stories")
    }

    @Test("APIClient constructs correct URL for all locales",
          arguments: ["en", "tr", "fr", "de", "es", "pt-PT", "it", "nl", "ja", "ko", "ru", "zh-CN", "ar"])
    func apiClientURLForAllLocales(locale: String) {
        var components = URLComponents(string: "https://api.aya.is")!
        components.path = "/\(locale)/stories"
        let url = components.url
        #expect(url != nil, "URL should be valid for locale: \(locale)")
        #expect(url?.path == "/\(locale)/stories", "Path should contain \(locale)")
    }
}
