/**
 * Convert a string to a URL-friendly slug.
 * - Normalizes accented characters (é -> e, ü -> u, etc.)
 * - Converts to lowercase
 * - Removes non-alphanumeric characters (except spaces and hyphens)
 * - Replaces spaces with hyphens
 * - Collapses multiple hyphens into one
 * - Trims leading/trailing hyphens
 */
export function slugify(text: string): string {
  return text
    .normalize("NFD") // Decompose accented characters (é -> e + ́)
    .replace(/[\u0300-\u036f]/g, "") // Remove combining diacritical marks
    .toLowerCase()
    .replace(/[^a-z0-9\s-]/g, "") // Remove non-alphanumeric (except spaces and hyphens)
    .replace(/\s+/g, "-") // Replace spaces with hyphens
    .replace(/-+/g, "-") // Collapse multiple hyphens
    .replace(/^-|-$/g, ""); // Trim leading/trailing hyphens
}

/**
 * Sanitize user input for a slug field.
 * - Normalizes accented characters (é -> e, ü -> u, etc.)
 * - Converts to lowercase
 * - Replaces invalid characters with hyphens
 * - Collapses multiple hyphens into one
 *
 * Unlike slugify(), this is for sanitizing direct user input in slug fields,
 * not for generating slugs from titles.
 */
export function sanitizeSlug(text: string): string {
  return text
    .normalize("NFD") // Decompose accented characters (é -> e + ́)
    .replace(/[\u0300-\u036f]/g, "") // Remove combining diacritical marks
    .toLowerCase()
    .replace(/[^a-z0-9-]/g, "-") // Replace invalid chars with hyphens
    .replace(/-+/g, "-"); // Collapse multiple hyphens
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
