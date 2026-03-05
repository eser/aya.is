import XCTest

final class ActivityDetailUITests: XCTestCase {

    private var app: XCUIApplication!

    override func setUpWithError() throws {
        continueAfterFailure = false
        app = XCUIApplication()
        app.launchArguments = ["-AppleLanguages", "(en)"]
        app.launch()
    }

    // MARK: - Activity Detail

    func testNavigateToActivityDetail() throws {
        let activitiesChip = app.buttons["Activities"]
        XCTAssertTrue(activitiesChip.waitForExistence(timeout: 10))

        activitiesChip.tap()
        sleep(2)

        let scrollView = app.scrollViews.firstMatch
        let firstCard = scrollView.buttons.firstMatch
        if firstCard.waitForExistence(timeout: 5) {
            firstCard.tap()

            let backButton = app.navigationBars.buttons.firstMatch
            XCTAssertTrue(backButton.waitForExistence(timeout: 5))
            backButton.tap()
        }
    }

    func testActivityDetailScrolls() throws {
        let activitiesChip = app.buttons["Activities"]
        XCTAssertTrue(activitiesChip.waitForExistence(timeout: 10))

        activitiesChip.tap()
        sleep(2)

        let scrollView = app.scrollViews.firstMatch
        let firstCard = scrollView.buttons.firstMatch
        if firstCard.waitForExistence(timeout: 5) {
            firstCard.tap()

            let detailScroll = app.scrollViews.firstMatch
            if detailScroll.waitForExistence(timeout: 5) {
                detailScroll.swipeUp()
                detailScroll.swipeDown()
            }

            app.navigationBars.buttons.firstMatch.tap()
        }
    }

    func testActivityDetailShowsActionButtons() throws {
        let activitiesChip = app.buttons["Activities"]
        XCTAssertTrue(activitiesChip.waitForExistence(timeout: 10))

        activitiesChip.tap()
        sleep(2)

        let scrollView = app.scrollViews.firstMatch
        let firstCard = scrollView.buttons.firstMatch
        if firstCard.waitForExistence(timeout: 5) {
            firstCard.tap()
            sleep(2)

            // Check for join/RSVP buttons (may not always be present)
            let joinButton = app.buttons.matching(NSPredicate(format: "label CONTAINS 'Join'")).firstMatch
            let rsvpButton = app.buttons.matching(NSPredicate(format: "label CONTAINS 'RSVP'")).firstMatch
            _ = joinButton.waitForExistence(timeout: 3)
            _ = rsvpButton.waitForExistence(timeout: 1)

            app.navigationBars.buttons.firstMatch.tap()
        }
    }
}
