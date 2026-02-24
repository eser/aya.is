import XCTest
import SwiftUI
import SnapshotTesting
@testable import AYAKit

@MainActor
final class RTLSnapshotTests: XCTestCase {

    override func invokeTest() {
        let original = NSTimeZone.default
        NSTimeZone.default = TimeZone(identifier: "UTC")!
        defer { NSTimeZone.default = original }
        withSnapshotTesting(record: .missing) {
            super.invokeTest()
        }
    }

    // MARK: - RTL Layout Tests

    func testStoryCard_rtl() {
        let card = StoryCard(
            title: "كيفية تقليل تأخير البلوتوث",
            summary: "ستهدف هذه المدونة إلى تقليل تأخير صوت البلوتوث من خلال تحديد برامج الترميز",
            authorName: "أحمد محمد",
            kind: "مقال"
        )
        assertSwiftUISnapshot(
            of: card.environment(\.layoutDirection, .rightToLeft),
            named: "rtl"
        )
    }

    func testActivityCard_rtl() {
        let card = ActivityCard(
            title: "ورشة عمل البرمجة",
            summary: "تعلم أساسيات البرمجة مع المحترفين",
            activityKind: "ورشة عمل",
            timeStart: "2026-03-15T10:00:00Z"
        )
        assertSwiftUISnapshot(
            of: card.environment(\.layoutDirection, .rightToLeft),
            named: "rtl"
        )
    }

    func testProfileCard_rtl() {
        let card = ProfileCard(
            title: "أحمد محمد",
            description: "مطور برمجيات في إسطنبول",
            kind: "individual",
            points: 150
        )
        assertSwiftUISnapshot(
            of: card.environment(\.layoutDirection, .rightToLeft),
            named: "rtl"
        )
    }

    func testProductCard_rtl() {
        let card = ProductCard(
            title: "منصة AYA",
            description: "منصة مجتمعية مفتوحة المصدر للتعاون ومشاركة المعرفة"
        )
        assertSwiftUISnapshot(
            of: card.environment(\.layoutDirection, .rightToLeft),
            named: "rtl"
        )
    }

    func testFeaturedStoryCard_rtl() {
        let card = FeaturedStoryCard(
            title: "مستقبل مجتمعات المصادر المفتوحة",
            summary: "كيف يعيد التعاون اللامركزي تشكيل تطوير البرمجيات في جميع أنحاء العالم",
            authorName: "أحمد محمد",
            kind: "مقال"
        )
        assertSwiftUISnapshot(
            of: card.environment(\.layoutDirection, .rightToLeft),
            named: "rtl",
            height: 320
        )
    }

    func testSearchBar_rtl() {
        let bar = AYASearchBar(
            text: .constant("بحث"),
            placeholder: "البحث عن مقالات، أشخاص، منتجات..."
        )
        assertSwiftUISnapshot(
            of: bar.environment(\.layoutDirection, .rightToLeft),
            named: "rtl"
        )
    }

    func testFilterChipBar_rtl() {
        let chips = FilterChipBar(
            chips: [
                FilterChip(id: "stories", label: "مقالات"),
                FilterChip(id: "activities", label: "أنشطة"),
                FilterChip(id: "people", label: "أشخاص"),
                FilterChip(id: "products", label: "منتجات"),
            ],
            selectedID: .constant("stories")
        )
        assertSwiftUISnapshot(
            of: chips.environment(\.layoutDirection, .rightToLeft),
            named: "rtl"
        )
    }

    func testUtilityViews_rtl() {
        let errorView = AYAErrorView(
            message: "حدث خطأ ما. يرجى التحقق من الاتصال والمحاولة مرة أخرى.",
            retryLabel: "إعادة المحاولة"
        ) {}
        assertSwiftUISnapshot(
            of: errorView.environment(\.layoutDirection, .rightToLeft),
            named: "rtl",
            height: 300
        )
    }
}

// MARK: - RTL Snapshot Helper

extension RTLSnapshotTests {
    func assertSwiftUISnapshot<V: View>(
        of view: V,
        named name: String? = nil,
        height: CGFloat = 200,
        file: StaticString = #filePath,
        testName: String = #function,
        line: UInt = #line
    ) {
        let localized = view.environment(\.locale, snapshotLocale)
        #if os(iOS)
        let platformSuffix = "iOS"
        let snapshotName = name.map { "\($0)-\(platformSuffix)" } ?? platformSuffix
        let vc = UIHostingController(rootView: localized)
        vc.view.frame = CGRect(x: 0, y: 0, width: 390, height: height)
        assertSnapshot(of: vc, as: .image(on: .init(safeArea: .zero, size: CGSize(width: 390, height: height), traits: .init())), named: snapshotName, file: file, testName: testName, line: line)
        #elseif os(macOS)
        let platformSuffix = "macOS"
        let snapshotName = name.map { "\($0)-\(platformSuffix)" } ?? platformSuffix
        let hostingView = NSHostingView(rootView: localized.frame(width: 390))
        hostingView.frame = NSRect(x: 0, y: 0, width: 390, height: height)
        assertSnapshot(of: hostingView, as: .image, named: snapshotName, file: file, testName: testName, line: line)
        #endif
    }
}
