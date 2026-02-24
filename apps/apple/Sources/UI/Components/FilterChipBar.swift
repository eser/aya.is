import SwiftUI

/// A data model representing a single selectable filter chip.
public struct FilterChip: Identifiable, Sendable {
    /// Unique identifier for the chip.
    public let id: String
    /// Display label shown inside the chip.
    public let label: String

    /// Creates a filter chip with the given identifier and label.
    public init(id: String, label: String) {
        self.id = id
        self.label = label
    }
}

/// A horizontally scrolling bar of selectable filter chips.
public struct FilterChipBar: View {
    let chips: [FilterChip]
    @Binding var selectedID: String?

    /// Creates a filter chip bar.
    public init(chips: [FilterChip], selectedID: Binding<String?>) {
        self.chips = chips
        self._selectedID = selectedID
    }

    /// The chip bar layout: a horizontal scroll view of toggle-style chip buttons.
    public var body: some View {
        ScrollView(.horizontal, showsIndicators: false) {
            HStack(spacing: AYASpacing.sm) {
                ForEach(chips) { chip in
                    let isSelected = selectedID == chip.id
                    Button {
                        withAnimation(AYAAnimation.filterSwitch) {
                            selectedID = isSelected ? nil : chip.id
                        }
                    } label: {
                        Text(chip.label)
                            .font(AYATypography.subheadline)
                            .fontWeight(isSelected ? .semibold : .regular)
                            .foregroundStyle(isSelected ? .white : AYAColors.textPrimary)
                            .padding(.horizontal, AYASpacing.md)
                            .padding(.vertical, AYASpacing.sm)
                            .background(isSelected ? AYAColors.accent : AYAColors.contentBackground)
                            .clipShape(Capsule())
                    }
                    .buttonStyle(.plain)
                }
            }
            .padding(.horizontal, AYASpacing.md)
        }
    }
}
