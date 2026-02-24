import SwiftUI

public struct AYASearchBar: View {
    @Binding var text: String
    var placeholder: String = "Search..."
    var onCommit: (() -> Void)?

    public init(text: Binding<String>, placeholder: String = "Search...", onCommit: (() -> Void)? = nil) {
        self._text = text
        self.placeholder = placeholder
        self.onCommit = onCommit
    }

    public var body: some View {
        HStack(spacing: AYASpacing.sm) {
            Image(systemName: "magnifyingglass")
                .foregroundStyle(AYAColors.textTertiary)
                .font(.body)

            TextField(placeholder, text: $text)
                .textFieldStyle(.plain)
                .font(AYATypography.body)
                .onSubmit { onCommit?() }
                #if os(iOS)
                .textInputAutocapitalization(.never)
                .autocorrectionDisabled()
                #endif

            if !text.isEmpty {
                Button {
                    withAnimation(AYAAnimation.filterSwitch) {
                        text = ""
                    }
                } label: {
                    Image(systemName: "xmark.circle.fill")
                        .foregroundStyle(AYAColors.textTertiary)
                }
                .buttonStyle(.plain)
            }
        }
        .padding(.horizontal, AYASpacing.md)
        .padding(.vertical, AYASpacing.sm + 2)
        .background(AYAColors.contentBackground)
        .clipShape(RoundedRectangle(cornerRadius: AYACornerRadius.xl))
    }
}
