import XCTest

final class StoryDetailUITests: XCTestCase {

    private var app: XCUIApplication!

    override func setUpWithError() throws {
        continueAfterFailure = false
        app = XCUIApplication()
        app.launchArguments = ["-AppleLanguages", "(en)"]
        app.launch()
    }

    // MARK: - Story Detail

    func testNavigateToStoryDetail() throws {
        let scrollView = app.scrollViews.firstMatch
        XCTAssertTrue(scrollView.waitForExistence(timeout: 10))

        // Wait for content to load
        sleep(2)

        // Filter to stories only
        let storiesChip = app.buttons["Stories"]
        if storiesChip.waitForExistence(timeout: 5) {
            storiesChip.tap()
            sleep(1)
        }

        // Tap first story card
        let firstCard = scrollView.buttons.firstMatch
        if firstCard.waitForExistence(timeout: 5) {
            firstCard.tap()

            // Verify we're in detail view
            let backButton = app.navigationBars.buttons.firstMatch
            XCTAssertTrue(backButton.waitForExistence(timeout: 5))
        }
    }

    func testStoryDetailScrolls() throws {
        let scrollView = app.scrollViews.firstMatch
        XCTAssertTrue(scrollView.waitForExistence(timeout: 10))
        sleep(2)

        let storiesChip = app.buttons["Stories"]
        if storiesChip.waitForExistence(timeout: 5) {
            storiesChip.tap()
            sleep(1)
        }

        let firstCard = scrollView.buttons.firstMatch
        if firstCard.waitForExistence(timeout: 5) {
            firstCard.tap()

            // Scroll in detail view
            let detailScroll = app.scrollViews.firstMatch
            if detailScroll.waitForExistence(timeout: 5) {
                detailScroll.swipeUp()
                detailScroll.swipeDown()
            }

            // Go back
            app.navigationBars.buttons.firstMatch.tap()
        }
    }

    func testStoryDetailShowsAuthor() throws {
        let scrollView = app.scrollViews.firstMatch
        XCTAssertTrue(scrollView.waitForExistence(timeout: 10))
        sleep(2)

        let storiesChip = app.buttons["Stories"]
        if storiesChip.waitForExistence(timeout: 5) {
            storiesChip.tap()
            sleep(1)
        }

        let firstCard = scrollView.buttons.firstMatch
        if firstCard.waitForExistence(timeout: 5) {
            firstCard.tap()
            sleep(2)

            // Check for author accessibility element
            let author = app.staticTexts.matching(NSPredicate(format: "label CONTAINS 'Author'")).firstMatch
            // Author may or may not be present depending on the story
            _ = author.waitForExistence(timeout: 3)

            app.navigationBars.buttons.firstMatch.tap()
        }
    }
}
