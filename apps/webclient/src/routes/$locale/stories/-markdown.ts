/**
 * Global stories domain - markdown generation utilities
 */
import { formatFrontmatter } from "@/lib/markdown";
import { registerMarkdownHandler } from "@/server/markdown-middleware";
import { formatDateShort, parseDateFromSlug } from "@/lib/date";
import { calculateReadingTime } from "@/lib/reading-time";
import { backend } from "@/modules/backend/backend";
import { supportedLocales, type SupportedLocaleCode, isValidLocale } from "@/config";
import type { StoryEx } from "@/modules/backend/types";

/**
 * Generate markdown for a global story
 */
export function generateGlobalStoryMarkdown(story: StoryEx, locale: string): string {
  // Try to get date from slug first, fall back to created_at
  const publishDate = story.slug !== null ? parseDateFromSlug(story.slug) : null;
  const dateToFormat = publishDate ?? new Date(story.created_at);

  const frontmatter = formatFrontmatter({
    title: story.title ?? "Untitled",
    author: story.author_profile?.title ?? "Unknown",
    publish_date: formatDateShort(dateToFormat, locale),
    reading_time: `${calculateReadingTime(story.content)} min`,
    kind: story.kind,
    status: story.publications.length > 0 ? "published" : "draft",
  });

  return `${frontmatter}\n\n${story.content}`;
}

/**
 * Generate markdown for the stories listing page
 */
export function generateStoriesListingMarkdown(
  stories: StoryEx[],
  locale: string,
): string {
  const frontmatter = formatFrontmatter({
    title: "Stories",
    generated_at: new Date().toISOString(),
  });

  const storyLinks = stories.map((story) => {
    const title = story.title ?? "Untitled";
    const slug = story.slug ?? "";
    const summary = story.summary ?? "";
    const author = story.author_profile?.title ?? "";
    const storyLocale = story.locale_code?.trim() ?? "";

    let line = `- [${title}](/${locale}/stories/${slug}.md)`;

    // Add locale badge when story language differs from viewer's locale
    if (storyLocale !== "" && storyLocale !== locale && isValidLocale(storyLocale)) {
      const localeData = supportedLocales[storyLocale as SupportedLocaleCode];
      line += ` [${localeData.englishName}]`;
    }

    if (author !== "") {
      line += ` by ${author}`;
    }
    if (summary !== "") {
      line += `\n  ${summary}`;
    }
    return line;
  });

  return `${frontmatter}\n\n${storyLinks.join("\n\n")}`;
}

/**
 * Register markdown handler for global stories listing
 * Pattern: /$locale/stories
 */
export function registerGlobalStoriesListingHandler(): void {
  registerMarkdownHandler("$locale/stories", async (_params, locale, _searchParams) => {
    // Match the HTML route: only articles
    const stories = await backend.getStoriesByKinds(locale, [
      "article",
    ]);

    if (stories === null) {
      return null;
    }

    return generateStoriesListingMarkdown(stories, locale);
  });
}

/**
 * Register markdown handler for individual global stories
 * Pattern: /$locale/stories/$storyslug
 */
export function registerGlobalStoryMarkdownHandler(): void {
  registerMarkdownHandler("$locale/stories/$storyslug", async (params, locale, _searchParams) => {
    const { storyslug } = params;

    if (storyslug === undefined) {
      return null;
    }

    const story = await backend.getStory(locale, storyslug);

    if (story === null) {
      return null;
    }

    return generateGlobalStoryMarkdown(story, locale);
  });
}
