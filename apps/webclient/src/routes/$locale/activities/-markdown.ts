/**
 * Activities domain - markdown generation utilities
 */
import { formatFrontmatter, formatStoryListItem } from "@/lib/markdown";
import { registerMarkdownHandler } from "@/server/markdown-middleware";
import { backend } from "@/modules/backend/backend";
import type { StoryEx } from "@/modules/backend/types";

/**
 * Generate markdown for the activities listing page
 */
export function generateActivitiesListingMarkdown(
  activities: StoryEx[],
  locale: string,
): string {
  const frontmatter = formatFrontmatter({
    title: "Activities",
    generated_at: new Date().toISOString(),
  });

  if (activities.length === 0) {
    return `${frontmatter}\n\nNo activities available yet.`;
  }

  const activityLinks = activities.map((activity) => formatStoryListItem(activity, locale, "activities"));

  return `${frontmatter}\n\n${activityLinks.join("\n\n")}`;
}

/**
 * Register markdown handler for activities listing
 * Pattern: /$locale/activities
 */
export function registerActivitiesListingHandler(): void {
  registerMarkdownHandler("$locale/activities", async (_params, locale, _searchParams) => {
    const activities = await backend.getActivities(locale);

    if (activities === null) {
      return null;
    }

    return generateActivitiesListingMarkdown(activities, locale);
  });
}
