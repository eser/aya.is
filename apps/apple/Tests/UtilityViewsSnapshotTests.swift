import XCTest
import SwiftUI
import SnapshotTesting
@testable import AYAKit

@MainActor
final class UtilityViewsSnapshotTests: XCTestCase {

    override func invokeTest() {
        withSnapshotTesting(record: .missing) {
            super.invokeTest()
        }
    }

    func testLoadingView() {
        assertSwiftUISnapshot(of: AYALoadingView(), height: 300)
    }

    func testErrorView() {
        let view = AYAErrorView(
            message: "Something went wrong. Please check your connection and try again."
        ) {}
        assertSwiftUISnapshot(of: view, height: 300)
    }

    func testEmptyView_default() {
        assertSwiftUISnapshot(of: AYAEmptyView(title: "No items found"), height: 300)
    }

    func testEmptyView_customIcon() {
        assertSwiftUISnapshot(of: AYAEmptyView(title: "No stories yet", systemImage: "newspaper"), height: 300)
    }

    func testMarkdownView() {
        let view = MarkdownView(content: "**Bold text** and *italic text* with a [link](https://example.com)")
        assertSwiftUISnapshot(of: view)
    }
}
