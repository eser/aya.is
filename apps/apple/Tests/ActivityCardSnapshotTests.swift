import XCTest
import SwiftUI
import SnapshotTesting
@testable import AYAKit

@MainActor
final class ActivityCardSnapshotTests: XCTestCase {

    override func invokeTest() {
        withSnapshotTesting(record: .missing) {
            super.invokeTest()
        }
    }

    func testActivityCard_complete() {
        let card = ActivityCard(
            title: "Swift Workshop: Building Apps",
            summary: "Learn how to build modern Swift apps",
            imageUrl: nil,
            activityKind: "workshop",
            timeStart: "2025-03-15T14:00:00Z",
            timeEnd: "2025-03-15T16:00:00Z"
        )
        assertSwiftUISnapshot(of: card)
    }

    func testActivityCard_minimal() {
        let card = ActivityCard(title: "Simple Activity")
        assertSwiftUISnapshot(of: card)
    }

    func testActivityCard_meetup() {
        let card = ActivityCard(
            title: "Community Meetup",
            activityKind: "meetup",
            timeStart: "2025-06-20T18:30"
        )
        assertSwiftUISnapshot(of: card)
    }

    func testActivityCard_conference() {
        let card = ActivityCard(
            title: "Annual Tech Conference 2025",
            summary: "The biggest tech event of the year",
            activityKind: "conference",
            timeStart: "2025-09-01T09:00:00Z",
            timeEnd: "2025-09-03T17:00:00Z"
        )
        assertSwiftUISnapshot(of: card)
    }
}
