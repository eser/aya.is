import SwiftUI

public struct FilterChip: Identifiable, Sendable {
    public let id: String
    public let label: String

    public init(id: String, label: String) {
        self.id = id
        self.label = label
    }
}

public struct FilterChipBar: View {
    let chips: [FilterChip]
    @Binding var selectedID: String?

    public init(chips: [FilterChip], selectedID: Binding<String?>) {
        self.chips = chips
        self._selectedID = selectedID
    }

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
