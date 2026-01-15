/**
 * Calculate reading time for text content
 * Assumes average reading speed of 200 words per minute
 */
export function calculateReadingTime(content: string): number {
  // Remove HTML tags and markdown syntax
  const cleanText = content
    .replace(/<[^>]*>/g, "") // Remove HTML tags
    .replace(/\*\*([^*]+)\*\*/g, "$1") // Remove bold markdown
    .replace(/\*([^*]+)\*/g, "$1") // Remove italic markdown
    .replace(/\[([^\]]+)\]\([^)]+\)/g, "$1") // Remove markdown links
    .replace(/```[\s\S]*?```/g, "") // Remove code blocks
    .replace(/`([^`]+)`/g, "$1") // Remove inline code
    .replace(/#{1,6}\s+/g, "") // Remove markdown headers
    .replace(/\n/g, " ") // Replace newlines with spaces
    .trim();

  // Count words (split by whitespace and filter empty strings)
  const words = cleanText.split(/\s+/).filter((word) => word.length > 0);
  const wordCount = words.length;

  // Calculate reading time (200 words per minute, minimum 1 minute)
  const readingTimeMinutes = Math.max(1, Math.ceil(wordCount / 200));

  return readingTimeMinutes;
}
