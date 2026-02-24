import SwiftUI
import SnapshotTesting
@testable import AYAKit

/// Fixed locale used across all snapshot tests to ensure consistent
/// date/number formatting regardless of host system locale.
let snapshotLocale = Locale(identifier: "en_US")

/// Fixed reference date used across all snapshot tests so that relative
/// date strings (e.g. "1 mo. ago") remain deterministic over time.
let snapshotReferenceDate: Date = {
    var components = DateComponents()
    components.year = 2026
    components.month = 2
    components.day = 1
    components.hour = 0
    components.minute = 0
    components.second = 0
    components.timeZone = TimeZone(identifier: "UTC")
    return Calendar(identifier: .gregorian).date(from: components)!
}()

#if os(iOS)
import UIKit

private let platformSuffix = "iOS"

@MainActor
func assertSwiftUISnapshot<V: View>(
    of view: V,
    named name: String? = nil,
    width: CGFloat = 350,
    height: CGFloat? = nil,
    file: StaticString = #filePath,
    testName: String = #function,
    line: UInt = #line
) {
    let snapshotName = name.map { "\($0)-\(platformSuffix)" } ?? platformSuffix
    let localized = view
        .environment(\.locale, snapshotLocale)
        .environment(\.referenceDate, snapshotReferenceDate)
    if let height {
        let layout = SwiftUISnapshotLayout.fixed(width: width, height: height)
        assertSnapshot(of: localized, as: .image(layout: layout), named: snapshotName, file: file, testName: testName, line: line)
    } else {
        let wrappedView = localized.frame(width: width)
        assertSnapshot(of: wrappedView, as: .image(layout: .sizeThatFits), named: snapshotName, file: file, testName: testName, line: line)
    }
}

#elseif os(macOS)
import AppKit

private let platformSuffix = "macOS"

@MainActor
func assertSwiftUISnapshot<V: View>(
    of view: V,
    named name: String? = nil,
    width: CGFloat = 350,
    height: CGFloat? = nil,
    file: StaticString = #filePath,
    testName: String = #function,
    line: UInt = #line
) {
    let snapshotName = name.map { "\($0)-\(platformSuffix)" } ?? platformSuffix
    let localized = view
        .environment(\.locale, snapshotLocale)
        .environment(\.referenceDate, snapshotReferenceDate)
    let hostingView = NSHostingView(rootView: localized)
    let fittingSize: CGSize
    if let height {
        fittingSize = CGSize(width: width, height: height)
    } else {
        hostingView.frame = NSRect(x: 0, y: 0, width: width, height: 10000)
        let fitted = hostingView.fittingSize
        fittingSize = CGSize(width: width, height: max(fitted.height, 1))
    }
    hostingView.frame = NSRect(origin: .zero, size: fittingSize)
    assertSnapshot(
        of: hostingView,
        as: .image(size: fittingSize),
        named: snapshotName,
        file: file,
        testName: testName,
        line: line
    )
}
#endif
