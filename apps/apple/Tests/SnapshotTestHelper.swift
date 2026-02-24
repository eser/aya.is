import SwiftUI
import SnapshotTesting

#if os(macOS)
import AppKit

@MainActor
func assertSwiftUISnapshot<V: View>(
    of view: V,
    width: CGFloat = 350,
    height: CGFloat? = nil,
    file: StaticString = #file,
    testName: String = #function,
    line: UInt = #line
) {
    let hostingView = NSHostingView(rootView: view)
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
        file: file,
        testName: testName,
        line: line
    )
}

#elseif os(iOS)
import UIKit

@MainActor
func assertSwiftUISnapshot<V: View>(
    of view: V,
    width: CGFloat = 350,
    height: CGFloat? = nil,
    file: StaticString = #file,
    testName: String = #function,
    line: UInt = #line
) {
    if let height {
        let layout = SwiftUISnapshotLayout.fixed(width: width, height: height)
        assertSnapshot(of: view, as: .image(layout: layout), file: file, testName: testName, line: line)
    } else {
        let wrappedView = view.frame(width: width)
        assertSnapshot(of: wrappedView, as: .image(layout: .sizeThatFits), file: file, testName: testName, line: line)
    }
}
#endif
