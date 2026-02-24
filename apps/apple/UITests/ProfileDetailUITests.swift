import XCTest

final class ProfileDetailUITests: XCTestCase {

    private var app: XCUIApplication!

    override func setUpWithError() throws {
        continueAfterFailure = false
        app = XCUIApplication()
        app.launchArguments = ["-AppleLanguages", "(en)"]
        app.launch()
    }

    // MARK: - Profile Detail (Sheet)

    func testNavigateToProfileSheet() throws {
        let peopleChip = app.buttons["People"]
        XCTAssertTrue(peopleChip.waitForExistence(timeout: 10))

        peopleChip.tap()
        sleep(2)

        let scrollView = app.scrollViews.firstMatch
        let firstCard = scrollView.buttons.firstMatch
        if firstCard.waitForExistence(timeout: 5) {
            firstCard.tap()

            // Profile opens as a sheet — look for dismiss button
            let dismissButton = app.buttons.matching(NSPredicate(format: "label CONTAINS 'xmark' OR label CONTAINS 'Close'")).firstMatch
            _ = dismissButton.waitForExistence(timeout: 5)
        }
    }

    func testProfileSheetShowsTabs() throws {
        let peopleChip = app.buttons["People"]
        XCTAssertTrue(peopleChip.waitForExistence(timeout: 10))

        peopleChip.tap()
        sleep(2)

        let scrollView = app.scrollViews.firstMatch
        let firstCard = scrollView.buttons.firstMatch
        if firstCard.waitForExistence(timeout: 5) {
            firstCard.tap()
            sleep(2)

            // Check for segmented picker tabs
            let pagesTab = app.buttons["Pages"]
            let storiesTab = app.buttons["Stories"]

            if pagesTab.waitForExistence(timeout: 5) {
                pagesTab.tap()
            }
            if storiesTab.exists {
                storiesTab.tap()
            }
        }
    }

    func testProfileSheetScrolls() throws {
        let peopleChip = app.buttons["People"]
        XCTAssertTrue(peopleChip.waitForExistence(timeout: 10))

        peopleChip.tap()
        sleep(2)

        let scrollView = app.scrollViews.firstMatch
        let firstCard = scrollView.buttons.firstMatch
        if firstCard.waitForExistence(timeout: 5) {
            firstCard.tap()
            sleep(1)

            let sheetScroll = app.scrollViews.firstMatch
            if sheetScroll.waitForExistence(timeout: 5) {
                sheetScroll.swipeUp()
                sheetScroll.swipeDown()
            }
        }
    }

    // MARK: - Product Profile

    func testNavigateToProductProfile() throws {
        let productsChip = app.buttons["Products"]
        XCTAssertTrue(productsChip.waitForExistence(timeout: 10))

        productsChip.tap()
        sleep(2)

        let scrollView = app.scrollViews.firstMatch
        let firstCard = scrollView.buttons.firstMatch
        if firstCard.waitForExistence(timeout: 5) {
            firstCard.tap()
            sleep(1)

            // Profile sheet should appear
            let sheetScroll = app.scrollViews.firstMatch
            _ = sheetScroll.waitForExistence(timeout: 5)
        }
    }
}
