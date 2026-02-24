import XCTest
import SwiftUI
import SnapshotTesting
@testable import AYAKit

@MainActor
final class FeaturedStoryCardSnapshotTests: XCTestCase {

    override func invokeTest() {
        withSnapshotTesting(record: .missing) {
            super.invokeTest()
        }
    }

    func testFeaturedStoryCard_complete() {
        let card = FeaturedStoryCard(
            title: "The Future of Open Source Communities",
            summary: "How decentralized collaboration is reshaping software development worldwide.",
            imageUrl: nil,
            authorName: "Alp Özcan",
            kind: "article"
        )
        assertSwiftUISnapshot(of: card, height: 320)
    }

    func testFeaturedStoryCard_minimal() {
        let card = FeaturedStoryCard(title: "Minimal Featured Story")
        assertSwiftUISnapshot(of: card, height: 320)
    }

    func testFeaturedStoryCard_longTitle() {
        let card = FeaturedStoryCard(
            title: "This Is an Extremely Long Title That Should Demonstrate How the Featured Card Handles Multi-Line Text Overflow",
            summary: "A brief summary underneath the long title.",
            kind: "announcement"
        )
        assertSwiftUISnapshot(of: card, height: 320)
    }
}
