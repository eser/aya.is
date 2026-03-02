import SwiftUI

/// A view that parses a Markdown string and renders headings, blockquotes, code blocks, lists, images, and paragraphs.
public struct RichContentView: View {
    let content: String

    /// Creates a rich content view from raw Markdown text.
    public init(content: String) {
        self.content = content
    }

    /// The rendered Markdown content as a vertical stack of block-level elements.
    public var body: some View {
        VStack(alignment: .leading, spacing: AYASpacing.md) {
            ForEach(Array(parseBlocks(content).enumerated()), id: \.offset) { _, block in
                blockView(block)
            }
        }
        .frame(maxWidth: .infinity, alignment: .leading)
        .textSelection(.enabled)
    }

    // MARK: - Block parsing

    private enum Block {
        case heading(Int, String)
        case blockquote(String)
        case codeBlock(String, String?)
        case orderedItem(Int, String)
        case unorderedItem(String)
        case image(String, String?)
        case embed(String)
        case horizontalRule
        case paragraph(String)
    }

    private func parseBlocks(_ text: String) -> [Block] {
        var blocks: [Block] = []
        let lines = text.components(separatedBy: "\n")
        var i = 0

        while i < lines.count {
            let line = lines[i]
            let trimmed = line.trimmingCharacters(in: .whitespaces)

            if trimmed.isEmpty {
                i += 1
                continue
            }

            if trimmed.hasPrefix("```") {
                let lang = String(trimmed.dropFirst(3)).trimmingCharacters(in: .whitespaces)
                var codeLines: [String] = []
                i += 1
                while i < lines.count && !lines[i].trimmingCharacters(in: .whitespaces).hasPrefix("```") {
                    codeLines.append(lines[i])
                    i += 1
                }
                if i < lines.count { i += 1 }
                blocks.append(.codeBlock(codeLines.joined(separator: "\n"), lang.isEmpty ? nil : lang))
                continue
            }

            if trimmed.isHorizontalRule {
                blocks.append(.horizontalRule)
                i += 1
                continue
            }

            if let headingMatch = trimmed.headingLevel() {
                blocks.append(.heading(headingMatch.level, headingMatch.text))
                i += 1
                continue
            }

            if trimmed.hasPrefix(">") {
                var quoteLines: [String] = []
                while i < lines.count {
                    let l = lines[i].trimmingCharacters(in: .whitespaces)
                    if l.hasPrefix(">") {
                        let stripped = String(l.dropFirst()).trimmingCharacters(in: .whitespaces)
                        quoteLines.append(stripped)
                        i += 1
                    } else if l.isEmpty {
                        i += 1
                        break
                    } else {
                        break
                    }
                }
                blocks.append(.blockquote(quoteLines.joined(separator: " ")))
                continue
            }

            // Embed: %[url] syntax (Hashnode-compatible)
            if let embedRange = trimmed.range(of: #"^%\[([^\]]+)\]$"#, options: .regularExpression) {
                let inner = String(trimmed[embedRange])
                let url = String(inner.dropFirst(2).dropLast())
                blocks.append(.embed(url))
                i += 1
                continue
            }

            if trimmed.hasPrefix("![") {
                if let range = trimmed.range(of: #"!\[(.*?)\]\((.*?)\)"#, options: .regularExpression) {
                    let match = String(trimmed[range])
                    let parts = match.components(separatedBy: "](")
                    let alt = String(parts[0].dropFirst(2))
                    let url = String(parts[1].dropLast())
                    blocks.append(.image(url, alt.isEmpty ? nil : alt))
                    i += 1
                    continue
                }
            }

            if let orderedMatch = trimmed.orderedListItem() {
                blocks.append(.orderedItem(orderedMatch.num, orderedMatch.text))
                i += 1
                continue
            }

            if trimmed.hasPrefix("- ") || trimmed.hasPrefix("* ") || trimmed.hasPrefix("+ ") {
                let text = String(trimmed.dropFirst(2))
                blocks.append(.unorderedItem(text))
                i += 1
                continue
            }

            var paraLines: [String] = []
            while i < lines.count {
                let l = lines[i].trimmingCharacters(in: .whitespaces)
                if l.isEmpty || l.hasPrefix("#") || l.hasPrefix(">") || l.hasPrefix("```") || l.hasPrefix("- ") || l.hasPrefix("* ") || l.hasPrefix("![") || l.hasPrefix("%[") || l.isHorizontalRule {
                    break
                }
                if l.orderedListItem() != nil { break }
                if lines[i].hasSuffix("  ") {
                    paraLines.append(l + "\n")
                } else {
                    paraLines.append(l)
                }
                i += 1
            }
            if !paraLines.isEmpty {
                blocks.append(.paragraph(paraLines.joined(separator: " ")))
            }
        }

        return blocks
    }

    // MARK: - Block rendering

