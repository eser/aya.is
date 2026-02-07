/**
 * Render inline markdown (bold, italic, links, code, strikethrough)
 * as HTML. Escapes raw HTML first to prevent XSS, then applies
 * markdown formatting.
 */

function escapeHtml(text: string): string {
  return text
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;");
}

export function renderInlineMarkdown(text: string): string {
  let html = escapeHtml(text);

  // Bold: **text** or __text__
  html = html.replace(/\*\*([^*]+)\*\*/g, "<strong>$1</strong>");
  html = html.replace(/__([^_]+)__/g, "<strong>$1</strong>");

  // Italic: *text* or _text_
  html = html.replace(/\*([^*]+)\*/g, "<em>$1</em>");
  html = html.replace(/(?<!\w)_([^_]+)_(?!\w)/g, "<em>$1</em>");

  // Strikethrough: ~~text~~
  html = html.replace(/~~([^~]+)~~/g, "<del>$1</del>");

  // Inline code: `text`
  html = html.replace(/`([^`]+)`/g, "<code>$1</code>");

  // Links: [text](url) — only allow http/https URLs
  html = html.replace(
    /\[([^\]]+)\]\((https?:\/\/[^)]+)\)/g,
    '<a href="$2" rel="noopener noreferrer">$1</a>',
  );

  return html;
}

type InlineMarkdownProps = {
  content: string;
  className?: string;
  as?: "p" | "span" | "div";
};

export function InlineMarkdown(props: InlineMarkdownProps) {
  const Tag = props.as ?? "p";
  const html = renderInlineMarkdown(props.content);

  return (
    <Tag
      className={props.className}
      dangerouslySetInnerHTML={{ __html: html }}
    />
  );
}
