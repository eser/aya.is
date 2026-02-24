import XCTest

final class FeedUITests: XCTestCase {

    private var app: XCUIApplication!

    override func setUpWithError() throws {
        continueAfterFailure = false
        app = XCUIApplication()
        app.launchArguments = ["--uitesting", "-AppleLanguages", "(en)"]
        app.launch()
    }

    // MARK: - Feed Loading

    func testFeedLoadsAndDisplaysContent() throws {
        let title = app.navigationBars["AYA"]
        XCTAssertTrue(title.waitForExistence(timeout: 10), "Navigation title should appear")

        let scrollView = app.scrollViews.firstMatch
        XCTAssertTrue(scrollView.waitForExistence(timeout: 10), "Feed scroll view should appear")
    }

    func testFeedScrollsVertically() throws {
        let scrollView = app.scrollViews.firstMatch
        XCTAssertTrue(scrollView.waitForExistence(timeout: 10))

        scrollView.swipeUp()
        scrollView.swipeUp()
        scrollView.swipeDown()
        scrollView.swipeDown()
    }

    // MARK: - Filter Chips

    func testFilterChipsExist() throws {
        let storiesChip = app.buttons["Stories"]
        XCTAssertTrue(storiesChip.waitForExistence(timeout: 10), "Stories filter chip should exist")

        let activitiesChip = app.buttons["Activities"]
        XCTAssertTrue(activitiesChip.exists, "Activities filter chip should exist")

        let peopleChip = app.buttons["People"]
        XCTAssertTrue(peopleChip.exists, "People filter chip should exist")

        let productsChip = app.buttons["Products"]
        XCTAssertTrue(productsChip.exists, "Products filter chip should exist")
    }

    func testFilterChipToggle() throws {
        let storiesChip = app.buttons["Stories"]
        XCTAssertTrue(storiesChip.waitForExistence(timeout: 10))

        storiesChip.tap()
        // Tap again to deselect
        storiesChip.tap()
    }

    func testFilterChipSwitching() throws {
        let storiesChip = app.buttons["Stories"]
        let activitiesChip = app.buttons["Activities"]
        XCTAssertTrue(storiesChip.waitForExistence(timeout: 10))

        storiesChip.tap()
        activitiesChip.tap()

        let peopleChip = app.buttons["People"]
        peopleChip.tap()

        let productsChip = app.buttons["Products"]
        productsChip.tap()
    }

    // MARK: - Search

    func testSearchBarExists() throws {
        let searchField = app.textFields.firstMatch
        XCTAssertTrue(searchField.waitForExistence(timeout: 10), "Search bar should exist")
    }

    func testSearchTypingAndClearing() throws {
        let searchField = app.textFields.firstMatch
        XCTAssertTrue(searchField.waitForExistence(timeout: 10))

        searchField.tap()
        searchField.typeText("test")

        let clearButton = app.buttons["Clear search"]
        if clearButton.waitForExistence(timeout: 3) {
            clearButton.tap()
        }
    }

    func testSearchShowsResults() throws {
        let searchField = app.textFields.firstMatch
        XCTAssertTrue(searchField.waitForExistence(timeout: 10))

        searchField.tap()
        searchField.typeText("AYA")

        // Wait for search results to appear (debounce + network)
        sleep(2)

        let scrollView = app.scrollViews.firstMatch
        XCTAssertTrue(scrollView.exists, "Search results should be scrollable")
    }

    // MARK: - Search Navigation

