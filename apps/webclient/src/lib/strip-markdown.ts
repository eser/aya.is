/**
 * Strip all markdown formatting from text, returning plain text.
 * Useful for rendering titles that may contain accidental markdown.
 */
export function stripMarkdown(text: string): string {
  return text
    .replace(/!\[([^\]]*)\]\([^)]+\)/g, "$1") // Remove images, keep alt
    .replace(/\[([^\]]+)\]\([^)]+\)/g, "$1") // Remove links, keep text
    .replace(/```[\s\S]*?```/g, "") // Remove code blocks
    .replace(/`([^`]+)`/g, "$1") // Remove inline code
    .replace(/\*\*([^*]+)\*\*/g, "$1") // Remove bold
    .replace(/__([^_]+)__/g, "$1") // Remove bold (underscore)
    .replace(/\*([^*]+)\*/g, "$1") // Remove italic
    .replace(/_([^_]+)_/g, "$1") // Remove italic (underscore)
    .replace(/~~([^~]+)~~/g, "$1") // Remove strikethrough
    .replace(/#{1,6}\s+/g, "") // Remove heading markers
    .replace(/>\s?/g, "") // Remove blockquotes
    .replace(/[-*+]\s/g, "") // Remove unordered list markers
    .replace(/\d+\.\s/g, "") // Remove ordered list markers
    .trim();
}
