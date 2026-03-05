import XCTest
import SwiftUI
import SnapshotTesting
@testable import AYAKit

@MainActor
final class ProfileCardSnapshotTests: XCTestCase {

    override func invokeTest() {
        let original = NSTimeZone.default
        NSTimeZone.default = TimeZone(identifier: "UTC")!
        defer { NSTimeZone.default = original }
        withSnapshotTesting(record: .missing) {
            super.invokeTest()
        }
    }

    func testProfileCard_organization() {
        let card = ProfileCard(
            title: "AYA Foundation",
            description: "Building open-source tools for community empowerment and education.",
            imageUrl: nil,
            kind: "organization",
            points: 1250
        )
        assertSwiftUISnapshot(of: card)
    }

    func testProfileCard_individual() {
        let card = ProfileCard(
            title: "Jane Developer",
            description: "Full-stack engineer passionate about Swift and open source.",
            imageUrl: nil,
            kind: "individual",
            points: 350
        )
        assertSwiftUISnapshot(of: card)
    }

    func testProfileCard_noPoints() {
        let card = ProfileCard(
            title: "New Member",
            description: "Just joined the community.",
            imageUrl: nil,
            kind: "individual",
            points: 0
        )
        assertSwiftUISnapshot(of: card)
    }

    func testProfileCard_minimal() {
        let card = ProfileCard(title: "Minimal Profile")
        assertSwiftUISnapshot(of: card)
    }
}
