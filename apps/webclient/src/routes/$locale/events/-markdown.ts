/**
 * Events domain - markdown generation utilities
 * Events is currently a placeholder with no data
 */
import { formatFrontmatter } from "@/lib/markdown";
import { registerMarkdownHandler } from "@/server/markdown-middleware";

/**
 * Generate markdown for the events listing page
 * Currently a placeholder as there's no events data
 */
export function generateEventsListingMarkdown(_locale: string): string {
  const frontmatter = formatFrontmatter({
    title: "Events",
    generated_at: new Date().toISOString(),
  });

  return `${frontmatter}\n\nNo events available yet.`;
}

/**
 * Register markdown handler for events listing
 * Pattern: /$locale/events
 */
export function registerEventsListingHandler(): void {
  registerMarkdownHandler("$locale/events", (_params, locale, _searchParams) => {
    return Promise.resolve(generateEventsListingMarkdown(locale));
  });
}
