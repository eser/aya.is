import XCTest
import SwiftUI
import SnapshotTesting
@testable import AYAKit

@MainActor
final class ProductCardSnapshotTests: XCTestCase {

    override func invokeTest() {
        withSnapshotTesting(record: .missing) {
            super.invokeTest()
        }
    }

    func testProductCard_complete() {
        let card = ProductCard(
            title: "AYA Platform",
            description: "An open-source community platform for collaboration and knowledge sharing.",
            imageUrl: nil
        )
        assertSwiftUISnapshot(of: card)
    }

    func testProductCard_noDescription() {
        let card = ProductCard(
            title: "Simple Product",
            imageUrl: nil
        )
        assertSwiftUISnapshot(of: card)
    }

    func testProductCard_longContent() {
        let card = ProductCard(
            title: "Product With a Very Long Name That Extends Beyond One Line",
            description: "This is a very detailed description that goes on for quite a while, explaining every feature and benefit of this amazing product in great detail.",
            imageUrl: nil
        )
        assertSwiftUISnapshot(of: card)
    }
}
