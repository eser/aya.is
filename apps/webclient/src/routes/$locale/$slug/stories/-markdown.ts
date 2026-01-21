/**
 * Profile stories domain - markdown generation utilities
 */
import { formatFrontmatter } from "@/lib/markdown";
import { registerMarkdownHandler } from "@/server/markdown-middleware";
import { formatDate, parseDateFromSlug } from "@/lib/date";
import { calculateReadingTime } from "@/lib/reading-time";
import { backend } from "@/modules/backend/backend";
import type { StoryEx } from "@/modules/backend/types";

export function generateProfileStoryMarkdown(
  story: StoryEx,
  locale: string,
): string {
  // Try to get date from slug first, fall back to created_at
  const publishDate = story.slug !== null
    ? parseDateFromSlug(story.slug)
    : null;
  const dateToFormat = publishDate ?? new Date(story.created_at);

  const frontmatter = formatFrontmatter({
    title: story.title ?? "Untitled",
    author: story.author_profile?.title ?? "Unknown",
    publish_date: formatDate(dateToFormat, locale),
    reading_time: `${calculateReadingTime(story.content)} min`,
    kind: story.kind,
    status: story.status,
  });

  return `${frontmatter}\n\n${story.content}`;
}

/**
 * Register markdown handler for profile stories
 * Pattern: /$locale/$slug/stories/$storyslug
 */
export function registerProfileStoryMarkdownHandler(): void {
  registerMarkdownHandler(
    "$locale/$slug/stories/$storyslug",
    async (params, locale, _searchParams) => {
      const { slug, storyslug } = params;

      if (slug === undefined || storyslug === undefined) {
        return null;
      }

      const story = await backend.getProfileStory(locale, slug, storyslug);

      if (story === null) {
        return null;
      }

      return generateProfileStoryMarkdown(story, locale);
    },
  );
}
