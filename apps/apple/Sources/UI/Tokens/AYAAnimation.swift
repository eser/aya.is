import SwiftUI

public enum AYAAnimation {
    public static let cardTap: Animation = .spring(response: 0.3, dampingFraction: 0.7)
    public static let sheetPresent: Animation = .spring(response: 0.4, dampingFraction: 0.85)
    public static let filterSwitch: Animation = .easeInOut(duration: 0.25)
    public static let shimmer: Animation = .easeInOut(duration: 1.2).repeatForever(autoreverses: false)
    public static let fadeIn: Animation = .easeIn(duration: 0.2)
    public static let contentTransition: Animation = .spring(response: 0.35, dampingFraction: 0.9)
}
