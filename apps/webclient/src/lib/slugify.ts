/**
 * Convert a string to a URL-friendly slug.
 * - Converts to lowercase
 * - Removes non-alphanumeric characters (except spaces and hyphens)
 * - Replaces spaces with hyphens
 * - Collapses multiple hyphens into one
 * - Trims leading/trailing hyphens
 */
export function slugify(text: string): string {
  return text
    .toLowerCase()
    .replace(/[^a-z0-9\s-]/g, "")
    .replace(/\s+/g, "-")
    .replace(/-+/g, "-")
    .trim();
}

/**
 * Get current date as YYYYMMDD- prefix for story slugs.
 */
export function getDatePrefix(): string {
  const now = new Date();
  const year = now.getFullYear();
  const month = String(now.getMonth() + 1).padStart(2, "0");
  const day = String(now.getDate()).padStart(2, "0");
  return `${year}${month}${day}-`;
}
