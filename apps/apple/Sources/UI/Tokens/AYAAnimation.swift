import SwiftUI

/// Design-token namespace for standard AYA animations.
public enum AYAAnimation {
    /// Snappy spring used when a card is tapped.
    public static let cardTap: Animation = .spring(response: 0.3, dampingFraction: 0.7)
    /// Smooth spring for presenting sheets and modals.
    public static let sheetPresent: Animation = .spring(response: 0.4, dampingFraction: 0.85)
    /// Quick ease-in-out for toggling filter selections.
    public static let filterSwitch: Animation = .easeInOut(duration: 0.25)
    /// Continuously repeating shimmer used by skeleton loaders.
    public static let shimmer: Animation = .easeInOut(duration: 1.2).repeatForever(autoreverses: false)
    /// Subtle fade-in for appearing content.
    public static let fadeIn: Animation = .easeIn(duration: 0.2)
    /// General-purpose spring for content transitions.
    public static let contentTransition: Animation = .spring(response: 0.35, dampingFraction: 0.9)
}