    func testSearchResultNavigatesToStoryDetail() throws {
        let searchField = app.textFields.firstMatch
        XCTAssertTrue(searchField.waitForExistence(timeout: 10))

        searchField.tap()
        searchField.typeText("scaling")

        // Wait for search results (debounce 400ms + response)
        let firstResult = app.buttons.matching(NSPredicate(format: "label CONTAINS 'Scaling Community Platforms'")).firstMatch
        XCTAssertTrue(firstResult.waitForExistence(timeout: 10), "Search results should appear")
        firstResult.tap()

        // Verify detail screen loads without error
        let backButton = app.navigationBars.buttons.firstMatch
        XCTAssertTrue(backButton.waitForExistence(timeout: 10), "Detail screen should load")

        // Verify no error view is shown
        let errorText = app.staticTexts.matching(NSPredicate(format: "label CONTAINS 'Failed to decode'"))
        XCTAssertEqual(errorText.count, 0, "Detail screen should not show a decode error")

        // Verify content loaded (scroll view exists with content)
        let detailScroll = app.scrollViews.firstMatch
        XCTAssertTrue(detailScroll.exists, "Detail content should be scrollable")

        backButton.tap()
    }

    func testSearchResultNavigatesToProfileDetail() throws {
        let searchField = app.textFields.firstMatch
        XCTAssertTrue(searchField.waitForExistence(timeout: 10))

        searchField.tap()
        searchField.typeText("eser")

        // Tap a profile result
        let firstResult = app.buttons.matching(NSPredicate(format: "label CONTAINS 'Eser Ozvataf'")).firstMatch
        if firstResult.waitForExistence(timeout: 10) {
            firstResult.tap()

            // Wait for detail to load
            sleep(2)

            // Verify no decode error
            let errorText = app.staticTexts.matching(NSPredicate(format: "label CONTAINS 'Failed to decode'"))
            XCTAssertEqual(errorText.count, 0, "Profile detail should not show a decode error")
        }
    }

    // MARK: - Navigation

    func testTapStoryCardNavigatesToDetail() throws {
        // Wait for feed content to load
        let firstCard = app.buttons.matching(NSPredicate(format: "label CONTAINS 'Scaling Community Platforms'")).firstMatch
        if firstCard.waitForExistence(timeout: 10) {
            firstCard.tap()

            // Should navigate to detail - check for back button
            let backButton = app.navigationBars.buttons.firstMatch
            XCTAssertTrue(backButton.waitForExistence(timeout: 5), "Back button should appear in detail view")

            backButton.tap()
        }
    }

    // MARK: - Toolbar

    func testAppearanceToggleExists() throws {
        let toggleButton = app.buttons["Toggle appearance"]
        XCTAssertTrue(toggleButton.waitForExistence(timeout: 10), "Appearance toggle should exist")
    }

    func testAppearanceToggleTap() throws {
        let toggleButton = app.buttons["Toggle appearance"]
        XCTAssertTrue(toggleButton.waitForExistence(timeout: 10))
        toggleButton.tap()
        // Toggle back
        toggleButton.tap()
    }

    func testLanguageMenuExists() throws {
        let languageMenu = app.buttons.matching(NSPredicate(format: "label CONTAINS 'Change language'")).firstMatch
        XCTAssertTrue(languageMenu.waitForExistence(timeout: 10), "Language menu should exist")
    }

    // MARK: - Pull to Refresh

    func testPullToRefresh() throws {
        let scrollView = app.scrollViews.firstMatch
        XCTAssertTrue(scrollView.waitForExistence(timeout: 10))

        // Wait for initial load
        sleep(2)

        // Pull to refresh
        let start = scrollView.coordinate(withNormalizedOffset: CGVector(dx: 0.5, dy: 0.3))
        let end = scrollView.coordinate(withNormalizedOffset: CGVector(dx: 0.5, dy: 0.8))
        start.press(forDuration: 0.1, thenDragTo: end)
    }

    // MARK: - Accessibility

    func testFeedElementsAreAccessible() throws {
        let scrollView = app.scrollViews.firstMatch
        XCTAssertTrue(scrollView.waitForExistence(timeout: 10))

        // Verify filter chips have accessibility labels
        XCTAssertTrue(app.buttons["Stories"].exists)
        XCTAssertTrue(app.buttons["Activities"].exists)
        XCTAssertTrue(app.buttons["People"].exists)
        XCTAssertTrue(app.buttons["Products"].exists)

        // Verify toolbar buttons are accessible
        XCTAssertTrue(app.buttons["Toggle appearance"].exists)
    }
}
