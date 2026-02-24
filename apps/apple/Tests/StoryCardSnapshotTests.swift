import XCTest
import SwiftUI
import SnapshotTesting
@testable import AYAKit

@MainActor
final class StoryCardSnapshotTests: XCTestCase {

    override func invokeTest() {
        withSnapshotTesting(record: .missing) {
            super.invokeTest()
        }
    }

    func testStoryCard_fullFields() {
        let card = StoryCard(
            title: "Building a Better Tomorrow",
            summary: "An exploration of how sustainable technology can reshape our communities and create lasting change for future generations.",
            imageUrl: nil,
            authorName: "Jane Doe",
            authorImageUrl: nil,
            date: "2025-01-15T10:00:00Z",
            kind: "story"
        )
        assertSwiftUISnapshot(of: card)
    }

    func testStoryCard_minimal() {
        let card = StoryCard(title: "Minimal Story")
        assertSwiftUISnapshot(of: card)
    }

    func testStoryCard_noImage() {
        let card = StoryCard(
            title: "Story Without Image",
            summary: "This story has a summary but no image.",
            imageUrl: nil,
            authorName: "Author",
            date: "2025-06-01T12:00:00Z",
            kind: "article"
        )
        assertSwiftUISnapshot(of: card)
    }

    func testStoryCard_longContent() {
        let card = StoryCard(
            title: "This is a very long story title that should be truncated after two lines of text to keep the card compact",
            summary: "This is a very long summary that goes on and on describing the story in great detail, covering many aspects and providing extensive context about the topic at hand.",
            imageUrl: nil,
            authorName: "A Very Long Author Name That Might Overflow",
            authorImageUrl: nil,
            date: "2025-12-31T23:59:59Z",
            kind: "essay"
        )
        assertSwiftUISnapshot(of: card)
    }
}
