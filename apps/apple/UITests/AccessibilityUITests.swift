import XCTest

@MainActor
final class AccessibilityUITests: XCTestCase {

    private var app: XCUIApplication!

    override func setUpWithError() throws {
        continueAfterFailure = false
        app = XCUIApplication()
        app.launchArguments = ["-AppleLanguages", "(en)"]
        app.launch()
    }

    // MARK: - VoiceOver Accessibility Labels

    func testToolbarButtonsHaveAccessibilityLabels() throws {
        let toggleButton = app.buttons["Toggle appearance"]
        XCTAssertTrue(toggleButton.waitForExistence(timeout: 10), "Appearance toggle should have accessibility label")

        let languageMenu = app.buttons.matching(NSPredicate(format: "label CONTAINS 'Change language'")).firstMatch
        XCTAssertTrue(languageMenu.exists, "Language menu should have accessibility label")
    }

    func testFilterChipsHaveAccessibilityLabels() throws {
        let storiesChip = app.buttons["Stories"]
        XCTAssertTrue(storiesChip.waitForExistence(timeout: 10))
        XCTAssertTrue(app.buttons["Activities"].exists)
        XCTAssertTrue(app.buttons["People"].exists)
        XCTAssertTrue(app.buttons["Products"].exists)
    }

    func testSearchBarIsAccessible() throws {
        let searchField = app.textFields.firstMatch
        XCTAssertTrue(searchField.waitForExistence(timeout: 10), "Search field should be accessible")

        searchField.tap()
        searchField.typeText("test")

        let clearButton = app.buttons["Clear search"]
        XCTAssertTrue(clearButton.waitForExistence(timeout: 5), "Clear button should have accessibility label")
    }

    // MARK: - Dynamic Type

    func testAppLaunchesWithLargerText() throws {
        let largeTextApp = XCUIApplication()
        largeTextApp.launchArguments = ["-AppleLanguages", "(en)", "-UIPreferredContentSizeCategoryName", "UICTContentSizeCategoryAccessibilityExtraLarge"]
        largeTextApp.launch()

        let title = largeTextApp.navigationBars["AYA"]
        XCTAssertTrue(title.waitForExistence(timeout: 10), "App should load with larger text")

        let scrollView = largeTextApp.scrollViews.firstMatch
        XCTAssertTrue(scrollView.waitForExistence(timeout: 10))
    }

    // MARK: - Color Contrast (Dark Mode)

    func testAppWorksInDarkMode() throws {
        let toggleButton = app.buttons["Toggle appearance"]
        XCTAssertTrue(toggleButton.waitForExistence(timeout: 10))

        // Switch to dark mode
        toggleButton.tap()

        // Verify app still functions
        let scrollView = app.scrollViews.firstMatch
        XCTAssertTrue(scrollView.exists, "Feed should still be visible in dark mode")

        let storiesChip = app.buttons["Stories"]
        XCTAssertTrue(storiesChip.exists, "Filter chips should be visible in dark mode")

        // Switch back
        toggleButton.tap()
    }

    // MARK: - RTL Layout

    func testAppLaunchesInArabic() throws {
        let arabicApp = XCUIApplication()
        arabicApp.launchArguments = ["-AppleLanguages", "(ar)"]
        arabicApp.launch()

        let scrollView = arabicApp.scrollViews.firstMatch
        XCTAssertTrue(scrollView.waitForExistence(timeout: 10), "App should load in Arabic RTL layout")
    }
}
