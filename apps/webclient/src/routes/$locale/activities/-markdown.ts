/**
 * Activities domain - markdown generation utilities
 * Activities is currently a placeholder with no data
 */
import { formatFrontmatter } from "@/lib/markdown";
import { registerMarkdownHandler } from "@/server/markdown-middleware";

/**
 * Generate markdown for the activities listing page
 * Currently a placeholder as there's no activities data
 */
export function generateActivitiesListingMarkdown(_locale: string): string {
  const frontmatter = formatFrontmatter({
    title: "Activities",
    generated_at: new Date().toISOString(),
  });

  return `${frontmatter}\n\nNo activities available yet.`;
}

/**
 * Register markdown handler for activities listing
 * Pattern: /$locale/activities
 */
export function registerActivitiesListingHandler(): void {
  registerMarkdownHandler("$locale/activities", (_params, locale, _searchParams) => {
    return Promise.resolve(generateActivitiesListingMarkdown(locale));
  });
}