    @ViewBuilder
    private func blockView(_ block: Block) -> some View {
        switch block {
        case .heading(let level, let text):
            inlineMarkdownText(text)
                .font(headingFont(level))
                .fontWeight(.bold)
                .foregroundStyle(AYAColors.textPrimary)
                .padding(.top, level == 1 ? AYASpacing.sm : AYASpacing.xs)

        case .blockquote(let text):
            HStack(spacing: AYASpacing.sm) {
                Rectangle()
                    .fill(AYAColors.accent.opacity(0.6))
                    .frame(width: 3)

                inlineMarkdownText(text)
                    .font(AYATypography.body)
                    .italic()
                    .foregroundStyle(AYAColors.textSecondary)
            }
            .padding(.vertical, AYASpacing.xs)

        case .codeBlock(let code, _):
            ScrollView(.horizontal, showsIndicators: false) {
                Text(code)
                    .font(AYATypography.code)
                    .foregroundStyle(AYAColors.textPrimary)
                    .padding(AYASpacing.md)
            }
            .background(AYAColors.contentBackground)
            .clipShape(RoundedRectangle(cornerRadius: AYACornerRadius.md))

        case .orderedItem(let num, let text):
            HStack(alignment: .top, spacing: AYASpacing.sm) {
                Text("\(num).")
                    .font(AYATypography.body)
                    .foregroundStyle(AYAColors.textTertiary)
                    .frame(width: 24, alignment: .trailing)
                inlineMarkdownText(text)
                    .font(AYATypography.body)
                    .foregroundStyle(AYAColors.textPrimary)
            }

        case .unorderedItem(let text):
            HStack(alignment: .top, spacing: AYASpacing.sm) {
                Text("\u{2022}")
                    .font(AYATypography.body)
                    .foregroundStyle(AYAColors.accent)
                    .frame(width: 16)
                inlineMarkdownText(text)
                    .font(AYATypography.body)
                    .foregroundStyle(AYAColors.textPrimary)
            }

        case .image(let url, let alt):
            VStack(spacing: AYASpacing.xs) {
                if let imageURL = URL(string: url) {
                    AsyncImage(url: imageURL) { phase in
                        switch phase {
                        case .success(let image):
                            image
                                .resizable()
                                .aspectRatio(contentMode: .fit)
                                .frame(maxWidth: .infinity)
                                .clipShape(RoundedRectangle(cornerRadius: AYACornerRadius.md))
                        case .failure:
                            EmptyView()
                        default:
                            RoundedRectangle(cornerRadius: AYACornerRadius.md)
                                .fill(AYAColors.contentBackground)
                                .frame(height: 200)
                                .overlay { ProgressView() }
                        }
                    }
                }
                if let alt {
                    Text(alt)
                        .font(AYATypography.caption)
                        .foregroundStyle(AYAColors.textTertiary)
                        .italic()
                }
            }

        case .embed(let url):
            EmbedView(url: url)

        case .horizontalRule:
            Rectangle()
                .fill(AYAColors.textTertiary.opacity(0.3))
                .frame(height: 1)
                .padding(.vertical, AYASpacing.sm)

        case .paragraph(let text):
            inlineMarkdownText(text)
                .font(AYATypography.body)
                .foregroundStyle(AYAColors.textPrimary)
                .lineSpacing(4)
        }
    }

    private func inlineMarkdownText(_ text: String) -> Text {
        if let attributed = try? AttributedString(markdown: text, options: .init(
            interpretedSyntax: .inlineOnlyPreservingWhitespace,
            failurePolicy: .returnPartiallyParsedIfPossible
        )) {
            return Text(attributed)
        }
        return Text(text)
    }

    private func headingFont(_ level: Int) -> Font {
        switch level {
        case 1: AYATypography.title
        case 2: AYATypography.title2
        case 3: AYATypography.title3
        default: AYATypography.headline
        }
    }
}

// MARK: - String helpers

private extension String {
    var isHorizontalRule: Bool {
        let stripped = filter { $0 != " " }
        return stripped.count >= 3 && stripped.allSatisfy({ $0 == "-" || $0 == "*" || $0 == "_" })
            && Set(stripped).count == 1
    }

    func headingLevel() -> (level: Int, text: String)? {
        var count = 0
        for ch in self {
            if ch == "#" { count += 1 }
            else { break }
        }
        guard count >= 1 && count <= 6 else { return nil }
        let rest = String(dropFirst(count)).trimmingCharacters(in: .whitespaces)
        guard !rest.isEmpty else { return nil }
        return (count, rest)
    }

    func orderedListItem() -> (num: Int, text: String)? {
        guard let dotIndex = firstIndex(of: ".") else { return nil }
        let prefix = self[startIndex..<dotIndex]
        guard let num = Int(prefix) else { return nil }
        let rest = String(self[index(after: dotIndex)...]).trimmingCharacters(in: .whitespaces)
        guard !rest.isEmpty else { return nil }
        return (num, rest)
    }
}
