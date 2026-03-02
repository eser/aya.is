import SwiftUI

private struct ReferenceDateKey: EnvironmentKey {
    static let defaultValue: Date = Date()
}

extension EnvironmentValues {
    var referenceDate: Date {
        get { self[ReferenceDateKey.self] }
        set { self[ReferenceDateKey.self] = newValue }
    }
}
